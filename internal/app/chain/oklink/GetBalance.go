package oklink

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TokenBalanceDetail struct {
	Token                string `json:"token"`
	HoldingAmount        string `json:"holdingAmount"`
	TotalTokenValue      string `json:"totalTokenValue"`
	Change24h            string `json:"change24h"`
	PriceUsd             string `json:"priceUsd"`
	TokenContractAddress string `json:"tokenContractAddress"`
}
type TokenBalanceDetailList struct {
	Page           string               `json:"page"`
	Limit          string               `json:"limit"`
	TotalPage      string               `json:"totalPage"`
	ChainFullName  string               `json:"chainFullName"`
	ChainShortName string               `json:"chainShortName"`
	TokenList      []TokenBalanceDetail `json:"tokenList"`
}
type TokenBalanceResponse struct {
	Code string                   `json:"code"`
	Msg  string                   `json:"msg"`
	Data []TokenBalanceDetailList `json:"data"`
}

type MainBalanceDetail struct {
	Balance       string `json:"balance"`
	TokenAmount   string `json:"tokenAmount"`
	Address       string `json:"address"`
	BalanceSymbol string `json:"balanceSymbol"`
	ChainFullName string `json:"chainFullName"`
}
type MainBalanceResponse struct {
	Code string              `json:"code"`
	Msg  string              `json:"msg"`
	Data []MainBalanceDetail `json:"data"`
}
type BalanceResp struct {
	MainBalanceData  []BlanceRespDetail
	TokenBalanceData TokenBalance
}
type TokenBalance struct {
	Page      string
	TotalPage string
	Tokenlist []BlanceRespDetail
}
type BlanceRespDetail struct {
	Symbol          string
	Balance         string
	IsNative        bool
	Name            string
	Address         string
	ContractAddress string
}

func GetBalance(address string, chainShortName string) (data BalanceResp, err error) {
	baseURL := OKLINK_ENPOINT + "/api/v5/explorer/address/address-summary"
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	// Construct request URL
	url := fmt.Sprintf("%s?chainShortName=%s&address=%s",
		baseURL, chainShortName, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Add("Ok-Access-Key", OKLINK_ACCESS_KEY)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to Oklink balance: %v", err)
	}
	defer resp.Body.Close()
	var result MainBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	// Check for success
	if result.Code != "0" {
		log.Printf("Oklink GetBalance API Error: %s", result.Msg)
	}
	for _, detail := range result.Data {
		tmpData := BlanceRespDetail{}
		tmpData.Symbol = detail.BalanceSymbol
		tmpData.Address = detail.Address
		tmpData.Balance = detail.Balance
		tmpData.IsNative = true
		tmpData.Name = detail.ChainFullName
		data.MainBalanceData = append(data.MainBalanceData, tmpData)
		tokenAmount, _ := strconv.ParseInt(detail.TokenAmount, 10, 64)
		if tokenAmount > 0 {
			tokenData := getTokenBalance(detail.Address, detail.ChainFullName)
			data.TokenBalanceData = tokenData
		}
	}
	return data, nil
}

func getTokenBalance(address string, chainShortName string) (data TokenBalance) {
	baseURL := OKLINK_ENPOINT + "/api/v5/explorer/address/address-balance-fills"
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	// Construct request URL
	url := fmt.Sprintf("%s?chainShortName=%s&address=%s",
		baseURL, chainShortName, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Add("Ok-Access-Key", OKLINK_ACCESS_KEY)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to Oklink balance: %v", err)
	}
	defer resp.Body.Close()
	var result TokenBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	// Check for success
	if result.Code != "0" {
		log.Printf("Oklink GetBalance API Error: %s", result.Msg)
	}
	for _, detail := range result.Data[0].TokenList {
		tmpData := BlanceRespDetail{}
		tmpData.Symbol = detail.Token
		tmpData.Address = ""
		tmpData.Balance = detail.HoldingAmount
		tmpData.IsNative = false
		tmpData.Name = detail.Token
		tmpData.ContractAddress = detail.TokenContractAddress
		data.Tokenlist = append(data.Tokenlist, tmpData)
	}
	data.Page = result.Data[0].Page
	data.TotalPage = result.Data[0].TotalPage

	return data

}
