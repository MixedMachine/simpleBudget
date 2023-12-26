package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type AllocationServiceInterface interface {
	GetAllAllocations() error
	GetSum() float64
	DeleteAll() error
}

type AllocationService struct {
	MonetaryService[models.Allocation]
	allocations *[]models.Allocation
}

func NewAllocationService(repo *store.SqlDB, allocations *[]models.Allocation) AllocationServiceInterface {
	return &AllocationService{
		MonetaryService: *NewMonetaryService[models.Allocation](repo, models.Income{}, allocations),
		allocations:     allocations,
	}
}

func (s *AllocationService) GetAllAllocations() error {
	err := store.GetAll(s.repo, &s.allocations)
	if err != nil {
		return err
	}
	return nil
}
