package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var ErrDeviceNotFound = errors.New("device: not found")
var ErrNoDeviceInstalled = errors.New("device: no device installed")

type emu struct{ bins SystemTool }

func NewEmulator(bins SystemTool) Devicer { return emu{bins} }

func (e emu) GetArch() (string, error) {
	out, err := exec.Command(e.bins.Adb(), "shell", "getprop ro.product.cpu.abi").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("device: exec query arch failed: %q %w", out, err)
	}
	return strings.TrimSuffix(string(out), "\n"), nil
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

func (e emu) List() ([]string, error) {
	var list []string
	out, err := exec.Command(e.bins.Emulator(), "-list-avds").CombinedOutput()
	if err != nil {
		return list, fmt.Errorf("device: exec -list-avds failed: %w", err)
	}
	out = bytes.TrimRight(out, "\n")
	list = strings.Split(string(out), "\n")
	return stripSpace(list), nil
}

func (e emu) Find(name string) error {
	dvs, err := e.List()
	if err != nil {
		return err
	}

	if len(dvs) == 0 {
		return ErrNoDeviceInstalled
	}
	for _, n := range dvs {
		if n == name {
			return nil
		}
	}
	return ErrDeviceNotFound
}

func (e emu) SwithToRoot() error {
	_, err := exec.Command(e.bins.Adb(), "wait-for-device", "root").CombinedOutput()
	return err
}

func (e emu) PushAndExecute(src, dst string) error {
	if out, err := exec.Command(e.bins.Adb(), "push", src, dst).CombinedOutput(); err != nil {
		return fmt.Errorf("device: push failed: %q %w", out, err)
	}

	_ = exec.Command(e.bins.Adb(), "shell", "killall frida-server").Run()

	if out, err := exec.Command(e.bins.Adb(), "shell", "chmod 0755 "+dst).CombinedOutput(); err != nil {
		return fmt.Errorf("device: chmod failed: %q %w", out, err)
	}

	// Execute the program and write to pid file
	if err := exec.Command(e.bins.Adb(), "shell", dst).Start(); err != nil {
		return fmt.Errorf("device: exec failed:%w", err)
	}

	time.Sleep(2 * time.Second)
	out, err := exec.Command(e.bins.FridaPS(), "-U").CombinedOutput()
	if err != nil {
		return fmt.Errorf("device: frida-server is down: %q %w", out, err)
	}

	return nil
}
