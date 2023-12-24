package main

import (
	"fmt"
	"github.com/mixedmachine/simple-budget-app/internal/services"
	"github.com/mixedmachine/simple-budget-app/internal/utils"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mixedmachine/simple-budget-app/internal/components"
	. "github.com/mixedmachine/simple-budget-app/internal/components"
	. "github.com/mixedmachine/simple-budget-app/internal/models"
	"github.com/mixedmachine/simple-budget-app/internal/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	AppName = "Simple Budget"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		println("Could not load configs")
	}
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
	myWindow := myApp.NewWindow(AppName)

	setIcon(myWindow)

	dbLocation = filepath.Join(myApp.Storage().RootURI().Path(), store.SQLITE_FILE)

	incomes := NewIncomes()
	expenses := NewExpenses()
	allocations := NewAllocations()
	notes := NewNotes()

	repo := store.NewSqlDB(store.InitializeSQL(store.SQLITE, dbLocation))

	incomeService := services.NewIncomeService(repo, incomes)
	err = incomeService.GetAllIncomes()
	utils.HandleErr(myWindow, err)

	expenseService := services.NewExpenseService(repo, expenses)
	err = expenseService.GetAllExpenses()
	utils.HandleErr(myWindow, err)

	allocationService := services.NewAllocationService(repo, allocations)
	err = allocationService.GetAllAllocations()
	utils.HandleErr(myWindow, err)

	noteService := services.NewNoteService(repo)
	*notes, err = noteService.GetNotes()

	incomeTotal := incomeService.GetSum()
	incomeAllocated := allocationService.GetSum()

	incomeTotalLabel := canvas.NewText(fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
		strconv.FormatFloat(incomeTotal, 'f', 2, 64),
		strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
		strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64)),
		theme.ForegroundColor())
	expenseTotalLabel := canvas.NewText(fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
		expenseService.GetSum(),
		expenseService.GetSum()-allocationService.GetSum()),
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

	incomeHeader := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			incomeLabel,
			incomeTotalLabel,
		),
		container.NewBorder(
			nil, nil, nil,
			widget.NewButton("clear", func() {
				dialogPopUp := dialog.NewConfirm(
					"Clear Income",
					"Are you sure you want to clear all income?",
					func(ok bool) {
						if ok {
							err := incomeService.DeleteAll()
							utils.HandleErr(myWindow, err)
							incomes = NewIncomes()
							budget["incomeList"].Refresh()
						}
					},
					myWindow,
				)
				dialogPopUp.SetDismissText("Cancel")
				dialogPopUp.SetConfirmText("Clear")
				dialogPopUp.Show()
			}),
			nil,
		),
	)

	expenseHeader := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			expenseLabel,
			expenseTotalLabel,
		),
		container.NewBorder(
			nil, nil, nil,
			widget.NewButton("clear", func() {
				dialogPopUp := dialog.NewConfirm(
					"Clear Expenses",
					"Are you sure you want to clear all expenses?",
					func(ok bool) {
						if ok {
							err := expenseService.DeleteAll()
							utils.HandleErr(myWindow, err)
							expenses = NewExpenses()
							budget["expenseList"].Refresh()
						}
					},
					myWindow,
				)
				dialogPopUp.SetDismissText("Cancel")
				dialogPopUp.SetConfirmText("Clear")
				dialogPopUp.Show()
			}),
			nil,
		),
	)

	cols := 2
	if fyne.CurrentDevice().IsMobile() {
		cols = 1
	}

	transactionsTab := container.New(layout.NewGridLayout(cols),
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

	allocationsHeader := container.NewBorder(
		nil, nil,
		canvas.NewText(fmt.Sprintf("Total: $%.2f", allocationService.GetSum()), color.White),
		widget.NewButton("clear", func() {
			dialogPopUp := dialog.NewConfirm(
				"Clear Allocations",
				"Are you sure you want to clear all allocations?",
				func(ok bool) {
					if ok {
						err := allocationService.DeleteAll()
						utils.HandleErr(myWindow, err)
						allocations = NewAllocations()
						budget["allocationList"].Refresh()
					}
				}, myWindow,
			)
			dialogPopUp.SetDismissText("Cancel")
			dialogPopUp.SetConfirmText("Clear")
			dialogPopUp.Show()
		}),
		nil,
	)

	allocationsTab := container.NewBorder(
		allocationsHeader,
		nil, nil, nil,
		budget["allocationList"],
	)

	notesTab := components.CreateNotesComponent(
		myWindow,
		repo,
		notes,
	)

	centerContainer := container.NewAppTabs(
		container.NewTabItem("Transactions",
			transactionsTab,
		),
		container.NewTabItem("Allocations",
			allocationsTab,
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

func setIcon(window fyne.Window) {
	iconPath := ""
	if _, err := os.Stat("assets/icon.png"); os.IsExist(err) {
		iconPath = "assets/icon.png"
	} else if _, err := os.Stat("icon.png"); os.IsExist(err) {
		iconPath = "icon.png"
	}
	if iconPath != "" {
		if resourceIconPng, err := fyne.LoadResourceFromPath("assets/icon.png"); err == nil {
			window.SetIcon(resourceIconPng)
		}
	}
}
