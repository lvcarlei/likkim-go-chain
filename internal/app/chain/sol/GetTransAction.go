package sol

import (
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"log"
	"time"
)

func GetTransaction(address string, before string) ([]map[string]interface{}, error) {
	c := GetClient()
	// 配置查询参数
	config := client.GetSignaturesForAddressConfig{
		Limit:  5, // 限制返回的记录数为10
		Before: before,
	}
	signatures, err := c.GetSignaturesForAddressWithConfig(ctx,
		address, config)
	if err != nil {
		log.Println("failed to GetSignaturesForAddressWithConfig : ", err)
	}
	var dataList []map[string]interface{}
	//redisClient := db.GetClient()
	for _, item := range signatures {
		//log.Println(item.BlockTime, item.Signature)
		transaction, err := c.GetTransaction(ctx, item.Signature)
		if err != nil {
			log.Println("failed to GetTransaction : ", err)
			continue
		}
		preBalances := transaction.Meta.PreBalances
		postBalances := transaction.Meta.PostBalances
		var totalAmountTransferred uint64
		for i := range preBalances {
			// Calculate the difference for each account
			diff := int64(preBalances[i]) - int64(postBalances[i])
			if diff > 0 {
				totalAmountTransferred += uint64(diff)
			}
		}

		fmt.Printf("sig %s Total Amount Transferred: %s lamports\n", item.Signature, LamportsToSOL(totalAmountTransferred))

		// Get token metadata from Redis
		//dataKey := "sol:token:" + preBalance.Mint
		//tokenResult, err := redisClient.HGetAll(ctx, dataKey).Result()
		//if err != nil {
		//log.Printf("failed to get fields from Redis: %v", err)
		//}

		// Extract token symbol and balance
		//symbol := tokenResult["symbol"]
		//dataList = append(dataList, map[string]interface{}{
		//	"blockTime": item.BlockTime,
		//	"signature": item.Signature,
		//})
		time.Sleep(1 * time.Second)
	}
	return dataList, nil
}
