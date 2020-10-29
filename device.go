package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type emu struct{ bins SystemTool }

func NewEmulator(bins SystemTool) Devicer { return emu{bins} }

func (e emu) GetArch() (string, error) {
	out, err := exec.Command(e.bins.Adb(), "shell", "getprop ro.product.cpu.abi").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("device: exec query arch failed: %q %w", out, err)
	}
	out = bytes.TrimRight(out, "\n")
	return string(out), nil
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
