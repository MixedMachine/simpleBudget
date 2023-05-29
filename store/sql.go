package store

import (
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

func Update[T any](s *SqlDB, id uint, element T) error {
	result := s.db.Model(element).Where("id = ?", id).Updates(element)
	return result.Error
}

func Delete[T any](s *SqlDB, id uint, element T) error {
	result := s.db.Delete(element, id)
	return result.Error
}
