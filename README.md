# go-kiss-tnc
A vibe-coded KISS TNC GUI for all platforms written in Go


go-kiss-tnc — starter scaffold

Build:
  go build ./cmd/kissterm

Run (Linux):
  KISS_PORT=/dev/ttyUSB0 ./kissterm

KISS over TCP with socat example (forward /dev/ttyUSB0 to TCP):
  socat -d -d TCP-LISTEN:8001,reuseaddr,fork FILE:/dev/ttyUSB0,raw,echo=0,nonblock

Notes:
- This scaffold is a minimal starting point: AX.25 parsing is simplified.
- Improve address parsing, error handling, concurrency and add UI features as needed.
