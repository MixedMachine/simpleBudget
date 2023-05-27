package components

import (
	"errors"
	"strconv"

	"github.com/mixedmachine/simple-budget-app/models"
	repo "github.com/mixedmachine/simple-budget-app/repository"
	"github.com/mixedmachine/simple-budget-app/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	log "github.com/sirupsen/logrus"
)

func CreateAddButtons(
	myWindow *fyne.Window,
	ic, ec, ac *repo.Collection,
	incomes *[]models.Income, expenses *[]models.Expense, allocations *[]models.Allocation,
	listComponents map[string]*(widget.List),
) map[string]*(widget.Button) {
	incomeList := listComponents["incomeList"]
	expenseList := listComponents["expenseList"]
	allocationList := listComponents["allocationList"]

	addIncome := widget.NewButton("Add Income", func() {

		entryName := widget.NewEntry()
		entryAmount := widget.NewEntry()
		entryDate := widget.NewEntry()

		nameForm := widget.NewFormItem("Name", entryName)
		amountForm := widget.NewFormItem("Amount", entryAmount)
		dateForm := widget.NewFormItem("Date", entryDate)

		entryAmount.Validator = func(s string) error {
			_, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return errors.New("invalid amount")
			}
			return nil
		}
		dateForm.HintText = "YYYY-MM-DD"
		entryDate.Validator = func(s string) error {
			if !utils.ValidateDate(s) {
				return errors.New("invalid date")
			}
			return nil
		}

		formItems := []*widget.FormItem{nameForm, amountForm, dateForm}

		dialogAdd := dialog.NewForm("Add Income", "Add", "Cancel", formItems, func(b bool) {
			if b {
				i := models.Income{
					Name:      entryName.Text,
					Amount:    entryAmount.Text,
					Date:      entryDate.Text,
					Allocated: "0",
				}

				if err := repo.Create(ic, &i); err != nil {
					log.Fatal(err)
				}
				repo.GetAll(ic, incomes)
				incomeList.Refresh()
			}
		}, *myWindow)

		dialogAdd.Resize(fyne.NewSize(500, 300))

		dialogAdd.Show()

	})

	addExpense := widget.NewButton("Add Expense", func() {

		entryName := widget.NewEntry()
		entryAmount := widget.NewEntry()
		entryDate := widget.NewEntry()

		nameForm := widget.NewFormItem("Name", entryName)
		amountForm := widget.NewFormItem("Amount", entryAmount)
		dateForm := widget.NewFormItem("Date", entryDate)

		entryAmount.Validator = func(s string) error {
			_, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return errors.New("invalid amount")
			}
			return nil
		}
		dateForm.HintText = "YYYY-MM-DD"
		entryDate.Validator = func(s string) error {
			if !utils.ValidateDate(s) {
				return errors.New("invalid date")
			}
			return nil
		}

		formItems := []*widget.FormItem{nameForm, amountForm, dateForm}

		dialogAdd := dialog.NewForm("Add Expense", "Add", "Cancel", formItems, func(b bool) {
			if b {
				e := models.Expense{
					Name:   entryName.Text,
					Amount: entryAmount.Text,
					Date:   entryDate.Text,
				}

				if err := repo.Create(ec, &e); err != nil {
					log.Fatal(err)
				}
				repo.GetAll(ec, expenses)
				expenseList.Refresh()
			}
		}, *myWindow)

		dialogAdd.Resize(fyne.NewSize(500, 300))

		dialogAdd.Show()

	})

	addAllocation := widget.NewButton("Add Allocation", func() {

		entryFromIncomeID := widget.NewSelect(models.GetIncomeNames(incomes), func(s string) {})
		entryToExpenseID := widget.NewSelect(models.GetExpenseNames(expenses), func(s string) {
			filterDate := models.GetExpenseByName(expenses, s).Date
			filterBy := func(inc models.Income) bool { return utils.CompareDates(inc.Date, filterDate) == -1 }
			entryFromIncomeID.Options = models.GetIncomeNames(models.Filter(incomes, filterBy))
		})
		entryAmount := widget.NewEntry()

		fromIncomeIDForm := widget.NewFormItem("From Income ID", entryFromIncomeID)
		toExpenseIDForm := widget.NewFormItem("To Expense ID", entryToExpenseID)
		amountForm := widget.NewFormItem("Amount", entryAmount)

		entryAmount.Validator = func(s string) error {
			_, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return errors.New("invalid amount")
			}
			return nil
		}

		formItems := []*widget.FormItem{toExpenseIDForm, fromIncomeIDForm, amountForm}

		dialogAdd := dialog.NewForm("Add Allocation", "Add", "Cancel", formItems, func(b bool) {
			if b {
				fromIncome := models.GetIncomeByName(incomes, entryFromIncomeID.Selected)

				toExpense := models.GetExpenseByName(expenses, entryToExpenseID.Selected)

				a := models.AllocatFunds(
					&fromIncome,
					&toExpense,
					entryAmount.Text,
				)

				if a == nil {
					return
				}

				repo.Update(ic, fromIncome.ID, fromIncome)

				if err := repo.Create(ac, &a); err != nil {
					log.Fatal(err)
				}
				repo.GetAll(ac, allocations)
				repo.GetAll(ic, incomes)
				allocationList.Refresh()
				incomeList.Refresh()
			}
		}, *myWindow)

		dialogAdd.Resize(fyne.NewSize(500, 300))

		dialogAdd.Show()

	})

	return map[string]*(widget.Button){
		"addIncome":     addIncome,
		"addExpense":    addExpense,
		"addAllocation": addAllocation,
	}
}
