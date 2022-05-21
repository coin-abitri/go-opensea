package opensea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/cheekybits/is"
)

func TestOpensea_GetListingWithContext(t *testing.T) {
	o, _ := NewOpensea(os.Getenv("OPENSEA_API_KEY"))

	got, err := o.GetListingWithContext(context.Background(), GetListingOpts{
		AssetContractAddress: "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
		TokenId:              "8520",
		Limit:                20,
	})
	is.New(t).NoErr(err)

	data, _ := json.Marshal(got)
	t.Log(string(data))
}
