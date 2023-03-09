package epoch

import (
	"context"
	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/models"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (m *Module) SaveEpoch(ctx context.Context, block *tmctypes.ResultBlock) error {
	return m.db.SaveEpoch(ctx, &models.Epoch{
		ID:          1,
		BlockHeight: block.Block.Height,
		BlockHash:   common.HexToHash(block.BlockID.Hash.String()),
		UpdateTime:  block.Block.Time.Unix(),
	})
}

func (m *Module) GetEpoch(ctx context.Context) (*models.Epoch, error) {
	return m.db.GetEpoch(ctx)
}
