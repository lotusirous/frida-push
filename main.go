package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var versionRegex = regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)$`)

var ErrFridaNotFound = errors.New("frida not found in your python environment")

func cacheDir() string {
	const base = "frida-push"
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(os.Getenv("HOME"), ".cache", base)
}

func fidraVersion() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("python", "-c", "import frida; print(frida.__version__)")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	v := strings.Trim(buf.String(), "\n")
	if matched := versionRegex.MatchString(v); !matched {
		return "", ErrFridaNotFound
	}
	return v, nil
}

func main() {
	var (
		device string
	)
	flag.StringVar(&device, "d", "pixel_2_api_281", "device name")
	flag.Parse()

	bins, err := LoadBinaries()
	if err != nil {
		log.Fatalln("failed to load binary:", err)
	}

	emu := NewEmulator(bins)
	adb := NewPusher(bins)

	if err := emu.Find(device); err != nil {
		log.Fatalln("list device failed:", err)
	}

	arch, err := adb.GetArch()
	if err != nil {
		log.Fatalln("main: get arch failed:", err)
	}

	version, err := fidraVersion()
	if err != nil {
		log.Fatalln("main: find version failed:", err)
	}

	fridaFile := fmt.Sprintf("frida-server-%s-android-%s.xz", version, arch)
	log.Printf("download version: %s\n", fridaFile)

	serverBin, err := adb.Download(cacheDir(), version, fridaFile)
	if err != nil {
		log.Println("main: download failed", err)
	}

	// push to device
	if err := adb.Push(serverBin); err != nil {
		log.Fatalln("push to device failed:", err)
	}
	log.Println("DONE: ", serverBin)

}
