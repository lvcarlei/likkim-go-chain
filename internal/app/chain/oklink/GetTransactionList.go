package oklink

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TransactionResponse struct {
	Code string                `json:"code"`
	Msg  string                `json:"msg"`
	Data []TransactionRespData `json:"data"`
}
type Transaction struct {
	TxId            string `json:"txId"`
	BlockHash       string `json:"blockHash"`
	Height          string `json:"height"`
	TransactionTime string `json:"transactionTime"`
	From            string `json:"from"`
	To              string `json:"to"`
	Amount          string `json:"amount"`
	//IsToContract         bool   `json:"isToContract"`
	TransactionSymbol    string `json:"transactionSymbol"`
	TxFee                string `json:"txFee"`
	State                string `json:"state"`
	TokenContractAddress string `json:"tokenContractAddress"`
	//ChallengeStatus      string `json:"challengeStatus"`
}
type TransactionRespData struct {
	Page            string        `json:"page"`
	Limit           string        `json:"limit"`
	TotalPage       string        `json:"totalPage"`
	ChainFullName   string        `json:"chainFullName"`
	ChainShortName  string        `json:"chainShortName"`
	TransactionList []Transaction `json:"transactionLists"`
}

func GetTransactionList(chainShortName string, address string, page string,
	tokenContractAddress string, protocolType string) (OklinkResp TransactionRespData) {
	baseURL := OKLINK_ENPOINT + "/api/v5/explorer/address/transaction-list"
	limit := 20 // Adjust size as needed
	accessKey := OKLINK_ACCESS_KEY
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	// Construct request URL
	url := fmt.Sprintf("%s?page=%s&limit=%d&chainShortName=%s&address=%s&tokenContractAddress=%s&protocolType=%s",
		baseURL, page, limit, chainShortName, address, tokenContractAddress, protocolType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Add("Ok-Access-Key", accessKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch token list: %v", err)
	}
	defer resp.Body.Close()
	var result TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	// Check for success
	if result.Code != "0" {
		log.Printf("Oklink GetTransactionList API Error: %s", result.Msg)
		return TransactionRespData{}
	}
	return result.Data[0]
}
