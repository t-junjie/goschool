package main

import (
	"fmt"
	"strconv"
)

type iteminfo struct {
	category int
	quantity int
	unitCost float64
}

func getItemValue(info iteminfo) (c, q int, u float64) {
	return info.category, info.quantity, info.unitCost
}

func addItem(c []string, l *list) (newlist list) {
	var name, cName string
	var unitcost float64
	var quantity, convertcat int
	var ok bool

	fmt.Println("What is the name of your item?")
	name = getInput()

	fmt.Println("What category does it belong to?")
	cName = getInput()

	fmt.Println("How many units are there?")
	quantity, _ = strconv.Atoi(getInput())

	fmt.Println("How much does it cost per unit?")
	unitcost, _ = strconv.ParseFloat(getInput(), 64)

	convertcat, ok = getCategoryIndex(c, cName)
	lst := *l

	if ok {
		iinfo := iteminfo{convertcat, quantity, unitcost}
		lst[name] = iinfo
		return lst
	} else {
		fmt.Println("Invalid category")
		return lst
	}
}

func deleteItem(l *list) (newlist list) {
	var name string

	fmt.Println(deletePrompt)
	name = getInput()

	lst := *l
	_, ok := lst[name]

	if ok {
		delete(lst, name)
		fmt.Println("Deleted", name)
		return lst
	} else {
		fmt.Println("Item not found. Nothing to delete!")
		return lst
	}
}

//assumes new quantity/unitcost cannot be 0 (this breaks the check for detecting Enter)
func modifyItem(l *list) (newlist list) {
	var modifyMsg, currentName, newName, newCategory string
	var currentUnitCost, newUnitCost float64
	var currentQuantity, newQuantity, currentCategory, newIntCategory int
	var ok bool

	lst := *l

	fmt.Println("What item do you wish to modify?")
	currentName = getInput()

	itminfo, found := lst[currentName]

	if !found {
		fmt.Println("Item not found in shopping list. Returning to main menu")
		return lst
	}

	currentCategory, currentQuantity, currentUnitCost = getItemValue(itminfo)
	cName, _ := getCategoryName(categories, currentCategory)
	fmt.Printf("Current item name is %s - Category is %v - Quantity is %d - Unit Cost %.2f\n", currentName, cName, currentQuantity, currentUnitCost)

	fmt.Println("Enter new Name. Enter for no change")
	newName = getInput()

	if len(newName) == 0 {
		modifyMsg += "No changes to name made.\n"
		newName = currentName
	}

	fmt.Println("Enter new Category. Enter for no change")
	newCategory = getInput()
	//fmt.Printf("' %T, %v'", newCategory, newCategory)
	if len(newCategory) == 0 {
		modifyMsg += "No changes to category made.\n"
		newIntCategory = currentCategory
		ok = true
	} else {
		newIntCategory, ok = getCategoryIndex(categories, newCategory)
	}

	fmt.Println("Enter new Quantity. Enter for no change")
	newQuantity, _ = strconv.Atoi(getInput())
	if newQuantity == 0 {
		modifyMsg += "No changes to quantity made.\n"
		newQuantity = currentQuantity
	} else if newQuantity < 0 {
		fmt.Println("Quantity must be more than 0. Returning to main menu.")
		return lst
	}

	fmt.Println("Enter new Unit Cost. Enter for no change")
	newUnitCost, _ = strconv.ParseFloat(getInput(), 64)
	if newUnitCost == 0 {
		modifyMsg += "No changes to unit cost made.\n"
		newUnitCost = currentUnitCost
	} else if newUnitCost < 0 {
		fmt.Println("Unit Cost must be more than 0. Returning to main menu.")
		return lst
	}

	if ok {
		iinfo := iteminfo{newIntCategory, newQuantity, newUnitCost}
		delete(lst, currentName)
		lst[newName] = iinfo
		fmt.Println(modifyMsg)
		return lst
	} else {
		fmt.Println("Invalid category")
		return lst
	}
}
