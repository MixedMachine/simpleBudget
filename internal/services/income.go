package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type IncomeServiceInterface[T models.Income] interface {
	GetAllIncomes() error
	GetSum() float64
	GetFilteredSum(query string, args ...interface{}) float64
	DeleteAll() error
	GetItems() *[]T
	CreateItem(item T) error
	UpdateItem(item T) error
	DeleteItem(item T) error
	GetSortedIncomes() []models.Income
}

type IncomeService struct {
	MonetaryService[models.Income]
}

func NewIncomeService(repo *store.SqlDB, incomes *[]models.Income) *IncomeService {
	return &IncomeService{
		MonetaryService: *NewMonetaryService[models.Income](repo, models.Income{}, incomes),
	}
}

func (s *IncomeService) GetAllIncomes() error {
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
