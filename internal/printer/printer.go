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
// macOS/Linux: pipes raw ESC/POS to the CUPS `lp` command, and on ANY failure
// (lp missing, no printer configured, printer offline) falls back to saving the
// receipt to a file so it is never lost.
// Windows: raw spooling via the print API (winspool WritePrinter) is the
// intended implementation and slots in here without touching callers; until it
// is added, the file fallback records each receipt so nothing is lost.
func New() Printer {
	switch runtime.GOOS {
	case "linux", "darwin":
		return newLPPrinter()
	default:
		return &filePrinter{}
	}
}

// lpPrinter prints via CUPS `lp`, with a file fallback. The `run` and `fallback`
// fields are injectable so the fallback behaviour can be unit-tested without a
// real printer.
type lpPrinter struct {
	run      func(data []byte, printerName string) error
	fallback Printer
}

func newLPPrinter() *lpPrinter {
	return &lpPrinter{run: runLP, fallback: &filePrinter{}}
}

// Print attempts to print; if that fails for ANY reason it saves the raw receipt
// to a file so nothing is ever lost. It only returns an error if even the file
// fallback fails.
func (p *lpPrinter) Print(data []byte, printerName string) error {
	if err := p.run(data, printerName); err != nil {
		// Printing failed (no printer, offline, lp missing…). Don't lose the
		// receipt — persist it to the receipts folder instead.
		return p.fallback.Print(data, printerName)
	}
	return nil
}

// runLP pipes the bytes to CUPS `lp -o raw`. Returns an error if lp is
// unavailable or the print job is rejected (e.g. no destination configured).
func runLP(data []byte, printerName string) error {
	lp, err := exec.LookPath("lp")
	if err != nil {
		return err
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
		_ = cmd.Wait()
		return err
	}
	if err := stdin.Close(); err != nil {
		_ = cmd.Wait()
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
