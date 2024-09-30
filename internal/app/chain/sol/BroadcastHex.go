package sol

import (
	"context"
	"log"

	"github.com/gagliardetto/solana-go"
)

func BroadcastHex(hex string) (map[string]interface{}, error) {
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

	data := map[string]interface{}{
		"code": code,
		"txid": txSignature,
	}
	return data, err
}
