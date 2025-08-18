package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}

	fName := os.Args[1]
	fmt.Println("Editing file:", fName)

	// Open (or create) the file in read-write mode, append if exists
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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

	// Initialize Fyne app
	a := app.New()
	w := a.NewWindow("Text Editor")
	w.SetFullScreen(true)

	// Use canvas.Text instead of Label to avoid bouncing
	message := canvas.NewText(string(data), theme.ForegroundColor())
	message.TextStyle.Monospace = true

	// Scrollable area for the text
	scroll := container.NewVScroll(message)
	scroll.SetMinSize(fyne.NewSize(800, 600))

	updateMessage := func() {
		message.Text = string(data)
		message.Refresh() // redraw text without re-layout
	}

	// Handle typed runes
	w.Canvas().SetOnTypedRune(func(r rune) {
		data = append(data, byte(r))
		if _, err := file.WriteString(string(r)); err != nil {
			fmt.Println("Error writing to file:", err)
		}
		updateMessage()
	})

	// Handle special keys
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		switch e.Name {
		case fyne.KeyEscape:
			w.Close()

		case fyne.KeyBackspace:
			if len(data) > 0 {
				data = data[:len(data)-1]

				// Truncate file and rewrite remaining data
				if err := file.Truncate(0); err != nil {
					fmt.Println("Error truncating file:", err)
				}
				if _, err := file.WriteAt(data, 0); err != nil {
					fmt.Println("Error writing file:", err)
				}

				updateMessage()
			}

		case fyne.KeySpace:
			data = append(data, ' ')
			if _, err := file.WriteString(" "); err != nil {
				fmt.Println("Error writing space:", err)
			}
			updateMessage()

		case fyne.KeyReturn:
			data = append(data, '\n')
			if _, err := file.WriteString("\n"); err != nil {
				fmt.Println("Error writing newline:", err)
			}
			updateMessage()
		}
	})

	w.SetContent(container.NewVBox(
		canvas.NewText("Type (ESC to exit):", theme.ForegroundColor()),
		scroll,
	))

	w.ShowAndRun()
}
