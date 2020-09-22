package main

import (
	"flag"
	"log"
)

// DefaultDownloadPath is temp folder for persist server file.
const DefaultDownloadPath = "./frida-push"

func main() {
	var (
		device string
		force  bool
	)
	flag.StringVar(&device, "d", "", "device name")
	flag.BoolVar(&force, "f", false, "force download")
	flag.Parse()

	emu, err := NewEmulator("")
	if err != nil {
		log.Fatalln("main: init emulator failed:", err)
	}

	dvs, err := emu.ListDevices()
	if err != nil {
		log.Fatalln("list device failed:", err)
	}

	if len(dvs) == 0 {
		log.Fatalln("No installed devices")
	}

}
