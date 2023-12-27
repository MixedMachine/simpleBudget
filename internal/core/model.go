package core

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/services"
	"github.com/mixedmachine/simple-budget-app/internal/store"
	"github.com/mixedmachine/simple-budget-app/internal/utils"
	"os"
	"path/filepath"
)

type SimpleBudget struct {
	App               fyne.App
	Window            fyne.Window
	Repo              *store.SqlDB
	IncomeService     services.IncomeServiceInterface[models.Income]
	ExpenseService    services.ExpenseServiceInterface[models.Expense]
	AllocationService services.AllocationServiceInterface[models.Allocation]
	NoteService       services.NoteServiceInterface
	Notes             models.Notes
	ListComponents    map[string]*widget.List
	LabelComponents   map[string]*canvas.Text
}

func (s *SimpleBudget) SetIcon() {
	iconPath := ""
	if _, err := os.Stat("assets/icon.png"); os.IsExist(err) {
		iconPath = "assets/icon.png"
	} else if _, err := os.Stat("icon.png"); os.IsExist(err) {
		iconPath = "icon.png"
	}
	if iconPath != "" {
		if resourceIconPng, err := fyne.LoadResourceFromPath("assets/icon.png"); err == nil {
			s.Window.SetIcon(resourceIconPng)
		}
	}
}

func (s *SimpleBudget) SetUpRepo() {
	s.Repo = store.NewSqlDB(
		store.InitializeSQL(
			store.SQLITE, filepath.Join(
				s.App.Storage().RootURI().Path(), store.SQLITE_FILE,
			),
		),
	)
}

func (s *SimpleBudget) SetUpServices() {
	var err error

	s.IncomeService = services.NewIncomeService(s.Repo, models.NewIncomes())
	err = s.IncomeService.GetAllIncomes()
	utils.HandleErr(s.Window, err)

	s.ExpenseService = services.NewExpenseService(s.Repo, models.NewExpenses())
	err = s.ExpenseService.GetAllExpenses()
	utils.HandleErr(s.Window, err)

	s.AllocationService = services.NewAllocationService(s.Repo, models.NewAllocations())
	err = s.AllocationService.GetAllAllocations()
	utils.HandleErr(s.Window, err)

	s.NoteService = services.NewNoteService(s.Repo)
	s.Notes, err = s.NoteService.GetNotes()
	utils.HandleErr(s.Window, err)
}
