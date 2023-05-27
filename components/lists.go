package components

import (
	"github.com/mixedmachine/simple-budget-app/models"
	repo "github.com/mixedmachine/simple-budget-app/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

func CreateListComponents(
	myWindow *fyne.Window,
	ic, ec, ac *(repo.Collection),
	incomes *[]models.Income, expenses *[]models.Expense, allocations *[]models.Allocation,
) map[string]*(widget.List) {
	var err error
	var incomeList *widget.List
	var expenseList *widget.List
	var allocationList *widget.List

	incomeList = widget.NewList(
		func() int { return len(*incomes) },
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabel("Name")
			allocatedLabel := widget.NewLabel("Allocated")
			amountLabel := widget.NewLabel("Amount")
			dateLabel := widget.NewLabel("Date")
			incomeContainer := container.NewGridWithColumns(3, nameLabel, allocatedLabel, amountLabel, dateLabel)
			edtb := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
			delb := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			buttonContainer := container.NewHBox(edtb, delb)
			return container.NewBorder(nil, nil, nil, buttonContainer, incomeContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)

			contactContainer := c.Objects[0].(*fyne.Container)
			buttonContainer := c.Objects[1].(*fyne.Container)

			nameLabel := contactContainer.Objects[0].(*widget.Label)
			amountLabel := contactContainer.Objects[1].(*widget.Label)
			dateLabel := contactContainer.Objects[2].(*widget.Label)

			edtb := buttonContainer.Objects[0].(*widget.Button)
			delb := buttonContainer.Objects[1].(*widget.Button)

			incomeID := (*incomes)[i].ID

			nameLabel.SetText((*incomes)[i].Name)
			amountLabel.SetText((*incomes)[i].Amount)
			dateLabel.SetText((*incomes)[i].Date)

			edtb.OnTapped = func() {
				incomeEntryName := widget.NewEntry()
				incomeEntryAmount := widget.NewEntry()
				incomeEntryDate := widget.NewEntry()

				incomeFormName := widget.NewFormItem("Name", incomeEntryName)
				incomeFormAmount := widget.NewFormItem("Amount", incomeEntryAmount)
				incomeFormDate := widget.NewFormItem("Date", incomeEntryDate)

				incomeFormItems := []*widget.FormItem{incomeFormName, incomeFormAmount, incomeFormDate}

				dialogBox := dialog.NewForm("Edit Income", "Save", "Cancel", incomeFormItems, func(b bool) {
					if b {
						income := models.Income{
							ID:     incomeID,
							Name:   incomeEntryName.Text,
							Amount: incomeEntryAmount.Text,
							Date:   incomeEntryDate.Text,
						}

						if err := repo.Update(ic, incomeID, income); err != nil {
							log.Fatal(err)
						}

						if err = repo.GetAll(ic, incomes); err != nil {
							log.Fatal(err)
						}
					}
					incomeList.Refresh()
				}, *myWindow)

				incomeEntryName.SetText((*incomes)[i].Name)
				incomeEntryAmount.SetText((*incomes)[i].Amount)
				incomeEntryDate.SetText((*incomes)[i].Date)

				dialogBox.Resize(fyne.NewSize(500, 300))

				dialogBox.Show()
			}

			delb.OnTapped = func() {

				dialogBox := dialog.NewConfirm(
					"Delete Income",
					"Do you wish to delete this income?",
					func(b bool) {
						if b {
							if err := repo.Delete(ic, incomeID); err != nil {
								log.Fatal(err)
							}
							if err = repo.GetAll(ic, incomes); err != nil {
								log.Fatal(err)
							}
						}
						incomeList.Refresh()
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
			nameLabel := widget.NewLabel("Name")
			amountLabel := widget.NewLabel("Amount")
			dateLabel := widget.NewLabel("Date")
			expenseContainer := container.NewGridWithColumns(3, nameLabel, amountLabel, dateLabel)
			edtb := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
			delb := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			buttonContainer := container.NewHBox(edtb, delb)
			return container.NewBorder(nil, nil, nil, buttonContainer, expenseContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)

			contactContainer := c.Objects[0].(*fyne.Container)
			buttonContainer := c.Objects[1].(*fyne.Container)

			nameLabel := contactContainer.Objects[0].(*widget.Label)
			amountLabel := contactContainer.Objects[1].(*widget.Label)
			dateLabel := contactContainer.Objects[2].(*widget.Label)

			edtb := buttonContainer.Objects[0].(*widget.Button)
			delb := buttonContainer.Objects[1].(*widget.Button)

			expenseID := (*expenses)[i].ID

			nameLabel.SetText((*expenses)[i].Name)
			amountLabel.SetText((*expenses)[i].Amount)
			dateLabel.SetText((*expenses)[i].Date)

			edtb.OnTapped = func() {

				expenseEntryName := widget.NewEntry()
				expenseEntryAmount := widget.NewEntry()
				expenseEntryDate := widget.NewEntry()

				expenseFormName := widget.NewFormItem("Name", expenseEntryName)
				expenseFormAmount := widget.NewFormItem("Amount", expenseEntryAmount)
				expenseFormDate := widget.NewFormItem("Date", expenseEntryDate)

				expensesFormItems := []*widget.FormItem{expenseFormName, expenseFormAmount, expenseFormDate}

				dialogBox := dialog.NewForm("Edit Expense", "Save", "Cancel", expensesFormItems, func(b bool) {
					if b {
						expense := models.Expense{
							ID:     expenseID,
							Name:   expenseEntryName.Text,
							Amount: expenseEntryAmount.Text,
							Date:   expenseEntryDate.Text,
						}
						expense.ID = primitive.NewObjectID()
						if err := repo.Update(ic, expense.ID, &expense); err != nil {
							log.Fatal(err)
						}
						if err = repo.GetAll(ec, expenses); err != nil {
							log.Fatal(err)
						}
					}
					expenseList.Refresh()
				}, *myWindow)

				expenseEntryName.SetText((*expenses)[i].Name)
				expenseEntryAmount.SetText((*expenses)[i].Amount)
				expenseEntryDate.SetText((*expenses)[i].Date)

				dialogBox.Resize(fyne.NewSize(500, 300))

				dialogBox.Show()

			}

			delb.OnTapped = func() {

				dialogBox := dialog.NewConfirm(
					"Delete Expense",
					"Do you wish to delete this expense?",
					func(b bool) {
						if b {
							if err := repo.Delete(ec, expenseID); err != nil {
								log.Fatal(err)
							}
							if err = repo.GetAll(ec, expenses); err != nil {
								log.Fatal(err)
							}
						}
						expenseList.Refresh()
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
			amountLabel.SetText((*allocations)[i].Amount)

			edtb.OnTapped = func() {

				allocationEntryFromIncomeID := widget.NewLabel(models.GetIncomeByID(incomes, (*allocations)[i].FromIncomeID).Name)
				allocationEntryToExpenseID := widget.NewLabel(models.GetExpenseByID(expenses, (*allocations)[i].ToExpenseID).Name)
				allocationEntryAmount := widget.NewEntry()

				allocationFormFromIncomeID := widget.NewFormItem("From", allocationEntryFromIncomeID)
				allocationFormToExpenseID := widget.NewFormItem("To", allocationEntryToExpenseID)
				allocationFormAmount := widget.NewFormItem("Amount", allocationEntryAmount)

				allocationFormItems := []*widget.FormItem{
					allocationFormFromIncomeID,
					allocationFormToExpenseID,
					allocationFormAmount,
				}

				dialogBox := dialog.NewForm("Edit Allocation", "Save", "Cancel", allocationFormItems, func(b bool) {
					if b {
						fromIncome := models.GetIncomeByName(incomes, allocationEntryFromIncomeID.Text)
						toExpense := models.GetExpenseByName(expenses, allocationEntryToExpenseID.Text)

						a := models.ReallocatFunds(
							&fromIncome,
							&toExpense,
							(*allocations)[i].Amount,
							allocationEntryAmount.Text,
						)

						if a == nil {
							return
						}

						a.ID = allocationID

						if err := repo.Update(ic, fromIncome.ID, fromIncome); err != nil {
							log.Fatal(err)
						}
						if err := repo.Update(ac, a.ID, &a); err != nil {
							log.Fatal(err)
						}
						repo.GetAll(ac, allocations)
						repo.GetAll(ic, incomes)
						allocationList.Refresh()
						incomeList.Refresh()

					}
					allocationList.Refresh()
				}, *myWindow)

				allocationEntryAmount.SetText((*allocations)[i].Amount)

				dialogBox.Resize(fyne.NewSize(500, 300))
				dialogBox.Show()

			}

			delb.OnTapped = func() {

				dialogBox := dialog.NewConfirm(
					"Delete Allocation",
					"Do you wish to delete this allocation?",
					func(b bool) {
						if b {
							fromIncome := models.GetIncomeByID(incomes, (*allocations)[i].FromIncomeID)
							toExpense := models.GetExpenseByID(expenses, (*allocations)[i].ToExpenseID)

							models.DeallocatFunds(
								&fromIncome,
								&toExpense,
								(*allocations)[i].Amount,
							)

							if err := repo.Update(ic, fromIncome.ID, fromIncome); err != nil {
								log.Fatal(err)
							}
							if err := repo.Delete(ac, allocationID); err != nil {
								log.Fatal(err)
							}
							repo.GetAll(ac, allocations)
							repo.GetAll(ic, incomes)
						}
						allocationList.Refresh()
						incomeList.Refresh()
					}, *myWindow,
				)

				dialogBox.Resize(fyne.NewSize(300, 200))
				dialogBox.Show()

			}
		},
	)

	return map[string]*(widget.List){
		"incomeList":     incomeList,
		"expenseList":    expenseList,
		"allocationList": allocationList,
	}
}
