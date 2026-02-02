package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const DO_PRINT = false

func println(a ...any) {
	if DO_PRINT {
		fmt.Println(a...)
	}
}

func loadHostFile() string {
	content, err := os.ReadFile("C:/Windows/System32/drivers/etc/hosts")
	if err != nil {
		println("Error opening hosts file:", err)
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
		println("Error saving hosts file:", err)
		return -1, -1
	}
	println("put rerouter tags")
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
	println("Parsing host file...")
	var overrides []override
	startIndex, endIndex := prepareHostFile()
	println("Start index:", startIndex, "End index:", endIndex)
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
		} else {
			thisOverride.enabled = true
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			thisOverride.override = parts[0]
			thisOverride.original = parts[1]
			overrides = append(overrides, thisOverride)
		}
	}
	println("Parsed overrides:", overrides)
	return overrides
}

func resizeTable(table *widget.Table, width int) {

	println("Resizing table to width:", width)

	remainingWidth := width

	table.SetColumnWidth(2, 40)
	remainingWidth -= 40 + 100

	table.SetColumnWidth(0, float32(remainingWidth/2))
	table.SetColumnWidth(1, float32(remainingWidth/2))
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

func checkAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func main() {

	prepareHostFile()
	a := app.New()
	w := a.NewWindow("Host Rerouter")
	w.SetMaster()

	if checkAdmin() {
		w.SetTitle("Host Rerouter (Admin)")
	} else {
		w.SetTitle("Host Rerouter (SEM ADMIN - APENAS LEITURA)")
	}

	selectedEntryIndex := -1
	madeChange := false

	overrides := parseHostFile(loadHostFile())
	println(loadHostFile())
	println(overrides)

	var table *widget.Table
	var saveBtn *widget.Button
	var discardBtn *widget.Button

	updateButtons := func() {
		println("update buttons")
		println("madeChange:", madeChange)
		if madeChange {
			saveBtn.Importance = widget.HighImportance
			discardBtn.Importance = widget.MediumImportance
			discardBtn.Text = "Descartar"
		} else {
			saveBtn.Importance = widget.MediumImportance
			//discardBtn.Importance = widget.LowImportance
			discardBtn.Text = "Atualizar"
		}
		saveBtn.Refresh()
		discardBtn.Refresh()
	}

	table = widget.NewTable(
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
					madeChange = true
					updateButtons()
				}

				o.(*fyne.Container).Add(entry)
			case 1:
				entry := widget.NewEntry()
				entry.SetText(thisOverride.override)
				entry.OnChanged = func(s string) {
					overrides[row].override = s
					thisOverride.override = s
					madeChange = true
					updateButtons()
				}

				o.(*fyne.Container).Add(entry)
			case 2:
				checkbox := widget.NewCheck("", func(b bool) {
					overrides[row].enabled = b
					thisOverride.enabled = b
					madeChange = true
					updateButtons()
				})
				checkbox.Checked = thisOverride.enabled
				center := container.NewCenter(checkbox)

				o.(*fyne.Container).Add(center)

				//o.(*fyne.Container).Add(widget.NewLabel(fmt.Sprintf("%v", thisOverride.enabled)))
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
		//label := widget.NewLabel("")

		// btn := widget.NewButton("", nil)
		// o.(*fyne.Container).Objects = nil
		// o.(*fyne.Container).Add(btn)

		if id.Row == -1 && id.Col == -1 { //canto
			label := widget.NewLabel("")
			o.(*fyne.Container).Objects = nil
			o.(*fyne.Container).Add(label)
			return
		} else if id.Row == -1 { //topo
			texts := [3]string{"Original", "Override", "Enabled"}
			label := widget.NewLabel(texts[id.Col])
			o.(*fyne.Container).Objects = nil
			o.(*fyne.Container).Add(label)
			return
		} else { //esquerda
			btn := widget.NewButton(fmt.Sprintf("%d", id.Row+1), func() {
				fmt.Printf("Clicou no botao da linha %d\n", id.Row+1)
				selectedEntryIndex = id.Row
				table.Refresh()
			})

			if id.Row == selectedEntryIndex {
				btn.Importance = widget.HighImportance
			} else {
				btn.Importance = widget.MediumImportance
			}

			o.(*fyne.Container).Objects = nil
			o.(*fyne.Container).Add(btn)
			return
		}
	}

	saveBtn = widget.NewButton("Salvar", func() {
		println("Salvando alterações...")
		err := saveHostFile(overrides)
		if err != nil {
			println("Erro ao salvar o arquivo hosts:", err)
		} else {
			println("Arquivo hosts salvo com sucesso.")
		}
		madeChange = false
		updateButtons()

		if !checkAdmin() {
			dialog.ShowInformation("Aviso de Permissão", "O aplicativo não está sendo executado como administrador.\nAs alterações não puderam ser salvas no arquivo hosts.", w)
		}
		overrides = parseHostFile(loadHostFile())
		table.Refresh()
	})

	discardBtn = widget.NewButton("Descartar", func() {
		println("Descartando alterações...")
		overrides = parseHostFile(loadHostFile())
		table.Refresh()
		madeChange = false
		updateButtons()
	})
	updateButtons()

	btnContainer := container.NewVBox(
		saveBtn,
		discardBtn,
		widget.NewButton("Novo", func() {
			overrides = append(overrides, override{original: "example.com", override: "127.0.0.1", enabled: true})
			selectedEntryIndex = len(overrides) - 1
			madeChange = true
			updateButtons()
			table.Refresh()
			println("Nova entrada adicionada.")
		}),
		widget.NewButton("Apagar", func() {
			//TODO: Pegar o indice selecionado da tabela (atualmente hardcoded)
			if selectedEntryIndex == -1 {
				println("Nenhuma entrada selecionada para apagar.")
				return
			}

			selecionado := selectedEntryIndex
			println("Apagando entrada selecionada:", selecionado)
			if selecionado >= 0 && selecionado < len(overrides) {
				overrides = append(overrides[:selecionado], overrides[selecionado+1:]...)
				table.Refresh()
			}
			selectedEntryIndex = -1
			madeChange = true
			updateButtons()
		}),
		widget.NewButton("Arquivo", func() {
			//usa o cmd para abrir o arquivo hosts no notepad
			println("Abrindo arquivo hosts no notepad...")
			cmd := exec.Command("notepad.exe", "C:/Windows/System32/drivers/etc/hosts")
			err := cmd.Start()
			if err != nil {
				println("Erro ao abrir o arquivo hosts:", err)
			}

		}),
		widget.NewButton("Editar", func() {
			//Cria uma nova janela com um textarea com o conteudo do arquivo hosts
			editWindow := a.NewWindow("Editar manualmente hosts")
			hostContent := loadHostFile()
			textArea := widget.NewMultiLineEntry()
			textArea.SetText(hostContent)
			saveManualBtn := widget.NewButton("Salvar", func() {
				newContent := textArea.Text
				newContentLines := strings.Split(newContent, "\n")

				foundStart := false
				foundEnd := false
				for _, line := range newContentLines {
					if strings.Contains(line, "# Host Rerouter Start") {
						foundStart = true
					}
					if strings.Contains(line, "# Host Rerouter End") {
						foundEnd = true
					}
				}
				if !foundStart || !foundEnd {
					dialog.ShowError(fmt.Errorf("As tags '# Host Rerouter Start' e '# Host Rerouter End' devem estar presentes no arquivo."), editWindow)
					return
				}

				err := saveHostFileRaw(newContent)
				if err != nil {
					println("Erro ao salvar o arquivo hosts:", err)
				} else {
					println("Arquivo hosts salvo com sucesso.")
				}

				if !checkAdmin() {
					dialog.ShowInformation("Aviso de Permissão", "O aplicativo não está sendo executado como administrador.\nAs alterações não puderam ser salvas no arquivo hosts.", w)
				}

				overrides = parseHostFile(newContent)
				madeChange = false
				updateButtons()
				table.Refresh()
				editWindow.Close()

			})
			closeManualBtn := widget.NewButton("Fechar", func() {
				editWindow.Close()
				table.Refresh()
			})
			content := container.NewBorder(nil, container.NewHBox(saveManualBtn, closeManualBtn), nil, nil, container.NewScroll(textArea))
			editWindow.SetContent(content)
			editWindow.Resize(fyne.NewSize(600, 400))
			editWindow.Show()
			editWindow.RequestFocus()
			textArea.FocusGained()
		}),
		widget.NewButton("FAQ e Ajuda", func() {
			helpWindow := a.NewWindow("Ajuda")

			content, err := os.ReadFile("help.md")
			if err != nil {
				helpWindow.SetContent(widget.NewLabel("Erro ao carregar ajuda: " + err.Error()))
			} else {
				rich := widget.NewRichTextFromMarkdown(string(content))
				rich.Wrapping = fyne.TextWrapWord
				helpWindow.SetContent(container.NewScroll(rich))
			}

			helpWindow.Resize(fyne.NewSize(500, 400))
			helpWindow.Show()
		}),
		widget.NewButton("Creditos", func() {
			//abre o github do autor
			url := "https://github.com/Laggh"
			println("Abrindo pagina do autor no github...")
			cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
			err := cmd.Start()
			if err != nil {
				println("Erro ao abrir a pagina do autor:", err)
			}
		}),
	)

	// resizeTable(table, int(w.Canvas().Size().Width))
	content := container.NewBorder(nil, nil, nil, btnContainer, &tableWrapper{Table: table})

	if !checkAdmin() {
		dialog.ShowInformation("Aviso de Permissão", "O aplicativo não está sendo executado como administrador.\nAs alterações não poderão ser salvas no arquivo hosts.", w)
	}

	w.Resize(fyne.NewSize(600, 400))
	w.SetContent(content)
	w.ShowAndRun()

}
