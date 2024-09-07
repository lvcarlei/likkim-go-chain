package oklink

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type UTXO struct {
	Txid          string `json:"txid"`
	Height        string `json:"height"`
	BlockTime     string `json:"blockTime"`
	Address       string `json:"address"`
	UnspentAmount string `json:"unspentAmount"`
	Index         string `json:"index"`
}
type UTXOResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Page      string `json:"page"`
		Limit     string `json:"limit"`
		TotalPage string `json:"TotalPage"`
		UTXOList  []UTXO `json:"utxoList"`
	} `json:"data"`
}

func GetUTXO(chainShortName string, address string, page string) (map[string]interface{}, error) {
	baseURL := OKLINK_ENPOINT + "/api/v5/explorer/address/utxo"
	limit := 20 // Adjust size as needed
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	// Construct request URL
	url := fmt.Sprintf("%s?page=%d&limit=%d&chainShortName=%s&address=%s",
		baseURL, page, limit, chainShortName, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Add("Ok-Access-Key", OKLINK_ACCESS_KEY)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch token list: %v", err)
	}
	defer resp.Body.Close()
	var result UTXOResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	// Check for success
	utxoList := []UTXO{}
	code := "0"
	// Check for success
	if result.Code == "0" {
		utxoList = result.Data[0].UTXOList
	} else {
		log.Printf("Broadcast Oklink API Error: %s, code: %s", result.Msg, result.Code)
		code = "1"
	}
	return map[string]interface{}{
		"code": code,
		"data": utxoList,
	}, err
}
