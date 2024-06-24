package opplasmairys

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

func newirysclient() (*IrysClient, error) {
	// test private key, don't use it in production
	var privatekeyHex = "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	return NewIrysClient("ethereum", "https://rpc.sepolia.org", privatekeyHex, true)
}

func TestClientUpload(t *testing.T) {
	c, err := newirysclient()
	if err != nil {
		t.Error(err)
		return
	}
	testData, _ := hexutil.Decode("0x500024cf833b1ee1b2a2d8abc1c9288000f2000000007855dada38c9719bd34aa627eb7b333ef67d7fe2d468119b73af6b86e5ed994aeed9266ff5a41f88930f6feee1f2513c5ff752e7a706c3bc3f9c9807591cd0042310a19413006061116264700518011400000000ffff0001")
	id, _ := hexutil.Decode("0x286a251d38aff7755e7d3ff97d1e630844eb3eb5e5fdf3147aa46cb3a9fbb940")
	err = c.Upload(context.Background(), id, testData)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("owner addr: ", c.ethereumAddr, c.arweaveAddr)
}

func TestClientDownload(t *testing.T) {
	c, err := newirysclient()
	if err != nil {
		t.Error(err)
		return
	}
	id, _ := hexutil.Decode("0x286a251d38aff7755e7d3ff97d1e630844eb3eb5e5fdf3147aa46cb3a9fbb940")
	_, err = c.Download(context.Background(), id)
	if err != nil {
		t.Error(err)
		return
	}
}
