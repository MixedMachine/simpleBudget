package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"

	"github.com/mixedmachine/simple-budget-app/components"
	. "github.com/mixedmachine/simple-budget-app/components"
	. "github.com/mixedmachine/simple-budget-app/models"
	"github.com/mixedmachine/simple-budget-app/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	APP_NAME = "Simple Budget App"
)

func init() {
	godotenv.Load()
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
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

	incomes := NewIncomes()
	expenses := NewExpenses()
	allocations := NewAllocations()
	notes := NewNotes()

	repo := store.NewSqlDB(store.InitializeSQL(store.SQLITE, dbLocation))

	err = store.GetAll(repo, incomes)
	if err != nil {
		log.Error(err)
		errBox := dialog.NewError(err, myWindow)
		errBox.Show()
	}
	err = store.GetAll(repo, expenses)
	if err != nil {
		log.Error(err)
		errBox := dialog.NewError(err, myWindow)
		errBox.Show()
	}
	err = store.GetAll(repo, allocations)
	if err != nil {
		log.Error(err)
		errBox := dialog.NewError(err, myWindow)
		errBox.Show()
	}
	err = store.GetAll(repo, notes)
	if err != nil {
		log.Error(err)
		errBox := dialog.NewError(err, myWindow)
		errBox.Show()
	}
	if len(*notes) == 0 {
		*notes = append(*notes, Notes{Content: ""})
		err = store.Create(repo, &(*notes)[0])
		if err != nil {
			log.Error(err)
			errBox := dialog.NewError(err, myWindow)
			errBox.Show()
		}
	}

	incomeTotalLabel := canvas.NewText("Total: $"+
		strconv.FormatFloat(store.GetSum(repo, incomes, "amount"), 'f', 2, 64)+
		" \t Allocated: $"+
		strconv.FormatFloat(store.GetSum(repo, allocations, "amount"), 'f', 2, 64),
		theme.ForegroundColor())
	expenseTotalLabel := canvas.NewText(fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
		store.GetSum(repo, Expense{}, "amount"),
		store.GetSum(repo, Expense{}, "amount")-store.GetSum(repo, Allocation{}, "amount")),
		theme.ForegroundColor())

	budget := CreateListComponents(
		&myWindow,
		repo, incomeTotalLabel, expenseTotalLabel,
		incomes, expenses, allocations,
	)

	addButtons := CreateAddButtons(
		&myWindow,
		repo, incomeTotalLabel, expenseTotalLabel,
		incomes, expenses, allocations,
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

	cols := 2
	if fyne.CurrentDevice().IsMobile() {
		cols = 1
	}

	transactions := container.New(layout.NewGridLayout(cols),
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

	notesTab := components.CreateNotesComponent(
		myWindow,
		repo,
		notes,
	)

	centerContainer := container.NewAppTabs(
		container.NewTabItem("Transactions",
			transactions,
		),
		container.NewTabItem("Allocations",
			budget["allocationList"],
		),
		container.NewTabItem("Notes",
			notesTab,
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
