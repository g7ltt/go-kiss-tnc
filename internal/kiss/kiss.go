package kiss

// Simple KISS framing and deframing for common TNCs.
// KISS special characters
const (
    FEND byte = 0xC0
    FESC byte = 0xDB
    TFEND byte = 0xDC
    TFESC byte = 0xDD
)

// Frame type 0 = data
func Frame(data []byte) []byte {
    out := []byte{FEND, 0x00} // port 0, type data
    for _, b := range data {
        switch b {
        case FEND:
            out = append(out, FESC, TFEND)
        case FESC:
            out = append(out, FESC, TFESC)
        default:
            out = append(out, b)
        }
    }
    out = append(out, FEND)
    return out
}

func Deframe(stream []byte) ([][]byte, []byte) {
    var frames [][]byte
    buf := make([]byte, 0, len(stream))
    inFrame := false
    esc := false

    for _, b := range stream {
        if !inFrame {
            if b == FEND {
                inFrame = true
                buf = buf[:0]
            }
            continue
        }
        if esc {
            if b == TFEND {
                buf = append(buf, FEND)
            } else if b == TFESC {
                buf = append(buf, FESC)
            } else {
                // unknown escape, append raw
                buf = append(buf, b)
            }
            esc = false
            continue
        }
        if b == FESC {
            esc = true
            continue
        }
        if b == FEND {
            // end frame
            // skip kiss header byte if present (first byte)
            if len(buf) > 0 {
                frames = append(frames, append([]byte(nil), buf[1:]...))
            }
            inFrame = false
            buf = buf[:0]
            continue
        }
        buf = append(buf, b)
    }
    return frames, buf
}
