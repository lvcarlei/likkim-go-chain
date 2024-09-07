package controller

import (
	"fmt"
	"go-wallet/internal/app/chain/oklink"
	"go-wallet/internal/app/chain/sol"
	"go-wallet/internal/app/chain/tron"

	"github.com/kataras/iris/v12"
)

type CoinsController struct{}

// 定义响应结构体
type Response struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data map[string]string `json:"data"`
}

// 获取交易记录

func (c *CoinsController) GetTransaction(ctx iris.Context) {
	chainShortName := ctx.URLParam("chainShortName")
	address := ctx.URLParam("address")
	page := ctx.URLParam("page")
	contractAddress := ctx.URLParam("contractAddress")
	protocolType := ctx.URLParam("protocolType")
	if chainShortName == "" || address == "" {
		ctx.JSON(Response{Code: 400, Msg: "param error"})
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
		ctx.JSON(Response{Code: 400, Msg: "param error"})
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
		ctx.JSON(Response{Code: 400, Msg: fmt.Sprintf("%s is not supporting", chainShortName)})
		return
	}
	if len(resultData) == 0 {
		ctx.JSON(Response{Code: 500, Msg: "server error"})
		return
	}
	ctx.JSON(map[string]interface{}{
		"code": 0,
		"msg":  "",
		"data": resultData,
	})

}
