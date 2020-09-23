package main

import (
	"os/exec"
)

// System centralizes requirements binary
type System interface {
	Adb() string
	Emulator() string
	UnXZ() string
}

func LoadBinaries() (System, error) {
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

	return emuSys{
		emulator: emu,
		adb:      adb,
		unxz:     uxz,
	}, nil
}

type emuSys struct {
	emulator string
	adb      string
	unxz     string
}

func (s emuSys) Adb() string      { return s.adb }
func (s emuSys) Emulator() string { return s.emulator }
func (s emuSys) UnXZ() string     { return s.unxz }
