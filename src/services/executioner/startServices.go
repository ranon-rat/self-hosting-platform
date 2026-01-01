package executionerServices

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/ranon-rat/self-hosting-manager/src/domain"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executioner"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

var runningProjects = domain.NewSecureMap[int, *executioner.RunningProject]()

// esto es util para el futuro websocket
var OutputChannels = domain.NewSecureMap[int, chan string]()

func StartServices() {
	projects, err := pRepo.Search("")
	if err != nil {
		log.Println(err) // deberia aqui de decir que hubo un error en este caso detendria todo?
		return
	}

	for _, project := range projects {
		fmt.Println("Starting project", project.Name)
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
	cmd.Cancel()
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
	project, err := pRepo.GetByID(id)
	if err != nil {
		return err
	}
	if project.IsPaused {
		return fmt.Errorf("user: you cannot restart a paused project")
	}

	cmd, exists := runningProjects.Get(id)
	if exists {
		cmd.Cancel()
	}
	Executioner(project)
	return nil
}
func Executioner(project *projectsD.Project) {
	if project.IsPaused {
		return
	}
	if rp, e := runningProjects.Get(project.ID); e {
		rp.Cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "bash", "-lc", fmt.Sprintf("cd \"%s\";", project.Dir)+project.Command)
	cmd.Dir = project.Dir
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
	go OutReader(project.ID, project.Name, stdout, channel)
	go ErrReader(project.ID, project.Name, stderr, channel)
	runningProjects.Set(project.ID, &executioner.RunningProject{
		Cmd:    cmd,
		Cancel: cancel,
	})

	cmd.Start()
	go func() {
		defer close(channel)
		err := cmd.Wait()
		OutputChannels.Delete(project.ID)
		runningProjects.Delete(project.ID)

		if ctx.Err() == context.Canceled {
			return
		}
		if err != nil {
			RestartProject(project.ID)
		}
	}()
}
func OutReader(id int, name string, buf io.ReadCloser, channel chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in ErrReader for %s: %v\n", name, r)
		}
	}()
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		// aqui podriamos decir
		output := scanner.Text()
		fmt.Println(name, output)
		logRepo.Create(&executionlogs.NewLog{
			IdProject: id,
			Content:   output,
		})
		select {
		case channel <- output:
		default:
		}
	}
}

func ErrReader(id int, name string, buf io.ReadCloser, channel chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in ErrReader for %s: %v\n", name, r)
		}
	}()
	scanner := bufio.NewScanner(buf)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		output := scanner.Text()
		// i just want to know what is happening this is for testing
		fmt.Println(name, output)
		logRepo.Create(&executionlogs.NewLog{
			IdProject: id,
			Content:   output,
		})
		select {
		case channel <- output:
		default:
		}

	}
}
