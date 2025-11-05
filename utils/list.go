package utils

import (
	"io"
	"net/http"

	"github.com/google/certificate-transparency-go/loglist3"
)

var (
	CTLogList *loglist3.LogList
)

func ScrapeCTLogList() error {
	req, err := http.NewRequest(http.MethodGet, loglist3.LogListURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// Read the response body as a string
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	CTLogList, err = loglist3.NewFromJSON(bytes)
	return err
}
