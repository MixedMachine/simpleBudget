package main

import (
	"fmt"
	"image/color"
	"path/filepath"

	. "github.com/mixedmachine/simple-budget-app/components"
	. "github.com/mixedmachine/simple-budget-app/models"
	"github.com/mixedmachine/simple-budget-app/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	APP_NAME = "Simple Budget App"
)

func init() {
	godotenv.Load()
}

func main() {
	var err error
	var dbLocation string

	myApp := app.NewWithID("com.mixedmachine.simplebudgetapp")
	myWindow := myApp.NewWindow(APP_NAME)
	if resourceIconPng, err := fyne.LoadResourceFromPath("assets/icon.png"); err == nil {
		myWindow.SetIcon(resourceIconPng)
	}
	if resourceIconPng, err := fyne.LoadResourceFromPath("icon.png"); err == nil {
		myWindow.SetIcon(resourceIconPng)
	}

	dbLocation = filepath.Join(myApp.Storage().RootURI().Path(), store.SQLITE_FILE)

	income := NewIncomes()
	expense := NewExpenses()
	allocation := NewAllocations()

	repo := store.NewSqlDB(store.InitializeSQL(store.SQLITE, dbLocation))

	err = store.GetAll(repo, income)
	if err != nil {
		log.Fatal(err)
	}
	err = store.GetAll(repo, expense)
	if err != nil {
		log.Fatal(err)
	}
	err = store.GetAll(repo, allocation)
	if err != nil {
		log.Fatal(err)
	}

	incomeTotalLabel := canvas.NewText(fmt.Sprintf("Total: $%.2f", store.GetSum(repo, income, "amount")), color.White)
	expenseTotalLabel := canvas.NewText(fmt.Sprintf("Total: $%.2f", store.GetSum(repo, expense, "amount")), color.White)

	budget := CreateListComponents(
		&myWindow,
		repo, incomeTotalLabel, expenseTotalLabel,
		income, expense, allocation,
	)

	addButtons := CreateAddButtons(
		&myWindow,
		repo, incomeTotalLabel, expenseTotalLabel,
		income, expense, allocation,
		budget,
	)

	footerContainerAdds := container.New(layout.NewHBoxLayout(),
		addButtons["addIncome"],
		addButtons["addExpense"],
		addButtons["addAllocation"],
	)

	footerContainerAdds.Resize(fyne.NewSize(1000, 100))

	footerContainer := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewCenterLayout(), footerContainerAdds),
	)

	incomeLabel := canvas.NewText("Income", color.White)
	incomeLabel.TextSize = 20
	incomeLabel.TextStyle = fyne.TextStyle{Bold: true}

	expenseLabel := canvas.NewText("Expenses", color.White)
	expenseLabel.TextSize = 20
	expenseLabel.TextStyle = fyne.TextStyle{Bold: true}

	incomeHeader := container.New(layout.NewHBoxLayout(),
		incomeLabel,
		incomeTotalLabel,
	)

	expenseHeader := container.New(layout.NewHBoxLayout(),
		expenseLabel,
		expenseTotalLabel,
	)

	transactions := container.New(layout.NewGridLayout(1),
		container.NewBorder(
			incomeHeader,
			nil,
			nil,
			nil,
			budget["incomeList"],
		),
		container.NewBorder(
			expenseHeader,
			nil,
			nil,
			nil,
			budget["expenseList"],
		),
	)

	centerContainer := container.NewAppTabs(
		container.NewTabItem("Transactions",
			transactions,
		),
		container.NewTabItem("Allocations",
			budget["allocationList"],
		),
	)

	myWindow.SetContent(
		container.NewBorder(
			nil,
			footerContainer,
			nil,
			nil,
			centerContainer,
		),
	)

	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.SetMaster()
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()

}
