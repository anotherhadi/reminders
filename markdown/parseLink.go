package markdown

import (
	"bufio"
	"os"
	"strings"
)

func parseLink(item TodoItem) TodoItem {
	file, err := os.Open(item.Link)
	if err != nil {
		return item
	}
	defer file.Close()

	inMetadata := false
	firstLine := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "---") && firstLine {
			inMetadata = true
			continue
		}
		firstLine = false
		if strings.HasPrefix(strings.TrimSpace(line), "---") {
			inMetadata = false
			continue
		}
		if inMetadata {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "importance: ") {
				item.Importance = strings.TrimPrefix(line, "importance: ")
				if item.Importance == "low" {
					item.Importance = importance_low
				} else if item.Importance == "medium" {
					item.Importance = importance_medium
				} else if item.Importance == "high" {
					item.Importance = importance_high
				} else {
					item.Importance = importance_none
				}
			}
			continue
		}
		item.Content += line + "\n"

	}

	return item
}
