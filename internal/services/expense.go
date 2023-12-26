package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type ExpenseServiceInterface interface {
	GetAllExpenses() error
	GetSum() float64
	DeleteAll() error
}

type ExpenseService struct {
	MonetaryService[models.Expense]
	expenses *[]models.Expense
}

func NewExpenseService(repo *store.SqlDB, expenses *[]models.Expense) ExpenseServiceInterface {
	return &ExpenseService{
		MonetaryService: *NewMonetaryService[models.Expense](repo, models.Income{}, expenses),
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
