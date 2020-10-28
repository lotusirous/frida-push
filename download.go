package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
)

func DownloadAndExtract(tools SystemTool, cache, version, arch string) (string, error) {
	svr := fmt.Sprintf("frida-server-%s-android-%s", version, arch)
	serverPath := path.Join(cache, svr)
	xzPath := serverPath + ".xz"

	// https://github.com/frida/frida/releases/download/12.11.17/frida-server-12.11.17-android-x86_64.xz
	targetUrl := fmt.Sprintf("https://github.com/frida/frida/releases/download/%s/%s.xz", version, svr)

	if _, err := os.Stat(serverPath); err == nil {
		return serverPath, nil
	}

	wc, err := prepareDirectory(xzPath)
	if err != nil {
		return "", err
	}
	defer wc.Close()

	if err := download(targetUrl, wc); err != nil {
		return "", err
	}

	// extract file
	defer os.Remove(xzPath)
	if out, err := exec.Command(tools.UnXZ(), xzPath).CombinedOutput(); err != nil {
		return "", fmt.Errorf("extract: %q %w", out, err)
	}

	return serverPath, nil
}

func prepareDirectory(dst string) (io.WriteCloser, error) {
	dir := path.Dir(dst)
	if _, err := os.Stat(dir); os.IsNotExist(err) { // mkdir if not exists
		_ = os.Mkdir(dir, os.ModePerm)
	}
	f, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func download(target string, wc io.WriteCloser) error {
	resp, err := http.DefaultClient.Get(target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(wc, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
