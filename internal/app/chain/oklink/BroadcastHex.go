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

func BroadcastHex(chainShortName string, hex string) (map[string]interface{}, error) {
	url := OKLINK_ENPOINT + "/api/v5/explorer/transaction/publish-tx"
	accessKey := "4783cafe-1710-48b4-ad2c-447a319a9f89"
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
	req.Header.Add("Ok-Access-Key", accessKey)

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
	data := map[string]interface{}{
		"txid": txid,
		"code": code,
	}

	return data, err

}
