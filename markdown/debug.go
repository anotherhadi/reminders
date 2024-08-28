package markdown

import "fmt"

// For debugging purposes
func (md Md) PrintLists() {
	fmt.Println("Title-" + md.Title + ":")
	for _, list := range md.Lists {
		fmt.Println("\nList-" + list.Title + ":")
		for _, item := range list.TodoItems {
			fmt.Println("Item-" + item.Label)
		}
	}
}
