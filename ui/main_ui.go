package ui

import (
    "fmt"
    "log"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/you/go-kiss-tnc/internal/backend"
)

var be *backend.Backend

func Run(b *backend.Backend) {
    be = b
    a := app.NewWithID("go-kiss-tnc")
    w := a.NewWindow("go-kiss-tnc")
    txEntry := widget.NewEntry()
    txEntry.SetPlaceHolder("Type text to send via KISS/AX.25...")
    sendBtn := widget.NewButton("Send Raw", func() {
        data := []byte(txEntry.Text)
        if err := be.WriteRawKISS(data); err != nil {
            log.Println("send error:", err)
        }
    })
    logArea := widget.NewMultiLineEntry()
    logArea.Wrapping = fyne.TextWrapWord
    // start UI update goroutine
    go func() {
        for {
            select {
            case l := <-be.Log:
                a.SendNotification(&fyne.Notification{Title: "go-kiss-tnc", Content: l})
                logArea.SetText(logArea.Text + l + "\n")
            case rx := <-be.RX:
                s := fmt.Sprintf("[%s] %s -> %s : %s\n", rx.Timestamp.Format("15:04:05"), rx.Src, rx.Dst, rx.InfoText())
                logArea.SetText(logArea.Text + s)
            }
        }
    }()
    // controls
    openBtn := widget.NewButton("Open Serial", func() {
        if err := be.OpenSerial(); err != nil {
            log.Println("open error:", err)
        } else {
            be.StartTCP()
        }
    })
    closeBtn := widget.NewButton("Close", func() { be.Close() })
    top := container.NewVBox(container.NewHBox(txEntry, sendBtn), container.NewHBox(openBtn, closeBtn))
    content := container.NewBorder(top, nil, nil, nil, logArea)
    w.SetContent(content)
    w.Resize(fyne.NewSize(800, 600))
    w.ShowAndRun()
}
