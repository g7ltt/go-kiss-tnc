package serial

import (
    "io"
    "time"

    "github.com/tarm/serial"
)

type Port struct {
    c *serial.Port
}

func Open(name string, baud int) (*Port, error) {
    cfg := &serial.Config{Name: name, Baud: baud, ReadTimeout: time.Millisecond * 500}
    p, err := serial.OpenPort(cfg)
    if err != nil {
        return nil, err
    }
    return &Port{c: p}, nil
}

func (p *Port) Read(b []byte) (int, error) { return p.c.Read(b) }
func (p *Port) Write(b []byte) (int, error){ return p.c.Write(b) }
func (p *Port) Close() error { return p.c.Close() }
