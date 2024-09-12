package sol

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

// TokenInfo 结构体表示每个 Token 的信息
type TokenInfo struct {
	Address  string  `json:"address"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	ChainId  float64 `json:"chainId"`
	Decimals float64 `json:"decimals"`
	LogoURI  string  `json:"logoURI"`
}

var baseSolKey = "tokenInfo:SOL:"

func GetchTokenList() {
	// 连接到 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "likkim2024",     // Redis 密码（如果没有设置密码，则为空）
		DB:       0,                // 使用的 Redis 数据库编号
	})

	// 读取本地的 solana.tokenlist.json 文件
	tokenList, err := loadTokenListFromFile("solana.tokenlist.json")
	if err != nil {
		log.Fatalf("failed to load token list: %v", err)
	}

	// 将 Token List 以 HashMap 形式存储到 Redis 中
	err = storeTokenListInRedisAsHash(rdb, tokenList)
	if err != nil {
		log.Fatalf("failed to store token list in redis as hash: %v", err)
	}

	log.Println("Token list stored successfully in Redis as HashMap")

	// 示例：根据地址从 Redis 中查找 Token 信息
	address := "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB" // 替换为你要查找的 address
	tokenInfo, err := getTokenInfoFromRedis(rdb, address)
	if err != nil {
		log.Fatalf("failed to get token info from redis: %v", err)
	}

	fmt.Printf("Token info for address %s: %v\n", address, tokenInfo)
}

// loadTokenListFromFile 从本地文件读取 Token List 并解析为 TokenInfo 列表
func loadTokenListFromFile(filename string) ([]TokenInfo, error) {
	data, err := os.ReadFile(filename) // 使用 os.ReadFile 代替 ioutil.ReadFile
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	tokens := result["tokens"].([]interface{})
	tokenList := make([]TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		tokenMap := token.(map[string]interface{})

		tokenList = append(tokenList, TokenInfo{
			Address:  tokenMap["address"].(string),
			Name:     tokenMap["name"].(string),
			Symbol:   tokenMap["symbol"].(string),
			ChainId:  tokenMap["chainId"].(float64),
			Decimals: tokenMap["decimals"].(float64),
			LogoURI:  tokenMap["logoURI"].(string),
		})
	}

	return tokenList, nil
}

// storeTokenListInRedisAsHash 将 Token List 存储为 Redis 的 HashMap
func storeTokenListInRedisAsHash(rdb *redis.Client, tokenList []TokenInfo) error {
	for _, token := range tokenList {
		// 将每个 Token 的信息存储为 HashMap
		tokenKey := baseSolKey + token.Address
		err := rdb.HSet(ctx, tokenKey, map[string]interface{}{
			"name":     token.Name,
			"symbol":   token.Symbol,
			"chainId":  token.ChainId,
			"decimals": token.Decimals,
			"logoURI":  token.LogoURI,
		}).Err()
		if err != nil {
			return err
		}
		//key2 := baseSolKey + token.Symbol
		//_, err = rdb.Do(context.Background(), "COPY", tokenKey, key2).Result()
		//if err != nil {
		//	log.Fatalf("Error copying key: %v", err)
		//}
	}
	return nil
}

// getTokenInfoFromRedis 从 Redis 中获取指定地址的 Token 信息
func getTokenInfoFromRedis(rdb *redis.Client, address string) (map[string]string, error) {
	tokenKey := baseSolKey + address
	result, err := rdb.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
