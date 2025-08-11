package main

import (
    "fmt"
    "log"

    "github.com/eiannone/keyboard" // using for demo
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
)

func main() {
    // Create the Fyne app and window
    a := app.New()
    w := a.NewWindow("Fyne + Keyboard")
    label := widget.NewLabel("Press ESC to quit keyboard reader")
    w.SetContent(label)
    w.Resize(fyne.NewSize(400, 300))

    // Start keyboard listener in a goroutine
    go func() {
        if err := keyboard.Open(); err != nil {
            log.Fatal(err)
        }
        defer keyboard.Close()

        for {
            char, key, err := keyboard.GetKey()
            if err != nil {
                log.Fatal(err)
            }

            fmt.Printf("Key pressed: %q, Special: %v\n", char, key)

            if key == keyboard.KeyEsc {
                fmt.Println("ESC pressed — stopping keyboard reader")
                return
            }
        }
    }()

    // Run Fyne app (blocks until window closes)
    w.ShowAndRun()
}
