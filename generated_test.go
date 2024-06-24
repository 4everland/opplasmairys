package opplasmairys

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"net/http"
	"testing"
)

func TestGetTx(t *testing.T) {
	client := graphql.NewClient(GraphqlMainNet, http.DefaultClient)

	tx, err := QueryTx(context.Background(), client, []string{"0x55b65D6971A0b3737A01be3095cDf1e01645DF68"}, "0x6ea56d1d9530d2d1678276ca1ebf308864605e4c17fb428782321de0ff42331d")
	if err != nil {
		t.Errorf("GetTx failed: %v", err)
	}
	if len(tx.Transactions.Edges) >= 1 {
		fmt.Println(tx.Transactions.Edges[0].Node.Id)
	}
}
