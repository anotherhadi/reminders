package markdown

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var (
	checkbox_pending string = "󰄮"
	checkbox_done           = "󰄲"
	checkbox_cancel         = "󱋬"
	checkbox_empty          = "󰄱"
	checkbox_outline        = "󰡖"

	importance_none   = ""
	importance_low    = "󰈅"
	importance_medium = "󱈸"
	importance_high   = ""

	defaultTitle = "Reminders"
)

func Markdown(path string) (Md, error) {
	directory := filepath.Dir(path)
	filename := filepath.Base(path)
	if filepath.Ext(path) == "" {
		filename += ".md"
	}

	err := os.Chdir(directory)
	if err != nil {
		return Md{}, err
	}

	md := Md{
		DirPath: directory,
		File:    filename,
		Title:   defaultTitle,
	}

	md, err = md.LoadMarkdown()
	if err != nil {
		return Md{}, err
	}

	return md, nil
}

func (md Md) LoadMarkdown() (Md, error) {
	_ = os.Chdir(md.DirPath)
	file, err := os.Open(md.File)
	if err != nil {
		return md, err
	}
	defer file.Close()

	currentList := List{}
	md.Lists = []List{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Change the default title with the first h1
		if strings.HasPrefix(line, "# ") && md.Title == defaultTitle {
			md.Title = line[2:]
			continue
		}

		// Start a new list if a new h2 is found
		if strings.HasPrefix(line, "## ") {
			if currentList.Title != "" {
				md.Lists = append(md.Lists, currentList)
			}
			currentList = List{Title: line[3:]}
			continue
		}

		// Skip if the line is not a todo item
		if !strings.HasPrefix(line, "- [") {
			continue
		}

		item := TodoItem{}
		itemLabel := line[len("- [ ] "):]

		// If contains a link
		if strings.Contains(itemLabel, "](") {
			linkStart := strings.Index(itemLabel, "](")
			linkEnd := strings.Index(itemLabel, ")")
			item.Label = itemLabel[1:linkStart]
			item.Link = itemLabel[linkStart+2 : linkEnd]
			item = parseLink(item)
		} else { // Else
			item.Label = itemLabel
		}

		// Set the state of the item
		if strings.HasPrefix(line, "- [ ]") {
			item.State = checkbox_empty
		} else if strings.HasPrefix(line, "- [x]") {
			item.State = checkbox_done
		} else if strings.HasPrefix(line, "- [-]") {
			item.State = checkbox_cancel
		} else if strings.HasPrefix(line, "- [?]") {
			item.State = checkbox_outline
		} else if strings.HasPrefix(line, "- [~]") {
			item.State = checkbox_pending
		}

		// Add the item to the current list
		currentList.TodoItems = append(currentList.TodoItems, item)
	}

	// Add the last list
	if currentList.Title != "" {
		md.Lists = append(md.Lists, currentList)
	}

	if err := scanner.Err(); err != nil {
		return Md{}, err
	}
	return md, nil
}
