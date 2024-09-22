package exchanges

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	base "github.com/hideyk/crypto-oms/base"
)

type BybitExchange base.Exchange

func (e BybitExchange) GetRequest(client *http.Client, url string, endpoint string, recv_window int, params string) []byte {
	now := time.Now()
	timestamp := now.UnixNano() / 1000000
	hmac256 := hmac.New(sha256.New, []byte(e.ApiKey))
	hmac256.Write([]byte(strconv.FormatInt(timestamp, 10) + e.ApiKey + strconv.Itoa(recv_window) + params))
	signature := hex.EncodeToString(hmac256.Sum(nil))
	request, err := http.NewRequestWithContext(context.Background(), "GET", url+endpoint+"?"+params, nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BAPI-API-KEY", e.ApiKey)
	request.Header.Set("X-BAPI-SIGN", signature)
	request.Header.Set("X-BAPI-TIMESTAMP", strconv.FormatInt(timestamp, 10))
	request.Header.Set("X-BAPI-SIGN-TYPE", "2")
	request.Header.Set("X-BAPI-RECV-WINDOW", strconv.Itoa(recv_window))
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Request Dump:\n%s", string(reqDump))
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	elapsed := time.Since(now).Seconds()
	fmt.Printf("\n%s took %v seconds \n", url+endpoint, elapsed)
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return body
}

func (e BybitExchange) PostRequest(client *http.Client, url string, endpoint string, recv_window int, request_body interface{}) []byte {
	now := time.Now()
	timestamp := now.UnixNano() / 1000000
	jsonData, err := json.Marshal(request_body)
	if err != nil {
		log.Fatal(err)
	}
	hmac256 := hmac.New(sha256.New, []byte(e.ApiKey))
	hmac256.Write([]byte(strconv.FormatInt(timestamp, 10) + e.ApiKey + strconv.Itoa(recv_window) + string(jsonData[:])))
	signature := hex.EncodeToString(hmac256.Sum(nil))
	request, err := http.NewRequestWithContext(context.Background(), "POST", url+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-BAPI-API-KEY", e.ApiKey)
	request.Header.Set("X-BAPI-SIGN", signature)
	request.Header.Set("X-BAPI-TIMESTAMP", strconv.FormatInt(timestamp, 10))
	request.Header.Set("X-BAPI-SIGN-TYPE", "2")
	request.Header.Set("X-BAPI-RECV-WINDOW", strconv.Itoa(recv_window))
	reqDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Request Dump:\n%s", string(reqDump))
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	elapsed := time.Since(now).Seconds()
	fmt.Printf("\n%s took %v seconds \n", url+endpoint, elapsed)
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return body
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
