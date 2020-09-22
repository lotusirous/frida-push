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
	"time"
)

const DefaultRemotePath = "/data/local/tmp/frida-server"

type Pusher interface {
	// Push a given binary to emulator server
	Push(string) error

	// GetArch return device OS
	GetArch() (string, error)

	// Download fetchs the binary and persists to given directory.
	DownloadAndExtract(store, version, arch string) (string, error)
}

type adb struct {
	bin        string
	httpClient *http.Client
}

func NewPusher(path string) (Pusher, error) {
	bin := "adb"
	if path != "" {
		bin = path
	}
	p, err := exec.LookPath(bin)
	if err != nil {
		return nil, err
	}
	return adb{bin: p, httpClient: &http.Client{Timeout: 10 * time.Second}}, nil
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
	return fmt.Sprintf("https://github.com/frida/frida/releases/download/%s/%s", version, fname)
}

func (a adb) DownloadAndExtract(dir, version, arch string) (string, error) {
	if err := os.RemoveAll(dir); err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, os.ModePerm)
	}
	fname := fmt.Sprintf("frida-server-%s-android-%s.xz", version, arch)
	req, err := http.NewRequest("GET", a.downloadURL(version, fname), nil)
	if err != nil {
		return "", err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	outfile := path.Join(dir, fname)
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
	cmd := exec.Command("unxz", path)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(path, ".xz"), nil
}

func (a adb) Push(binfile string) error {
	// prepare frida-server file

	var out bytes.Buffer
	fmt.Println(binfile, DefaultRemotePath)
	cmd := exec.Command("adb", "push", binfile, DefaultRemotePath)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	os.Stdout.Write(out.Bytes())
	return nil
}
