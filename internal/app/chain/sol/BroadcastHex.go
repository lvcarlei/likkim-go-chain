package sol

import (
	"context"
	"go-wallet/internal/app/chain/oklink"
	"log"

	"github.com/gagliardetto/solana-go"
)

func BroadcastHex(hex string) (oklink.BroadcastResult, error) {
	signedTransaction, err := solana.TransactionFromBase58(hex)
	if err != nil {
		log.Fatalf("Failed to decode signed transaction: %v", err)
	}
	txSignature, err := GetClient().SendTransaction(context.TODO(), signedTransaction)
	code := "0"
	if err != nil {
		log.Fatalf("SOL Failed to send transaction: %v", err)
		code = "1"
	}
	data := oklink.BroadcastResult{}
	data.Txid = txSignature.String()
	data.Code = code
	return data, err
}
