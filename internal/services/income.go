package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type IncomeServiceInterface[T models.Income] interface {
	GetAllIncomes() error
	GetSum() float64
	DeleteAll() error
	GetItems() []T
	GetSortedIncomes() []models.Income
}

type IncomeService struct {
	MonetaryService[models.Income]
	incomes *[]models.Income
}

func NewIncomeService(repo *store.SqlDB, incomes *[]models.Income) *IncomeService {
	return &IncomeService{
		MonetaryService: *NewMonetaryService[models.Income](repo, models.Income{}, incomes),
		incomes:         incomes,
	}
}

func (s *IncomeService) GetAllIncomes() error {
	err := store.GetAll(s.GetRepo(), &s.incomes)
	if err != nil {
		return err
	}
	return nil
}

func (s *IncomeService) GetSortedIncomes() []models.Income {
	sortedIncomes := s.GetItems()
	models.SortIncomeByDate(&sortedIncomes)
	return sortedIncomes
}
