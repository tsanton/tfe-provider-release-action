package action

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func ReadFileToByteArray(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileContents, nil
}

func CreateRequestFromByteArray(method string, url string, fileBytes []byte, fileName string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(fileBytes))
	if err != nil {
		panic(err)
	}

	return req, nil
}
