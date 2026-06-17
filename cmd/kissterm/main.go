package main

import (
    "log"
    "os"

    "github.com/you/go-kiss-tnc/internal/backend"
    "github.com/you/go-kiss-tnc/ui"
)

func main() {
    cfg := backend.DefaultConfig()
    // allow serial port from env for quick testing
    if p := os.Getenv("KISS_PORT"); p != "" {
        cfg.Port = p
    }
    be, err := backend.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    ui.Run(be)
}
