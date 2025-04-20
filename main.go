package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type Expense struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

var expenses []Expense
var nextID int

const dataFile = "expenses.json"

func loadExpenses() {
	file, err := ioutil.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(file, &expenses)
		if len(expenses) > 0 {
			nextID = expenses[len(expenses)-1].ID + 1
		}
	}
}

func saveExpenses() {
	data, _ := json.MarshalIndent(expenses, "", "  ")
	ioutil.WriteFile(dataFile, data, 0644)
}

func addExpense(description string, amount float64) {
	expense := Expense{
		ID:          nextID,
		Date:        time.Now().Format("2006-01-02"),
		Description: description,
		Amount:      amount,
	}
	expenses = append(expenses, expense)
	nextID++
	saveExpenses()
	fmt.Printf("Expense added successfully (ID: %d)\n", expense.ID)
}

func deleteExpense(id int) {
	for i, expense := range expenses {
		if expense.ID == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			saveExpenses()
			fmt.Println("Expense deleted successfully")
			return
		}
	}
	fmt.Println("Expense not found")
}

func listExpenses() {
	fmt.Println("ID  Date       Description  Amount")
	for _, expense := range expenses {
		fmt.Printf("%d   %s  %s        $%.2f\n", expense.ID, expense.Date, expense.Description, expense.Amount)
	}
}

func summary(month int) {
	total := 0.0
	for _, expense := range expenses {
		if month == 0 || time.Now().Month() == time.Month(month) {
			total += expense.Amount
		}
	}
	if month == 0 {
		fmt.Printf("Total expenses: $%.2f\n", total)
	} else {
		fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(month).String(), total)
	}
}

func main() {
	loadExpenses()

	app := &cli.App{
		Name:  "expense-tracker",
		Usage: "Manage your expenses",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a new expense",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "description",
						Required: true,
					},
					&cli.Float64Flag{
						Name:     "amount",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					addExpense(c.String("description"), c.Float64("amount"))
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete an expense by ID",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "id",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					deleteExpense(c.Int("id"))
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List all expenses",
				Action: func(c *cli.Context) error {
					listExpenses()
					return nil
				},
			},
			{
				Name:  "summary",
				Usage: "Show summary of expenses",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "month",
						Value: 0,
					},
				},
				Action: func(c *cli.Context) error {
					summary(c.Int("month"))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
