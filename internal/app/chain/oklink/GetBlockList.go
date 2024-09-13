package oklink

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Block struct {
	Hash      string `json:"hash"`
	Height    string `json:"height"`
	BlockTime string `json:"blockTime"`
	State     string `json:"state"`
	NetWork   string `json:"netWork"`
}
type BlockListData struct {
	Page           string  `json:"page"`
	Limit          string  `json:"limit"`
	TotalPage      string  `json:"totalPage"`
	ChainFullName  string  `json:"chainFullName"`
	ChainShortName string  `json:"chainShortName"`
	BlockList      []Block `json:"blockList"`
}
type BlockListResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []BlockListData `json:"data"`
}

func GetBlockList(chainShortName string) (oklinkRep BlockListData) {
	baseUrl := OKLINK_ENPOINT + "/api/v5/explorer/block/block-list"
	page := 1
	limit := 20
	url := fmt.Sprintf("%s?page=%d&limit=%d&chainShortName=%s", baseUrl, page, limit, chainShortName)
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置请求超时为5秒
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}
	req.Header.Add("Ok-Access-Key", OKLINK_ACCESS_KEY)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch token list: %v", err)
	}
	defer resp.Body.Close()
	// Parse the response
	var result BlockListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	}
	if result.Code != "0" {
		log.Printf("Oklink GetBlockList API Error: %s", result.Msg)
		return BlockListData{}
	}
	return result.Data[0]
}
