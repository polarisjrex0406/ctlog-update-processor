package utils

import (
	"context"
	"fmt"

	"bitbucket.org/xoduxcrt/ctlog-update-processor/logger"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/msg"
	"github.com/CaliDog/certstream-go"
	json "github.com/goccy/go-json"
)

type (
	CTLogEntry struct {
		UpdateType UpdateType    `json:"update_type"`
		CertIndex  int64         `json:"cert_index,omitempty"`
		CertLink   string        `json:"cert_link,omitempty"`
		LeafCert   Certificate   `json:"leaf_cert"`
		Chain      []Certificate `json:"chain,omitempty"`
		Seen       float64       `json:"seen"`
		Source     Source        `json:"source"`
	}

	Source struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	Certificate struct {
		Subject            SubjectIssuer  `json:"subject"`
		Issuer             *SubjectIssuer `json:"issuer,omitempty"`
		Extensions         Extensions     `json:"extensions"`
		NotBefore          float64        `json:"not_before"`
		NotAfter           float64        `json:"not_after"`
		SerialNumber       string         `json:"serial_number"`
		Fingerprint        string         `json:"fingerprint"`
		AsDer              string         `json:"as_der,omitempty"`
		AllDomains         []string       `json:"all_domains,omitempty"`
		SignatureAlgorithm string         `json:"signature_algorithm,omitempty"`
	}

	SubjectIssuer struct {
		Aggregated   string  `json:"aggregated"`
		C            *string `json:"C,omitempty"`
		ST           *string `json:"ST,omitempty"`
		L            *string `json:"L,omitempty"`
		O            *string `json:"O,omitempty"`
		OU           *string `json:"OU,omitempty"`
		CN           *string `json:"CN,omitempty"`
		EmailAddress *string `json:"emailAddress,omitempty"`
	}

	Extensions struct {
		KeyUsage               string `json:"keyUsage,omitempty"`
		ExtendedKeyUsage       string `json:"extendedKeyUsage,omitempty"`
		BasicConstraints       string `json:"basicConstraints,omitempty"`
		SubjectKeyIdentifier   string `json:"subjectKeyIdentifier,omitempty"`
		AuthorityKeyIdentifier string `json:"authorityKeyIdentifier,omitempty"`
		AuthorityInfoAccess    string `json:"authorityInfoAccess,omitempty"`
		SubjectAltName         string `json:"subjectAltName,omitempty"`
		CertificatePolicies    string `json:"certificatePolicies,omitempty"`
		CRLDistributionPoints  string `json:"crlDistributionPoints,omitempty"`
		CTLPoisonByte          *bool  `json:"ctlPoisonByte,omitempty"`
	}

	UpdateType string
)

const (
	UpdateTypeX509LogEntry    = UpdateType("X509LogEntry")
	UpdateTypePrecertLogEntry = UpdateType("PrecertLogEntry")
)

func RTUpdate(ctx context.Context) {
	stream, errStream := certstream.CertStreamEventStream(false)
	for {
		select {
		case jq := <-stream:
			messageType, err := jq.String("message_type")
			if err != nil || messageType != "certificate_update" {
				continue
			}

			rawData, err := jq.Object("data")
			if err != nil {
				continue
			}

			bytes, err := json.Marshal(rawData)
			if err != nil {
				continue
			}

			var entry CTLogEntry
			if err = json.Unmarshal(bytes, &entry); err != nil {
				continue
			}

		case err := <-errStream:
			fmt.Println(err)
		case <-ctx.Done(): // Respond to graceful shutdown requests.
			msg.ShutdownWG.Done()
			logger.Logger.Info("Stopped LogConfigSyncer")
			return
		}
	}
}
