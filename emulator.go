package main

import (
	"bytes"
	"os/exec"
	"strings"
)

// DefaultEmulatorBin defines a default emulator binary from PATH.
const DefaultEmulatorBin = "emulator"

type Emulator interface {
	// ListDevices shows all devices in system.
	ListDevices() ([]string, error)

	// GetExec returns the executable file.
	GetExec() string
}

type emu struct {
	bin string
}

func NewEmulator(path string) (Emulator, error) {
	bin := DefaultEmulatorBin
	if path != "" {
		bin = path
	}
	p, err := exec.LookPath(bin)
	if err != nil {
		return nil, err
	}
	return emu{bin: p}, nil
}

func (e emu) ListDevices() ([]string, error) {
	var out []string
	cmd := exec.Command("ls", "-alh")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return out, err
	}
	out = strings.Split(buf.String(), "\n")
	return out, nil
}

func (e emu) GetExec() string { return e.bin }
