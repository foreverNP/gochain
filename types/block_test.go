package types

import (
	"testing"

	"github.com/foreverNP/gochain/crypto"
	"github.com/foreverNP/gochain/util"
	"gotest.tools/assert"
)

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)

	assert.Equal(t, len(hash), 32)
}

func TestSignBlock(t *testing.T) {
	var (
		block = util.RandomBlock()
		privK = crypto.GeneratePrivateKey()
		pubK  = privK.Public()
	)

	sig := SignBlock(privK, block)

	assert.Equal(t, 64, len(sig.Bytes()))
	assert.Equal(t, sig.Verify(pubK, HashBlock(block)), true)

	// Test that the signature is different for different blocks
	block2 := util.RandomBlock()

	assert.Equal(t, sig.Verify(pubK, HashBlock(block2)), false)
}
