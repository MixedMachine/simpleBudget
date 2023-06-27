package store

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"gorm.io/gorm"
)

type SqlDB struct {
	db *gorm.DB
}

func NewSqlDB(db *gorm.DB) *SqlDB {
	return &SqlDB{
		db: db,
	}
}

func Create[T any](s *SqlDB, element T) error {
	result := s.db.Create(element)
	return result.Error
}

func GetAll[T any](s *SqlDB, elements T) error {
	result := s.db.Find(elements)
	return result.Error
}

func GetAllWhere[T any](s *SqlDB, elements T, query string, args ...interface{}) error {
	result := s.db.Where(query, args...).Find(elements)
	return result.Error
}

func Get[T any](s *SqlDB, id uint, element T) error {
	result := s.db.First(element, id)
	return result.Error
}

func GetWhere[T any](s *SqlDB, element T, query string, args ...interface{}) error {
	result := s.db.Where(query, args...).First(element)
	return result.Error
}

func GetSum[T any](s *SqlDB, element T, column string) float64 {
	var sum float64
	s.db.Model(element).Select("total(" + column + ")").Scan(&sum)
	return sum
}

func GetSumWhere[T any](s *SqlDB, element T, column, query string, args ...interface{}) float64 {
	var sum float64
	s.db.Model(element).Where(query, args...).Select("total(" + column + ")").Scan(&sum)
	return sum
}

func Update[T any](s *SqlDB, id uint, element T) error {
	result := s.db.Model(element).Where("id = ?", id).Updates(element)
	return result.Error
}

func Delete[T any](s *SqlDB, id uint, element T) error {
	result := s.db.Delete(element, id)
	return result.Error
}

func DeleteAllAllocations(repo *SqlDB, allocations *[]models.Allocation) error {
	for _, allocation := range *allocations {
		err := Delete(repo, allocation.ID, &models.Allocation{})
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteAllIncomes(repo *SqlDB, incomes *[]models.Income) error {
	for _, income := range *incomes {
		err := Delete(repo, income.ID, &models.Income{})
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteAllExpenses(repo *SqlDB, expenses *[]models.Expense) error {
	for _, expense := range *expenses {
		err := Delete(repo, expense.ID, &models.Expense{})
		if err != nil {
			return err
		}
	}
	return nil
}
