package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type IncomeServiceInterface interface {
	GetAllIncomes() error
	GetSum() float64
	DeleteAll() error
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
