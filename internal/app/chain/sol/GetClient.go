package sol

import (
	"context"
	"github.com/gagliardetto/solana-go/rpc"
	"sync"
)

var (
	once      sync.Once
	ctx       = context.Background()
	rpcClient *rpc.Client
)

func initialize() {
	once.Do(func() {
		//endpoint := rpc.MainNetBeta_RPC
		endpoint := "https://solana-mainnet.g.alchemy.com/v2/LJ6in8xgKPw5QyF69bftR3_kH7ChpoF4"
		rpcClient = rpc.New(endpoint)
	})
}

func GetClient() *rpc.Client {
	if rpcClient == nil {
		initialize()
	}
	return rpcClient
}

func DefaultCtx() context.Context {
	return context.Background()
}

func LamportsToSOL(lamports uint64) float64 {
	const lamportsPerSOL = 1000000000 // 1 SOL = 10^9 lamports
	return float64(lamports) / float64(lamportsPerSOL)
}
