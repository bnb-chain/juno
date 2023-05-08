package parser

import (
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/types"
)

// FindValidatorByAddr finds a validator by a consensus address given a set of
// Tendermint validators for a particular block. If no validator is found, nil
// is returned.
func FindValidatorByAddr(consAddr string, vals *tmctypes.ResultValidators) *tmtypes.Validator {
	for _, val := range vals.Validators {
		if consAddr == sdk.ConsAddress(val.Address).String() {
			return val
		}
	}
	return nil
}

// SumGasTxs returns the total gas consumed by a set of transactions.
func SumGasTxs(txs []*types.Tx) uint64 {
	var totalGas uint64

	for _, tx := range txs {
		totalGas += uint64(tx.GasUsed)
	}

	return totalGas
}
