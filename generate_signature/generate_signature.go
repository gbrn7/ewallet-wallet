package generate_signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"ewallet-wallet/constants"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GenerateSignature(endpoint, clientID string, now time.Time, method string, strPayload string) string {
	secretKey := constants.MappingClient[clientID]
	timestamp := time.Now().Format(time.RFC3339)

	if method != http.MethodGet {
		re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
		strPayload = re.ReplaceAllString(strPayload, "")
		strPayload = strings.ToLower(strPayload) + timestamp + endpoint
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(strPayload))
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	fmt.Println(generatedSignature)

	return generatedSignature
}
