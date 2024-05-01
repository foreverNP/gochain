package crypto

import (
	"fmt"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	if len(privKey.Bytes()) != privKeyLen {
		t.Errorf("Expected %d, but got %d", privKeyLen, len(privKey.Bytes()))
	}

	pubKey := privKey.Public()
	if len(pubKey.Bytes()) != pubKeyLen {
		t.Errorf("Expected %d, but got %d", pubKeyLen, len(pubKey.Bytes()))
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("message")

	sig := privKey.Sign(msg)

	if !sig.Verify(pubKey, msg) {
		t.Errorf("Expected %t, but got %t", true, false)
	}

	if sig.Verify(pubKey, []byte("foo")) {
		t.Errorf("Expected %t, but got %t", false, true)
	}
}

func TestPrivateKey_Address(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	fmt.Println(address)

	if len(address.Bytes()) != addressLen {
		t.Errorf("Expected %d, but got %d", addressLen, len(address.Bytes()))
	}
}

func TestNewPrivateKetFromString(t *testing.T) {
	var (
		seed       = "a3fe7f380d886c117097c452309a9f2eec869df9fdac3f9754f3273bae208a2a"
		privKey    = NewPrivateKetFromString(seed)
		addressStr = "e4a0a55c14f3fb7870d2b8ae4024844aee24acf1"
	)
	address := privKey.Public().Address()

	if address.String() != addressStr {
		t.Errorf("Expected %s, but got %s", addressStr, address.String())
	}
}
