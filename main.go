package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}

	fName := os.Args[1]
	fmt.Println("Editing file:", fName)

	// Open file
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read initial content
	data, err := os.ReadFile(fName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	a := app.New()
	w := a.NewWindow("Text Editor")
	w.Resize(fyne.NewSize(800, 600))

	textColor := color.RGBA{R: 211, G: 198, B: 170, A: 255} // #D3C6AA

	lines := []fyne.CanvasObject{}
	for _, line := range strings.Split(string(data), "\n") {
		t := canvas.NewText(line, textColor)
		t.TextStyle.Monospace = true
		lines = append(lines, t)
	}
messageBox := container.NewVBox(lines...)

// Scrollable view of the VBox
scroll := container.NewVScroll(messageBox)
	cursor := 0
	cursorLabel := canvas.NewText(fmt.Sprintf("Cursor: %d", cursor), color.White)

	// Update functionc
	updateMessage := func() {
		lines := []fyne.CanvasObject{}
		for _, line := range strings.Split(string(data), "\n") {
			t := canvas.NewText(line, textColor)
			t.TextStyle.Monospace = true
			lines = append(lines, t)
		}
		messageBox.Objects = lines
		messageBox.Refresh()
		cursorLabel.Text = fmt.Sprintf("Cursor: %d", cursor)
		cursorLabel.Refresh()
	}


	// Insert character at cursor
	insertRune := func(r rune) {
		data = append(data[:cursor], append([]byte{byte(r)}, data[cursor:]...)...)
		cursor++
		writeFile(file, data)
		updateMessage()
	}

	// Delete character before cursor
	deleteRune := func() {
		if len(data) > 0 && cursor > 0 {
			data = append(data[:cursor-1], data[cursor:]...)
			cursor--
			writeFile(file, data)
			updateMessage()
		}
	}

	// Handle typing
	w.Canvas().SetOnTypedRune(func(r rune) {
		insertRune(r)
	})

	// Handle special keys
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		switch e.Name {
		case fyne.KeyEscape:
			w.Close()
			panic("exited by user")
		case fyne.KeyBackspace:
			deleteRune()
		case fyne.KeyReturn:
			insertRune('\n')
		case fyne.KeyLeft:
			if cursor > 0 {
				cursor--
				updateMessage()
			}
		case fyne.KeyRight:
			if cursor < len(data) {
				cursor++
				updateMessage()
			}
		}
	})

	// Layout: text with cursor counter at bottom
	content := container.NewBorder(nil, cursorLabel, nil, nil, scroll)

	w.SetContent(content)
	w.Canvas().Refresh(content)
	scroll.ScrollToTop()
	w.ShowAndRun()
}

// Helper to overwrite file with new data
func writeFile(file *os.File, data []byte) {
	if err := file.Truncate(0); err != nil {
		fmt.Println("Error truncating file:", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
	}
	if _, err := file.Write(data); err != nil {
		fmt.Println("Error writing file:", err)
	}
}
