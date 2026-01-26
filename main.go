package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func loadHostFile() string {
	content, err := os.ReadFile("C:/Windows/System32/drivers/etc/hosts")
	if err != nil {
		fmt.Println("Error opening hosts file:", err)
		return ""
	}
	return string(content)
}

func prepareHostFile() (int, int) {
	hostStr := loadHostFile()
	lines := strings.Split(hostStr, "\n")
	startIndex := -1
	endIndex := -1

	for i, line := range lines {
		if strings.Contains(line, "# Host Rerouter Start") {
			startIndex = i
		}
		if strings.Contains(line, "# Host Rerouter End") {
			endIndex = i
		}
	}

	if startIndex != -1 && endIndex != -1 {
		return startIndex, endIndex
	}

	if startIndex == -1 {
		lines = append(lines, "# Host Rerouter Start")
		startIndex = len(lines) - 1
	}

	if endIndex == -1 {
		lines = append(lines, "# Host Rerouter End")
		endIndex = len(lines) - 1
	}

	newHostStr := strings.Join(lines, "\n")
	err := saveHostFileRaw(newHostStr)
	if err != nil {
		fmt.Println("Error saving hosts file:", err)
		return -1, -1
	}
	fmt.Println("put rerouter tags")
	return startIndex, endIndex
}

func saveHostFile(overrides []override) error {
	hostStr := loadHostFile()
	lines := strings.Split(hostStr, "\n")
	startIndex, endIndex := prepareHostFile()

	if startIndex == -1 || endIndex == -1 {
		return fmt.Errorf("Could not find Host Rerouter tags")
	}
	var newLines []string
	newLines = append(newLines, lines[:startIndex+1]...)
	for _, ovr := range overrides {
		line := ""
		if !ovr.enabled {
			line += "# "
		}
		line += fmt.Sprintf("%s %s", ovr.override, ovr.original)
		newLines = append(newLines, line)
	}
	newLines = append(newLines, lines[endIndex:]...)
	newHostStr := strings.Join(newLines, "\n")
	return saveHostFileRaw(newHostStr)
}

func saveHostFileRaw(content string) error {
	return os.WriteFile("C:/Windows/System32/drivers/etc/hosts", []byte(content), 0644)
}

func parseHostFile(content string) []override {
	fmt.Println("Parsing host file...")
	var overrides []override
	startIndex, endIndex := prepareHostFile()
	fmt.Println("Start index:", startIndex, "End index:", endIndex)
	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return overrides
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines[startIndex+1 : endIndex] {
		thisOverride := override{original: "", override: "", enabled: false}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			thisOverride.enabled = false
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimSpace(line)
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			thisOverride.override = parts[0]
			thisOverride.original = parts[1]
			overrides = append(overrides, thisOverride)
		}
	}
	fmt.Println("Parsed overrides:", overrides)
	return overrides
}

func resizeTable(table *widget.Table, width int) {

	fmt.Println("Resizing table to width:", width)

	remainingWidth := width

	table.SetColumnWidth(2, 40)
	remainingWidth -= 40 + 100

	table.SetColumnWidth(0, float32(remainingWidth/2))
	table.SetColumnWidth(1, float32(remainingWidth/2))

	return
}

type override struct {
	original string
	override string
	enabled  bool
}
type tableWrapper struct {
	*widget.Table
}

func (t *tableWrapper) Resize(s fyne.Size) {
	t.Table.Resize(s)
	resizeTable(t.Table, int(s.Width))
}

func littleTestFunction() {
	fmt.Println("Little test function executed")
	fmt.Println(loadHostFile())
	prepareHostFile()
	fmt.Println("After preparing hosts file:")
	fmt.Println(loadHostFile())
}

func main() {

	prepareHostFile()
	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetTitle("Host Rerouter")

	//overrides := []override{{original: "site.com", override: "127.0.0.1", enabled: false}}
	overrides := parseHostFile(loadHostFile())
	fmt.Println(loadHostFile())
	fmt.Println(overrides)

	table := widget.NewTable(
		//lenght
		func() (int, int) {
			return len(overrides), 3
		},
		//create
		func() fyne.CanvasObject {
			c := container.NewStack()
			c.Add(widget.NewLabel("fuck niggers"))
			return c
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			row, col := i.Row, i.Col
			o.(*fyne.Container).Objects = nil
			thisOverride := overrides[row]

			switch col {
			case 0:
				entry := widget.NewEntry()
				entry.SetText(thisOverride.original)
				entry.OnChanged = func(s string) {
					overrides[row].original = s
					thisOverride.original = s
				}

				o.(*fyne.Container).Add(entry)
			case 1:
				o.(*fyne.Container).Add(widget.NewLabel(thisOverride.override))
			case 2:
				o.(*fyne.Container).Add(widget.NewLabel(fmt.Sprintf("%v", thisOverride.enabled)))
			}
		},
	)

	table.ShowHeaderRow = true
	table.ShowHeaderColumn = true
	table.CreateHeader = func() fyne.CanvasObject {
		return container.NewStack(widget.NewLabel("000"))
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {

		//label := o.(*widget.Label)
		label := widget.NewLabel("")
		o.(*fyne.Container).Objects = nil
		o.(*fyne.Container).Add(label)

		if id.Row == -1 && id.Col == -1 {
			label.SetText("")
		} else if id.Row == -1 {
			texts := [3]string{"Original", "Override", "Enabled"}

			label.SetText(texts[id.Col])
		} else {
			label.SetText(fmt.Sprintf("%d", id.Row+1))
		}
	}

	btnContainer := container.NewVBox(
		widget.NewButton("Salvar", func() {
			fmt.Println("Salvando alterações...")
			err := saveHostFile(overrides)
			if err != nil {
				fmt.Println("Erro ao salvar o arquivo hosts:", err)
			} else {
				fmt.Println("Arquivo hosts salvo com sucesso.")
			}
		}),
		widget.NewButton("Novo", func() {
			overrides = append(overrides, override{original: "example.com", override: "127.0.0.1", enabled: true})
		}),
		widget.NewButton("Apagar", func() {
			//TODO: Pegar o indice selecionado da tabela (atualmente hardcoded)
			selecionado := 1
			if selecionado >= 0 && selecionado < len(overrides) {
				overrides = append(overrides[:selecionado], overrides[selecionado+1:]...)
				table.Refresh()
			}
		}),
		widget.NewButton("Arquivo", func() {
			//usa o cmd para abrir o arquivo hosts no notepad
			fmt.Println("Abrindo arquivo hosts no notepad...")
			cmd := exec.Command("notepad.exe", "C:/Windows/System32/drivers/etc/hosts")
			err := cmd.Start()
			if err != nil {
				fmt.Println("Erro ao abrir o arquivo hosts:", err)
			}

		}),
	)
	// resizeTable(table, int(w.Canvas().Size().Width))
	content := container.NewBorder(nil, nil, nil, btnContainer, &tableWrapper{Table: table})

	w.Resize(fyne.NewSize(600, 400))
	w.SetContent(content)
	w.ShowAndRun()
}
