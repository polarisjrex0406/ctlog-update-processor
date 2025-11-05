package utils

import (
	"bytes"
	"crypto/sha256"
)

func ExtractLogName(keyID [sha256.Size]byte) string {
	for _, operator := range CTLogList.Operators {
		for _, log := range operator.Logs {
			if bytes.Equal(keyID[:], log.LogID[:]) {
				return log.Description
			}
		}
	}

	return ""
}
