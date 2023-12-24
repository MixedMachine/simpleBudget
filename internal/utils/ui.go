package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	log "github.com/sirupsen/logrus"
)

func HandleErr(window fyne.Window, err error) {
	if err != nil {
		log.Error(err)
		errBox := dialog.NewError(err, window)
		errBox.Show()
	}
}
