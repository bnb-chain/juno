package epoch

import (
	"context"
	"github.com/forbole/juno/v4/log"
)

func (m *Module) IsProcessed(height uint64) (bool, error) {
	ep, err := m.db.GetEpoch(context.Background())
	if err != nil {
		return false, err
	}
	log.Infof("epoch height:%d, cur height: %d", ep.BlockHeight, height)
	return ep.BlockHeight > int64(height), nil
}
