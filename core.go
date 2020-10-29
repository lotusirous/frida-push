package main

type (

	// Devicer represents a connected device.
	Devicer interface {
		// SwithToRoot changes current adb user to root user.
		SwithToRoot() error

		// PushAndExecute sends a binary to connected device and execute it.
		PushAndExecute(src, dst string) error

		// GetArch retrieves device's CPU architecture.
		GetArch() (string, error)
	}

	// SystemTool centralizes required binaries from system.
	SystemTool interface {
		// Adb defines a absolute path to `adb` tool.
		Adb() string

		// Emulator defines a absolute path to `emulator` tool.
		Emulator() string

		// UnXZ defines a absolute path to `unxz` tool.
		UnXZ() string

		// GetFridaToolVersion retrieves frida version from current PYTHONPATH.
		GetFridaToolVersion() (string, error)

		// FridaPS defines a absolute path to `frida-ps` tool.
		FridaPS() string
	}
)
