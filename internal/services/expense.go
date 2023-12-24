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
	MonetaryService
	expenses *[]models.Expense
}

func NewExpenseService(repo *store.SqlDB, expenses *[]models.Expense) ExpenseServiceInterface {
	return &ExpenseService{
		MonetaryService: *NewMonetaryService(repo, models.Income{}),
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

func (s *ExpenseService) DeleteAll() error {
	for _, expense := range *s.expenses {
		err := store.Delete(s.GetRepo(), expense.ID, s.GetElementType())
		if err != nil {
			return err
		}
	}
	return nil
}
