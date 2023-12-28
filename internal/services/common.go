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

type MonetaryServiceInterface[T models.MonetaryItemInterface] interface {
	GetElementExample() any
	GetItems() []T
	GetRepo() *store.SqlDB
	GetSum() float64
	DeleteAll() error
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

func (s *Service) GetElementExample() any {
	return &s.element
}

func (s *MonetaryService[T]) GetItems() *[]T {
	return s.items
}

func (s *MonetaryService[T]) CreateItem(item T) error {
	if err := store.Create(s.GetRepo(), item); err != nil {
		return err
	}
	if err := store.GetAll(s.GetRepo(), s.GetItems()); err != nil {
		return err
	}
	return nil
}

func (s *MonetaryService[T]) UpdateItem(item T) error {
	if err := store.Update(s.GetRepo(), item.GetID(), item); err != nil {
		return err
	}
	if err := store.GetAll(s.GetRepo(), s.GetItems()); err != nil {
		return err
	}
	return nil
}

func (s *MonetaryService[T]) DeleteItem(item T) error {
	if err := store.Delete(s.GetRepo(), item.GetID(), s.GetElementExample()); err != nil {
		return err
	}
	if err := store.GetAll(s.GetRepo(), s.GetItems()); err != nil {
		return err
	}
	return nil
}

func (s *MonetaryService[T]) GetSum() float64 {
	return store.GetSum(s.GetRepo(), s.GetElementExample(), "amount")
}

func (s *MonetaryService[T]) GetFilteredSum(query string, args ...interface{}) float64 {
	return store.GetSumWhere(s.GetRepo(), s.GetElementExample(), "amount", query, args)
}

func (s *MonetaryService[T]) DeleteAll() error {
	for _, item := range *s.items {
		err := store.Delete(s.GetRepo(), item.GetID(), s.GetElementExample())
		if err != nil {
			return err
		}
	}
	return nil
}
