package controller

import (
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/chain/oklink"
	"go-wallet/internal/app/chain/sol"
)

type TransactionController struct{}

// 获取交易记录

func (c *TransactionController) GetTransaction(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")
	address := ctx.URLParam("address")
	page := ctx.URLParam("page")
	contractAddress := ctx.URLParam("contractAddress")
	protocolType := ctx.URLParam("protocolType")
	before := ctx.URLParam("before")
	if chainShortName == "" || address == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	var data oklink.TransactionRespData
	switch chainShortName {
	case "SOL":
		data = sol.GetTransaction(address, before)
	default:
		data = oklink.GetTransactionList(chainShortName, address, page, contractAddress, protocolType)
	}
	ctx.JSON(map[string]interface{}{
		"code":      0,
		"msg":       "",
		"data":      data.TransactionList,
		"totalPage": data.TotalPage,
		"page":      data.Page,
	})

}
