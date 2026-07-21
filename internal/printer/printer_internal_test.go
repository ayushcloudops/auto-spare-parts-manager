package printer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// printerFunc adapts a function to the Printer interface for testing.
type printerFunc func(data []byte, name string) error

func (f printerFunc) Print(data []byte, name string) error { return f(data, name) }

// TestPrintFallsBackWhenPrintingFails is the regression test for the bug where a
// Mac with `lp` installed but no printer configured lost the receipt: printing
// fails, so we MUST fall back to the file instead of erroring.
func TestPrintFallsBackWhenPrintingFails(t *testing.T) {
	fallbackUsed := false
	p := &lpPrinter{
		run: func([]byte, string) error { return errors.New("lp: no destination") },
		fallback: printerFunc(func([]byte, string) error {
			fallbackUsed = true
			return nil
		}),
	}

	if err := p.Print([]byte("receipt"), ""); err != nil {
		t.Fatalf("expected nil (fell back to file), got %v", err)
	}
	if !fallbackUsed {
		t.Fatal("expected the file fallback to be used when printing fails")
	}
}

// TestPrintDoesNotFallBackOnSuccess: when printing succeeds, the fallback must
// NOT run (we don't want a stray file for every successful print).
func TestPrintDoesNotFallBackOnSuccess(t *testing.T) {
	fallbackUsed := false
	p := &lpPrinter{
		run: func([]byte, string) error { return nil },
		fallback: printerFunc(func([]byte, string) error {
			fallbackUsed = true
			return nil
		}),
	}

	if err := p.Print([]byte("receipt"), ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fallbackUsed {
		t.Fatal("fallback should not run when printing succeeds")
	}
}

// TestFilePrinterWritesReceipt verifies the fallback actually persists the bytes.
// HOME is redirected to a temp dir so we don't write into the real app data dir.
func TestFilePrinterWritesReceipt(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)                       // macOS/Linux: UserConfigDir derives from HOME
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))

	if err := (&filePrinter{}).Print([]byte("hello-receipt"), ""); err != nil {
		t.Fatalf("file printer: %v", err)
	}

	// Find the written .bin somewhere under the temp home.
	var found string
	_ = filepath.Walk(tmp, func(path string, info os.FileInfo, _ error) error {
		if info != nil && !info.IsDir() && filepath.Ext(path) == ".bin" {
			found = path
		}
		return nil
	})
	if found == "" {
		t.Fatal("expected a receipt .bin file to be written")
	}
	data, _ := os.ReadFile(found)
	if string(data) != "hello-receipt" {
		t.Fatalf("receipt content mismatch: %q", data)
	}
}
