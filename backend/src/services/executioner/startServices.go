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

var runningProjects domain.SecureMap[int, *executioner.RunningProject]

// esto es util para el futuro websocket
var OutputChannels domain.SecureMap[int, chan string]

func StartServices() {
	projects, err := pRepo.Search("")
	if err != nil {
		log.Println(err) // deberia aqui de decir que hubo un error en este caso detendria todo?
	}

	for _, project := range projects {
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

	cmd := exec.CommandContext(ctx, "bash", "-c", project.Command)
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
	go OutReader(project.ID, stdout, channel)
	go ErrReader(project.ID, stderr, channel)
	runningProjects.Set(project.ID, &executioner.RunningProject{
		Cmd:    cmd,
		Cancel: cancel,
	})

	cmd.Start()
	go func() {
		err := cmd.Wait()
		OutputChannels.Delete(project.ID)
		runningProjects.Delete(project.ID)
		close(channel)

		if ctx.Err() == context.Canceled {
			return
		}
		if err != nil {
			RestartProject(project.ID)
		}
	}()
}
func OutReader(id int, buf io.ReadCloser, channel chan string) {
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		// aqui podriamos decir
		output := scanner.Text()
		fmt.Println(output)
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

func ErrReader(id int, buf io.ReadCloser, channel chan string) {
	scanner := bufio.NewScanner(buf)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		output := scanner.Text()
		fmt.Println(output)
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
