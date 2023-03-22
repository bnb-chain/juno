package models

import (
	"encoding/hex"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/juno/v4/common"
)

type Tx struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"-"`

	Hash    common.Hash `gorm:"column:hash;type:BINARY(32);not null;uniqueIndex:idx_hash"`
	Height  uint64      `gorm:"column:height;not null;uniqueIndex:idx_height_tx_index,priority:1"`
	TxIndex uint32      `gorm:"column:tx_index;not null;uniqueIndex:idx_height_tx_index,priority:2"`

	Success     bool   `gorm:"column:success"`
	Messages    string `gorm:"column:messages;type:json;not null;default:(JSON_ARRAY())"`
	Memo        string `gorm:"column:memo"`
	Signatures  string `gorm:"column:signatures"`
	SignerInfos string `gorm:"column:signer_infos;type:json;not null;default:(JSON_ARRAY())"`
	Fee         string `gorm:"column:fee;type:json;not null;default:(JSON_ARRAY())"`

	GasWanted uint64 `gorm:"column:gas_wanted"`
	GasUsed   uint64 `gorm:"column:gas_used"`
	RawLog    string `gorm:"column:raw_log"`
	Logs      string `gorm:"column:logs;type:json;not null;default:(JSON_ARRAY())"`

	Timestamp uint64 `gorm:"timestamp"` // refer block.header.timestamp
}

func (*Tx) TableName() string {
	return "txs"
}

func (t *Tx) ToTmTx() *ResultTx {
	txResult := ResponseDeliverTx{
		Log:       t.Logs,
		GasWanted: int64(t.GasWanted),
		GasUsed:   int64(t.GasUsed),
		Messages:  t.Messages,
	}

	if t.Success {
		txResult.Code = 0
	} else {
		txResult.Code = 1
	}

	txHash, _ := hex.DecodeString(t.Hash.Hex()[2:])

	return &ResultTx{
		Hash:     txHash,
		Height:   int64(t.Height),
		Index:    t.TxIndex,
		Time:     time.Unix(0, int64(t.Timestamp)).UTC(),
		TxResult: txResult,
	}
}

type ResponseDeliverTx struct {
	Code      uint32       `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	Data      []byte       `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log       string       `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string       `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasWanted int64        `protobuf:"varint,5,opt,name=gas_wanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed   int64        `protobuf:"varint,6,opt,name=gas_used,proto3" json:"gas_used,omitempty"`
	Events    []abci.Event `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"`
	Codespace string       `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
	Messages  string       `json:"messages,omitempty"`
}

// ResultTx Result of querying for a tx
type ResultTx struct {
	Hash     bytes.HexBytes    `json:"hash"`
	Height   int64             `json:"height"`
	Index    uint32            `json:"index"`
	Time     time.Time         `json:"time"`
	TxResult ResponseDeliverTx `json:"tx_result"`
	Tx       tmtypes.Tx        `json:"tx"`
	Proof    tmtypes.TxProof   `json:"proof,omitempty"`
}
