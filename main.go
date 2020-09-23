package main

import (
	"flag"
	"log"
	"os"
)

// DefaultDownloadPath is temp folder for persist server file.
const DefaultDownloadPath = "./frida_push_cache"

var DefaultVersion = "12.11.17"

func init() {
	// Allow user to overwrite version from env.
	if v := os.Getenv("FRIDA_VERSION"); v != "" {
		DefaultVersion = v
	}
}

func main() {
	var (
		device  string
		version string
		force   bool
	)
	flag.StringVar(&device, "d", "pixel_2_api_281", "device name")
	flag.BoolVar(&force, "f", false, "force download")
	flag.StringVar(&version, "version", DefaultVersion, "frida version")
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

	log.Printf("use frida-server: %s-%s\n", arch, version)

	// Download
	outfile, err := adb.DownloadAndExtract(DefaultDownloadPath, version, arch)
	if err != nil {
		log.Println("download failed:", err)
	}

	// push to device
	if err := adb.Push(outfile); err != nil {
		log.Fatalln("push to device failed:", err)
	}
	log.Println("FINSHED", outfile)

}
