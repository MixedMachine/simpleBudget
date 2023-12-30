package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
	"time"
)

type ExpenseServiceInterface[T models.Expense] interface {
	RefreshExpenses() error
	GetSum() float64
	GetFilteredSum(query string, args ...interface{}) float64
	DeleteAll() error
	GetItems() *[]T
	CreateItem(item T) error
	UpdateItem(item T) error
	DeleteItem(item T) error
	GetSortedExpenses() []models.Expense
	FilterExpensesAfterDate(date time.Time) []models.Expense
	GetExpenseByName(expenseName string) models.Expense
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

func (s *ExpenseService) RefreshExpenses() error {
	err := store.GetAll(s.repo, s.GetItems())
	if err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetSortedExpenses() []models.Expense {
	sortedExpenses := s.GetItems()
	models.SortExpenseByDate(sortedExpenses)
	return *sortedExpenses
}

func (s *ExpenseService) FilterExpensesAfterDate(date time.Time) []models.Expense {
	var filteredExpenses []models.Expense
	for _, expense := range *s.GetItems() {
		if expense.Date.After(date) {
			filteredExpenses = append(filteredExpenses, expense)
		}
	}

	return filteredExpenses
}

func (s *ExpenseService) GetExpenseByName(expenseName string) models.Expense {
	for _, expense := range *s.GetItems() {
		if expense.Name == expenseName {
			return expense
		}
	}
	return models.Expense{}
}
