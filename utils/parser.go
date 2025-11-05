package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	x509 "github.com/google/certificate-transparency-go/x509"
)

func IsPrecertificate(cert *x509.Certificate) bool {
	for _, ext := range cert.Extensions {
		if x509.OIDExtensionCTPoison.Equal(ext.Id) && ext.Critical {
			return true // Precertificate.
		}
	}

	return false
}

func ParseCertString(strCert string) (*x509.Certificate, error) {
	// Step 1: Trim whitespace and special characters
	trimmed := strings.TrimSpace(strCert)

	// Step 2: Add missing padding if needed
	if missing := len(trimmed) % 4; missing != 0 {
		trimmed += strings.Repeat("=", 4-missing)
	}

	// Decode base64 certificate (DER format)
	certDER, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to base64 decode certificate"))
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to parse certificate"))
	}

	return cert, nil
}
