package main

import (
	"fmt"
	"github.com/mixedmachine/simple-budget-app/internal/utils"
	"image/color"
	"strconv"

	"github.com/mixedmachine/simple-budget-app/internal/components"
	. "github.com/mixedmachine/simple-budget-app/internal/components"
	"github.com/mixedmachine/simple-budget-app/internal/models"

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
	var simpleBudget models.SimpleBudget

	simpleBudget.App = app.NewWithID("com.mixedmachine.simplebudgetapp")
	simpleBudget.Window = simpleBudget.App.NewWindow(AppName)

	simpleBudget.SetIcon()
	simpleBudget.SetUpRepo()
	simpleBudget.SetUpServices()

	incomeTotal := simpleBudget.IncomeService.GetSum()
	incomeAllocated := simpleBudget.AllocationService.GetSum()

	simpleBudget.LabelComponents = map[string]*canvas.Text{
		"incomeTotal": canvas.NewText(fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
			strconv.FormatFloat(incomeTotal, 'f', 2, 64),
			strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
			strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64)),
			theme.ForegroundColor()),
		"expenseTotal": canvas.NewText(fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
			simpleBudget.ExpenseService.GetSum(),
			simpleBudget.ExpenseService.GetSum()-simpleBudget.AllocationService.GetSum()),
			theme.ForegroundColor()),
	}

	simpleBudget.ListComponents = CreateListComponents(&simpleBudget)

	addButtons := CreateAddButtons(&simpleBudget)

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
			simpleBudget.LabelComponents["incomeTotal"],
		),
		container.NewBorder(
			nil, nil, nil,
			widget.NewButton("clear", func() {
				dialogPopUp := dialog.NewConfirm(
					"Clear Income",
					"Are you sure you want to clear all income?",
					func(ok bool) {
						if ok {
							err := simpleBudget.IncomeService.DeleteAll()
							utils.HandleErr(simpleBudget.Window, err)
							simpleBudget.ListComponents["income"].Refresh()
						}
					},
					simpleBudget.Window,
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
			simpleBudget.LabelComponents["expenseTotal"],
		),
		container.NewBorder(
			nil, nil, nil,
			widget.NewButton("clear", func() {
				dialogPopUp := dialog.NewConfirm(
					"Clear Expenses",
					"Are you sure you want to clear all expenses?",
					func(ok bool) {
						if ok {
							err := simpleBudget.ExpenseService.DeleteAll()
							utils.HandleErr(simpleBudget.Window, err)
							simpleBudget.ListComponents["expense"].Refresh()
						}
					},
					simpleBudget.Window,
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
			simpleBudget.ListComponents["income"],
		),
		container.NewBorder(
			expenseHeader,
			nil,
			nil,
			nil,
			simpleBudget.ListComponents["expense"],
		),
	)

	allocationsHeader := container.NewBorder(
		nil, nil,
		canvas.NewText(fmt.Sprintf("Total: $%.2f", simpleBudget.AllocationService.GetSum()), color.White),
		widget.NewButton("clear", func() {
			dialogPopUp := dialog.NewConfirm(
				"Clear Allocations",
				"Are you sure you want to clear all allocations?",
				func(ok bool) {
					if ok {
						err := simpleBudget.AllocationService.DeleteAll()
						utils.HandleErr(simpleBudget.Window, err)
						simpleBudget.ListComponents["allocation"].Refresh()
					}
				}, simpleBudget.Window,
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
		simpleBudget.ListComponents["allocation"],
	)

	notesTab := components.CreateNotesComponent(
		simpleBudget.Window,
		simpleBudget.Repo,
		&simpleBudget.Notes,
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

	simpleBudget.Window.SetContent(
		container.NewBorder(
			nil,
			footerContainer,
			nil,
			nil,
			centerContainer,
		),
	)

	simpleBudget.Window.Resize(fyne.NewSize(1000, 600))
	simpleBudget.Window.SetMaster()
	simpleBudget.Window.CenterOnScreen()
	simpleBudget.Window.ShowAndRun()

}
