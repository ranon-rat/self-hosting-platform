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

	"github.com/ranon-rat/self-hosting-manager/src/domain"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executioner"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

var runningProjects = domain.NewSecureMap[int, *executioner.RunningProject]()

// esto es util para el futuro websocket
var OutputChannels = domain.NewSecureMap[int, chan string]()
var deletingAll = false

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
		fmt.Println("dir:\n", project.Dir)
		fmt.Println("command:", project.Command)
		fmt.Println(divider)
		Executioner(&project)
	}

}
func StopProject(id int) error {
	if err := pRepo.PauseProject(true, id); err != nil {
		return err
	}
	cmd, exists := runningProjects.Get(id)
	if !exists {
		return nil
	}
	// Primero cancela el contexto (cierre gracioso)
	stopCmd(cmd.Cmd)
	runningProjects.Delete(id)
	return nil
}
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
func RestartProject(id int) error {
	if deletingAll {
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
		stopCmd(rp.Cmd)
		runningProjects.Delete(project.ID)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "bash", "-lc" /*"trap 'kill 0' SIGTERM; "+*/, project.Command)
	cmd.Dir = project.Dir
	cmd.Env = executableEnv

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
	channel := make(chan string, executioner.MAX_CHANNEL_BUFFER)
	OutputChannels.Set(project.ID, channel)
	lastErrOutput := domain.NewSecureStrContainer()
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
	runningProjects.Set(project.ID, &executioner.RunningProject{
		Cmd:    cmd,
		Cancel: cancel,
	})

	if err := cmd.Start(); err != nil {
		stopCmd(cmd)
		return
	}
	go func() {
		err := cmd.Wait()
		runningProjects.Delete(project.ID)
		wg.Wait()
		lastErrStr := strings.ToLower(lastErrOutput.Content())
		fmt.Println(lastErrStr)
		if ctx.Err() == context.Canceled {
			SaveAndSend(channel, err.Error(), project.ID)
			goto justClean
		}
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
		close(channel)
		OutputChannels.Delete(project.ID)
		RestartProject(project.ID)
		return
	justClean:
		close(channel)
		OutputChannels.Delete(project.ID)
	}()
}
func OutReader(id int, name string, buf io.ReadCloser, channel chan string, lastErr *domain.SecureStringContainer) {

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
		lastErr.FromStdout()

	}
}

func ErrReader(id int, name string, buf io.ReadCloser, channel chan string, lastErr *domain.SecureStringContainer) {
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
		fmt.Println("from error", output)
		lastErr.AppendValue(output)
		lastErr.FromStderr()
		SaveAndSend(channel, output, id)

		// i just want to know what is happening this is for testing
		//		fmt.Println(name, output)

	}
}
func SaveAndSend(channel chan string, output string, projectID int) {
	logRepo.Create(&executionlogs.NewLog{
		IdProject: projectID,
		Content:   output,
	})
	select {
	case channel <- output:
	default:
	}
}

func StoppingAll() {
	deletingAll = true
	runningProjects.Range(func(i int, rp *executioner.RunningProject) bool {
		stopCmd(rp.Cmd)
		return true
	})
}
