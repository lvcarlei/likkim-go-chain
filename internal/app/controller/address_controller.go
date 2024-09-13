package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/chain/sol"
	"go-wallet/internal/app/chain/tron"
)

type AddressController struct{}

// 获取钱包余额

func (c *AddressController) GetTokenBalance(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")
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
