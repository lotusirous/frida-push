package main

type (
	Devicer interface {
		SwithToRoot() error

		// Find a given device
		Find(name string) error

		// ListDevices shows all devices in system.
		List() ([]string, error)

		// Push a given binary to emulator server
		PushAndExecute(src, dst string) error

		// GetArch return device OS
		GetArch() (string, error)
	}

	// SystemTool defines a requires binary from system
	SystemTool interface {
		// Adb specify adb tools from path
		Adb() string

		// Emulator specify adb tools from path
		Emulator() string

		// UnXZ extract the downloaded file
		UnXZ() string

		// FridaToolVersion
		GetFridaToolVersion() (string, error)

		// FridaPS to check whether frida-server is running
		FridaPS() string
	}
)
