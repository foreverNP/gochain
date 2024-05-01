package types

import (
	"crypto/sha256"

	"github.com/foreverNP/gochain/crypto"
	"github.com/foreverNP/gochain/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		pubKey := crypto.NewPublicKeyFromBytes(input.PublicKey)
		sig := crypto.NewSignatureFromBytes(input.Signature)
		// TODO: may be some problem after verify
		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}
	return true
}
