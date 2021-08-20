package main

import (
	"fmt"
	"strconv"
	"time"
)

/*
Add Items option assumes user will enter valid values for prompt
Modify/Delete Items assume user will enter valid values that exist in the shopping list
*/

const runtimePrompt string = `
Shopping List Application
===========================
1. View entire shopping list.
2. Generate Shopping List Report
3. Add Items
4. Modify Items
5. Delete Item
Select your choice:`

const deletePrompt string = `
Delete Item.
Enter item name to delete:`

const reportPrompt = `
Generate Report
1. Total Cost of each Category
2. List of Item by Category
3. Main Menu
Choose your report:
`

var categories = setCategory()
var shoppingList list = make(map[string]iteminfo)

func runMenuOptions() {

	fmt.Println(runtimePrompt)

	var input int
	fmt.Println("... Trying input method 1")
	inputStr := getInput()
	input, _ = strconv.Atoi(inputStr)
	fmt.Printf("Input value: %v, Input type: %T\n", input, input)

	fmt.Println("... Trying input method 2")
	fmt.Scanln(&input)
	fmt.Printf("Input value: %v, Input type: %T\n", input, input)

	switch input {
	case 1:
		shoppingList.print(categories)
	case 2:
		generateShoppingListReport(&shoppingList)
	case 3:
		newlist := addItem(categories, &shoppingList)
		shoppingList = newlist
	case 4:
		newlist := modifyItem(&shoppingList)
		shoppingList = newlist
	case 5:
		newlist := deleteItem(&shoppingList)
		shoppingList = newlist
	default:
		fmt.Println("Invalid option. Returning to main menu.")
	}
	time.Sleep(1000 * time.Millisecond)
}
