package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/lotusirous/frida-push/log"
)

func cacheDir() string {
	const base = "frida-push"
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(os.Getenv("HOME"), ".cache", base)
}

func main() {
	var (
		device     string
		remotePath string
	)
	flag.StringVar(&device, "d", "", "device name")
	flag.StringVar(&remotePath, "r", "/data/local/tmp/frida-server", "default remote path")
	flag.Parse()

	tools, err := LoadTools()
	if err != nil {
		log.Fatalln("main: load tools:", err)
	}

	dev := NewEmulator(tools)

	// load connected device
	if device != "" {
		device = "emulator-5554"
	}

	// if err := dev.Find(device); err != nil {
	// 	log.Fatalln("main: find device failed:", err)
	// }

	arch, err := dev.GetArch()
	if err != nil {
		log.Fatalln("main: get arch:", err)
	}

	version, err := tools.GetFridaToolVersion()
	if err != nil {
		log.Fatalln("main: find frida-tools version:", err)
	}
	log.Infoln("Found frida-tools version:", version)

	log.Infoln("Download and extract file to:", cacheDir())
	serverBin, err := DownloadAndExtract(tools, cacheDir(), version, arch)
	if err != nil {
		log.Fatalln("main: download and extract:", err)
	}
	log.Infoln("Downloaded path:", serverBin)

	log.Infoln("Switch to root")
	if err := dev.SwithToRoot(); err != nil {
		log.Fatalln("main: switch to root:", err)
	}
	log.Infoln("Push and execute")
	// push to device
	if err := dev.PushAndExecute(serverBin, remotePath); err != nil {
		log.Fatalln("main: push and execute:", err)
	}

	log.Infoln("DONE")

}
