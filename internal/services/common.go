package services

import "github.com/mixedmachine/simple-budget-app/internal/store"

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

type MonetaryService struct {
	Service
	elements []any
}

func NewMonetaryService(repo *store.SqlDB, element any) *MonetaryService {
	return &MonetaryService{
		Service: *NewService(repo, element),
	}
}

func (s *Service) GetRepo() *store.SqlDB {
	return s.repo
}

func (s *Service) GetElementType() any {
	return &s.element
}

func (s *MonetaryService) GetSum() float64 {
	return store.GetSum(s.GetRepo(), s.GetElementType(), "amount")
}
