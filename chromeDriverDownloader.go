package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func main() {

	const chromeDriverZip = "chromeDriver_mac64.zip"

	var chromeDriverLatestVersion = getChromeDriverLastestVersion()

	var fileURL = "https://chromedriver.storage.googleapis.com/" + chromeDriverLatestVersion + "/chromedriver_mac64.zip"

	downloadFile(chromeDriverZip, fileURL, chromeDriverLatestVersion)

	unzip(chromeDriverZip, getUserHomeDir())

	removeZipFile(chromeDriverZip)
}

func getChromeDriverLastestVersion() string {
	resp, err := http.Get("http://chromedriver.storage.googleapis.com/LATEST_RELEASE")
	if err != nil {
		panic(err)
	}

	versionText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	return string(versionText)
}

func downloadFile(filepath string, url string, chromeDriverLatestVersion string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, copyErr := io.Copy(out, resp.Body)
	if copyErr != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("ChromeDriver zip version %s downloaded", chromeDriverLatestVersion))
}

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func unzip(archive, target string) {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		panic(err)
	}
	fmt.Println("ChromeDriver unzipped")

	if err := os.MkdirAll(target, 0755); err != nil {
		panic(err)
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			panic(err)
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			panic(err)
		}
	}
	fmt.Println(fmt.Sprintf("ChromeDriver moved to $HOME dir: %s", target))
}

func removeZipFile(chromeDriverZip string) {
	if err := os.Remove(chromeDriverZip); err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("ChromeDriver zip file %s removed", chromeDriverZip))
}
