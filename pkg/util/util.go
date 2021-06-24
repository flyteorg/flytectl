package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type githubversion struct {
	TagName string `json:"tag_name"`
}

const allowedChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GetRequest will get the response of get request
func GetRequest(baseURL, url string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("%v%v", baseURL, url))
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

// ParseGithubTag will parse github tag from github response
func ParseGithubTag(data []byte) (string, error) {
	var result = githubversion{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}
	return result.TagName, nil
}

// WriteIntoFile will write content in a file
func WriteIntoFile(data []byte, file string) error {
	err := ioutil.WriteFile(file, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

//RandString will return a random string of n length
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = allowedChar[rand.Intn(len(allowedChar))]
	}
	return string(b)
}
