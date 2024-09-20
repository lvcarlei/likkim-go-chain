package sol

import (
	"encoding/binary"
	"fmt"
	"go-wallet/internal/app/chain/helper"
	"go-wallet/internal/app/chain/oklink"
	"log"
	"math/big"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetTransaction(mainAddress string, before string) (OklinkResp oklink.TransactionRespData) {
	c := GetClient()
	address := solana.MustPublicKeyFromBase58(mainAddress)
	// 配置查询参数

	var limit int = 100
	config := &rpc.GetSignaturesForAddressOpts{
		Limit: &limit,
	}
	if before != "" {
		config.Before = solana.MustSignatureFromBase58(before)
	}
	signatures, err := c.GetSignaturesForAddressWithOpts(ctx,
		address, config)
	if err != nil {
		log.Println("failed to GetSignaturesForAddressWithConfig : ", err)
	}
	var dataList []oklink.Transaction
	//var parsedData map[string]interface{}
	//redisClient := db.GetClient()
	for _, item := range signatures {
		getTransactionOpts := rpc.GetTransactionOpts{}
		transaction, err := c.GetTransaction(ctx, item.Signature, &getTransactionOpts)
		if err != nil {
			log.Println("failed to Get SOL Transaction : ", err)
			continue
		}
		// 检查交易状态
		var sender, receiver, status string
		if transaction.Meta.Err == nil {
			status = "success"
		} else {
			status = "fail"
			log.Printf("SOL交易失败: %v\n", transaction.Meta.Err)
		}
		data, err := transaction.Transaction.GetTransaction()
		if err != nil {
			log.Println("SOL failed to GetTransaction : ", err)
		}
		//检查每条指令，识别代币转账
		for _, instr := range data.Message.Instructions {
			var diff big.Int
			accountKey := data.Message.AccountKeys[instr.ProgramIDIndex]
			if accountKey.String() == solana.TokenProgramID.String() {
				// 检查是否是代币转账指令
				fmt.Printf("type %v, amount %v, decimals %v, sourceAccount %s, destinationAccount %s \n", instr.Data[0], binary.LittleEndian.Uint64(instr.Data[1:9]), instr.Data[9], data.Message.AccountKeys[instr.Accounts[0]], data.Message.AccountKeys[instr.Accounts[1]])
				if len(transaction.Meta.PreTokenBalances) > 0 && len(transaction.Meta.PostTokenBalances) > 0 {
					// 遍历 preTokenBalances 和 postTokenBalances
					for i, preBalance := range transaction.Meta.PreTokenBalances {
						postBalance := transaction.Meta.PostTokenBalances[i]
						// 确保 preBalance 和 postBalance 是同一个账户
						if preBalance.Owner.String() == postBalance.Owner.String() && preBalance.Mint == postBalance.Mint {
							preAmount := preBalance.UiTokenAmount.Amount
							postAmount := postBalance.UiTokenAmount.Amount
							var tmpDiff, tmp1, tmp2 big.Int
							bigintPreAmount, ok1 := tmp1.SetString(preAmount, 10)
							bigintPostAmount, ok2 := tmp2.SetString(postAmount, 10)
							if !ok1 || !ok2 {
								fmt.Println("字符串转换失败")
								continue
							}
							tmpDiff.Sub(bigintPreAmount, bigintPostAmount)
							if tmpDiff.Sign() < 0 {
								// 余额减少，表示是发送地址
								sender = preBalance.Owner.String()
							} else if tmpDiff.Sign() > 0 {
								receiver = preBalance.Owner.String()
							} else {
								continue
							}
							//fmt.Printf("代币信息,发送者%s,接受者%s,变化值 %s \n", sender, receiver, tmpDiff.String())
						}
					}
				}

			} else if accountKey.String() == solana.SystemProgramID.String() {
				// 处理主链（SOL）交易
				preBalances := transaction.Meta.PreBalances
				postBalances := transaction.Meta.PostBalances
				var totalAmountTransferred int64
				for i := range preBalances {
					diffInt64 := int64(preBalances[i]) - int64(postBalances[i])
					if diffInt64 > 0 {
						totalAmountTransferred += int64(diffInt64)
					}
				}
				diff.SetInt64(totalAmountTransferred)
				//fmt.Println("主链转账信息:")
			} else {
				continue
			}
			dataList = append(dataList, oklink.Transaction{
				TxId:            item.Signature.String(),
				Amount:          helper.ConvertToReadableAmount(diff.String(), 9),
				From:            sender,
				To:              receiver,
				TxFee:           helper.ConvertToReadableAmount(transaction.Meta.Fee, 9),
				State:           status,
				Height:          strconv.FormatUint(item.Slot, 10),
				TransactionTime: helper.DateTimeToUnix(item.BlockTime.String()),
			})

			//	// Get token metadata from Redis
			//	//dataKey := "sol:token:" + preBalance.Mint
			//	//tokenResult, err := redisClient.HGetAll(ctx, dataKey).Result()
			//	//if err != nil {
			//	//log.Printf("failed to get fields from Redis: %v", err)
			//	//}
			//
			//	// Extract token symbol and balance
			//	//symbol := tokenResult["symbol"]

		}
	}
	result := oklink.TransactionRespData{
		Page:            "1",
		TransactionList: dataList,
		TotalPage:       "2",
		ChainShortName:  "sol",
	}
	return result
}
