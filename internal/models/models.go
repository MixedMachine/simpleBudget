package models

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"unique"`
	Password string `gorm:"not null"`
}

type Income struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	Name   string    `gorm:"unique;not null"`
	Amount float64   `gorm:"type:decimal(10,2);default:0.00;not null"`
	Date   time.Time `gorm:"type:date;not null"`
}

type Expense struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	Name   string    `gorm:"unique;not null"`
	Amount float64   `gorm:"type:decimal(10,2);default:0.00;not null"`
	Date   time.Time `gorm:"type:date;not null"`
}

type Allocation struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	Amount       float64 `gorm:"type:decimal(10,2);default:0.00;not null"`
	FromIncomeID uint    `gorm:"index:idx_from_income_id;foreignKey:FromIncomeID"`
	ToExpenseID  uint    `gorm:"index:idx_to_expense_id;foreignKey:ToExpenseID"`
}

type Notes struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Content string `gorm:"not null"`
}

func NewIncomes() *[]Income {
	return &[]Income{}
}

func NewExpenses() *[]Expense {
	return &[]Expense{}
}

func NewAllocations() *[]Allocation {
	return &[]Allocation{}
}

func NewNotes() *Notes {
	return &Notes{}
}
