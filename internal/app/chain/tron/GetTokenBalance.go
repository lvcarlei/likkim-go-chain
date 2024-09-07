package tron

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// 定义 TRC-20 代币信息结构体

type TRC20Token struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

// 定义 TRC-10 代币信息结构体

type TRC10Token struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

// 定义主链信息结构体

type AccountResource struct {
	EnergyWindowOptimized                     bool  `json:"energy_window_optimized"`
	AcquiredDelegatedFrozenV2BalanceForEnergy int64 `json:"acquired_delegated_frozenV2_balance_for_energy"`
	EnergyUsage                               int64 `json:"energy_usage"`
	LatestConsumeTimeForEnergy                int64 `json:"latest_consume_time_for_energy"`
	EnergyWindowSize                          int64 `json:"energy_window_size"`
}

// 定义 TRC-20 信息结构体
type TRC20Response struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

// 定义主链响应结构体
type AccountResponse struct {
	TRC20           []map[string]string `json:"trc20"`
	AssetV2         []TRC10Token        `json:"assetV2"`
	AccountResource AccountResource     `json:"account_resource"`
	Balance         int64               `json:"Balance"`
}

func GetTokenBalance(mainAddress string) ([]map[string]interface{}, error) {
	url := MainnetRPCEndpoint + fmt.Sprintf("/v1/accounts/%s", mainAddress)
	var tokenList []map[string]interface{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to fetch token list: %v", err)
		return tokenList, nil
	}
	defer resp.Body.Close()
	// 创建 AccountResponse 实例
	var result struct {
		Data []AccountResponse `json:"data"`
	}
	// 解析 JSON 数据
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return tokenList, nil
	}
	trxBalance := result.Data[0].Balance
	tokenList = append(tokenList, map[string]interface{}{
		"symbol":          "TRX",
		"balance":         ConvertToReadableAmount(trxBalance, 6),
		"isNative":        true,
		"name":            "TRON",
		"protocolType":    "",
		"contractAddress": "",
	})
	for _, account := range result.Data {
		for _, token := range account.TRC20 {
			for address, balance := range token {
				tokenInfo := FetchTokenList(address, "TRC20")
				tokenList = append(tokenList, map[string]interface{}{
					"symbol":          tokenInfo["symbol"],
					"balance":         ConvertToReadableAmount(balance, tokenInfo["decimals"]),
					"isNative":        false,
					"name":            tokenInfo["name"],
					"protocolType":    tokenInfo["protocolType"],
					"contractAddress": tokenInfo["address"],
				})
			}
		}
		for _, trc10token := range account.AssetV2 {
			tokenInfo := FetchTokenList(trc10token.Key, "TRC10")
			tokenList = append(tokenList, map[string]interface{}{
				"symbol":          tokenInfo["symbol"],
				"balance":         ConvertToReadableAmount(trc10token.Value, tokenInfo["decimals"]),
				"isNative":        false,
				"name":            tokenInfo["name"],
				"protocolType":    tokenInfo["protocolType"],
				"contractAddress": tokenInfo["address"],
			})
		}
	}
	return tokenList, nil

}
