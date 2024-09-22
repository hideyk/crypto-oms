package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"strings"
	"time"

	base "github.com/hideyk/crypto-oms/base"
)

func NewBybitGetHeaders(apiKey string, params string) map[string]interface{} {
	timestamp := time.Now().UnixNano() / 1000000
	var sign_string strings.Builder

	sign_string.WriteString(strconv.FormatInt(timestamp, 10))
	sign_string.WriteString(apiKey)
	sign_string.WriteString(params)
	hmac256 := hmac.New(sha256.New, []byte(apiKey))

	return map[string]interface{}{
		"X-BAPI-API-KEY":   apiKey,
		"X-BAPI-TIMESTAMP": timestamp,
		"X-BAPI-SIGN":      "",
	}
}

func NewBybitPostHeaders(apiKey string, body interface{}) map[string]interface{} {
	timestamp := time.Now().UnixNano() / 1000000
	var sign_string strings.Builder

	sign_string.WriteString(strconv.FormatInt(timestamp, 10))
	sign_string.WriteString(apiKey)
	sign_string.WriteString(body)

	return map[string]interface{}{
		"X-BAPI-API-KEY":   apiKey,
		"X-BAPI-TIMESTAMP": timestamp,
		"X-BAPI-SIGN":      "",
	}
}

func NewBybitHeadersSecure(apiKey string, params map[string]interface{}, window time.Duration) map[string]interface{} {
	return map[string]interface{}{
		"X-BAPI-API-KEY":     apiKey,
		"X-BAPI-TIMESTAMP":   time.Now().UnixMilli(),
		"X-BAPI-SIGN":        "",
		"X-BAPI-RECV-WINDOW": "",
	}
}

func NewBybitTestClient(apiKey string) *base.Exchange {
	return &base.Exchange{
		Id:           "bybit",
		Name:         "Bybit",
		Countries:    []string{"Japan"},
		Version:      "v5",
		RateLimit:    20,
		Hostname:     "bybit.com",
		BaseEndpoint: "https://api-testnet.bybit.com",
		ApiKey:       apiKey,
	}
}

func NewBybitClient(apiKey string) *base.Exchange {
	return &base.Exchange{
		Id:           "bybit",
		Name:         "Bybit",
		Countries:    []string{"Japan"},
		Version:      "v5",
		RateLimit:    20,
		Hostname:     "bybit.com",
		BaseEndpoint: "https://api.bybit.com",
		ApiKey:       apiKey,
	}
}
