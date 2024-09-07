package oklink

import (
	"context"
	"encoding/json"
	"fmt"
	"go-wallet/db"
	"io/ioutil"
	"log"
	"net/http"
)

// Define the struct to match the API response

type BlockchainResponse struct {
	Code string       `json:"code"`
	Msg  string       `json:"msg"`
	Data []Blockchain `json:"data"`
}

type Blockchain struct {
	ChainFullName  string `json:"chainFullName"`
	ChainShortName string `json:"chainShortName"`
	Symbol         string `json:"symbol"`
}

// Store each blockchain's details in Redis
func storeBlockchainInRedis(blockchain Blockchain) error {
	// Create Redis key, e.g., "blockchain:ETH"
	redisKey := fmt.Sprintf("support:blockchain:%s", blockchain.ChainShortName)
	rdb := db.GetClient()
	// Option 1: Store as a hash (more structured and query-friendly)
	err := rdb.HSet(context.Background(), redisKey, map[string]interface{}{
		"chainFullName": blockchain.ChainFullName,
		"symbol":        blockchain.Symbol,
	}).Err()

	if err != nil {
		return err
	}
	return nil
}

func getSupportedBlockchains() (BlockchainResponse, error) {
	var result BlockchainResponse

	// OKLink API URL
	url := "https://www.oklink.com/api/v5/explorer/blockchain/summary"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Add("OK-ACCESS-KEY", OKLINK_ACCESS_KEY)
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal the JSON response into struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return result, nil
}

func HandleSupportChain() {
	blockchainResponse, err := getSupportedBlockchains()
	if err != nil {
		log.Fatalf("Error fetching blockchains: %v", err)
	}

	// Store each blockchain in Redis with a unique key
	for _, blockchain := range blockchainResponse.Data {
		err := storeBlockchainInRedis(blockchain)
		if err != nil {
			log.Fatalf("Error storing blockchain in Redis: %v", err)
		}
	}

	fmt.Println("Blockchains stored in Redis successfully")

}
