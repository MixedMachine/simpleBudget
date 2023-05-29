package main

import (
	. "github.com/mixedmachine/simple-budget-app/components"
	. "github.com/mixedmachine/simple-budget-app/models"
	"github.com/mixedmachine/simple-budget-app/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/joho/godotenv"

	"image/color"

	log "github.com/sirupsen/logrus"
)

const (
	APP_NAME = "Simple Budget App"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Info("Error loading .env file")
	}

	store.MONGO_URI = "mongodb+srv://MixedMachine:Eugene23@personalprojects.zpn6jfo.mongodb.net/?retryWrites=true&w=majority"

}

func main() {
	var err error

	myApp := app.New()
	myWindow := myApp.NewWindow(APP_NAME)
	resourceIconPng, err := fyne.LoadResourceFromPath("assets/icon.png")
	if err != nil {
		log.Info(err)
	}
	myWindow.SetIcon(resourceIconPng)

	income := NewIncomes()
	expense := NewExpenses()
	allocation := NewAllocations()

	repo := store.NewSqlDB(store.InitializeSQL(store.SQLITE))

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

	budget := CreateListComponents(
		&myWindow,
		repo,
		income, expense, allocation,
	)

	addButtons := CreateAddButtons(
		&myWindow,
		repo,
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

	transactions := container.New(layout.NewGridLayout(1),
		container.NewBorder(
			incomeLabel,
			nil,
			nil,
			nil,
			budget["incomeList"],
		),
		container.NewBorder(
			expenseLabel,
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
