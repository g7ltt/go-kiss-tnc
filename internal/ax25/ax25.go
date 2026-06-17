package ax25

import (
    "encoding/hex"
    "strings"
    "time"
)

type Frame struct {
    Timestamp time.Time
    Src string
    Dst string
    Info []byte
    Raw []byte
}

func Parse(raw []byte) *Frame {
    // Minimal, forgiving parser: find addresses and payload.
    // This does not implement full validation; it's a working stub.
    f := &Frame{Timestamp: time.Now(), Raw: append([]byte(nil), raw...)}
    if len(raw) < 15 {
        f.Info = raw
        return f
    }
    // addresses are 7 bytes each with shifted ASCII; decode simple
    dst := decodeCallsign(raw[0:7])
    src := decodeCallsign(raw[7:14])
    f.Dst = dst
    f.Src = src
    // find info after control/pid (usually at byte 16)
    if len(raw) > 16 {
        f.Info = raw[16:]
    } else {
        f.Info = raw[14:]
    }
    return f
}

func decodeCallsign(b []byte) string {
    var s []byte
    for i := 0; i < 6 && i < len(b); i++ {
        ch := b[i] >> 1
        if ch == 0x20 { break }
        s = append(s, ch)
    }
    // SSID in upper bits
    if len(b) >= 7 {
        ssid := (b[6] >> 1) & 0x0F
        if ssid != 0 {
            return strings.TrimSpace(string(s)) + "-" + strconv.Itoa(int(ssid))
        }
    }
    return strings.TrimSpace(string(s))
}

func (f *Frame) InfoText() string {
    if len(f.Info) == 0 {
        return ""
    }
    // show printable or hex
    printable := true
    for _, c := range f.Info {
        if c < 32 || c > 126 {
            printable = false
            break
        }
    }
    if printable {
        return string(f.Info)
    }
    return strings.ToUpper(hex.EncodeToString(f.Info))
}
