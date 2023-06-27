package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var Version string

func main() {
	config := NewSSHConfig()

	if Version == "" {
		// todo temp workaround for broken fyne ldflags
		Version = "1.0.0"
	}
	name := "Port Punch"

	a := app.NewWithID("com.nbembedded.portpunch")

	title := fmt.Sprintf("%s %s", name, Version)
	w := a.NewWindow(title)

	// context for ssh command
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// When the window is closing, cancel the context
	w.SetCloseIntercept(func() {
		cancel()
		w.Close()
	})

	msgContainer := container.NewVBox()

	form := NewForm(a, w, ctx, msgContainer, config)

	w.SetContent(container.NewVBox(
		form,
		msgContainer,
	))

	w.SetFixedSize(true)

	// Set the window dimensions
	w.Resize(fyne.NewSize(480.0, 300.0))

	w.ShowAndRun()
}
