package utils

import (
	"fmt"
	"os/exec"
)

type CTLog struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	EntryNumber int64    `json:"entry_number"`
	Timestamp   uint64   `json:"timestamp"` // milliseconds since epoch
	Certificate string   `json:"certificate"`
	Chain       []string `json:"chain"`
}

func Scrape(outputDir, url, filenamePrefix string, logId int, entryBeg, entrySize int64) error {
	args := []string{
		"--include-chains",
		"-o",
		fmt.Sprintf("%s%s_%d(%d-%d).json", outputDir, filenamePrefix, logId, entryBeg, entryBeg+entrySize-1),
		"-s",
		fmt.Sprintf("%d", entryBeg),
		"-n",
		fmt.Sprintf("%d", entrySize),
		url,
	}

	cmd := exec.Command("/root/_zzz_/scrape-ct-log/target/release/scrape-ct-log", args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
