package services

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
)

type ServiceInterface interface {
	IncomeService | ExpenseService | AllocationService | NoteService
}

type Service struct {
	repo    *store.SqlDB
	element any
}

func NewService(repo *store.SqlDB, element any) *Service {
	return &Service{
		repo:    repo,
		element: element,
	}
}

type MonetaryService[T models.MonetaryItemInterface] struct {
	Service
	items *[]T
}

func NewMonetaryService[T models.MonetaryItemInterface](repo *store.SqlDB, element any, items *[]T) *MonetaryService[T] {
	return &MonetaryService[T]{
		Service: *NewService(repo, element),
		items:   items,
	}
}

func (s *Service) GetRepo() *store.SqlDB {
	return s.repo
}

func (s *Service) GetElementType() any {
	return &s.element
}

func (s *MonetaryService[T]) GetSum() float64 {
	return store.GetSum(s.GetRepo(), s.GetElementType(), "amount")
}

func (s *MonetaryService[T]) DeleteAll() error {
	for _, item := range *s.items {
		err := store.Delete(s.GetRepo(), item.GetID(), s.GetElementType())
		if err != nil {
			return err
		}
	}
	return nil
}
