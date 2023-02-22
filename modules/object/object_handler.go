package object

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ObjectModule struct {
}

func (bucket *ObjectModule) HandleEvent(index int, event sdk.Event) error {
	return nil
}

func (bucket *ObjectModule) PrepareTables() {

}
