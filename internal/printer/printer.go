package printer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"autoshop/internal/config"
)

// Printer sends a raw byte stream to a physical or virtual printer.
type Printer interface {
	Print(data []byte, printerName string) error
}

// New returns the best available Printer for the current OS.
//
// macOS/Linux: pipes raw ESC/POS to the CUPS `lp` command.
// Windows: raw spooling via the print API (winspool WritePrinter) is the
// intended implementation and slots in here without touching callers; until it
// is added, the file fallback records each receipt so nothing is lost.
func New() Printer {
	switch runtime.GOOS {
	case "linux", "darwin":
		return &lpPrinter{}
	default:
		return &filePrinter{}
	}
}

// lpPrinter pipes bytes to CUPS `lp -o raw`. Falls back to a file if lp is
// missing (e.g. a machine with no printer configured).
type lpPrinter struct{}

func (p *lpPrinter) Print(data []byte, printerName string) error {
	lp, err := exec.LookPath("lp")
	if err != nil {
		return (&filePrinter{}).Print(data, printerName)
	}
	args := []string{"-o", "raw"}
	if printerName != "" {
		args = append([]string{"-d", printerName}, args...)
	}
	cmd := exec.Command(lp, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := stdin.Write(data); err != nil {
		_ = stdin.Close()
		return err
	}
	if err := stdin.Close(); err != nil {
		return err
	}
	return cmd.Wait()
}

// filePrinter writes the raw receipt to a .bin file under <appdata>/receipts.
// It is the universal fallback and useful for verifying the ESC/POS output.
type filePrinter struct{}

func (p *filePrinter) Print(data []byte, _ string) error {
	dir, err := config.AppDataDir()
	if err != nil {
		return err
	}
	recDir := filepath.Join(dir, "receipts")
	if err := os.MkdirAll(recDir, 0o755); err != nil {
		return err
	}
	name := fmt.Sprintf("receipt-%d.bin", time.Now().UnixNano())
	return os.WriteFile(filepath.Join(recDir, name), data, 0o644)
}
