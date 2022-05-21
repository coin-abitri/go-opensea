package opensea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

type Order struct {
	ID    int64 `json:"id"`
	Asset Asset `json:"asset"`
	// AssetBundle          interface{}          `json:"asset_bundle"`
	CreatedDate *TimeNano `json:"created_date"`
	ClosingDate *TimeNano `json:"closing_date"`
	// ClosingExtendable bool      `json:"closing_extendable"`
	ExpirationTime int64 `json:"expiration_time"`
	ListingTime    int64 `json:"listing_time"`
	// OrderHash            string               `json:"order_hash"`
	// Metadata Metadata `json:"metadata"`
	Exchange     Address `json:"exchange"`
	Maker        Account `json:"maker"`
	Taker        Account `json:"taker"`
	CurrentPrice Number  `json:"current_price"`
	// CurrentBounty        string               `json:"current_bounty"`
	// BountyMultiple       string               `json:"bounty_multiple"`
	MakerRelayerFee    Number    `json:"maker_relayer_fee"`
	TakerRelayerFee    Number    `json:"taker_relayer_fee"`
	MakerProtocolFee   Number    `json:"maker_protocol_fee"`
	TakerProtocolFee   Number    `json:"taker_protocol_fee"`
	MakerReferrerFee   Number    `json:"maker_referrer_fee"`
	FeeRecipient       Account   `json:"fee_recipient"`
	FeeMethod          FeeMethod `json:"fee_method"`
	Side               Side      `json:"side"` // 0 for buy orders and 1 for sell orders.
	SaleKind           SaleKind  `json:"sale_kind"`
	Target             Address   `json:"target"`
	HowToCall          HowToCall `json:"how_to_call"`
	Calldata           Bytes     `json:"calldata"`
	ReplacementPattern Bytes     `json:"replacement_pattern"`
	StaticTarget       Address   `json:"static_target"`
	StaticExtradata    Bytes     `json:"static_extradata"`
	PaymentToken       Address   `json:"payment_token"`
	// PaymentTokenContract PaymentTokenContract `json:"payment_token_contract"`
	BasePrice       Number `json:"base_price"`
	Extra           Number `json:"extra"`
	Quantity        string `json:"quantity"`
	Salt            Number `json:"salt"`
	V               *uint8 `json:"v"`
	R               *Bytes `json:"r"`
	S               *Bytes `json:"s"`
	ApprovedOnChain bool   `json:"approved_on_chain"`
	Cancelled       bool   `json:"cancelled"`
	Finalized       bool   `json:"finalized"`
	MarkedInvalid   bool   `json:"marked_invalid"`
	// PrefixedHash         string               `json:"prefixed_hash"`
}

func (o Order) IsPrivate() bool {
	if o.Taker.Address != NullAddress {
		return true
	}
	return false
}

type Side uint8

const (
	Buy Side = iota
	Sell
)

type SaleKind uint8

const (
	FixedOrMinBit SaleKind = iota // 0 for fixed-price sales or min-bid auctions
	DutchAuctions                 // 1 for declining-price Dutch Auctions
)

type HowToCall uint8

const (
	Call HowToCall = iota
	DelegateCall
)

type FeeMethod uint8

const (
	ProtocolFee FeeMethod = iota
	SplitFee
)

const (
	OrderByCreateDate = "created_date"
	OrderByPrice      = "eth_price"

	OrderDirectionAscend  = "asc"
	OrderDirectionDescend = "desc"
)

type GetOrderOpts struct {
	AssetContractAddress string   `json:"asset_contract_address,omitempty"`
	ListedAfter          string   `json:"listed_after,omitempty"`
	ListedBefore         string   `json:"listed_before,omitempty"`
	OrderBy              string   `json:"order_by,omitempty"`
	OrderDirection       string   `json:"order_direction,omitempty"`
	TokenId              string   `json:"token_id,omitempty"`
	TokenIds             []string `json:"token_ids,omitempty"`
	Limit                int32    `json:"limit,omitempty"`
	Offset               int32    `json:"offset,omitempty"`
}

func (o Opensea) GetOrders(opt GetOrderOpts) ([]*Order, error) {
	return o.GetOrdersWithContext(context.TODO(), opt)
}

func (o Opensea) GetOrdersWithContext(ctx context.Context, opt GetOrderOpts) (orders []*Order, err error) {
	if opt.Limit == 0 {
		opt.Limit = 100
	}
	offset := int32(0)

	q, err := query.Values(opt)
	if err != nil {
		return nil, err
	}
	orders = []*Order{}

	for true {
		q.Set("offset", fmt.Sprintf("%d", offset))
		path := "/wyvern/v1/orders?" + q.Encode()
		b, err := o.getPath(ctx, path)
		if err != nil {
			return nil, err
		}

		out := &struct {
			Count  int64    `json:"count"`
			Orders []*Order `json:"orders"`
		}{}

		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}
		orders = append(orders, out.Orders...)

		if int32(len(out.Orders)) < opt.Limit {
			break
		}
		offset += opt.Limit
	}

	return
}

type GetOfferOpts struct {
	AssetContractAddress string `json:"asset_contract_address,omitempty"`
	TokenId              string `json:"token_id,omitempty"`
	// max 50
	Limit int `json:"limit,omitempty"`
}

type GetListingOpts struct {
	AssetContractAddress string `json:"asset_contract_address,omitempty"`
	TokenId              string `json:"token_id,omitempty"`
	// max 50
	Limit int `json:"limit,omitempty"`
}

type GetOfferResponse struct {
	Offers []Offers `json:"offers"`
}

type FeeRecipient struct {
	User          int    `json:"user"`
	ProfileImgURL string `json:"profile_img_url"`
	Address       string `json:"address"`
	Config        string `json:"config"`
}

type PaymentTokenContract struct {
	ID       int         `json:"id"`
	Symbol   string      `json:"symbol"`
	Address  string      `json:"address"`
	ImageURL string      `json:"image_url"`
	Name     interface{} `json:"name"`
	Decimals int         `json:"decimals"`
	EthPrice string      `json:"eth_price"`
	UsdPrice interface{} `json:"usd_price"`
}

type Offers struct {
	CreatedDate          string               `json:"created_date"`
	ClosingDate          interface{}          `json:"closing_date"`
	ClosingExtendable    bool                 `json:"closing_extendable"`
	ExpirationTime       int                  `json:"expiration_time"`
	ListingTime          int                  `json:"listing_time"`
	OrderHash            string               `json:"order_hash"`
	Metadata             Metadata             `json:"metadata"`
	Exchange             string               `json:"exchange"`
	Maker                Maker                `json:"maker"`
	Taker                Taker                `json:"taker"`
	CurrentPrice         string               `json:"current_price"`
	CurrentBounty        string               `json:"current_bounty"`
	BountyMultiple       string               `json:"bounty_multiple"`
	MakerRelayerFee      string               `json:"maker_relayer_fee"`
	TakerRelayerFee      string               `json:"taker_relayer_fee"`
	MakerProtocolFee     string               `json:"maker_protocol_fee"`
	TakerProtocolFee     string               `json:"taker_protocol_fee"`
	MakerReferrerFee     string               `json:"maker_referrer_fee"`
	FeeRecipient         FeeRecipient         `json:"fee_recipient"`
	FeeMethod            int                  `json:"fee_method"`
	Side                 int                  `json:"side"`
	SaleKind             int                  `json:"sale_kind"`
	Target               string               `json:"target"`
	HowToCall            int                  `json:"how_to_call"`
	Calldata             string               `json:"calldata"`
	ReplacementPattern   string               `json:"replacement_pattern"`
	StaticTarget         string               `json:"static_target"`
	StaticExtradata      string               `json:"static_extradata"`
	PaymentToken         string               `json:"payment_token"`
	PaymentTokenContract PaymentTokenContract `json:"payment_token_contract"`
	BasePrice            string               `json:"base_price"`
	Extra                string               `json:"extra"`
	Quantity             string               `json:"quantity"`
	Salt                 string               `json:"salt"`
	V                    int                  `json:"v"`
	R                    string               `json:"r"`
	S                    string               `json:"s"`
	ApprovedOnChain      bool                 `json:"approved_on_chain"`
	Cancelled            bool                 `json:"cancelled"`
	Finalized            bool                 `json:"finalized"`
	MarkedInvalid        bool                 `json:"marked_invalid"`
	PrefixedHash         string               `json:"prefixed_hash"`
}

type GetListingResponse struct {
	Listings []Listings `json:"listings"`
}

type Asset struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type Metadata struct {
	Asset  Asset  `json:"asset"`
	Schema string `json:"schema"`
}

type Maker struct {
	User          int    `json:"user"`
	ProfileImgURL string `json:"profile_img_url"`
	Address       string `json:"address"`
	Config        string `json:"config"`
}
type Taker struct {
	User          int    `json:"user"`
	ProfileImgURL string `json:"profile_img_url"`
	Address       string `json:"address"`
	Config        string `json:"config"`
}

type Listings struct {
	CreatedDate          string               `json:"created_date"`
	ClosingDate          string               `json:"closing_date"`
	ClosingExtendable    bool                 `json:"closing_extendable"`
	ExpirationTime       int                  `json:"expiration_time"`
	ListingTime          int                  `json:"listing_time"`
	OrderHash            string               `json:"order_hash"`
	Metadata             Metadata             `json:"metadata"`
	Exchange             string               `json:"exchange"`
	Maker                Maker                `json:"maker"`
	Taker                Taker                `json:"taker"`
	CurrentPrice         string               `json:"current_price"`
	CurrentBounty        string               `json:"current_bounty"`
	BountyMultiple       string               `json:"bounty_multiple"`
	MakerRelayerFee      string               `json:"maker_relayer_fee"`
	TakerRelayerFee      string               `json:"taker_relayer_fee"`
	MakerProtocolFee     string               `json:"maker_protocol_fee"`
	TakerProtocolFee     string               `json:"taker_protocol_fee"`
	MakerReferrerFee     string               `json:"maker_referrer_fee"`
	FeeRecipient         FeeRecipient         `json:"fee_recipient"`
	FeeMethod            int                  `json:"fee_method"`
	Side                 int                  `json:"side"`
	SaleKind             int                  `json:"sale_kind"`
	Target               string               `json:"target"`
	HowToCall            int                  `json:"how_to_call"`
	Calldata             string               `json:"calldata"`
	ReplacementPattern   string               `json:"replacement_pattern"`
	StaticTarget         string               `json:"static_target"`
	StaticExtradata      string               `json:"static_extradata"`
	PaymentToken         string               `json:"payment_token"`
	PaymentTokenContract PaymentTokenContract `json:"payment_token_contract"`
	BasePrice            string               `json:"base_price"`
	Extra                string               `json:"extra"`
	Quantity             string               `json:"quantity"`
	Salt                 string               `json:"salt"`
	V                    interface{}          `json:"v"`
	R                    interface{}          `json:"r"`
	S                    interface{}          `json:"s"`
	ApprovedOnChain      bool                 `json:"approved_on_chain"`
	Cancelled            bool                 `json:"cancelled"`
	Finalized            bool                 `json:"finalized"`
	MarkedInvalid        bool                 `json:"marked_invalid"`
	PrefixedHash         string               `json:"prefixed_hash"`
}

func (o Opensea) GetOffersWithContext(ctx context.Context, opt GetOfferOpts) (*GetOfferResponse, error) {
	if opt.Limit == 0 {
		opt.Limit = 50
	}

	q := url.Values{}
	q.Set("limit", fmt.Sprintf("%d", opt.Limit))
	var offers = &GetOfferResponse{}

	// https://api.opensea.io/api/v1/asset/{asset_contract_address}/{token_id}/offers
	path := fmt.Sprintf("/api/v1/asset/%s/%s/offers", opt.AssetContractAddress, opt.TokenId)
	b, err := o.getPath(ctx, path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, offers)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (o Opensea) GetListingWithContext(ctx context.Context, opt GetListingOpts) (*GetListingResponse, error) {
	if opt.Limit == 0 {
		opt.Limit = 50
	}

	q := url.Values{}
	q.Set("limit", fmt.Sprintf("%d", opt.Limit))
	var listings = &GetListingResponse{}

	// https://api.opensea.io/api/v1/asset/{asset_contract_address}/{token_id}/offers
	path := fmt.Sprintf("/api/v1/asset/%s/%s/listings", opt.AssetContractAddress, opt.TokenId)
	b, err := o.getPath(ctx, path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, listings)
	if err != nil {
		return nil, err
	}
	return listings, nil
}
