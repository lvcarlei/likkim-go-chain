package sol

import (
	"context"
	"sync"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

var (
	once      sync.Once
	ctx       = context.Background()
	rpcClient *client.Client
)

func initialize() {
	once.Do(func() {
		rpcClient = client.NewClient(rpc.MainnetRPCEndpoint)
	})
}

func GetClient() *client.Client {
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
