package opplasmairys

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"github.com/Ja7ad/irys/types"
	ghl "github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"log/slog"
	"net/http"
)

const (
	Matic     = "matic"
	Bnb       = "bnb"
	Avalanche = "avalanche"
	Arbitrum  = "arbitrum"
	Fantom    = "fantom"
	Ethereum  = "ethereum"
)

const (
	GraphqlMainNet  = "https://arweave.mainnet.irys.xyz/graphql"
	GraphqlMainNet2 = "https://arweave.net/graphql"
)

type Config struct {
	NetworkName string
	PrivateKey  string
	NetWorkRpc  string

	FreeUpload bool
}

type IrysClient struct {
	c       irys.Irys
	key     string
	graphql ghl.Client

	enableFreeUpload bool

	ethereumAddr string
	arweaveAddr  string
}

// NewIrysClient
// support tokens ethereum,matic,bnb,avalanche,arbitrum,fantom
func NewIrysClient(networkName, rpc, privateKey string, enableFreeUpload bool) (*IrysClient, error) {
	prKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	pubkBytes := crypto.FromECDSAPub(&prKey.PublicKey)
	key := base64.RawURLEncoding.EncodeToString(pubkBytes)
	ethereumAddress := crypto.PubkeyToAddress(prKey.PublicKey).String()
	addr := sha256.Sum256(pubkBytes)
	arweaveAddr := base64.RawURLEncoding.EncodeToString(addr[:])
	currencyFunc := newProvider(networkName)
	if currencyFunc == nil {
		return nil, fmt.Errorf("unsupport network: %s", networkName)
	}
	matic, err := currencyFunc(privateKey, rpc)
	if err != nil {
		return nil, err
	}

	c, err := irys.New(irys.DefaultNode1, matic, false)
	if err != nil {
		return nil, err
	}
	graphql := ghl.NewClient(GraphqlMainNet, http.DefaultClient)

	return &IrysClient{
		c:       c,
		key:     key,
		graphql: graphql,

		enableFreeUpload: enableFreeUpload,

		ethereumAddr: ethereumAddress,
		arweaveAddr:  arweaveAddr,
	}, nil
}

func (c *IrysClient) Upload(ctx context.Context, id []byte, data []byte) error {
	var (
		tx  types.Transaction
		err error
	)
	if c.enableFreeUpload && len(data) > 1024*1024*100 {
		tx, err = c.c.BasicUpload(ctx, data, types.Tag{
			Name:  "OP-PLASMA-KEY",
			Value: hexutil.Encode(id),
		})
	} else {
		tx, err = c.c.Upload(ctx, data, types.Tag{
			Name:  "OP-PLASMA-KEY",
			Value: hexutil.Encode(id),
		})
	}
	if err != nil {
		return err
	}
	slog.Info("irys uploaded, ", "tx", tx.ID)
	return nil
}

func (c *IrysClient) Download(ctx context.Context, id []byte) (io.ReadCloser, error) {
	//query txid by op plasma key
	// filter owners:
	//use ethereum address for https://arweave.mainnet.irys.xyz/graphql
	//use arweave address for https://arweave.net/graphql
	tx, err := QueryTx(ctx, c.graphql, []string{c.ethereumAddr}, hexutil.Encode(id))
	if err != nil {
		return nil, err
	}
	edges := tx.Transactions.Edges
	if len(edges) == 0 {
		return nil, fmt.Errorf("not found")
	}
	txId := edges[0].Node.Id
	f, err := c.c.Download(ctx, txId)
	if err != nil {
		return nil, err
	}
	return f.Data, nil
}

func newProvider(name string) func(privateKey, rpc string) (currency.Currency, error) {
	switch name {
	case Bnb:
		return currency.NewBNB
	case Arbitrum:
		return currency.NewArbitrum
	case Avalanche:
		return currency.NewAvalanche
	case Fantom:
		return currency.NewFantom
	case Ethereum:
		return currency.NewEthereum
	case Matic:
		return currency.NewMatic
	}
	return nil
}
