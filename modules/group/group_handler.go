package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GroupModule struct {
}

func (bucket *GroupModule) HandleEvent(index int, event sdk.Event) error {
	return nil
}

func (bucket *GroupModule) PrepareTables() {

}
