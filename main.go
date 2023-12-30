package main

import (
	"fmt"
	"github.com/mixedmachine/simple-budget-app/internal/utils"
	"image/color"
	"strconv"

	"github.com/mixedmachine/simple-budget-app/internal/components"
	"github.com/mixedmachine/simple-budget-app/internal/components/buttons"
	"github.com/mixedmachine/simple-budget-app/internal/components/lists"
	"github.com/mixedmachine/simple-budget-app/internal/core"

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
	var simpleBudget core.SimpleBudget

	simpleBudget.App = app.NewWithID("com.mixedmachine.simplebudgetapp")
	simpleBudget.Window = simpleBudget.App.NewWindow(AppName)

	simpleBudget.SetIcon()
	simpleBudget.SetUpRepo()
	simpleBudget.SetUpServices()

	simpleBudget.LabelComponents = map[string]*canvas.Text{
		"income":  canvas.NewText("Income", theme.ForegroundColor()),
		"expense": canvas.NewText("Expenses", theme.ForegroundColor()),
		"incomeTotal": canvas.NewText(fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
			strconv.FormatFloat(simpleBudget.IncomeService.GetSum(), 'f', 2, 64),
			strconv.FormatFloat(simpleBudget.AllocationService.GetSum(), 'f', 2, 64),
			strconv.FormatFloat(simpleBudget.IncomeService.GetSum()-simpleBudget.AllocationService.GetSum(), 'f', 2, 64)),
			theme.ForegroundColor()),
		"expenseTotal": canvas.NewText(fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
			simpleBudget.ExpenseService.GetSum(),
			simpleBudget.ExpenseService.GetSum()-simpleBudget.AllocationService.GetSum()),
			theme.ForegroundColor()),
	}

	simpleBudget.ListComponents = lists.CreateListComponents(&simpleBudget)
	simpleBudget.ButtonComonents = buttons.CreateAddButtons(&simpleBudget)

	footerContainerAdds := container.New(layout.NewHBoxLayout(),
		simpleBudget.ButtonComonents["addIncome"],
		simpleBudget.ButtonComonents["addExpense"],
		simpleBudget.ButtonComonents["addAllocation"],
	)

	footerContainerAdds.Resize(fyne.NewSize(1000, 100))

	footerContainer := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewCenterLayout(), footerContainerAdds),
	)

	simpleBudget.LabelComponents["income"].TextSize = 20
	simpleBudget.LabelComponents["income"].TextStyle = fyne.TextStyle{Bold: true}

	simpleBudget.LabelComponents["expense"].TextSize = 20
	simpleBudget.LabelComponents["expense"].TextStyle = fyne.TextStyle{Bold: true}

	incomeHeader := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			simpleBudget.LabelComponents["income"],
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
							incomeTotal := simpleBudget.IncomeService.GetSum()
							incomeAllocated := simpleBudget.AllocationService.GetSum()
							simpleBudget.LabelComponents["incomeTotal"].Text = fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
								strconv.FormatFloat(incomeTotal, 'f', 2, 64),
								strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
								strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64))
							simpleBudget.LabelComponents["incomeTotal"].Refresh()
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
			simpleBudget.LabelComponents["expense"],
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
							simpleBudget.LabelComponents["expenseTotal"].Text = fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
								simpleBudget.ExpenseService.GetSum(),
								simpleBudget.ExpenseService.GetSum()-simpleBudget.AllocationService.GetSum())
							simpleBudget.LabelComponents["expenseTotal"].Refresh()
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

	simpleBudget.LabelComponents["allocationTotal"] = canvas.NewText(fmt.Sprintf("Total: $%.2f", simpleBudget.AllocationService.GetSum()), color.White)

	allocationsHeader := container.NewBorder(
		nil, nil,
		simpleBudget.LabelComponents["allocationTotal"],
		widget.NewButton("clear", func() {
			dialogPopUp := dialog.NewConfirm(
				"Clear Allocations",
				"Are you sure you want to clear all allocations?",
				func(ok bool) {
					if ok {
						err := simpleBudget.AllocationService.DeleteAll()
						utils.HandleErr(simpleBudget.Window, err)
						simpleBudget.ListComponents["allocation"].Refresh()
						incomeTotal := simpleBudget.IncomeService.GetSum()
						expenseTotal := simpleBudget.ExpenseService.GetSum()
						incomeAllocated := simpleBudget.AllocationService.GetSum()
						simpleBudget.LabelComponents["incomeTotal"].Text = fmt.Sprintf("Total: $%s\tAllocated: $%s\tDifference: $%s",
							strconv.FormatFloat(incomeTotal, 'f', 2, 64),
							strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
							strconv.FormatFloat(incomeTotal-incomeAllocated, 'f', 2, 64))
						simpleBudget.LabelComponents["incomeTotal"].Refresh()
						simpleBudget.LabelComponents["expenseTotal"].Text = fmt.Sprintf("Total: $%.2f \t Needed: $%.2f",
							strconv.FormatFloat(expenseTotal, 'f', 2, 64),
							strconv.FormatFloat(incomeAllocated, 'f', 2, 64),
							strconv.FormatFloat(expenseTotal-incomeAllocated, 'f', 2, 64))
						simpleBudget.LabelComponents["expenseTotal"].Refresh()
						simpleBudget.LabelComponents["allocationTotal"].Text = fmt.Sprintf("Total: $%.2f", simpleBudget.AllocationService.GetSum())
						simpleBudget.LabelComponents["allocationTotal"].Refresh()
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
