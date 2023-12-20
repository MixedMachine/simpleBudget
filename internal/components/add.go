package components

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"
	"github.com/mixedmachine/simple-budget-app/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	log "github.com/sirupsen/logrus"
)

func CreateAddButtons(
	myWindow *fyne.Window,
	repo *store.SqlDB, incomeTotalLabel, expenseTotalLabel *canvas.Text,
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

		dialogAdd := dialog.NewForm("Add Income", "Add", "Cancel", formItems, func(ok bool) {
			if ok {
				amount, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, *myWindow)
					errBox.Show()
				}
				date, err := time.Parse("2006-01-02", entryDate.Text)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, *myWindow)
					errBox.Show()

				}

				i := models.Income{
					Name:   entryName.Text,
					Amount: amount,
					Date:   date,
				}

				if err := store.Create(repo, &i); err != nil {
					log.Fatal(err)
				}
				err = store.GetAll(repo, incomes)
				if err != nil {
					return
				}
				incomeList.Refresh()
				incomeTotal := store.GetSum(repo, incomes, "amount")
				incomeAllocated := store.GetSum(repo, allocations, "amount")
				incomeTotalLabel.Text = fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
					strconv.FormatFloat(incomeTotal, 'f', 2, 64),
					strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
					strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64))
				incomeTotalLabel.Refresh()
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

		dialogAdd := dialog.NewForm("Add Expense", "Add", "Cancel", formItems, func(ok bool) {
			if ok {
				amount, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, *myWindow)
					errBox.Show()
				}
				date, err := time.Parse("2006-01-02", entryDate.Text)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, *myWindow)
					errBox.Show()
				}

				e := models.Expense{
					Name:   entryName.Text,
					Amount: amount,
					Date:   date,
				}

				if err := store.Create(repo, &e); err != nil {
					log.Fatal(err)
				}
				err = store.GetAll(repo, expenses)
				if err != nil {
					return
				}
				expenseList.Refresh()
				expenseTotalLabel.Text = fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
					store.GetSum(repo, models.Expense{}, "amount"),
					store.GetSum(repo, models.Expense{}, "amount")-
						store.GetSum(repo, models.Allocation{}, "amount"))
				expenseTotalLabel.Refresh()
			}
		}, *myWindow)

		dialogAdd.Resize(fyne.NewSize(500, 300))

		dialogAdd.Show()

	})

	addAllocation := widget.NewButton("Add Allocation", func() {

		entryAmount := widget.NewEntry()
		amountForm := widget.NewFormItem("Amount", entryAmount)
		entryFromIncomeID := widget.NewSelect(models.GetIncomeNames(incomes), func(incomeName string) {
			if incomeName != "" {
				hint := ""
				incomeID := models.GetIncomeByName(incomes, incomeName).ID

				amt, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					amt = 0.0
				}
				hint += "Avail: $" + strconv.FormatFloat(store.GetSumWhere(
					repo, incomes, "amount", "id = ?", incomeID,
				)-store.GetSumWhere(
					repo, allocations, "amount", "from_income_id = ?", incomeID,
				)-amt, 'f', 2, 64)
				amountForm.HintText = hint
				amountForm.Widget.Refresh()
			}
		})
		fromIncomeIDForm := widget.NewFormItem("From Income ID", entryFromIncomeID)
		entryToExpenseID := widget.NewSelect(models.GetExpenseNames(expenses), func(s string) {
			filterDate := models.GetExpenseByName(expenses, s).Date
			filterBy := func(inc models.Income) bool { return utils.CompareDates(inc.Date, filterDate) == -1 }
			entryFromIncomeID.Options = models.GetIncomeNames(models.Filter(incomes, filterBy))
		})
		toExpenseIDForm := widget.NewFormItem("To Expense ID", entryToExpenseID)

		entryAmount.OnChanged = func(amount string) {
			if entryFromIncomeID.Selected != "" {
				hint := ""
				amt, err := strconv.ParseFloat(amount, 64)
				if err != nil {
					amt = 0.0
				}

				incomeID := models.GetIncomeByName(incomes, entryFromIncomeID.Selected).ID
				hint += "Avail: $" + strconv.FormatFloat(store.GetSumWhere(
					repo, incomes, "amount", "id = ?", incomeID,
				)-store.GetSumWhere(
					repo, allocations, "amount", "from_income_id = ?", incomeID,
				)-amt, 'f', 2, 64)
				amountForm.HintText = hint
				amountForm.Widget.Refresh()
			}
		}

		entryAmount.Validator = func(s string) error {
			_, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return errors.New("invalid amount")
			}
			if entryFromIncomeID.Selected != "" {
				incomeID := models.GetIncomeByName(incomes, entryFromIncomeID.Selected).ID
				chosenIcomeAmount := store.GetSumWhere(repo, incomes, "amount", "id = ?", incomeID)
				amt, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, *myWindow)
					errBox.Show()
				}
				if chosenIcomeAmount < amt {
					log.Infof("chosenIcomeAmount: %f", chosenIcomeAmount)
					log.Infof("amt: %f", amt)
					return errors.New("amount is greater than income amount")
				}
			}
			return nil
		}

		formItems := []*widget.FormItem{toExpenseIDForm, fromIncomeIDForm, amountForm}

		dialogAdd := dialog.NewForm("Add Allocation", "Add", "Cancel", formItems, func(ok bool) {
			if ok {
				fromIncome := models.GetIncomeByName(incomes, entryFromIncomeID.Selected)

				toExpense := models.GetExpenseByName(expenses, entryToExpenseID.Selected)

				amount, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, *myWindow)
					errBox.Show()
				}

				a := models.AllocateFunds(
					&fromIncome,
					&toExpense,
					amount,
				)

				if a == nil {
					return
				}

				err = store.Update(repo, fromIncome.ID, fromIncome)
				if err != nil {
					return
				}

				if err = store.Create(repo, &a); err != nil {
					log.Fatal(err)
				}
				err = store.GetAll(repo, allocations)
				if err != nil {
					return
				}
				err = store.GetAll(repo, incomes)
				if err != nil {
					return
				}
				allocationList.Refresh()
				incomeList.Refresh()
				incomeTotal := store.GetSum(repo, incomes, "amount")
				incomeAllocated := store.GetSum(repo, allocations, "amount")
				incomeTotalLabel.Text = fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
					strconv.FormatFloat(incomeTotal, 'f', 2, 64),
					strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
					strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64))
				incomeTotalLabel.Refresh()
				expenseTotalLabel.Text = fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
					store.GetSum(repo, models.Expense{}, "amount"),
					store.GetSum(repo, models.Expense{}, "amount")-
						store.GetSum(repo, models.Allocation{}, "amount"))
				expenseTotalLabel.Refresh()
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
