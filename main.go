package main

import (
	"flag"
	"log"
)

// DefaultDownloadPath is temp folder for persist server file.
const DefaultDownloadPath = "./frida_push_cache"

func main() {
	var (
		device  string
		version string
		force   bool
	)
	flag.StringVar(&device, "d", "pixel_2_api_281", "device name")
	flag.BoolVar(&force, "f", false, "force download")
	flag.StringVar(&version, "version", "12.11.17", "frida version")
	flag.Parse()

	emu, err := NewEmulator("")
	if err != nil {
		log.Fatalln("main: init emulator failed:", err)
	}

	if err := emu.Find(device); err != nil {
		log.Fatalln("list device failed:", err)
	}

	adb, err := NewPusher("")
	if err != nil {
		log.Fatalln("main: init pusher failed:", err)
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
