package tron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-wallet/internal/app/chain/helper"
	"go-wallet/internal/app/chain/oklink"
	"log"
	"net/http"
	"time"
)

type BroadcastResponse struct {
	Result      bool   `json:"result"`
	Code        string `json:"code"`
	Txid        string `json:"txid"`
	Message     string `json:"message"`
	Transaction string `json:"transaction"`
}

type Transaction struct {
	RawData   RawData  `json:"raw_data"`
	Signature []string `json:"signature"`
}

type RawData struct {
	RefBlockBytes string     `json:"ref_block_bytes"`
	RefBlockHash  string     `json:"ref_block_hash"`
	Expiration    int64      `json:"expiration"`
	Contract      []Contract `json:"contract"`
}

type Contract struct {
	Type      string    `json:"type"`
	Parameter Parameter `json:"parameter"`
}

type Parameter struct {
	TypeURL string `json:"type_url"`
	Value   string `json:"value"`
}

func BroadcastHex(hex string) (oklink.BroadcastResult, error) {
	url := helper.MainnetRPCEndpoint + fmt.Sprintf("/wallet/broadcasthex")
	// 需要发送的JSON数据
	payload := map[string]interface{}{
		"transaction": hex,
	}
	// 将Go数据结构编码为JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
	}
	// 设置请求头，表明发送的数据是JSON
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	// 发送请求
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending tron broadcast request:", err)
	}
	defer resp.Body.Close()
	var result BroadcastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
		//return tokenList
	}
	code := "0"
	if !result.Result {
		log.Println("Failed to broadcast transaction:", result.Message, result.Code)
		code = "1"
	}
	data := oklink.BroadcastResult{}
	data.Txid = result.Txid
	data.Code = code
	return data, err
}
