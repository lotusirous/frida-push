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
	// https://github.com/frida/frida/releases/download/12.11.17/frida-server-12.11.17-android-x86_64.xz
	serverBin := fmt.Sprintf("frida-server-%s-android-%s", version, arch)
	absServerBin := path.Join(cache, serverBin)
	serverZipFile := serverBin + ".xz"
	targetUrl := fmt.Sprintf("https://github.com/frida/frida/releases/download/%s/%s", version, serverZipFile)

	// if exist
	if _, err := os.Stat(absServerBin); err == nil {
		return absServerBin, nil
	}

	// prepare temp dir to extract file
	outfile := path.Join(cache, serverZipFile)
	wc, err := prepareDirectory(outfile)
	if err != nil {
		return "", err
	}
	defer wc.Close()

	if err := download(targetUrl, wc); err != nil {
		return "", err
	}

	// unzip file
	defer os.Remove(outfile)
	if out, err := exec.Command(tools.UnXZ(), outfile).CombinedOutput(); err != nil {
		return "", fmt.Errorf("extract: %q %w", out, err)
	}

	return absServerBin, nil
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
