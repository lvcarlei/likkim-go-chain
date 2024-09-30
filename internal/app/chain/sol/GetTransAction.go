package sol

import (
	"go-wallet/internal/app/chain/helper"
	"go-wallet/internal/app/chain/oklink"
	"log"
	"strconv"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
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
	for _, item := range signatures {
		//getTransactionOpts := rpc.GetTransactionOpts{}
		out, err := c.GetTransaction(ctx, item.Signature, &rpc.GetTransactionOpts{
			Encoding: solana.EncodingBase64,
		},
		)
		if err != nil {
			log.Println("failed to Get SOL Transaction : ", err)
			continue
		}
		tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
		if err != nil {
			panic(err)
		}
		status := "fail"
		if out.Meta.Err == nil {
			status = "success"
		}
		txFee := helper.ConvertToReadableAmount(out.Meta.Fee, 9)
		for _, inst := range tx.Message.Instructions {
			progKey, err := tx.ResolveProgramIDIndex(inst.ProgramIDIndex)
			if err == nil {
				accounts, err := inst.ResolveInstructionAccounts(&tx.Message)
				if err != nil {
					panic(err)
				}
				switch progKey {
				case solana.TokenProgramID:
					decodedInstruction, err := token.DecodeInstruction(accounts, inst.Data)
					if err == nil {
						//spew.Dump(decodedInstruction)
						if transferChecked, ok := decodedInstruction.Impl.(*token.TransferChecked); ok {
							// 获取转账金额
							amount := *transferChecked.Amount
							decimals := *transferChecked.Decimals
							readableAmount := helper.ConvertToReadableAmount(amount, decimals)
							dataList = append(dataList, oklink.Transaction{
								TxId:                 item.Signature.String(),
								Amount:               readableAmount,
								From:                 transferChecked.GetOwnerAccount().PublicKey.String(),
								To:                   transferChecked.GetDestinationAccount().PublicKey.String(),
								TxFee:                txFee,
								State:                status,
								Height:               strconv.FormatUint(item.Slot, 10),
								TransactionTime:      helper.DateTimeToUnix(item.BlockTime.String()),
								TokenContractAddress: transferChecked.GetMintAccount().PublicKey.String(),
								TransactionSymbol:    token.InstructionIDToName(decodedInstruction.TypeID.Uint8()),
							})
						} else {
							log.Println("The instruction is not a token transfer")
						}

					} else {
						//fmt.Println("token decodedInstruction err ======", err)
						continue
					}
				case solana.SystemProgramID:
					systemInstruction, err := system.DecodeInstruction(accounts, inst.Data)
					if err != nil {
						panic(err)
					}
					if transfer, ok := systemInstruction.Impl.(*system.Transfer); ok {
						amount := *transfer.Lamports
						readableAmount := helper.ConvertToReadableAmount(amount, 9)
						sender := accounts[0].PublicKey.String()   // 发送者
						receiver := accounts[1].PublicKey.String() // 接收者
						dataList = append(dataList, oklink.Transaction{
							TxId:                 item.Signature.String(),
							Amount:               readableAmount,
							From:                 sender,
							To:                   receiver,
							TxFee:                txFee,
							State:                status,
							Height:               strconv.FormatUint(item.Slot, 10),
							TransactionTime:      helper.DateTimeToUnix(item.BlockTime.String()),
							TokenContractAddress: "",
							TransactionSymbol:    system.InstructionIDToName(systemInstruction.TypeID.Uint32()),
						})
						//fmt.Printf("system Transfer Amount: %d tokens\n", amount)
					}

				default:
					continue
				}
			} else {
				panic(err)
			}
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
