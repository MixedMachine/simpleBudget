package services

import (
	"errors"
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
	"gorm.io/gorm"
)

type NoteServiceInterface interface {
	GetNotes() (models.Notes, error)
}

type NoteService struct {
	Service
}

func NewNoteService(repo *store.SqlDB) NoteServiceInterface {
	return &NoteService{
		Service: *NewService(repo, models.Notes{}),
	}
}

func (s *NoteService) GetNotes() (models.Notes, error) {
	var notes models.Notes
	err := store.GetAll(s.GetRepo(), &notes)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = store.Create(s.GetRepo(), &notes)
		if err != nil {
			return models.Notes{}, err
		}
		err = store.GetAll(s.GetRepo(), &notes)
	}
	if err != nil {
		return models.Notes{}, err
	}
	return notes, nil
}
