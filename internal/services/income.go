package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
	"time"
)

type IncomeServiceInterface[T models.Income] interface {
	RefreshIncomes() error
	GetSum() float64
	GetFilteredSum(query string, args ...interface{}) float64
	DeleteAll() error
	GetItems() *[]T
	CreateItem(item T) error
	UpdateItem(item T) error
	DeleteItem(item T) error
	GetSortedIncomes() []models.Income
	FilterIncomesBeforeDate(date time.Time) []models.Income
	GetIncomeByName(incomeName string) models.Income
}

type IncomeService struct {
	MonetaryService[models.Income]
}

func NewIncomeService(repo *store.SqlDB, incomes *[]models.Income) *IncomeService {
	return &IncomeService{
		MonetaryService: *NewMonetaryService[models.Income](repo, models.Income{}, incomes),
	}
}

func (s *IncomeService) RefreshIncomes() error {
	if err := store.GetAll(s.GetRepo(), s.GetItems()); err != nil {
		return err
	}
	return nil
}

func (s *IncomeService) GetSortedIncomes() []models.Income {
	sortedIncomes := s.GetItems()
	models.SortIncomeByDate(sortedIncomes)
	return *sortedIncomes
}

func (s *IncomeService) FilterIncomesBeforeDate(date time.Time) []models.Income {
	var filteredIncome []models.Income
	for _, income := range *s.GetItems() {
		if income.Date.Before(date) {
			filteredIncome = append(filteredIncome, income)
		}
	}
	return filteredIncome
}

func (s *IncomeService) GetIncomeByName(incomeName string) models.Income {
	for _, income := range *s.GetItems() {
		if income.Name == incomeName {
			return income
		}
	}
	return models.Income{}
}
