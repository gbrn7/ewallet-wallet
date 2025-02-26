package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"ewallet-wallet/constants"
	"ewallet-wallet/helpers"
	"ewallet-wallet/internal/handler/wallet"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ExternalDependency struct {
	External wallet.External
}

func (d *ExternalDependency) MiddlewareValidateToken(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		fmt.Println("authorization empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	tokenData, err := d.External.ValidateToken(c.Request.Context(), auth)
	if err != nil {
		fmt.Printf("%v", err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()

	}

	c.Set("token", tokenData)

	c.Next()
}

func (d *ExternalDependency) MiddlewareSignatureValidation(c *gin.Context) {

	clientID := c.Request.Header.Get("Client-id")
	if clientID == "" {
		log.Println("Client-id empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	secretKey := constants.MappingClient[clientID]
	if secretKey == "" {
		log.Println("invalid client id")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	timestamp := c.Request.Header.Get("Timestamp")
	if timestamp == "" {
		log.Println("Timestamp empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	requestTime, err := time.Parse(time.RFC3339, timestamp)
	now := time.Now()

	if err != nil || now.Sub(requestTime) > 5*time.Minute {
		log.Println("invalid timestamp request")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	signature := c.Request.Header.Get("Signature")
	if signature == "" {
		log.Println("Signature empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	strPayload := ""

	if c.Request.Method != http.MethodGet {
		byteData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("failed to read request body")
			helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}
		copyBody := io.NopCloser(bytes.NewBuffer(byteData))
		c.Request.Body = copyBody

		endpoint := c.Request.URL.Path
		strPayload = string(byteData)
		re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
		strPayload = re.ReplaceAllString(strPayload, "")
		strPayload = strings.ToLower(strPayload) + timestamp + endpoint
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(strPayload))
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	if signature != generatedSignature {
		log.Printf("invalid signature, requested: %s, generated: %s \n", signature, generatedSignature)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	c.Set("client_id", clientID)
	c.Next()
}
