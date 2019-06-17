package executor

import (
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common/db/table"
	pty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"gitlab.33.cn/chain33/chain33/common"
)

var (
	opt = table.Option{
		Prefix:  tokennoteLocalPre,
		Name:    "test",
		Primary: "hash",
		Index:   []string{"symbol", "isBurn"},
	}
	marketopt = table.Option{
		Prefix:  tokennoteLocalPre,
		Name:    "market",
		Primary: "hash",
		Index:   []string{"symbol","addr"},
	}
)

type TransactionRow struct {
	*types.Transaction
}

func NewTransactionRow() *TransactionRow {
	return &TransactionRow{Transaction: &types.Transaction{}}
}

func (tx *TransactionRow) CreateRow() *table.Row {
	return &table.Row{Data: &types.Transaction{}}
}

func (tx *TransactionRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*types.Transaction); ok {
		tx.Transaction = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *TransactionRow) Get(key string) ([]byte, error) {
	if key == "hash" {
		return tx.Hash(), nil
	} else if key == "symbol" {
		var action pty.TokennoteAction
		err := types.Decode(tx.Payload,&action)
		if err != nil {
			return nil,err
		}
		switch action.Ty {
		case pty.TokennoteActionMint:
			return []byte(action.GetTokennoteMint().GetSymbol()),nil
		case pty.TokennoteActionBurn:
			return []byte(action.GetTokennoteBurn().GetSymbol()),nil
		default:
			return nil,types.ErrNotFound
		}
	} else if key == "isBurn" {
		var action pty.TokennoteAction
		err := types.Decode(tx.Payload,&action)
		if err != nil {
			return nil,err
		}
		switch action.Ty {
		case pty.TokennoteActionMint:
			return []byte(string(pty.TokennoteActionMint)),nil
		case pty.TokennoteActionBurn:
			return []byte(string(pty.TokennoteActionBurn)),nil
		default:
			return nil,types.ErrNotFound
		}
	}
	return nil, types.ErrNotFound
}


type TokennoteMarketRow struct {
	*pty.Tokennote
}

func NewTokennoteMarketRow() *TokennoteMarketRow {
	return &TokennoteMarketRow{ Tokennote:&pty.Tokennote{}}
}

func (tx *TokennoteMarketRow) CreateRow() *table.Row {
	return &table.Row{Data: &pty.Tokennote{}}
}

func (tx *TokennoteMarketRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*pty.Tokennote); ok {
		tx.Tokennote = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *TokennoteMarketRow) Get(key string) ([]byte, error) {
	if key == "addr" {
		return []byte(tx.Issuer), nil
	} else if key == "symbol" {
		return []byte(tx.Currency),nil
	} else if key == "hash" {
		return common.Sha256(types.Encode(tx)),nil
	}
	return nil, types.ErrNotFound
}

type TokennoteAccountCreditRow struct {
	*pty.TokennoteAccountCredit
}

func NewTokennoteAccountCreditRow() *TokennoteAccountCreditRow {
	return &TokennoteAccountCreditRow{TokennoteAccountCredit:&pty.TokennoteAccountCredit{}}
}

func (tx *TokennoteAccountCreditRow) CreateRow() *table.Row {
	return &table.Row{Data: &pty.Tokennote{}}
}

func (tx *TokennoteAccountCreditRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*pty.TokennoteAccountCredit); ok {
		tx.TokennoteAccountCredit = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *TokennoteAccountCreditRow) Get(key string) ([]byte, error) {
	if key == "Addr" {
		return []byte(tx.Addr), nil
	}
	return nil, types.ErrNotFound
}