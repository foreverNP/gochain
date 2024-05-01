package types

import (
	"testing"

	"github.com/foreverNP/gochain/crypto"
	"github.com/foreverNP/gochain/proto"
	"github.com/foreverNP/gochain/util"
	"gotest.tools/assert"
)

// sender balance = 100
// out1: 5 -> to send
// out2: 95 -> to return to sender
func TestHashTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddr := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddr := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddr,
	}
	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddr,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}
	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	assert.Equal(t, VerifyTransaction(tx), true)
}
