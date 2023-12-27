package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type ExpenseServiceInterface[T models.Expense] interface {
	GetAllExpenses() error
	GetSum() float64
	DeleteAll() error
	GetItems() []T
	GetSortedExpenses() []models.Expense
}

type ExpenseService struct {
	MonetaryService[models.Expense]
	expenses *[]models.Expense
}

func NewExpenseService(repo *store.SqlDB, expenses *[]models.Expense) *ExpenseService {
	return &ExpenseService{
		MonetaryService: *NewMonetaryService[models.Expense](repo, models.Expense{}, expenses),
		expenses:        expenses,
	}
}

func (s *ExpenseService) GetAllExpenses() error {
	err := store.GetAll(s.repo, &s.expenses)
	if err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetSortedExpenses() []models.Expense {
	sortedExpenses := s.GetItems()
	models.SortExpenseByDate(&sortedExpenses)
	return sortedExpenses
}
