package controller

import (
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/chain/oklink"
)

type TokenController struct{}

// 获取代币信息
func (c *TokenController) GetTokenInfo(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")
	symbol := ctx.URLParam("symbol")
	protocolType := ctx.URLParam("protocolType")
	if chainShortName == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	resultData := oklink.GetTokenInfo(chainShortName, symbol, protocolType)
	ctx.JSON(map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": resultData,
	})
}

// 更新token信息

func (c *TokenController) UpdateTokenInfo(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")

	oklink.FetchTokenList(chainShortName, "")
}
