package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

var ErrDeviceNotReady = errors.New("device not ready")
var ErrDeviceNotFound = errors.New("device not found")
var ErrNoDeviceInstalled = errors.New("no device installed")

type Emulator interface {
	// ListDevices shows all devices in system.
	ListDevices() ([]string, error)

	// Find a given device
	Find(name string) error
}

type emu struct {
	bins System
}

func NewEmulator(bins System) Emulator {
	return emu{bins}
}

func (e emu) ListDevices() ([]string, error) {
	var out []string
	cmd := exec.Command(e.bins.Emulator(), "-list-avds")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return out, err
	}
	out = strings.Split(buf.String(), "\n")
	out = stripSpace(out) // remove empty string from output
	return out, nil
}

func (e emu) Find(name string) error {
	dvs, err := e.ListDevices()
	if err != nil {
		return err
	}
	if len(dvs) == 0 {
		return ErrNoDeviceInstalled
	}
	for _, dv := range dvs {
		if dv == name {
			return nil
		}
	}
	return ErrDeviceNotFound
}

func stripSpace(arr []string) []string {
	var r []string
	for _, e := range arr {
		if e != "" {
			r = append(r, e)
		}
	}
	return r
}
