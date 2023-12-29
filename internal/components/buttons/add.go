package buttons

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mixedmachine/simple-budget-app/internal/core"
	"github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	log "github.com/sirupsen/logrus"
)

func CreateAddButtons(simpleBudget *core.SimpleBudget) map[string]*widget.Button {
	incomeList := simpleBudget.ListComponents["income"]
	expenseList := simpleBudget.ListComponents["expense"]
	allocationList := simpleBudget.ListComponents["allocation"]
	var incomes []models.Income
	var expenses []models.Expense

	addIncome := widget.NewButton("Add Income", func() {
		incomes = simpleBudget.IncomeService.GetSortedIncomes()
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
					errBox := dialog.NewError(err, simpleBudget.Window)
					errBox.Show()
				}
				date, err := time.Parse("2006-01-02", entryDate.Text)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, simpleBudget.Window)
					errBox.Show()

				}

				i := models.Income{
					TransactionItem: models.TransactionItem{
						MonetaryItem: models.MonetaryItem{
							Amount: amount,
						},
						Name: entryName.Text,
						Date: date,
					},
				}

				if err := simpleBudget.IncomeService.CreateItem(i); err != nil {
					utils.HandleErr(simpleBudget.Window, err)
				}

				incomeList.Refresh()
				incomeTotal := simpleBudget.IncomeService.GetSum()
				incomeAllocated := simpleBudget.AllocationService.GetSum()
				simpleBudget.LabelComponents["incomeTotal"].Text = fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
					strconv.FormatFloat(incomeTotal, 'f', 2, 64),
					strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
					strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64))
				simpleBudget.LabelComponents["incomeTotal"].Refresh()
			}
		}, simpleBudget.Window)

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
					errBox := dialog.NewError(err, simpleBudget.Window)
					errBox.Show()
				}
				date, err := time.Parse("2006-01-02", entryDate.Text)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, simpleBudget.Window)
					errBox.Show()
				}

				e := models.Expense{
					TransactionItem: models.TransactionItem{
						MonetaryItem: models.MonetaryItem{
							Amount: amount,
						},
						Name: entryName.Text,
						Date: date,
					},
				}

				if err := simpleBudget.ExpenseService.CreateItem(e); err != nil {
					utils.HandleErr(simpleBudget.Window, err)
				}

				expenseList.Refresh()
				simpleBudget.LabelComponents["expenseTotal"].Text = fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
					simpleBudget.ExpenseService.GetSum(),
					simpleBudget.ExpenseService.GetSum()-simpleBudget.AllocationService.GetSum())
				simpleBudget.LabelComponents["expenseTotal"].Refresh()
			}
		}, simpleBudget.Window)

		dialogAdd.Resize(fyne.NewSize(500, 300))

		dialogAdd.Show()

	})

	addAllocation := widget.NewButton("Add Allocation", func() {
		incomes = *simpleBudget.IncomeService.GetItems()
		expenses = *simpleBudget.ExpenseService.GetItems()
		entryAmount := widget.NewEntry()
		amountForm := widget.NewFormItem("Amount", entryAmount)
		entryFromIncomeID := widget.NewSelect(models.GetIncomeNames(&incomes), func(incomeName string) {
			if incomeName != "" {
				hint := ""
				incomeID := models.GetIncomeByName(&incomes, incomeName).ID

				amt, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					amt = 0.0
				}
				hint += "Avail: $" + strconv.FormatFloat(
					simpleBudget.IncomeService.GetFilteredSum("id = ?", incomeID)-
						simpleBudget.AllocationService.GetFilteredSum("from_income_id = ?", incomeID)-
						amt, 'f', 2, 64)
				amountForm.HintText = hint
				amountForm.Widget.Refresh()
			}
		})
		fromIncomeIDForm := widget.NewFormItem("From Income ID", entryFromIncomeID)
		entryToExpenseID := widget.NewSelect(models.GetExpenseNames(&expenses), func(s string) {
			filterDate := models.GetExpenseByName(&expenses, s).Date
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

				incomeID := models.GetIncomeByName(&incomes, entryFromIncomeID.Selected).ID
				hint += "Avail: $" + strconv.FormatFloat(
					simpleBudget.IncomeService.GetFilteredSum("id = ?", incomeID)-
						simpleBudget.AllocationService.GetFilteredSum("from_income_id = ?", incomeID)-
						amt, 'f', 2, 64)
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
				incomeID := models.GetIncomeByName(&incomes, entryFromIncomeID.Selected).ID
				chosenIcomeAmount := simpleBudget.IncomeService.GetFilteredSum("id = ?", incomeID)
				amt, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, simpleBudget.Window)
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

		formItems := []*widget.FormItem{fromIncomeIDForm, toExpenseIDForm, amountForm}

		dialogAdd := dialog.NewForm("Add Allocation", "Add", "Cancel", formItems, func(ok bool) {
			if ok {
				fromIncome := models.GetIncomeByName(&incomes, entryFromIncomeID.Selected)

				toExpense := models.GetExpenseByName(&expenses, entryToExpenseID.Selected)

				amount, err := strconv.ParseFloat(entryAmount.Text, 64)
				if err != nil {
					log.Error(err)
					errBox := dialog.NewError(err, simpleBudget.Window)
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

				if err = simpleBudget.IncomeService.UpdateItem(fromIncome); err != nil {
					utils.HandleErr(simpleBudget.Window, err)
				}

				if err = simpleBudget.AllocationService.CreateItem(*a); err != nil {
					utils.HandleErr(simpleBudget.Window, err)
				}

				allocationList.Refresh()
				incomeList.Refresh()
				incomeTotal := simpleBudget.IncomeService.GetSum()
				incomeAllocated := simpleBudget.AllocationService.GetSum()
				simpleBudget.LabelComponents["incomeTotal"].Text = fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
					strconv.FormatFloat(incomeTotal, 'f', 2, 64),
					strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
					strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64))
				simpleBudget.LabelComponents["incomeTotal"].Refresh()
				simpleBudget.LabelComponents["expenseTotal"].Text = fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
					simpleBudget.ExpenseService.GetSum(),
					simpleBudget.ExpenseService.GetSum()-simpleBudget.AllocationService.GetSum())
				simpleBudget.LabelComponents["expenseTotal"].Refresh()
			}
		}, simpleBudget.Window)

		dialogAdd.Resize(fyne.NewSize(500, 300))

		dialogAdd.Show()

	})

	return map[string]*widget.Button{
		"addIncome":     addIncome,
		"addExpense":    addExpense,
		"addAllocation": addAllocation,
	}
}
