package cbt

import (
	"cbt/token/nep17"

	"github.com/nspcc-dev/neo-go/pkg/interop"
	"github.com/nspcc-dev/neo-go/pkg/interop/lib/address"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const (
	decimals   = 8
	multiplier = 100000000
	total      = 600000000
)

var (
	owner  = address.ToHash160("NYE7BDVfMNiGuNR1KLdRqtV3ryuczGhwJg")
	token  nep17.Token
	ctx    storage.Context
	tratio = 1 * multiplier
	nratio = 1 * multiplier
)

// init initializes Token Interface and storage context for the Smart
// Contract to operate with
func init() {
	token = nep17.Token{
		Name:           "cbt",
		Symbol:         "$Card",
		Decimals:       decimals,
		Owner:          owner,
		Multiplier:     multiplier,
		TotalSupply:    total * multiplier,
		CirculationKey: "TokenCirculation",
	}
	ctx = storage.GetContext()
}

// Symbol returns the token symbol
func Symbol() string {
	return token.Symbol
}

// Decimals returns the token decimals
func Decimals() int {
	return token.Decimals
}

// TotalSupply returns the token total supply value
func TotalSupply() int {
	return token.GetSupply(storage.GetReadOnlyContext())
}

// BalanceOf returns the amount of token on the specified address
func BalanceOf(holder interop.Hash160) int {
	return token.BalanceOf(storage.GetReadOnlyContext(), holder)
}

// Transfer token from one user to another
func Transfer(from interop.Hash160, to interop.Hash160, amount int, data any) bool {
	return token.Transfer(ctx, from, to, amount, data)
}

// Mint initial supply of tokens
func Mint(to interop.Hash160, amount int) bool {
	return token.Mint(ctx, to, amount)
}

// set mint is once
func SetOnce(v bool) bool {
	return token.SetOnce(ctx, v)
}

// exchange gas to token
func Exchange(to interop.Hash160, amount int) bool {
	return token.Exchange(ctx, to, amount)
}

// sell token to gas
func Sale(to interop.Hash160, amount int) bool {
	return token.Sale(ctx, to, amount)
}

// grant nft to address
func GrantNFT(to interop.Hash160, nftType int, id []byte, name string, owner, image string, atk, def, hp string) bool {
	return token.GrantNFT(ctx, to, nftType, id, name, owner, image, atk, def, hp)
}

func SetNFTPrice(v int) bool {
	return token.SetNFTPrice(ctx, v)
}

// change gas/token exchange ratio
func ChangeTR(ratio int) bool {
	return token.ChangeTR(ratio)
}
