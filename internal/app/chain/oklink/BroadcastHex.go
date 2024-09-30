package oklink

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type BroadcastResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []BroadcastData `json:"data"`
}
type BroadcastData struct {
	ChainFullName  string `json:"chainFullName"`
	ChainShortName string `json:"chainShortName"`
	Txid           string `json:"txid"`
}
type BroadcastResult struct {
	Txid string
	Code string
}

func BroadcastHex(chainShortName string, hex string) (BroadcastResult, error) {
	url := OKLINK_ENPOINT + "/api/v5/explorer/transaction/publish-tx"
	payload := map[string]interface{}{
		"signedTx":       hex,
		"chainShortName": chainShortName,
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
	req.Header.Add("Ok-Access-Key", OKLINK_ACCESS_KEY)

	// 发送请求
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending oklink broadcast request:", err)
	}
	defer resp.Body.Close()
	// Parse the response
	var result BroadcastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	txid := ""
	code := "0"
	// Check for success
	if result.Code == "0" {
		txid = result.Data[0].Txid
	} else {
		log.Printf("Broadcast Oklink API Error: %s, code: %s", result.Msg, result.Code)
		code = "1"
	}
	data := BroadcastResult{}
	data.Txid = txid
	data.Code = code
	return data, err

}
