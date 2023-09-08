package nep17

import (
	"github.com/nspcc-dev/neo-go/pkg/interop"
	"github.com/nspcc-dev/neo-go/pkg/interop/contract"
	"github.com/nspcc-dev/neo-go/pkg/interop/lib/address"
	"github.com/nspcc-dev/neo-go/pkg/interop/native/gas"
	"github.com/nspcc-dev/neo-go/pkg/interop/native/management"
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

const (
	CardNFT = 1
	HeroNFT = 2
)

var (
	IsOnce   = []byte("isOnce")
	NFTPrice = []byte("nftPrice")
	Hero     = address.ToHash160("NQ3cNE4bErfBXysYKgUZXXFg34x5G2vHKV")
	Card     = address.ToHash160("NNvZdSXuvTk8K5yYdAmdec7rWTDHF9VaB8")
	tratio   = int(1) // token/gas
	nratio   = int(1) // token/nft
)

// Token holds all token info
type Token struct {
	// Token name
	Name string
	// Ticker symbol
	Symbol string
	// Amount of decimals
	Decimals int
	// Token owner address
	Owner []byte
	// Total tokens * multiplier
	TotalSupply int
	// Storage key for circulation value
	CirculationKey string
	Multiplier     int
}

// getIntFromDB is a helper that checks for nil result of storage.Get and returns
// zero as the default value.
func getIntFromDB(ctx storage.Context, key []byte) int {
	var res int
	val := storage.Get(ctx, key)
	if val != nil {
		res = val.(int)
	}
	return res
}

// GetSupply gets the token totalSupply value from VM storage
func (t Token) GetSupply(ctx storage.Context) int {
	return getIntFromDB(ctx, []byte(t.CirculationKey))
}

// BalanceOf gets the token balance of a specific address
func (t Token) BalanceOf(ctx storage.Context, holder []byte) int {
	return getIntFromDB(ctx, holder)
}

// Transfer token from one user to another
func (t Token) Transfer(ctx storage.Context, from, to interop.Hash160, amount int, data any) bool {
	amountFrom := t.CanTransfer(ctx, from, to, amount)
	if amountFrom == -1 {
		return false
	}

	if amountFrom == 0 {
		storage.Delete(ctx, from)
	}

	if amountFrom > 0 {
		diff := amountFrom - amount
		storage.Put(ctx, from, diff)
	}

	amountTo := getIntFromDB(ctx, to)
	totalAmountTo := amountTo + amount
	if totalAmountTo != 0 {
		storage.Put(ctx, to, totalAmountTo)
	}

	runtime.Notify("Transfer", from, to, amount)
	if to != nil && management.GetContract(to) != nil {
		contract.Call(to, "onNEP17Payment", contract.All, from, amount, data)
	}
	return true
}

// CanTransfer returns the amount it can transfer
func (t Token) CanTransfer(ctx storage.Context, from []byte, to []byte, amount int) int {
	if len(to) != 20 || !IsUsableAddress(from) {
		return -1
	}

	amountFrom := getIntFromDB(ctx, from)
	if amountFrom < amount {
		return -1
	}

	// Tell Transfer the result is equal - special case since it uses Delete
	if amountFrom == amount {
		return 0
	}

	// return amountFrom value back to Transfer, reduces extra Get
	return amountFrom
}

// IsUsableAddress checks if the sender is either the correct Neo address or SC address
func IsUsableAddress(addr []byte) bool {
	if len(addr) == 20 {
		if runtime.CheckWitness(addr) {
			return true
		}

		// Check if a smart contract is calling scripthash
		callingScriptHash := runtime.GetCallingScriptHash()
		if callingScriptHash.Equals(addr) {
			return true
		}
	}

	return false
}

func (t Token) SetOnce(ctx storage.Context, v bool) bool {
	if !IsUsableAddress(t.Owner) {
		return false
	}

	storage.Put(ctx, IsOnce, v)
	return true
}

// Mint initial supply of tokens
func (t Token) Mint(ctx storage.Context, to interop.Hash160, amount int) bool {
	if !IsUsableAddress(t.Owner) {
		return false
	}

	isOnce := storage.Get(ctx, IsOnce)
	if isOnce != nil && isOnce.(bool) {
		minted := storage.Get(ctx, []byte("minted"))
		if minted != nil && minted.(bool) == true {
			return false
		}
	}

	storage.Put(ctx, []byte("minted"), true)

	v := getIntFromDB(ctx, to)
	v += amount
	storage.Put(ctx, to, v)

	v = getIntFromDB(ctx, []byte(t.CirculationKey))
	v += amount
	storage.Put(ctx, []byte(t.CirculationKey), v)

	// from is nil
	var from interop.Hash160
	runtime.Notify("Transfer", from, to, amount)
	return true
}

func (t Token) Exchange(ctx storage.Context, to interop.Hash160, amount int) bool {
	if !runtime.CheckWitness(to) {
		return false
	}

	gasMount := amount / tratio
	amount = gasMount * tratio

	if gas.BalanceOf(to) < gasMount {
		return false
	}

	if !gas.Transfer(to, t.Owner, gasMount, nil) {
		return false
	}

	if !t.Transfer(ctx, t.Owner, to, amount, nil) {
		return false
	}
	return true
}

func (t Token) Sale(ctx storage.Context, to interop.Hash160, amount int) bool {
	if !runtime.CheckWitness(to) {
		return false
	}

	gasMount := amount / tratio
	amount = gasMount * tratio

	if gas.BalanceOf(t.Owner) < gasMount {
		return false
	}

	if !gas.Transfer(t.Owner, to, gasMount, nil) {
		return false
	}

	if !t.Transfer(ctx, to, t.Owner, amount, nil) {
		return false
	}
	return true
}

func (t Token) GrantNFT(ctx storage.Context, to interop.Hash160, nftType int, id []byte, name string, owner string, image string, atk, def, hp string) bool {
	switch nftType {
	case CardNFT, HeroNFT:
	default:
		return false
	}

	price := getIntFromDB(ctx, NFTPrice)
	if price == 0 {
		price = t.Multiplier
	}

	if !t.Transfer(ctx, to, t.Owner, price, nil) {
		return false
	}

	switch nftType {
	case CardNFT:
		return contract.Call(Card, "mint", contract.All, to, id, name, image, atk, def, hp).(bool)
	}
	return contract.Call(Hero, "mint", contract.All, to, id, name, image, atk, def, hp).(bool)
}

func (t Token) SetNFTPrice(ctx storage.Context, v int) bool {
	if IsUsableAddress(t.Owner) || v < 0 {
		return false
	}

	storage.Put(ctx, NFTPrice, v)
	return true
}

func (t Token) ChangeTR(ratio int) bool {
	if !runtime.CheckWitness(t.Owner) {
		return false
	}
	tratio = ratio
	return true
}
