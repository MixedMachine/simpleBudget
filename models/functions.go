package models

import (
	log "github.com/sirupsen/logrus"
)

type transaction interface {
	Income | Expense
}

func GetIncomeByID(incomes *[]Income, id uint) Income {
	for _, income := range *incomes {
		if income.ID == id {
			return income
		}
	}
	return Income{}
}

func GetExpenseByID(expenses *[]Expense, id uint) Expense {
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

func AllocatFunds(income *Income, expense *Expense, amount float64) *Allocation {
	prevAmount := income.Allocated
	allocated := amount
	maxAmount := income.Amount

	allocated += prevAmount

	if allocated > maxAmount {
		log.Error("Cannot allocate more than income amount")
		return nil
	}

	income.Allocated = allocated

	return &Allocation{
		Amount:       amount,
		FromIncomeID: income.ID,
		ToExpenseID:  expense.ID,
	}
}

func DeallocatFunds(income *Income, expense *Expense, amount float64) {
	prevAmount := income.Allocated

	amount -= prevAmount

	if amount < 0 {
		log.Error("Cannot deallocate more than allocated")
	}

	income.Allocated = amount
}

func ReallocatFunds(income *Income, expense *Expense, originalAmount, newAmount float64) *Allocation {
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

func (i *Income) AmountLeftToAllocate() float64 {
	return i.Amount - i.Allocated
}

func SortIncomeByDate(income *[]Income) {
	for i := 0; i < len(*income); i++ {
		for j := 0; j < len(*income)-1; j++ {
			if (*income)[j].Date.After((*income)[j+1].Date) {
				(*income)[j], (*income)[j+1] = (*income)[j+1], (*income)[j]
			}
		}
	}
}

func SortExpenseByDate(expense *[]Expense) {
	for i := 0; i < len(*expense); i++ {
		for j := 0; j < len(*expense)-1; j++ {
			if (*expense)[j].Date.After((*expense)[j+1].Date) {
				(*expense)[j], (*expense)[j+1] = (*expense)[j+1], (*expense)[j]
			}
		}
	}
}
