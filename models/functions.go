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

func ReallocateFunds(income *Income, expense *Expense, allocatedToIncome, prevAmount, amount float64) *Allocation {
	maxAmount := income.Amount

	newAmount := allocatedToIncome - prevAmount + amount

	if newAmount > maxAmount {
		log.Error("Cannot allocate more than income amount")
		return nil
	}

	return &Allocation{
		Amount:       amount,
		FromIncomeID: income.ID,
		ToExpenseID:  expense.ID,
	}
}

func AllocateFunds(income *Income, expense *Expense, amount float64) *Allocation {
	maxAmount := income.Amount

	if amount > maxAmount {
		log.Error("Cannot allocate more than income amount")
		return nil
	}

	return &Allocation{
		Amount:       amount,
		FromIncomeID: income.ID,
		ToExpenseID:  expense.ID,
	}
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
