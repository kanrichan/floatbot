package xianqu

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
)

func Download(url, path string) (err error) {
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Set("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Set("Net-Type", "Wifi")
	resp, err := client.Do(reqest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)
	return nil
}

func DecodeBase64(hash, path string) (err error) {
	data, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(data)
	return nil
}
