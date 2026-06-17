package backend

import (
    "errors"
    "io"
    "net"
    "sync"
    "time"

    "github.com/you/go-kiss-tnc/internal/kiss"
    "github.com/you/go-kiss-tnc/internal/ax25"
    "github.com/you/go-kiss-tnc/internal/serial"
)

type Config struct {
    Port string
    Baud int
    TCPListen string
}

func DefaultConfig() *Config {
    return &Config{Port: "/dev/ttyUSB0", Baud: 9600, TCPListen: ":8001"}
}

type Backend struct {
    cfg *Config
    conn io.ReadWriteCloser
    mu sync.Mutex
    stop chan struct{}

    // channels for UI
    RX chan *ax25.Frame
    Log chan string
}

func New(cfg *Config) (*Backend, error) {
    b := &Backend{cfg: cfg, RX: make(chan *ax25.Frame, 256), Log: make(chan string, 512)}
    return b, nil
}

func (b *Backend) OpenSerial() error {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.conn != nil { return errors.New("already open") }
    p, err := serial.Open(b.cfg.Port, b.cfg.Baud)
    if err != nil { return err }
    b.conn = p
    b.stop = make(chan struct{})
    go b.readLoop()
    return nil
}

func (b *Backend) Close() error {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.conn == nil { return nil }
    close(b.stop)
    err := b.conn.Close()
    b.conn = nil
    return err
}

func (b *Backend) readLoop() {
    buf := make([]byte, 4096)
    var acc []byte
    for {
        select {
        case <-b.stop:
            return
        default:
        }
        n, err := b.conn.Read(buf)
        if err != nil {
            if err == io.EOF { time.Sleep(100*time.Millisecond); continue }
            b.Log <- "read error: " + err.Error()
            return
        }
        if n == 0 { continue }
        acc = append(acc, buf[:n]...)
        frames, rem := kiss.Deframe(acc)
        acc = rem
        for _, fr := range frames {
            axf := ax25.Parse(fr)
            b.RX <- axf
            b.Log <- "RX: " + axf.InfoText()
        }
    }
}

func (b *Backend) WriteRawKISS(data []byte) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.conn == nil { return errors.New("not open") }
    out := kiss.Frame(data)
    _, err := b.conn.Write(out)
    return err
}

// TCP server: accept and proxy to serial (simple single-client)
func (b *Backend) StartTCP() error {
    if b.cfg.TCPListen == "" { return nil }
    ln, err := net.Listen("tcp", b.cfg.TCPListen)
    if err != nil { return err }
    go func() {
        for {
            conn, err := ln.Accept()
            if err != nil { break }
            b.Log <- "TCP client connected"
            go b.handleTCP(conn)
        }
    }()
    return nil
}

func (b *Backend) handleTCP(c net.Conn) {
    defer c.Close()
    // simple loop: copy from serial to tcp and tcp to serial
    rw := struct {
        io.Reader
        io.Writer
    }{Reader: c, Writer: c}
    // copy tcp->serial
    go func() {
        io.Copy(b.conn, c)
    }()
    // copy serial->tcp
    io.Copy(c, b.conn)
}
