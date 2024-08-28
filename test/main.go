package main

import "reminders/markdown"

func main() {
	md, err := markdown.Markdown("md_files/md.md")
	if err != nil {
		panic(err)
	}

	md.ChangeState("List 2", "Item 1", " ")
	md.ChangeLabel("List 2", "Item 2", "Item 1.2")

}
