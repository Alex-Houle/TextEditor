package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fmt"
	"os"
)

func main() {

	// read the file name 
	var fName string
	if len(os.Args) > 1 {
		fName = os.Args[1]
		fmt.Println("File name:", fName)
	} else {
		fmt.Println("No file name provided")
		os.Exit(1)
	}

	// load data from the file
	data, err := readFile(fName) 
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	


	// Create a new Fyne application
	a := app.New()
	w := a.NewWindow("Keyboard Input")
	w.Resize(fyne.NewSize(400, 300))

	message := widget.NewLabel(data)

	// 
	w.Canvas().SetOnTypedRune(func(r rune) {
		data += string(r)
		message.SetText(data)
	})

	// Handle special keys
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		switch e.Name {
		case fyne.KeyEscape:

			w.Close()
		case fyne.KeyBackspace:
			if len(data) > 0 {
				data = data[:len(data)-1]
				message.SetText(data)
			}
		case fyne.KeySpace:
			data += " "
			message.SetText(data)
		case fyne.KeyReturn:
			data += "\n"
			message.SetText(data)
		}
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Type anything (ESC to quit):"),
		message,
	))
	w.ShowAndRun()
}



/* 
This function reads the content of a file and returns it as a string.
If the file cannot be opened or read, it returns an error.
*/ 
func readFile(fName string) (string, error) {
	file, err := os.Open(fName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := os.ReadFile(fName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}