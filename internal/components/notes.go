package components

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

func CreateNotesComponent(
	myWindow fyne.Window,
	repo *(store.SqlDB),
	notes *models.Notes,
) *fyne.Container {
	autoSave := false
	notesTitle := canvas.NewText("Notes", theme.ForegroundColor())
	notesTitle.TextSize = 20
	notesTitle.TextStyle = fyne.TextStyle{Bold: true}

	notesHeader := container.New(layout.NewHBoxLayout(),
		notesTitle,
		widget.NewCheckWithData("Auto Save", binding.BindBool(&autoSave)),
	)

	notesEntry := widget.NewMultiLineEntry()

	notesEntry.OnChanged = func(s string) {
		if autoSave {
			(*notes).Content = s
			err := store.Update(repo, (*notes).ID, (*notes))
			if err != nil {
				log.Error(err)
				errBox := dialog.NewError(err, myWindow)
				errBox.Show()
			}
		}
	}

	notesEntry.SetText((*notes).Content)

	saveNoteBtn := widget.NewButton("Save", func() {
		(*notes).Content = notesEntry.Text
		err := store.Update(repo, (*notes).ID, *notes)
		if err != nil {
			log.Error(err)
			errBox := dialog.NewError(err, myWindow)
			errBox.Show()
		}
	})

	return container.NewBorder(
		notesHeader,
		saveNoteBtn,
		nil,
		nil,
		notesEntry,
	)

}
