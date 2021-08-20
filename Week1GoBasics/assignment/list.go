package main

import (
	"fmt"
	"strconv"
)

type list map[string]iteminfo

const displayMsg = `
Shopping List Contents`

func (l list) print(categories []string) {
	fmt.Println(displayMsg)
	if len(l) == 0 {
		fmt.Println("Oops, the list is empty!")
	}
	for name, info := range l {
		c, q, u := getItemValue(info)
		cat, _ := getCategoryName(categories, c)
		fmt.Printf("Category: %s - Item: %s Quantity: %d Unit Cost: %.2f\n", cat, name, q, u)
	}
}

func generateShoppingListReport(l *list) {
	var input int
	fmt.Println(reportPrompt)
	input, _ = strconv.Atoi(getInput())

	switch input {
	//report by category total cost
	case 1:
		getReportByCategoryTotalCost(l)
	//report by category list
	case 2:
		getReportByCategoryList(l)
	case 3:
		fmt.Println("Returning to main menu")
	default:
		fmt.Println("Invalid option. Returning to main menu")
	}
}
