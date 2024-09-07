package sol

import (
	"context"
	"log"
	"strconv"

	"go-wallet/db"

	"github.com/blocto/solana-go-sdk/common"
)

func GetTokenBalance(mainAddress string) ([]map[string]interface{}, error) {
	// Create a new Solana client
	c := GetClient()
	// Get balance in lamports
	balance, err := c.GetBalance(
		context.TODO(),
		mainAddress,
	)
	// Initialize the result slice
	var tokenList []map[string]interface{}
	if err != nil {
		log.Println("failed to get balance", err)
	} else {
		lamports := uint64(balance)
		lamportsPerSol := uint64(1000000000) // 1 SOL = 1,000,000,000 Lamports
		solBalance := float64(lamports) / float64(lamportsPerSol)
		solBalanceStr := strconv.FormatFloat(solBalance, 'f', -1, 64)
		tokenList = append(tokenList, map[string]interface{}{
			"symbol":   "sol",
			"balance":  solBalanceStr,
			"isNative": true,
		})
	}

	// Get token accounts
	tokenAccounts, err := c.GetTokenAccountsByOwnerWithContextByProgram(context.TODO(), mainAddress, common.TokenProgramID.String())
	if err != nil {
		log.Println("failed to get token account by owner: ", err)
	}

	// Get Redis client
	redisClient := db.GetClient()

	// Iterate over token accounts
	for _, tokenAccount := range tokenAccounts.Value {
		// Get token account balance
		tokenBalance, err := c.GetTokenAccountBalance(context.TODO(), tokenAccount.PublicKey.String())
		if err != nil {
			log.Printf("failed to get token account balance: %v", err)
			continue
		}

		// Get token metadata from Redis
		dataKey := "sol:token:" + tokenAccount.Mint.String()
		tokenResult, err := redisClient.HGetAll(context.TODO(), dataKey).Result()
		if err != nil {
			log.Printf("failed to get fields from Redis: %v", err)
		}

		// Extract token symbol and balance
		symbol := tokenResult["symbol"]
		balanceStr := tokenBalance.UIAmountString

		// Append token balance information
		tokenList = append(tokenList, map[string]interface{}{
			"symbol":   symbol,
			"balance":  balanceStr,
			"isNative": false,
		})
	}
	return tokenList, nil
}
