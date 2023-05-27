package models

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type transaction interface {
	Income | Expense
}

func GetIncomeByID(incomes *[]Income, id primitive.ObjectID) Income {
	for _, income := range *incomes {
		if income.ID == id {
			return income
		}
	}
	return Income{}
}

func GetExpenseByID(expenses *[]Expense, id primitive.ObjectID) Expense {
	for _, expense := range *expenses {
		if expense.ID == id {
			return expense
		}
	}
	return Expense{}
}

func GetIncomeByName(incomes *[]Income, name string) Income {
	for _, income := range *incomes {
		if income.Name == name {
			return income
		}
	}
	return Income{}
}

func GetExpenseByName(expenses *[]Expense, name string) Expense {
	for _, expense := range *expenses {
		if expense.Name == name {
			return expense
		}
	}
	return Expense{}
}

func GetIncomeNames(incomes *[]Income) []string {
	var names []string
	for _, income := range *incomes {
		names = append(names, income.Name)
	}
	return names
}

func GetExpenseNames(expenses *[]Expense) []string {
	var names []string
	for _, expense := range *expenses {
		names = append(names, expense.Name)
	}
	return names
}

func AllocatFunds(income *Income, expense *Expense, amount string) *Allocation {
	prevAmount, err := strconv.ParseFloat(income.Allocated, 64)
	if err != nil {
		log.Error(err)
		return nil
	}

	allocated, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Error(err)
		return nil
	}

	maxAmount, err := strconv.ParseFloat(income.Amount, 64)
	if err != nil {
		log.Error(err)
		return nil
	}

	allocated += prevAmount

	if allocated > maxAmount {
		log.Error("Cannot allocate more than income amount")
		return nil
	}

	newAmount := fmt.Sprintf("%.2f", allocated)

	income.Allocated = newAmount

	return &Allocation{
		Amount:       amount,
		FromIncomeID: income.ID,
		ToExpenseID:  expense.ID,
	}
}

func DeallocatFunds(income *Income, expense *Expense, amount string) {
	prevAmount, err := strconv.ParseFloat(income.Allocated, 64)
	if err != nil {
		log.Fatal(err)
	}

	deallocated, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatal(err)
	}

	deallocated = prevAmount - deallocated

	if deallocated < 0 {
		log.Fatal("Cannot deallocate more than allocated")
	}

	newAmount := fmt.Sprintf("%.2f", deallocated)

	income.Allocated = newAmount
}

func ReallocatFunds(income *Income, expense *Expense, originalAmount, newAmount string) *Allocation {
	var allocation *Allocation

	DeallocatFunds(income, expense, originalAmount)
	allocation = AllocatFunds(income, expense, newAmount)

	return allocation
}

func Filter[T transaction](transactions *[]T, filterFunc func(t T) bool) *[]T {
	var filtered []T
	for _, t := range *transactions {
		if filterFunc(t) {
			filtered = append(filtered, t)
		}
	}
	return &filtered
}

func MapExpenseAllocations(expenses *[]Expense, allocations *[]Allocation) map[string]string {
	allocationsMap := make(map[string]string)
	for _, expense := range *expenses {
		allocationsMap[expense.ID.Hex()] = "0"
	}
	for _, allocation := range *allocations {
		allocationsMap[allocation.ToExpenseID.Hex()] = allocation.Amount
	}
	return allocationsMap
}

func (i *Income) AmountLeftToAllocate() string {
	allocated, err := strconv.ParseFloat(i.Allocated, 64)
	if err != nil {
		log.Fatal(err)
	}

	amount, err := strconv.ParseFloat(i.Amount, 64)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%.2f", amount-allocated)
}
