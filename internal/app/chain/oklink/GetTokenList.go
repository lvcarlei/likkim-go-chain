package oklink

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

type Token struct {
	ContractAddress string `json:"tokenContractAddress"`
	Name            string `json:"tokenFullName"`
	Symbol          string `json:"token"`
	Decimals        string `json:"precision"`
	LogoURI         string `json:"logoUrl"`
	ProtocolType    string `json:"protocolType"`
}
type TokenListData struct {
	Page           string  `json:"page"`
	Limit          string  `json:"limit"`
	TotalPage      string  `json:"totalPage"`
	ChainFullName  string  `json:"chainFullName"`
	ChainShortName string  `json:"chainShortName"`
	TokenList      []Token `json:"tokenList"`
}
type TokenListResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []TokenListData `json:"data"`
}

func GetTokenInfo(chainShortName string, symbol string, protocolType string) map[string]string {
	key := fmt.Sprintf("%s%s:%s:%s", TokenInfoKey, chainShortName, symbol, protocolType)
	result, _ := getTokenInfoFromRedis(key)
	return result
}

func FetchTokenList(chainShortName string, tokenContractAddress string) {
	// Define base URL and initial page
	baseURL := "https://www.oklink.com/api/v5/explorer/token/token-list"
	page := 1
	limit := 50 // Adjust size as needed
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
		req.Header.Add("Ok-Access-Key", OKLINK_ACCESS_KEY)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to fetch token list: %v", err)
		}
		defer resp.Body.Close()

		// Parse the response
		var result TokenListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Fatalf("Failed to decode response: %v", err)
		}
		// Check for success
		if result.Code != "0" {
			log.Fatalf("Oklink Tocken list API Error: %s", result.Msg)
		}
		// Store tokens in Redis
		for _, token := range result.Data[0].TokenList {
			storeTokenListInRedisAsHash(chainShortName, token)
		}
		if strconv.Itoa(page) == result.Data[0].TotalPage || result.Data[0].TotalPage == "" {
			break
		}
		page++
		time.Sleep(100 * time.Millisecond)
	}

}
func storeTokenListInRedisAsHash(chainShortName string, token Token) {
	rdb := db.GetClient()
	key := fmt.Sprintf("%s%s:%s", TokenInfoKey, chainShortName, token.ContractAddress)
	err := rdb.HSet(context.Background(), key, map[string]interface{}{
		"name":            token.Name,
		"symbol":          token.Symbol,
		"decimals":        token.Decimals,
		"logoURI":         token.LogoURI,
		"contractAddress": token.ContractAddress,
		"protocolType":    token.ProtocolType,
	}).Err()
	if err != nil {
		log.Fatalf("Failed to store token in Redis: %v", err)
	}
	//key2 := fmt.Sprintf("%s%s:%s:%s", TokenInfoKey, chainShortName, token.Symbol, token.ProtocolType)
	//_, err = rdb.Do(context.Background(), "COPY", key, key2).Result()
	//if err != nil {
	//	log.Fatalf("Error copying key: %v", err)
	//}
}

func getTokenInfoFromRedis(tokenKey string) (map[string]string, error) {
	rdb := db.GetClient()
	result, err := rdb.HGetAll(context.Background(), tokenKey).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
