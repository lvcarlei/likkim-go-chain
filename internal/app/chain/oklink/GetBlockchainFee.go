package oklink

import (
	"context"
	"encoding/json"
	"fmt"
	"go-wallet/db"
	"log"
	"net/http"
	"time"
)

type FeeResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data []FeeDetail `json:"data"`
}

type FeeDetail struct {
	//ChainFullName         string `json:"chainFullName"`
	//ChainShortName        string `json:"chainShortName"`
	Symbol                string `json:"symbol"`
	BestTransactionFee    string `json:"bestTransactionFee"`
	BestTransactionFeeSat string `json:"bestTransactionFeeSat"`
	RecommendedGasPrice   string `json:"recommendedGasPrice"`
	RapidGasPrice         string `json:"rapidGasPrice"`
	StandardGasPrice      string `json:"standardGasPrice"`
	SlowGasPrice          string `json:"slowGasPrice"`
	BaseFee               string `json:"baseFee"`
	GasUsedRatio          string `json:"gasUsedRatio"`
}

func GetBlockchainFee(chainShortName string) (map[string]interface{}, error) {
	// 查询是否在缓存中
	rdb := db.GetClient()
	supportKey := fmt.Sprintf("support:blockchain:%s", chainShortName)
	chain, _ := rdb.HGetAll(context.Background(), supportKey).Result()
	if len(chain) == 0 {
		return map[string]interface{}{
			"code": "1",
			"msg":  "blockchain not support",
		}, nil
	}
	feeKey := fmt.Sprintf("fee:oklink:%s", chainShortName)
	redisData, err := rdb.HGetAll(context.Background(), feeKey).Result()

	if err != nil {
		log.Printf("Failed to get fee from redis: %v", err)
	}
	if len(redisData) > 0 {
		return map[string]interface{}{
			"code": "0",
			"data": redisData,
		}, nil
	}
	baseURL := OKLINK_ENPOINT + "/api/v5/explorer/blockchain/fee"
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	// Construct request URL
	url := fmt.Sprintf("%s?chainShortName=%s",
		baseURL, chainShortName)
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
	var result FeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	// Check for success
	if result.Code != "0" {
		log.Printf("Oklink Blockchain Fee API Error: %s", result.Msg)
	}
	// Check for success
	feeDetail := FeeDetail{}
	code := "0"
	// Check for success
	if result.Code == "0" {
		feeDetail = result.Data[0]
	} else {
		log.Printf("Chain fee Oklink API Error: %s, code: %s", result.Msg, result.Code)
		code = "1"
	}
	redisErr := rdb.HSet(context.Background(), feeKey, map[string]interface{}{
		"bestTransactionFee":    feeDetail.BestTransactionFee,
		"symbol":                feeDetail.Symbol,
		"bestTransactionFeeSat": feeDetail.BestTransactionFeeSat,
		"recommendedGasPrice":   feeDetail.RecommendedGasPrice,
		"rapidGasPrice":         feeDetail.RecommendedGasPrice,
		"standardGasPrice":      feeDetail.StandardGasPrice,
		"slowGasPrice":          feeDetail.SlowGasPrice,
		"gasUsedRatio":          feeDetail.GasUsedRatio,
		"baseFee":               feeDetail.BaseFee,
	}).Err()

	if redisErr != nil {
		log.Fatalf("Failed to store chain fee in Redis: %v", redisErr)
	}
	rdb.Expire(context.Background(), feeKey, 24*7*time.Hour)
	return map[string]interface{}{
		"code": code,
		"data": feeDetail,
	}, err

}
