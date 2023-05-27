package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstname,omitempty"`
	LastName  string             `bson:"lastname,omitempty"`
	Email     string             `bson:"email,omitempty"`
	Phone     string             `bson:"phone,omitempty"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type Income struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Amount    string             `bson:"amount,omitempty"`
	Date      string             `bson:"date,omitempty"`
	Allocated string             `bson:"allocated,omitempty"`
}

type Expense struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name,omitempty"`
	Amount string             `bson:"amount,omitempty"`
	Date   string             `bson:"date,omitempty"`
}

type Allocation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Amount       string             `bson:"amount,omitempty"`
	FromIncomeID primitive.ObjectID `bson:"fromincomeid,omitempty"`
	ToExpenseID  primitive.ObjectID `bson:"toexpenseid,omitempty"`
}

type BudgetElement interface {
	Income | Expense | Allocation
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
