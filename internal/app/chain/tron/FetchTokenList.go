package tron

import (
	"context"
	"encoding/json"
	"fmt"
	"go-wallet/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

var baseTronKey = "tron:token:"

//TokenInfo 结构体表示每个 Token 的信息

type Token struct {
	Address      string `json:"tokenContractAddress"`
	Name         string `json:"tokenFullName"`
	Symbol       string `json:"token"`
	Decimals     string `json:"precision"`
	LogoURI      string `json:"logoUrl"`
	ProtocolType string `json:"protocolType"`
}
type TokenListData struct {
	Page           string  `json:"page"`
	Limit          string  `json:"limit"`
	TotalPage      string  `json:"totalPage"`
	ChainFullName  string  `json:"chainFullName"`
	ChainShortName string  `json:"chainShortName"`
	TokenList      []Token `json:"tokenList"`
}
type ApiResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []TokenListData `json:"data"`
}

func FetchTokenList(address string, protocolType string) map[string]string {
	result, _ := getTokenInfoFromRedis(address)
	if len(result) == 0 { // 处理错误
		if protocolType == "TRC10" {
			curlToken10(address)
		} else if protocolType == "TRC20" {
			curlOkLink(address)
		}
	}
	result, _ = getTokenInfoFromRedis(address)
	return result
}

func curlToken10(id interface{}) {
	// 定义 TRC-10 代币的数据结构
	type trc10Token struct {
		ID           int64  `json:"id"`
		Abbr         string `json:"abbr"`
		Description  string `json:"description"`
		Name         string `json:"name"`
		Num          int64  `json:"num"`
		Precision    int64  `json:"precision"`
		URL          string `json:"url"`
		TotalSupply  int64  `json:"total_supply"`
		TrxNum       int64  `json:"trx_num"`
		VoteScore    int64  `json:"vote_score"`
		OwnerAddress string `json:"owner_address"`
		StartTime    int64  `json:"start_time"`
		EndTime      int64  `json:"end_time"`
	}

	// 定义响应的 Meta 信息
	type Meta struct {
		At       int64 `json:"at"`
		PageSize int64 `json:"page_size"`
	}

	// 定义完整的响应结构体
	type Response struct {
		Data    []trc10Token `json:"data"`
		Meta    Meta         `json:"meta"`
		Success bool         `json:"success"`
	}
	url := fmt.Sprintf("https://api.trongrid.io/v1/assets/%s", id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}
	if response.Success { // 处理成功响应
		var redisToken = Token{}
		for _, token := range response.Data {
			redisToken.Address = strconv.FormatInt(token.ID, 10)
			redisToken.Decimals = strconv.FormatInt(token.Precision, 10)
			redisToken.LogoURI = ""
			redisToken.Name = token.Name
			redisToken.Symbol = token.Abbr
			redisToken.ProtocolType = "TRC10"
			storeTokenListInRedisAsHash(redisToken)
		}
	}

}

func curlOkLink(tokenContractAddress string) {
	// Define base URL and initial page
	baseURL := "https://www.oklink.com/api/v5/explorer/token/token-list"
	page := 1
	limit := 50 // Adjust size as needed
	chainShortName := "TRON"
	accessKey := "4783cafe-1710-48b4-ad2c-447a319a9f89"
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	for {
		// Construct request URL
		url := fmt.Sprintf("%s?page=%d&limit=%d&chainShortName=%s&tokenContractAddress=%s", baseURL, page, limit, chainShortName, tokenContractAddress)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Add("Ok-Access-Key", accessKey)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to fetch token list: %v", err)
		}
		defer resp.Body.Close()

		// Parse the response
		var result ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Fatalf("Failed to decode response: %v", err)
		}
		// Check for success
		if result.Code != "0" {
			log.Fatalf("API Error: %s", result.Msg)
		}
		// Store tokens in Redis
		for _, token := range result.Data[0].TokenList {
			storeTokenListInRedisAsHash(token)
		}
		if strconv.Itoa(page) == result.Data[0].TotalPage || result.Data[0].TotalPage == "" {
			break
		}
		page++
		time.Sleep(100 * time.Millisecond)
	}

}
func storeTokenListInRedisAsHash(token Token) {
	rdb := db.GetClient()
	key := baseTronKey + token.Address
	err := rdb.HSet(context.Background(), key, map[string]interface{}{
		"name":         token.Name,
		"symbol":       token.Symbol,
		"decimals":     token.Decimals,
		"logoURI":      token.LogoURI,
		"address":      token.Address,
		"protocolType": token.ProtocolType,
	}).Err()
	if err != nil {
		log.Fatalf("Failed to store token in Redis: %v", err)
	}
}

func getTokenInfoFromRedis(address string) (map[string]string, error) {
	tokenKey := baseTronKey + address
	rdb := db.GetClient()
	result, err := rdb.HGetAll(context.Background(), tokenKey).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
