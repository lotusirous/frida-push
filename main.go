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
		remotePath string
	)
	flag.StringVar(&remotePath, "r", "/data/local/tmp/frida-server", "default remote path")
	flag.Parse()

	tools, err := LoadTools()
	if err != nil {
		log.Fatalln("main: load tools:", err)
	}

	dev := NewEmulator(tools)

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
	if err := dev.PushAndExecute(serverBin, remotePath); err != nil {
		log.Fatalln("main: push and execute:", err)
	}

	log.Infoln("DONE")

}
