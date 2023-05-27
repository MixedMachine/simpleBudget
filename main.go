package main

import (
	. "github.com/mixedmachine/simple-budget-app/components"
	. "github.com/mixedmachine/simple-budget-app/models"
	repo "github.com/mixedmachine/simple-budget-app/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/joho/godotenv"

	"image/color"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	APP_NAME = "Budget Buddy"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo.MONGO_URI = os.Getenv("MONGO_URI")
}

func main() {
	var err error

	myApp := app.New()
	myWindow := myApp.NewWindow(APP_NAME)

	income := NewIncomes()
	expense := NewExpenses()
	allocation := NewAllocations()

	ctx, client := repo.InitializeDB()
	defer client.Disconnect(*ctx)

	collections := repo.CreateCollections(ctx, client)

	err = repo.GetAll(collections["income"], income)
	if err != nil {
		log.Fatal(err)
	}
	err = repo.GetAll(collections["expense"], expense)
	if err != nil {
		log.Fatal(err)
	}
	err = repo.GetAll(collections["allocation"], allocation)
	if err != nil {
		log.Fatal(err)
	}

	budget := CreateListComponents(
		&myWindow,
		collections["income"], collections["expense"], collections["allocation"],
		income, expense, allocation,
	)

	addButtons := CreateAddButtons(
		&myWindow,
		collections["income"], collections["expense"], collections["allocation"],
		income, expense, allocation,
		budget,
	)

	footerContainerAdds := container.New(layout.NewHBoxLayout(),
		container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.Black), addButtons["addIncome"]),
		addButtons["addExpense"],
		addButtons["addAllocation"],
	)

	footerContainerAdds.Resize(fyne.NewSize(1000, 100))

	footerContainer := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewCenterLayout(), footerContainerAdds),
	)

	incomeLabel := canvas.NewText("Income", color.White)
	incomeLabel.Alignment = fyne.TextAlignCenter
	incomeLabel.TextSize = 20
	incomeLabel.TextStyle = fyne.TextStyle{Bold: true}
	incomeLabel.SetMinSize(fyne.NewSize(100, 50))
	expenseLabel := canvas.NewText("Expenses", color.White)
	expenseLabel.Alignment = fyne.TextAlignCenter
	expenseLabel.TextSize = 20
	expenseLabel.TextStyle = fyne.TextStyle{Bold: true}
	expenseLabel.SetMinSize(fyne.NewSize(100, 50))

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
