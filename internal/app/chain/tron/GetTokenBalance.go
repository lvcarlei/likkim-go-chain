package tron

import (
	"encoding/json"
	"fmt"
	"go-wallet/internal/app/chain/helper"
	"go-wallet/internal/app/chain/oklink"
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

func GetTokenBalance(mainAddress string) (respData oklink.BalanceResp, err error) {
	url := helper.MainnetRPCEndpoint + fmt.Sprintf("/v1/accounts/%s", mainAddress)
	//var tokenList []map[string]interface{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to fetch token list: %v", err)
		return respData, nil
	}
	defer resp.Body.Close()
	// 创建 AccountResponse 实例
	var result struct {
		Data []AccountResponse `json:"data"`
	}
	// 解析 JSON 数据
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return respData, nil
	}
	if len(result.Data) == 0 {
		return respData, nil
	}
	var mainBalanceDeatil oklink.BlanceRespDetail
	trxBalance := result.Data[0].Balance
	mainBalanceDeatil.Address = mainAddress
	mainBalanceDeatil.Balance = helper.ConvertToReadableAmount(trxBalance, 6)
	mainBalanceDeatil.ContractAddress = ""
	mainBalanceDeatil.IsNative = true
	mainBalanceDeatil.Symbol = "TRX"
	respData.MainBalanceData = append(respData.MainBalanceData, mainBalanceDeatil)
	for _, account := range result.Data {
		for _, token := range account.TRC20 {
			for address, balance := range token {
				tokenInfo := FetchTokenList(address, "TRC20")
				var tokenDetail oklink.BlanceRespDetail
				tokenDetail.Address = ""
				tokenDetail.Symbol = tokenInfo["symbol"]
				tokenDetail.Balance = helper.ConvertToReadableAmount(balance, tokenInfo["decimals"])
				tokenDetail.Name = tokenInfo["name"]
				tokenDetail.ContractAddress = tokenInfo["address"]
				tokenDetail.IsNative = false
				respData.TokenBalanceData.Tokenlist = append(respData.TokenBalanceData.Tokenlist, tokenDetail)
			}
		}
		for _, trc10token := range account.AssetV2 {
			tokenInfo := FetchTokenList(trc10token.Key, "TRC10")
			var tokenDetail oklink.BlanceRespDetail
			tokenDetail.Address = ""
			tokenDetail.Symbol = tokenInfo["symbol"]
			tokenDetail.Balance = helper.ConvertToReadableAmount(trc10token.Value, tokenInfo["decimals"])
			tokenDetail.Name = tokenInfo["name"]
			tokenDetail.ContractAddress = tokenInfo["address"]
			tokenDetail.IsNative = false
			respData.TokenBalanceData.Tokenlist = append(respData.TokenBalanceData.Tokenlist, tokenDetail)
		}
	}
	respData.TokenBalanceData.Page = "1"
	respData.TokenBalanceData.TotalPage = "1"
	return respData, nil

}
