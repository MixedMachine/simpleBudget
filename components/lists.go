package components

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mixedmachine/simple-budget-app/models"
	"github.com/mixedmachine/simple-budget-app/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

func CreateListComponents(
	myWindow *fyne.Window,
	repo *(store.SqlDB), incomeTotalLabel, expenseTotalLabel *canvas.Text,
	incomes *[]models.Income, expenses *[]models.Expense, allocations *[]models.Allocation,
) map[string]*(widget.List) {
	var err error
	var incomeList *widget.List
	var expenseList *widget.List
	var allocationList *widget.List

	incomeList = widget.NewList(
		func() int { return len(*incomes) },
		func() fyne.CanvasObject {
			var cols int
			models.SortIncomeByDate(incomes)
			nameLabel := widget.NewLabel("Name")
			allocatedLabel := widget.NewLabel("Allocated")
			amountLabel := widget.NewLabel("Amount")
			dateLabel := widget.NewLabel("Date")
			if fyne.CurrentDevice().IsMobile() {
				cols = 2
			} else {
				cols = 4
			}
			incomeContainer := container.NewGridWithColumns(cols, dateLabel, nameLabel, allocatedLabel, amountLabel)
			edtb := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
			delb := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			buttonContainer := container.NewHBox(edtb, delb)
			return container.NewBorder(nil, nil, nil, buttonContainer, incomeContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)

			contactContainer := c.Objects[0].(*fyne.Container)
			buttonContainer := c.Objects[1].(*fyne.Container)

			dateLabel := contactContainer.Objects[0].(*widget.Label)
			nameLabel := contactContainer.Objects[1].(*widget.Label)
			allocatedLabel := contactContainer.Objects[2].(*widget.Label)
			amountLabel := contactContainer.Objects[3].(*widget.Label)

			edtb := buttonContainer.Objects[0].(*widget.Button)
			delb := buttonContainer.Objects[1].(*widget.Button)

			incomeID := (*incomes)[i].ID

			nameLabel.SetText((*incomes)[i].Name)
			allocatedLabel.SetText(fmt.Sprintf("allocated: $ %.2f", store.GetSumWhere(repo, allocations, "amount", "from_income_id = ?", incomeID)))
			amountLabel.SetText(fmt.Sprintf("total: $ %.2f", (*incomes)[i].Amount))
			dateLabel.SetText((*incomes)[i].Date.Format("2006-01-02"))

			edtb.OnTapped = func() {
				incomeEntryName := widget.NewEntry()
				incomeEntryAmount := widget.NewEntry()
				incomeEntryDate := widget.NewEntry()

				incomeFormName := widget.NewFormItem("Name", incomeEntryName)
				incomeFormAmount := widget.NewFormItem("Amount", incomeEntryAmount)
				incomeFormDate := widget.NewFormItem("Date", incomeEntryDate)

				incomeFormItems := []*widget.FormItem{incomeFormName, incomeFormAmount, incomeFormDate}

				dialogBox := dialog.NewForm("Edit Income", "Save", "Cancel", incomeFormItems, func(ok bool) {
					if ok {
						amount, err := strconv.ParseFloat(incomeEntryAmount.Text, 64)
						if err != nil {
							log.Fatal(err)
						}
						date, error := time.Parse("2006-01-02", incomeEntryDate.Text)

						if error != nil {
							log.Fatal(error)
						}

						income := &models.Income{
							ID:     incomeID,
							Name:   incomeEntryName.Text,
							Amount: amount,
							Date:   date,
						}

						if err := store.Update(repo, incomeID, income); err != nil {
							log.Fatal(err)
						}

						if err = store.GetAll(repo, incomes); err != nil {
							log.Fatal(err)
						}
					}
					incomeList.Refresh()
					incomeTotalLabel.Text = "Total: $" + strconv.FormatFloat(store.GetSum(repo, incomes, "amount"), 'f', 2, 64)
					incomeTotalLabel.Refresh()
				}, *myWindow)

				incomeEntryName.SetText((*incomes)[i].Name)
				incomeEntryAmount.SetText(fmt.Sprintf("%.2f", (*incomes)[i].Amount))
				incomeEntryDate.SetText((*incomes)[i].Date.Format("2006-01-02"))

				dialogBox.Resize(fyne.NewSize(500, 300))

				dialogBox.Show()
			}

			delb.OnTapped = func() {

				dialogBox := dialog.NewConfirm(
					"Delete Income",
					"Do you wish to delete this income?",
					func(ok bool) {
						if ok {
							if err := store.Delete(repo, incomeID, &models.Income{}); err != nil {
								log.Fatal(err)
							}
							if err = store.GetAll(repo, incomes); err != nil {
								log.Fatal(err)
							}
						}
						incomeList.Refresh()
						incomeTotalLabel.Text = "Total: $" + strconv.FormatFloat(store.GetSum(repo, incomes, "amount"), 'f', 2, 64)
						incomeTotalLabel.Refresh()
					}, *myWindow,
				)

				dialogBox.Resize(fyne.NewSize(300, 200))
				dialogBox.Show()

			}

		},
	)

	expenseList = widget.NewList(
		func() int { return len(*expenses) },
		func() fyne.CanvasObject {
			var cols int
			models.SortExpenseByDate(expenses)
			nameLabel := widget.NewLabel("Name")
			allocatedLabel := widget.NewLabel("Allocated")
			amountLabel := widget.NewLabel("Amount")
			dateLabel := widget.NewLabel("Date")
			if fyne.CurrentDevice().IsMobile() {
				cols = 2
			} else {
				cols = 4
			}
			expenseContainer := container.NewGridWithColumns(cols, dateLabel, nameLabel, allocatedLabel, amountLabel)
			edtb := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
			delb := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			buttonContainer := container.NewHBox(edtb, delb)
			return container.NewBorder(nil, nil, nil, buttonContainer, expenseContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)

			contactContainer := c.Objects[0].(*fyne.Container)
			buttonContainer := c.Objects[1].(*fyne.Container)

			dateLabel := contactContainer.Objects[0].(*widget.Label)
			nameLabel := contactContainer.Objects[1].(*widget.Label)
			allocatedLabel := contactContainer.Objects[2].(*widget.Label)
			amountLabel := contactContainer.Objects[3].(*widget.Label)

			edtb := buttonContainer.Objects[0].(*widget.Button)
			delb := buttonContainer.Objects[1].(*widget.Button)

			expenseID := (*expenses)[i].ID

			nameLabel.SetText((*expenses)[i].Name)
			allocatedLabel.SetText(
				fmt.Sprintf("allocated: $ %.2f",
					store.GetSumWhere(repo, allocations, "Amount", "to_expense_id = ?", expenseID),
				),
			)
			amountLabel.SetText(fmt.Sprintf("total: $ %.2f", (*expenses)[i].Amount))
			dateLabel.SetText((*expenses)[i].Date.Format("2006-01-02"))

			edtb.OnTapped = func() {

				expenseEntryName := widget.NewEntry()
				expenseEntryAmount := widget.NewEntry()
				expenseEntryDate := widget.NewEntry()

				expenseFormName := widget.NewFormItem("Name", expenseEntryName)
				expenseFormAmount := widget.NewFormItem("Amount", expenseEntryAmount)
				expenseFormDate := widget.NewFormItem("Date", expenseEntryDate)

				expensesFormItems := []*widget.FormItem{expenseFormName, expenseFormAmount, expenseFormDate}

				dialogBox := dialog.NewForm("Edit Expense", "Save", "Cancel", expensesFormItems, func(ok bool) {
					if ok {
						amount, err := strconv.ParseFloat(expenseEntryAmount.Text, 64)
						if err != nil {
							log.Fatal(err)
						}

						date, err := time.Parse("2006-01-02", expenseEntryDate.Text)
						if err != nil {
							log.Fatal(err)
						}

						expense := models.Expense{
							ID:     expenseID,
							Name:   expenseEntryName.Text,
							Amount: amount,
							Date:   date,
						}

						if err := store.Update(repo, expenseID, &expense); err != nil {
							log.Fatal(err)
						}
						if err = store.GetAll(repo, expenses); err != nil {
							log.Fatal(err)
						}
					}
					expenseList.Refresh()
					expenseTotalLabel.Text = "Total: $" + strconv.FormatFloat(store.GetSum(repo, expenses, "amount"), 'f', 2, 64)
					expenseTotalLabel.Refresh()
				}, *myWindow)

				expenseEntryName.SetText((*expenses)[i].Name)
				expenseEntryAmount.SetText(fmt.Sprintf("$ %.2f", (*expenses)[i].Amount))
				expenseEntryDate.SetText((*expenses)[i].Date.Format("2006-01-02"))

				dialogBox.Resize(fyne.NewSize(500, 300))

				dialogBox.Show()

			}

			delb.OnTapped = func() {

				dialogBox := dialog.NewConfirm(
					"Delete Expense",
					"Do you wish to delete this expense?",
					func(ok bool) {
						if ok {
							if err := store.Delete(repo, expenseID, &models.Expense{}); err != nil {
								log.Fatal(err)
							}
							if err = store.GetAll(repo, expenses); err != nil {
								log.Fatal(err)
							}
						}
						expenseList.Refresh()
						expenseTotalLabel.Text = "Total: $" + strconv.FormatFloat(store.GetSum(repo, expenses, "amount"), 'f', 2, 64)
						expenseTotalLabel.Refresh()
					}, *myWindow,
				)

				dialogBox.Resize(fyne.NewSize(300, 200))
				dialogBox.Show()

			}

		},
	)

	allocationList = widget.NewList(
		func() int { return len(*allocations) },
		func() fyne.CanvasObject {
			fromLabel := widget.NewLabel("FromIncomeID")
			toLabel := widget.NewLabel("ToExpenseID")
			amountLabel := widget.NewLabel("Amount")
			allocationContainer := container.NewGridWithColumns(3, fromLabel, toLabel, amountLabel)
			edtb := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
			delb := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			buttonContainer := container.NewHBox(edtb, delb)
			return container.NewBorder(nil, nil, nil, buttonContainer, allocationContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			c := o.(*fyne.Container)

			contactContainer := c.Objects[0].(*fyne.Container)
			buttonContainer := c.Objects[1].(*fyne.Container)

			fromLabel := contactContainer.Objects[0].(*widget.Label)
			toLabel := contactContainer.Objects[1].(*widget.Label)
			amountLabel := contactContainer.Objects[2].(*widget.Label)

			edtb := buttonContainer.Objects[0].(*widget.Button)
			delb := buttonContainer.Objects[1].(*widget.Button)

			allocationID := (*allocations)[i].ID

			toLabel.SetText(models.GetExpenseByID(expenses, (*allocations)[i].ToExpenseID).Name)
			fromLabel.SetText(models.GetIncomeByID(incomes, (*allocations)[i].FromIncomeID).Name)
			amountLabel.SetText(fmt.Sprintf("%.2f", (*allocations)[i].Amount))

			edtb.OnTapped = func() {

				allocationEntryFromIncomeID := widget.NewLabel(models.GetIncomeByID(incomes, (*allocations)[i].FromIncomeID).Name)
				allocationEntryToExpenseID := widget.NewLabel(models.GetExpenseByID(expenses, (*allocations)[i].ToExpenseID).Name)
				allocationEntryAmount := widget.NewEntry()

				allocationFormFromIncomeID := widget.NewFormItem("From", allocationEntryFromIncomeID)
				allocationFormToExpenseID := widget.NewFormItem("To", allocationEntryToExpenseID)
				allocationFormAmount := widget.NewFormItem("Amount", allocationEntryAmount)

				incomeID, err := strconv.ParseUint(allocationEntryFromIncomeID.Text, 10, 64)
				if err != nil {
					log.Error(err)
				}
				allocationFormAmount.HintText = fmt.Sprintf("Max: $%.2f", store.GetSumWhere(repo, incomes, "amount", "id = ?", incomeID))

				allocationFormItems := []*widget.FormItem{
					allocationFormFromIncomeID,
					allocationFormToExpenseID,
					allocationFormAmount,
				}

				dialogBox := dialog.NewForm("Edit Allocation", "Save", "Cancel", allocationFormItems, func(ok bool) {
					if ok {
						fromIncome := models.GetIncomeByName(incomes, allocationEntryFromIncomeID.Text)
						toExpense := models.GetExpenseByName(expenses, allocationEntryToExpenseID.Text)
						amount, err := strconv.ParseFloat(allocationEntryAmount.Text, 64)
						if err != nil {
							log.Fatal(err)
						}

						a := models.AllocatFunds(
							&fromIncome,
							&toExpense,
							store.GetSumWhere(
								repo,
								allocations,
								"amount",
								"from_income_id = ?",
								fromIncome.ID,
							),
							amount,
						)

						if a == nil {
							return
						}

						a.ID = allocationID

						if err := store.Update(repo, fromIncome.ID, fromIncome); err != nil {
							log.Fatal(err)
						}
						if err := store.Update(repo, a.ID, &a); err != nil {
							log.Fatal(err)
						}
						store.GetAll(repo, allocations)
						store.GetAll(repo, incomes)
						allocationList.Refresh()
						incomeList.Refresh()
						incomeTotalLabel.Text = "Total: $" + strconv.FormatFloat(store.GetSum(repo, incomes, "amount"), 'f', 2, 64)
						incomeTotalLabel.Refresh()
					}
					allocationList.Refresh()
				}, *myWindow)

				allocationEntryAmount.SetText(fmt.Sprintf("%.2f", (*allocations)[i].Amount))

				dialogBox.Resize(fyne.NewSize(500, 300))
				dialogBox.Show()

			}

			delb.OnTapped = func() {

				dialogBox := dialog.NewConfirm(
					"Delete Allocation",
					"Do you wish to delete this allocation?",
					func(ok bool) {
						if ok {
							if err := store.Delete(repo, allocationID, models.NewAllocations()); err != nil {
								log.Fatal(err)
							}
							store.GetAll(repo, allocations)
						}
						allocationList.Refresh()
						incomeList.Refresh()
						incomeTotalLabel.Text = "Total: $" + strconv.FormatFloat(store.GetSum(repo, incomes, "amount"), 'f', 2, 64)
						incomeTotalLabel.Refresh()
					}, *myWindow,
				)

				dialogBox.Resize(fyne.NewSize(300, 200))
				dialogBox.Show()

			}
		},
	)

	incomeList.OnSelected = func(id widget.ListItemID) {
		incomeList.UnselectAll()
	}

	expenseList.OnSelected = func(id widget.ListItemID) {
		expenseList.UnselectAll()
	}

	allocationList.OnSelected = func(id widget.ListItemID) {
		allocationList.UnselectAll()
	}

	return map[string]*(widget.List){
		"incomeList":     incomeList,
		"expenseList":    expenseList,
		"allocationList": allocationList,
	}
}
