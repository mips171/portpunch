package main

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func RunSSHCommand(ctx context.Context, a fyne.App, w fyne.Window, msgContainer *fyne.Container, config *sshConfig) error {
	cmd := fmt.Sprintf("-f -N -L %s:%s:%s %s@%s -f", config.localPort, config.targetHost, config.targetPort, config.hostUser, config.sshHost)

	var command *exec.Cmd
	var err error

	if runtime.GOOS == "windows" {
		command = exec.CommandContext(ctx, "cmd", "/C", "start /B ssh "+cmd)
	} else {
		args := strings.Split(cmd, " ")
		command = exec.CommandContext(ctx, "ssh", args...)
	}

	done := make(chan error, 1)

	go func() {
		err = command.Run()
		done <- err
	}()

	childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-childCtx.Done():
		err = childCtx.Err()
	case err = <-done:
	}

	if err != nil {
		return err
	}

	if config.targetPort == "22" {
		sshCommand := fmt.Sprintf("ssh root@localhost -p %s", config.localPort)
		sshEntry := widget.NewLabel(sshCommand)

		copyButton := widget.NewButton("Copy", func() {
			w.Clipboard().SetContent(sshCommand)
		})

		instructionLabel := widget.NewLabel("Paste that in a terminal")
		instructionLabel.Wrapping = fyne.TextWrapWord

		msgContainer.Objects = []fyne.CanvasObject{instructionLabel, sshEntry, copyButton}
		msgContainer.Refresh()
		return nil
	}

	proto := "http"
	if config.targetPort == "443" {
		proto = "https"
	}

	urlStr := fmt.Sprintf("%s://localhost:%s", proto, config.localPort)
	parsedUrl, err := url.Parse(urlStr)

	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
	}

	link := widget.NewHyperlink("Open in browser", parsedUrl)
	link.OnTapped = func() {
		err := a.OpenURL(parsedUrl)
		if err != nil {
			return
		}
	}

	msgContainer.Objects = []fyne.CanvasObject{link}
	msgContainer.Refresh()

	return nil
}
