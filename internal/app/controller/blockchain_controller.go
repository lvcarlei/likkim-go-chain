package controller

import (
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/chain/oklink"
	"go-wallet/internal/app/chain/tron"
)

type BlockchainController struct{}

// 定义响应结构体

type Response struct {
	Code      string                 `json:"code"`
	Msg       string                 `json:"msg"`
	Data      map[string]interface{} `json:"data"`
	TotalPage string                 `json:"totalPage"`
}

// 广播交易

func (c *BlockchainController) BroadcastHex(ctx iris.Context) {
	chainShortName := ctx.PostValue("chain")
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

func (c *BlockchainController) GetUTXO(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")
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

func (c *BlockchainController) GetBlockchainFee(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")
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

func (c *BlockchainController) UpdateSupportChain(ctx iris.Context) {
	oklink.HandleSupportChain()
}
