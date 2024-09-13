package controller

import (
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/chain/oklink"
)

type BlockController struct{}

func (c *BlockController) GetBlockList(ctx iris.Context) {
	chainShortName := ctx.URLParam("chain")
	if chainShortName == "" {
		ctx.JSON(Response{Code: "400", Msg: "param error"})
		return
	}
	data := oklink.GetBlockList(chainShortName)
	ctx.JSON(map[string]interface{}{
		"code":      0,
		"msg":       "",
		"data":      data.BlockList,
		"totalPage": data.TotalPage,
		"page":      data.Page,
	})

}
