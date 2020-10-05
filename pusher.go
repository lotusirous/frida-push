package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

const DefaultRemotePath = "/data/local/tmp/frida-server"

type Pusher interface {
	// Push a given binary to emulator server
	Push(string) error

	// GetArch return device OS
	GetArch() (string, error)

	// Download fetchs the binary and persists to given directory.
	Download(store, version, name string) (string, error)
}

type adb struct {
	bins System
}

func NewPusher(bins System) Pusher {
	return adb{bins}
}

func (a adb) GetArch() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("adb", "shell", "getprop ro.product.cpu.abi")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", ErrDeviceNotReady
	}
	return strings.TrimSuffix(out.String(), "\n"), nil
}

func (a adb) downloadURL(version, fname string) string {
	// https://github.com/frida/frida/releases/download/12.11.17/frida-server-12.11.17-android-x86_64.xz
	return fmt.Sprintf("https://github.com/frida/frida/releases/download/%s/%s", version, fname)
}

func (a adb) Download(dest, version, name string) (string, error) {
	if err := os.RemoveAll(dest); err != nil {
		return "", err
	}
	if _, err := os.Stat(dest); os.IsNotExist(err) { // mkdir if not exists
		_ = os.Mkdir(dest, os.ModePerm)
	}

	url := a.downloadURL(version, name)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	outfile := path.Join(dest, name)
	out, err := os.Create(outfile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	outfile, err = a.UnXZ(outfile)
	if err != nil {
		return "", err
	}
	return outfile, nil
}

func (a adb) UnXZ(path string) (string, error) {
	cmd := exec.Command(a.bins.UnXZ(), path)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(path, ".xz"), nil
}

func (a adb) Push(binfile string) error {
	// prepare frida-server file
	var out bytes.Buffer
	cmd := exec.Command("adb", "push", binfile, DefaultRemotePath)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	os.Stdout.Write(out.Bytes())
	return nil
}
