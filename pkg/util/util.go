package util

import (
	"encoding/json"
	"fmt"
	"math/big"

	"crypto/rand"
	"io/ioutil"
	"net/http"
)

type githubversion struct {
	TagName string `json:"tag_name"`
}

const allowedChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

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
func RandString(n int) (string, error) {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChar))))
		if err != nil {
			return "", err
		}
		b[i] = allowedChar[num.Int64()]
	}
	return string(b), nil
}
