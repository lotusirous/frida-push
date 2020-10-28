package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func extract(bin, path string) (string, error) {
	serverBin := strings.TrimSuffix(path, ".xz")
	if _, err := os.Stat(serverBin); !os.IsNotExist(err) {
		return serverBin, nil
	}

	if out, err := exec.Command(bin, path).CombinedOutput(); err != nil {
		return "", fmt.Errorf("extract: %q %w", out, err)
	}
	return serverBin, nil

}

func DownloadAndExtract(tools SystemTool, cache, version, arch string) (string, error) {
	// https://github.com/frida/frida/releases/download/12.11.17/frida-server-12.11.17-android-x86_64.xz
	serverFile := fmt.Sprintf("frida-server-%s-android-%s.xz", version, arch)
	targetUrl := fmt.Sprintf("https://github.com/frida/frida/releases/download/%s/%s", version, serverFile)

	// prepare temp dir
	outfile := path.Join(cache, serverFile)
	wc, err := prepareDirectory(outfile)
	if err != nil {
		return "", err
	}
	defer wc.Close()

	if err := download(targetUrl, wc); err != nil {
		return "", err
	}

	// extract
	defer os.Remove(outfile)
	bin, err := extract(tools.UnXZ(), outfile)
	if err != nil {
		return "", err
	}

	return bin, nil
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
