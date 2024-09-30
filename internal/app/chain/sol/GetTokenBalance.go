package sol

import (
	"context"
	"encoding/json"
	"fmt"
	"go-wallet/db"
	"go-wallet/internal/app/chain/oklink"
	"log"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTokenBalance(mainAddress string) (respData oklink.BalanceResp, err error) {
	// Create a new Solana client
	c := GetClient()
	var address solana.PublicKey = solana.MustPublicKeyFromBase58(mainAddress)
	balance, err := c.GetBalance(
		context.TODO(),
		address,
		"",
	)
	// Initialize the result slice
	//var tokenList []map[string]interface{}
	var mainBalanceDeatil oklink.BlanceRespDetail
	//var tokenBalance oklink.TokenBalance
	//var tokenBalanceDetail oklink.BlanceRespDetail
	if err != nil {
		log.Println("failed to get balance", err)
	} else {
		lamports := balance.Value
		lamportsPerSol := uint64(1000000000) // 1 SOL = 1,000,000,000 Lamports
		solBalance := float64(lamports) / float64(lamportsPerSol)
		solBalanceStr := strconv.FormatFloat(solBalance, 'f', -1, 64)
		mainBalanceDeatil.Address = mainAddress
		mainBalanceDeatil.Balance = solBalanceStr
		mainBalanceDeatil.ContractAddress = ""
		mainBalanceDeatil.IsNative = true
		mainBalanceDeatil.Symbol = "sol"
		respData.MainBalanceData = append(respData.MainBalanceData, mainBalanceDeatil)

	}

	// Get token accounts
	tokenAccounts, err := c.GetTokenAccountsByOwner(context.Background(), address, &rpc.GetTokenAccountsConfig{
		ProgramId: &solana.TokenProgramID},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingJSONParsed,
		})
	if err != nil {
		log.Println("failed to get SOL token account by owner: ", err)
	}

	// Get Redis client
	redisClient := db.GetClient()
	var parsedData map[string]interface{}
	// Iterate over token accounts
	for _, tokenAccount := range tokenAccounts.Value {
		// Get token metadata from Redis
		rawJSON := tokenAccount.Account.Data.GetRawJSON()
		if err != nil {
			log.Fatalf("SOL解析代币账户信息失败: %v", err)
		}
		if err := json.Unmarshal([]byte(rawJSON), &parsedData); err != nil {
			log.Fatalf("SOL解析 JSON 数据失败: %v", err)
		}
		data := parsedData["parsed"].(map[string]interface{})["info"]
		tokenAmount := data.(map[string]interface{})["tokenAmount"].(map[string]interface{})
		uiAmount := tokenAmount["uiAmount"].(float64)
		isNative := data.(map[string]interface{})["isNative"].(bool)
		mint := data.(map[string]interface{})["mint"].(string)
		//tokenAddress := data.(map[string]interface{})["owner"].(string)
		//tokenAddress := tokenAccount.Account.Owner.String()
		tokenAddress := tokenAccount.Pubkey.String()
		dataKey := baseSolKey + mint
		tokenResult, err := redisClient.HGetAll(context.TODO(), dataKey).Result()
		if err != nil {
			log.Printf("failed to get fields from Redis: %v", err)
		}
		symbol := tokenResult["symbol"]
		var tokenDetail oklink.BlanceRespDetail
		tokenDetail.Address = tokenAddress
		tokenDetail.Symbol = symbol
		tokenDetail.Balance = fmt.Sprintf("%f", uiAmount)
		tokenDetail.Name = tokenResult["name"]
		tokenDetail.ContractAddress = mint
		tokenDetail.IsNative = isNative
		respData.TokenBalanceData.Tokenlist = append(respData.TokenBalanceData.Tokenlist, tokenDetail)
	}
	respData.TokenBalanceData.Page = "1"
	respData.TokenBalanceData.TotalPage = "1"
	return respData, nil
}
