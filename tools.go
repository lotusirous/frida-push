package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

var versionRegex = regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)$`)

var ErrFridaNotFound = errors.New("frida-tools not found in your python environment")

func LoadTools() (SystemTool, error) {
	adb, err := exec.LookPath("adb")
	if err != nil {
		return nil, err
	}
	emu, err := exec.LookPath("emulator")
	if err != nil {
		return nil, err
	}
	uxz, err := exec.LookPath("unxz")
	if err != nil {
		return nil, err
	}

	fridaPS, err := exec.LookPath("frida-ps")
	if err != nil {
		return nil, err
	}

	return emuSys{
		emulator: emu,
		adb:      adb,
		unxz:     uxz,
		fridaPS:  fridaPS,
	}, nil
}

type emuSys struct {
	emulator string
	adb      string
	unxz     string
	fridaPS  string
}

func (s emuSys) Adb() string      { return s.adb }
func (s emuSys) Emulator() string { return s.emulator }
func (s emuSys) UnXZ() string     { return s.unxz }
func (s emuSys) FridaPS() string  { return s.fridaPS }

func (s emuSys) GetFridaToolVersion() (string, error) {
	v, err := exec.Command("python", "-c", "import frida; print(frida.__version__)").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tools: query frida client failed: %q %w", v, err)
	}
	out := bytes.TrimRight(v, "\n")
	version := string(out)

	if matched := versionRegex.MatchString(version); !matched {
		return "", ErrFridaNotFound
	}
	return string(out), nil
}
