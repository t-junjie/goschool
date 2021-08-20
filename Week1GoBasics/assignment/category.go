package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func getCategoryIndex(c []string, cName string) (index int, ok bool) {
	for index, name := range c {
		if name == cName {
			return index, true
		}
	}
	return -1, false
}

func getCategoryName(c []string, cIndex int) (name string, ok bool) {
	if cIndex < 0 || cIndex > len(c) {
		return "", false
	} else {
		return c[cIndex], true
	}
}

func setCategoryWithUserInput() (category []string) {
	fmt.Println("Please input categories (e.g. Household,Food,Drinks)")

	input := getInput()

	time.Sleep(1000 * time.Millisecond)

	categories := strings.Split(input, ",")
	return categories
}

func setCategory() (category []string) {
	return []string{"Household", "Food", "Drinks"}
}

/*
//unsorted print
func getReportByCategoryList(l *list) {
	lst := *l
	for idx, _ := range categories {
		for name, itm := range lst {
			if c, q, u := getItemValue(itm); c == idx {
				cat, _ := getCategoryName(categories, c)
				fmt.Printf("Category: %s - Item: %s Quantity: %d Unit Cost: %.2f\n", cat, name, q, u)
			}
		}
	}
}
*/

//sorted print
func getReportByCategoryList(l *list) {
	lst := *l
	newMap := make(list)
	for key, value := range lst {
		newMap[key] = value
	}

	keys := make([]string, 0, len(newMap))
	for idx, _ := range categories {
		subkeys := make([]string, 0, len(newMap))
		for name, itm := range newMap {
			if c, _, _ := getItemValue(itm); c == idx {
				subkeys = append(subkeys, name)
			}
		}
		sort.Strings(subkeys)
		keys = append(keys, subkeys...)
	}
	for _, k := range keys {
		c, q, u := getItemValue(newMap[k])
		cat, _ := getCategoryName(categories, c)
		fmt.Printf("Category: %s - Item: %s Quantity: %d Unit Cost: %.2f\n", cat, k, q, u)
	}

}

func getReportByCategoryTotalCost(l *list) {
	var totalCostPrompt string = "Total cost by Category.\n"

	lst := *l
	categoryCost := make(map[string]float64)

	for idx, cName := range categories {
		categoryName := cName
		var sum float64
		for _, itm := range lst {
			if c, q, u := getItemValue(itm); c == idx {
				sum += u * float64(q)
			}
		}
		categoryCost[categoryName] = sum
	}

	for cat, cost := range categoryCost {
		totalCostPrompt += fmt.Sprintf("%v cost: %v\n", cat, cost)
	}
	fmt.Println(totalCostPrompt)

}
