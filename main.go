package main

import (
	"github.com/kataras/iris/v12"
	"go-wallet/internal/app/controller"
)

func main() {
	//sol.GetTransaction("EB1CH72E8LubfuUuSPBmoVGkuNVNwSzS8kCgojQM7aQS", "")
	//tron.FetchTokenList("TEubZd9pstkp9kHRR5o7ab1TFhxw5day94")
	//tron.GetTokenBalance("TEN4KrL95t6cSWZwb71gaiXj5ZbadJuT3o")
	//oklink.HandleSupportChain()
	irisApp := iris.New()

	route := irisApp.Party("/api")
	{
		addressController := controller.AddressController{}
		blockchainController := controller.BlockchainController{}
		tokenController := controller.TokenController{}
		transactionController := controller.TransactionController{}
		blockController := controller.BlockController{}
		route.Use(iris.Compression)
		route.Get("/wallet/balance", addressController.GetTokenBalance)
		route.Get("/wallet/transaction", transactionController.GetTransaction)
		route.Get("/chain/utxo", blockchainController.GetUTXO)
		route.Get("/chain/blockchain-fee", blockchainController.GetBlockchainFee)
		route.Post("/wallet/broadcastHex", blockchainController.BroadcastHex)
		route.Get("/manual/update-support-chain", blockchainController.UpdateSupportChain)
		route.Get("/manual/update-token-info", tokenController.UpdateTokenInfo)
		route.Get("/chain/blocklist", blockController.GetBlockList)
	}

	irisApp.Listen(":8082")
}
