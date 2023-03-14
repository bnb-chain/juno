package models

import (
	"encoding/hex"
	"time"

	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	"github.com/forbole/juno/v4/common"
)

type PartSetHeader struct {
	Total uint32      `gorm:"-"`
	Hash  common.Hash `gorm:"-"`
}

type BlockID struct {
	Hash          common.Hash   `gorm:"column:hash;type:BINARY(32);not null;uniqueIndex:idx_hash"`
	PartSetHeader PartSetHeader `gorm:"-"`
}

type Header struct {
	Height             uint64      `gorm:"column:height;not null;uniqueIndex:idx_height"`
	LastCommitHash     common.Hash `gorm:"last_commit_hash;type:BINARY(32)"`     // commit from validators from the last block
	DataHash           common.Hash `gorm:"data_hash;type:BINARY(32)"`            // transactions
	ValidatorsHash     common.Hash `gorm:"validators_hash;type:BINARY(32)"`      // validators for the current block
	NextValidatorsHash common.Hash `gorm:"next_validators_hash;type:BINARY(32)"` // validators for the next block
	ConsensusHash      common.Hash `gorm:"consensus_hash;type:BINARY(32)"`       // consensus params for current block
	AppHash            common.Hash `gorm:"app_hash;type:BINARY(32)"`             // state after txs from the previous block
	// root hash of all results from the txs from the previous block
	// see `deterministicResponseDeliverTx` to understand which parts of a tx is hashed into here
	LastResultsHash common.Hash `gorm:"last_results_hash;type:BINARY(32)"`
	// consensus info
	EvidenceHash    common.Hash    `gorm:"evidence_hash;type:BINARY(32)"` // evidence included in the block
	ProposerAddress common.Address `gorm:"column:proposer_address;type:BINARY(20);index:idx_proposer_address"`

	Timestamp uint64 `gorm:"column:timestamp"`
}

type Block struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	BlockID
	Header

	NumTxs   uint64 `gorm:"column:num_txs"`
	TotalGas uint64 `gorm:"column:total_gas"`
}

func (*Block) TableName() string {
	return "blocks"
}

func (b *Block) ToTmBlock() *tmctypes.ResultBlock {
	blockID := tmtypes.BlockID{
		Hash: b.Hash.Bytes(),
	}
	header := tmtypes.Header{
		Version: tmversion.Consensus{Block: version.BlockProtocol},
		//ChainID: ,
		Height: int64(b.Height),
		Time:   time.Unix(int64(b.Timestamp), 0),
		//LastBlockID: ,
	}
	header.LastCommitHash, _ = hex.DecodeString(b.LastResultsHash.Hex()[2:])
	header.DataHash, _ = hex.DecodeString(b.DataHash.Hex()[2:])
	header.ValidatorsHash, _ = hex.DecodeString(b.DataHash.Hex()[2:])
	header.NextValidatorsHash, _ = hex.DecodeString(b.NextValidatorsHash.Hex()[2:])
	header.ConsensusHash, _ = hex.DecodeString(b.ConsensusHash.Hex()[2:])
	header.AppHash, _ = hex.DecodeString(b.AppHash.Hex()[2:])
	header.LastResultsHash, _ = hex.DecodeString(b.LastResultsHash.Hex()[2:])
	header.EvidenceHash, _ = hex.DecodeString(b.EvidenceHash.Hex()[2:])
	header.ProposerAddress, _ = hex.DecodeString(b.ProposerAddress.Hex()[2:])

	block := &tmtypes.Block{
		Header: header,
	}
	return &tmctypes.ResultBlock{
		BlockID: blockID,
		Block:   block,
	}
}

// NewBlockFromTmBlock builds a new Block instance from a given ResultBlock object
func NewBlockFromTmBlock(blk *tmctypes.ResultBlock, totalGas uint64) *Block {
	return &Block{
		BlockID: BlockID{
			Hash: common.HexToHash(blk.Block.Hash().String()),
		},
		Header: Header{
			uint64(blk.Block.Height),
			common.HexToHash(blk.Block.Header.LastCommitHash.String()),
			common.HexToHash(blk.Block.Header.DataHash.String()),
			common.HexToHash(blk.Block.Header.ValidatorsHash.String()),
			common.HexToHash(blk.Block.Header.NextValidatorsHash.String()),
			common.HexToHash(blk.Block.Header.ConsensusHash.String()),
			common.HexToHash(blk.Block.Header.AppHash.String()),
			common.HexToHash(blk.Block.Header.LastResultsHash.String()),
			common.HexToHash(blk.Block.Header.EvidenceHash.String()),
			common.HexToAddress(blk.Block.Header.ProposerAddress.String()),
			uint64(blk.Block.Time.Unix()),
		},
		NumTxs:   uint64(len(blk.Block.Txs)),
		TotalGas: totalGas,
	}
}

type Genesis struct {
	OneRowId      bool   `gorm:"one_row_id;primaryKey;default:true"`
	ChainID       string `gorm:"chain_id"`
	Timestamp     uint64 `gorm:"timestamp"`
	InitialHeight uint64 `gorm:"initial_height"`
}

func (*Genesis) TableName() string {
	return "geneses"
}
