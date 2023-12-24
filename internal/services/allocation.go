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
	MonetaryService
	allocations *[]models.Allocation
}

func NewAllocationService(repo *store.SqlDB, allocations *[]models.Allocation) AllocationServiceInterface {
	return &AllocationService{
		MonetaryService: *NewMonetaryService(repo, models.Income{}),
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

func (s *AllocationService) DeleteAll() error {
	for _, allocations := range *s.allocations {
		err := store.Delete(s.GetRepo(), allocations.ID, s.GetElementType())
		if err != nil {
			return err
		}
	}
	return nil
}
