package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	serverHost, ok := os.LookupEnv("SERVER_HOST")
	if !ok {
		serverHost = "localhost"
	}
	serverPort, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		serverPort = "3000"
	}
	saveFilePath, ok := os.LookupEnv("SAVE_FILE_PATH")
	if !ok {
		saveFilePath = "."
	}
	fmt.Println("saveFilePath: ", saveFilePath)

	url := fmt.Sprintf("http://%s:%s", serverHost, serverPort)
	fmt.Println("url: ", url)

	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		id := createNewPerson(url)
		getAllPersons(url)
		getPersonById(url, id)
		deletePersonById(url, id)
		getAllPersons(url)
		getRandomFile(url, saveFilePath)
		fmt.Println("======================================================")
	}
}

func createNewPerson(url string) string {

	var person = map[string]string{
		"name":   "mahmoud",
		"family": "kheyrati fard",
	}

	requestBody, err := json.Marshal(person)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/person", url), bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	request.Header.Set("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var resp map[string]string
	err = json.Unmarshal(all, &resp)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	id := resp["id"]
	fmt.Println("person created with id= ", id)
	return id
}
func getAllPersons(url string) {
	resp, err := http.Get(fmt.Sprintf("%s/person/all", url))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("all persons: ", string(all))

}
func getPersonById(url string, id string) {
	resp, err := http.Get(fmt.Sprintf("%s/person/%s", url, id))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("person by id: ", string(all))
}

func deletePersonById(url string, id string) {

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/person/%s", url, id), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resp map[string]string
	err = json.Unmarshal(all, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("person deleted: ", resp)
}

func getRandomFile(url string, saveFilePath string) {
	resp, err := http.Get(fmt.Sprintf("%s/randomFile", url))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("randomFile: ", string(all))
	checksum := resp.Header.Get("checksum")

	contentDisposition := resp.Header.Get("Content-Disposition")
	contentDisposition = strings.ReplaceAll(contentDisposition, "\"", "")
	a := strings.Split(contentDisposition, ";")
	b := strings.Split(a[1], "=")
	fileName := b[1]

	path := fmt.Sprintf("%s/%s", saveFilePath, fileName)
	err = ioutil.WriteFile(path, all, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	newFileCheckSum, err := getChecksumOfFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	checkSumResult := strings.Compare(newFileCheckSum, checksum)
	var result string
	if checkSumResult == 0 {
		result = "Matched"
	} else {
		result = "DO NOT Matched"
	}
	fmt.Println(fmt.Sprintf("file saved to %s. and checksums are %s", path, result))
}

func getChecksumOfFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(content)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}
