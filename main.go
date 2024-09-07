package main

import (
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/controller"
)

func main() {
	//sol.GetTransaction("EB1CH72E8LubfuUuSPBmoVGkuNVNwSzS8kCgojQM7aQS", "")
	//tron.FetchTokenList("TEubZd9pstkp9kHRR5o7ab1TFhxw5day94")
	//tron.GetTokenBalance("TEN4KrL95t6cSWZwb71gaiXj5ZbadJuT3o")
	irisApp := iris.New()

	route := irisApp.Party("/api")
	{
		conisController := controller.CoinsController{}
		route.Use(iris.Compression)
		route.Get("/wallet/get-balance", conisController.GetTokenBalance)
		route.Get("/wallet/get-transaction", conisController.GetTransaction)
		route.Post("/wallet/broadcastHex", conisController.BroadcastHex)
	}

	irisApp.Listen(":8082")
}
