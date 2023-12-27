package models

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/mixedmachine/simple-budget-app/internal/services"
	"github.com/mixedmachine/simple-budget-app/internal/store"
	"github.com/mixedmachine/simple-budget-app/internal/utils"
	"os"
	"path/filepath"
	"time"
)

type SimpleBudget struct {
	App      fyne.App
	Window   fyne.Window
	Repo     *store.SqlDB
	IncomeService services.IncomeServiceInterface[Income]
	ExpenseService services.ExpenseServiceInterface[Expense]
	AllocationService services.AllocationServiceInterface[Allocation]
	NoteService services.NoteServiceInterface
	Notes Notes
	ListComponents map[string]*widget.List
	LabelComponents map[string]*canvas.Text
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

	s.IncomeService = services.NewIncomeService(s.Repo, NewIncomes())
	err = s.IncomeService.GetAllIncomes()
	utils.HandleErr(s.Window, err)

	s.ExpenseService = services.NewExpenseService(s.Repo, NewExpenses())
	err = s.ExpenseService.GetAllExpenses()
	utils.HandleErr(s.Window, err)

	s.AllocationService = services.NewAllocationService(s.Repo, NewAllocations())
	err = s.AllocationService.GetAllAllocations()
	utils.HandleErr(s.Window, err)

	s.NoteService = services.NewNoteService(s.Repo)
	s.Notes, err = s.NoteService.GetNotes()
	utils.HandleErr(s.Window, err)
}


type MonetaryItemInterface interface {
	Income | Expense | Allocation
	GetID() uint
}
type MonetaryItem struct {
	ID     uint    `gorm:"primaryKey;autoIncrement"`
	Amount float64 `gorm:"type:decimal(10,2);default:0.00;not null"`
}

func (m MonetaryItem) GetID() uint {
	return m.ID
}

type TransactionItem struct {
	MonetaryItem
	Name string    `gorm:"unique;not null"`
	Date time.Time `gorm:"type:date;not null"`
}

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"unique"`
	Password string `gorm:"not null"`
}

type Income struct {
	TransactionItem
}

type Expense struct {
	TransactionItem
}

type Allocation struct {
	MonetaryItem
	FromIncomeID uint `gorm:"index:idx_from_income_id;foreignKey:FromIncomeID"`
	ToExpenseID  uint `gorm:"index:idx_to_expense_id;foreignKey:ToExpenseID"`
}

type Notes struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Content string `gorm:"not null"`
}

func NewIncomes() *[]Income {
	return &[]Income{}
}

func NewExpenses() *[]Expense {
	return &[]Expense{}
}

func NewAllocations() *[]Allocation {
	return &[]Allocation{}
}

func NewNotes() *Notes {
	return &Notes{}
}
