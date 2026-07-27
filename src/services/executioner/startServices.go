package executionerServices

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ranon-rat/self-hosting-manager/src/domain"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executioner"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

var runningProjects = domain.NewSecureMap[int, *executioner.RunningProject]()
var OutputChannels = domain.NewSecureMap[int, *domain.SafeChan[string]]()
var deletingAll atomic.Bool

// this is for starting the services that are stored in our database
func StartServices() {
	projects, err := pRepo.Search("")
	if err != nil {
		log.Println(err) // deberia aqui de decir que hubo un error en este caso detendria todo?
		return
	}
	divider := strings.Repeat("-", 50)
	fmt.Println(divider)
	for _, project := range projects {
		fmt.Println("Starting project", project.Name)
		fmt.Printf("dir:\n%s\n", project.Dir)
		fmt.Println("command:", project.Command)
		fmt.Println(divider)
		Executioner(&project)
	}

}
func DeleteProject(id int) error {
	StopProject(id)
	runningProjects.Delete(id)
	output, exist := OutputChannels.Get(id)
	if !exist {
		return nil
	}
	output.Close()
	return nil
}

// this stops the project directly
func StopProject(id int) error {
	if err := pRepo.PauseProject(true, id); err != nil {
		return err
	}
	cmd, exists := runningProjects.Get(id)
	if !exists {
		return nil
	}
	cmd.Cancel()
	stopCmd(cmd.Cmd)
	runningProjects.Delete(id)
	return nil
}

// this starts the project, it changes tha paused status and then it runs it
func StartProject(id int) error {
	if err := pRepo.PauseProject(false, id); err != nil {
		return err
	}
	_, exists := runningProjects.Get(id)
	if exists {
		return nil
	}
	project, err := pRepo.GetByID(id)
	if err != nil {
		return err
	}
	Executioner(project)
	return nil

}

// this is for restarting the project in case of any error
func RestartProject(id int) error {
	if deletingAll.Load() {
		return nil
	}
	project, err := pRepo.GetByID(id)
	if err != nil {
		return err
	}
	if project.IsPaused {
		return fmt.Errorf("user: you cannot restart a paused project")
	}

	Executioner(project)
	return nil
}

func Executioner(project *projectsD.Project) {
	if project.IsPaused {
		return
	}
	if rp, e := runningProjects.Get(project.ID); e {
		rp.Cmd.Cancel()
		stopCmd(rp.Cmd)
		runningProjects.Delete(project.ID)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "bash", "-lc" /*"trap 'kill 0' SIGTERM; "+*/, project.Command)
	cmd.Dir = project.Dir
	cmd.Env = executableEnv
	// i have to configure the system process attributes to avoid any weird behaviour
	// this function is empty on windows but on linux and unix derived systems
	// it sets the children process group id, for it to be used. later on the stop cmd
	setSysProcAttr(cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return
	}
	channel, exists := OutputChannels.Get(project.ID)
	if !exists {
		// since i am closing the channel, this allows me to avoid any weird behaviour
		channel = domain.NewSafeChan[string](executioner.MAX_CHANNEL_BUFFER)
		OutputChannels.Set(project.ID, channel)
	}
	// this is constantly being written and readed, so it uses a rwmutex
	// it also helps to store the values and to also know and debug what happened before
	// since golang sends the log.Println to the stderr
	// i need to use this to avoid any weird shit later on
	// specially with the port binding
	lastErrOutput := executioner.NewSecureErrContainer(executioner.MAX_LINES_DEFAULT)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		OutReader(project.ID, project.Name, stdout, channel, lastErrOutput)
		wg.Done()
	}()
	go func() {
		ErrReader(project.ID, project.Name, stderr, channel, lastErrOutput)
		wg.Done()
	}()
	if err := cmd.Start(); err != nil {
		cancel()
		stopCmd(cmd)
		return
	}

	runningProjects.Set(project.ID, &executioner.RunningProject{
		Cmd:    cmd,
		Cancel: cancel,
	})
	go func() {
		// in case of any interruption i am using this
		err := cmd.Wait()
		// it deletes the running project
		runningProjects.Delete(project.ID)
		wg.Wait()
		lastErrStr := strings.ToLower(lastErrOutput.Content())
		fmt.Println(lastErrStr)
		if ctx.Err() == context.Canceled {
			//
			SaveAndSend(channel, err.Error(), project.ID)
			goto justClean
		}
		// i am managing the port binding just after this. I should have something to manage any weird behaviour
		if err != nil {

			SaveAndSend(channel, err.Error(), project.ID)
			if strings.Contains(strings.ToLower(err.Error()), "bind") {
				log.Println("Port already in use, not restarting", project.Name)
				goto justClean
			}
		}
		if strings.Contains(lastErrStr, "bind") {
			log.Println("Port already in use, not restarting", project.Name)
			goto justClean
		}
		time.Sleep(1 * time.Second)
		RestartProject(project.ID)
		return
		// the just clean part is to avoid getting inside an infite loop
	justClean:
		// bueno aqui no deberia de haber un problema pero
		channel.Close()
		OutputChannels.Delete(project.ID)
	}()
}
func OutReader(id int, name string, buf io.ReadCloser, channel *domain.SafeChan[string], lastErr *executioner.SecureErrContainer) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in ErrReader for %s: %v\n", name, r)
		}
	}()
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		// aqui podriamos decir
		output := scanner.Text()
		//	fmt.Println(name, output)
		SaveAndSend(channel, output, id)
		// here i am just notifying that i came from the stdout
		lastErr.FromStdout()

	}
}

func ErrReader(id int, name string, buf io.ReadCloser, channel *domain.SafeChan[string], lastErr *executioner.SecureErrContainer) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in ErrReader for %s: %v\n", name, r)
		}
	}()
	scanner := bufio.NewScanner(buf)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		// i am just trying to check something i dont want to erase anything weird
		if lastErr.ComingFromStdOut() {
			lastErr.Clean()
		}
		output := scanner.Text()
		// so i am adding something new to the append value
		lastErr.AppendValue(output)
		SaveAndSend(channel, output, id)

		// i just want to know what is happening this is for testing

	}
}
func SaveAndSend(channel *domain.SafeChan[string], output string, projectID int) {
	logRepo.Create(&executionlogs.NewLog{
		IdProject: projectID,
		Content:   output,
	})
	channel.AppendSend(output)
}

// i need to manage the way i am executing this
func StoppingAll(ctx context.Context) {
	deletingAll.Store(true)
	var wg sync.WaitGroup
	runningProjects.Range(func(i int, rp *executioner.RunningProject) bool {
		wg.Add(1)
		defer wg.Done()
		rp.Cancel()
		stopCmd(rp.Cmd)
		return true
	})
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		log.Println("All processes stopped")
	case <-ctx.Done():
		log.Println("Shutdown interrupted by context cancel")
	}
}
