package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/chain/oklink"
	"go-wallet/internal/app/chain/sol"
	"go-wallet/internal/app/chain/tron"
)

type CoinsController struct{}

// 定义响应结构体

type Response struct {
	Code      string                 `json:"code"`
	Msg       string                 `json:"msg"`
	Data      map[string]interface{} `json:"data"`
	TotalPage string                 `json:"totalPage"`
}

// 获取交易记录

func (c *CoinsController) GetTransaction(ctx iris.Context) {
	chainShortName := ctx.URLParam("chainShortName")
	address := ctx.URLParam("address")
	page := ctx.URLParam("page")
	contractAddress := ctx.URLParam("contractAddress")
	protocolType := ctx.URLParam("protocolType")
	if chainShortName == "" || address == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	data := oklink.GetTransactionList(chainShortName, address, page, contractAddress, protocolType)
	ctx.JSON(map[string]interface{}{
		"code":      0,
		"msg":       "",
		"data":      data.TransactionList,
		"totalPage": data.TotalPage,
		"page":      data.Page,
	})

}

// 获取钱包余额

func (c *CoinsController) GetTokenBalance(ctx iris.Context) {
	chainShortName := ctx.URLParam("chainShortName")
	address := ctx.URLParam("address")

	if chainShortName == "" || address == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	var resultData []map[string]interface{}
	if chainShortName == "SOL" {
		info, _ := sol.GetTokenBalance(address)
		resultData = info
	} else if chainShortName == "TRON" {
		info, _ := tron.GetTokenBalance(address)
		resultData = info
	} else {
		ctx.JSON(Response{Code: "400", Msg: fmt.Sprintf("%s is not supporting", chainShortName)})
		return
	}
	if len(resultData) == 0 {
		ctx.JSON(Response{Code: "500", Msg: "server error"})
		return
	}
	ctx.JSON(map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": resultData,
	})

}

// 广播交易

func (c *CoinsController) BroadcastHex(ctx iris.Context) {
	chainShortName := ctx.PostValue("chainShortName")
	hex := ctx.PostValue("hex")
	if chainShortName == "" || hex == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	var resultData map[string]interface{}
	var err error
	if chainShortName == "TRON" {
		resultData, err = tron.BroadcastHex(hex)
	} else {
		resultData, err = oklink.BroadcastHex(chainShortName, hex)
	}
	if err != nil {
		ctx.JSON(Response{Code: "500", Msg: "server error"})
		return
	}
	ctx.JSON(Response{
		Code: resultData["code"].(string),
		Msg:  "",
		Data: map[string]interface{}{
			"txid": resultData["txid"].(string),
		},
	})
}

// 获取UTXO

func (c *CoinsController) GetUTXO(ctx iris.Context) {
	chainShortName := ctx.URLParam("chainShortName")
	address := ctx.URLParam("address")
	if chainShortName == "" || address == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	//var resultData map[string]interface{}
	//var err error
	resultData, err := oklink.GetUTXO(chainShortName, address, "")
	if err != nil {
		ctx.JSON(Response{Code: "500", Msg: "server error"})
		return
	}
	ctx.JSON(Response{
		Code:      resultData["code"].(string),
		TotalPage: "",
		Data: map[string]interface{}{
			"utxoList": resultData["data"],
		},
	})
}

// 获取手续费

func (c *CoinsController) GetBlockchainFee(ctx iris.Context) {
	chainShortName := ctx.URLParam("chainShortName")
	if chainShortName == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	resultData, err := oklink.GetBlockchainFee(chainShortName)
	if err != nil {
		ctx.JSON(Response{Code: "500", Msg: "server error"})
		return
	}
	ctx.JSON(map[string]interface{}{
		"code": "0",
		"msg":  resultData["msg"],
		"data": resultData["data"],
	})
}

// 更新支持的公链

func (c *CoinsController) UpdateSupportChain(ctx iris.Context) {
	oklink.HandleSupportChain()
}
