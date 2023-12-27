package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type AllocationServiceInterface[T models.Allocation] interface {
	GetAllAllocations() error
	GetSum() float64
	DeleteAll() error
	GetItems() []T
	GetSortedAllocations() []models.Allocation
}

type AllocationService struct {
	MonetaryService[models.Allocation]
	allocations *[]models.Allocation
}

func NewAllocationService(repo *store.SqlDB, allocations *[]models.Allocation) *AllocationService {
	return &AllocationService{
		MonetaryService: *NewMonetaryService[models.Allocation](repo, models.Allocation{}, allocations),
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

func (s *AllocationService) GetSortedAllocations() []models.Allocation {
	sortedAllocations := s.GetItems()
	models.SortAllocationsByAmount(&sortedAllocations)
	return sortedAllocations
}
