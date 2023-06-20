package main

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewForm(a fyne.App, w fyne.Window, ctx context.Context, msgContainer *fyne.Container, config *sshConfig) *widget.Form {

	prefs := a.Preferences()

	sshServerPref := prefs.String("preference.sshserver")
	hostUserPref := prefs.String("preference.hostuser")

	inputSShServer := widget.NewEntry()
	inputSShServer.SetPlaceHolder("Enter SSH Server")
	inputSShServer.SetText(sshServerPref)

	inputHostUser := widget.NewEntry()
	inputHostUser.SetPlaceHolder("Enter Host Username")
	inputHostUser.SetText(hostUserPref)

	inputTargetHost := newIPEntry()
	inputTargetHost.SetPlaceHolder("Enter Target Host")

	optionMap := map[string]string{
		"Web (HTTP - 80)":   "80",
		"Web (HTTPS - 443)": "443",
		"SSH (CLI - 22)":    "22",
	}

	options := make([]string, 0, len(optionMap))
	for option := range optionMap {
		options = append(options, option)
	}

	inputTargetPort := widget.NewSelect(options, func(label string) {
		config.targetPort = optionMap[label]
	})
	inputTargetPort.PlaceHolder = "Select Target Port"
	inputTargetPort.SetSelected("Web (HTTPS - 443)")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Jump Server Address", Widget: inputSShServer},
			{Text: "Target Host Username", Widget: inputHostUser},
			{Text: "Target Host IP", Widget: inputTargetHost},
			{Text: "Target Host Port", Widget: inputTargetPort},
		},
		SubmitText: "Connect",
		OnSubmit: func() {
			matched, _ := regexp.MatchString(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, inputTargetHost.Text)
			if !matched {
				errLabel := widget.NewLabel("Invalid IP address.")
				errLabel.Wrapping = fyne.TextWrapWord
				msgContainer.Objects = []fyne.CanvasObject{errLabel}
				msgContainer.Refresh()
				return
			}
			rand.NewSource(time.Now().UnixNano())
			config.localPort = fmt.Sprintf("%d", rand.Intn(60000-50000+1)+50000)

			config.targetHost = inputTargetHost.Text

			prefs.SetString("preference.sshserver", inputSShServer.Text)
			prefs.SetString("preference.hostuser", inputHostUser.Text)

			config.hostUser = inputHostUser.Text
			config.sshHost = inputSShServer.Text

			err := RunSSHCommand(ctx, a, w, msgContainer, config)
			if err != nil {
				errLabel := widget.NewLabel("Check that your SSH keys are set up correctly and that you have an internet connection. Error: " + err.Error())
				errLabel.Wrapping = fyne.TextWrapWord
				msgContainer.Objects = []fyne.CanvasObject{errLabel}
				msgContainer.Refresh()
			}
		},
	}

	return form
}
