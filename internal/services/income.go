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
	MonetaryService
	incomes *[]models.Income
}

func NewIncomeService(repo *store.SqlDB, incomes *[]models.Income) IncomeServiceInterface {
	return &IncomeService{
		MonetaryService: *NewMonetaryService(repo, models.Income{}),
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

func (s *IncomeService) DeleteAll() error {
	for _, income := range *s.incomes {
		err := store.Delete(s.GetRepo(), income.ID, s.GetElementType())
		if err != nil {
			return err
		}
	}
	return nil
}
