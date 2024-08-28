package markdown

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/rand"
)

func (md Md) ChangeState(section, label, state string) {
	_ = os.Chdir(md.DirPath)
	file, _ := os.Open(md.File)
	defer file.Close()
	randomString := strconv.Itoa(rand.Intn(100000))
	tempfile, _ := os.Create("/tmp/" + randomString + ".md")

	var sectionName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var toWrite string
		line := scanner.Text()
		toWrite = line
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			_, _ = tempfile.WriteString(toWrite + "\n")
			continue
		}

		// Skip if the line is not a todo item or a section
		if strings.HasPrefix(line, "## ") {
			sectionName = line[3:]
		} else if strings.HasPrefix(line, "- [") {
			if sectionName == section {
				itemLabel := line[len("- [ ] "):]
				if strings.Contains(itemLabel, "](") {
					linkStart := strings.Index(itemLabel, "](")
					itemLabel = itemLabel[1:linkStart]
				}
				if itemLabel == label {
					toWrite = strings.Replace(toWrite, "- [ ] ", "- ["+state+"] ", 1)
					toWrite = strings.Replace(toWrite, "- [x] ", "- ["+state+"] ", 1)
					toWrite = strings.Replace(toWrite, "- [-] ", "- ["+state+"] ", 1)
					toWrite = strings.Replace(toWrite, "- [~] ", "- ["+state+"] ", 1)
					toWrite = strings.Replace(toWrite, "- [?] ", "- ["+state+"] ", 1)
				}
			}
		}
		_, _ = tempfile.WriteString(toWrite + "\n")
	}

	// Replace the original file with the temporary File
	_ = os.Remove(md.File)
	_ = os.Rename("/tmp/"+randomString+".md", md.File)
}

func (md Md) ChangeLabel(section, oldLabel, newLabel string) {
	_ = os.Chdir(md.DirPath)
	file, _ := os.Open(md.File)
	defer file.Close()
	randomString := strconv.Itoa(rand.Intn(100000))
	tempfile, _ := os.Create("/tmp/" + randomString + ".md")

	var sectionName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var toWrite string
		line := scanner.Text()
		toWrite = line
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			_, _ = tempfile.WriteString(toWrite + "\n")
			continue
		}

		// Skip if the line is not a todo item or a section
		if strings.HasPrefix(line, "## ") {
			sectionName = line[3:]
		} else if strings.HasPrefix(line, "- [") {
			if sectionName == section {
				itemLabel := line[len("- [ ] "):]
				if strings.Contains(itemLabel, "](") {
					linkStart := strings.Index(itemLabel, "](")
					itemLabel = itemLabel[1:linkStart]
				}
				if itemLabel == oldLabel {
					toWrite = strings.Replace(toWrite, oldLabel, newLabel, 1)
				}
			}
		}
		_, _ = tempfile.WriteString(toWrite + "\n")
	}

	// Replace the original file with the temporary File
	_ = os.Remove(md.File)
	_ = os.Rename("/tmp/"+randomString+".md", md.File)
}
