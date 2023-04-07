package permission

import (
	"context"
	"encoding/json"
	"errors"

	permissiontypes "github.com/bnb-chain/greenfield/x/permission/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventPutPolicy    = proto.MessageName(&permissiontypes.EventPutPolicy{})
	EventDeletePolicy = proto.MessageName(&permissiontypes.EventDeletePolicy{})
)

var policyEvents = map[string]bool{
	EventPutPolicy:    true,
	EventDeletePolicy: true,
}

var actionTypeMap = map[permissiontypes.ActionType]int{
	permissiontypes.ACTION_TYPE_ALL:            0,
	permissiontypes.ACTION_UPDATE_BUCKET_INFO:  1,
	permissiontypes.ACTION_DELETE_BUCKET:       2,
	permissiontypes.ACTION_CREATE_OBJECT:       3,
	permissiontypes.ACTION_DELETE_OBJECT:       4,
	permissiontypes.ACTION_COPY_OBJECT:         5,
	permissiontypes.ACTION_GET_OBJECT:          6,
	permissiontypes.ACTION_EXECUTE_OBJECT:      7,
	permissiontypes.ACTION_LIST_OBJECT:         8,
	permissiontypes.ACTION_UPDATE_GROUP_MEMBER: 9,
	permissiontypes.ACTION_DELETE_GROUP:        10,
	//permissiontypes.ACTION_GROUP_MEMBER:        11,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, event sdk.Event) error {
	if !policyEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}
	switch event.Type {
	case EventPutPolicy:
		putPolicy, ok := typedEvent.(*permissiontypes.EventPutPolicy)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateObject", "event", typedEvent)
			return errors.New("put policy event assert error")
		}
		return m.handlePutPolicy(ctx, block, putPolicy)
	case EventDeletePolicy:
		deletePolicy, ok := typedEvent.(*permissiontypes.EventDeletePolicy)
		if !ok {
			log.Errorw("type assert error", "type", "EventCancelCreateObject", "event", typedEvent)
			return errors.New("cancel delete policy event assert error")
		}
		return m.handleDeletePolicy(ctx, block, deletePolicy)
	}

	return nil
}

func (m *Module) handlePutPolicy(ctx context.Context, block *tmctypes.ResultBlock, policy *permissiontypes.EventPutPolicy) error {
	p := &models.Permission{
		PrincipalType:   int32(policy.Principal.Type),
		PrincipalValue:  policy.Principal.Value,
		ResourceType:    policy.ResourceType.String(),
		ResourceID:      common.BigToHash(policy.ResourceId.BigInt()),
		PolicyID:        common.BigToHash(policy.PolicyId.BigInt()),
		CreateTimestamp: block.Block.Time.Unix(),
	}

	statements := make([]*models.Statements, 0, 0)
	for _, statement := range policy.Statements {
		actionValue := 0
		for _, action := range statement.Actions {
			value, ok := actionTypeMap[action]
			if !ok {
				return errors.New("unknown action type action")
			}
			actionValue |= value
		}
		s := &models.Statements{
			PolicyID:       common.HexToHash(policy.PolicyId.String()),
			Effect:         statement.Effect.String(),
			ActionValue:    actionValue,
			ExpirationTime: statement.ExpirationTime.UTC().Unix(),
			LimitSize:      statement.LimitSize.Value,
		}
		if len(statement.Resources) != 0 {
			resources, err := json.Marshal(statement.Resources)
			if err != nil {
				return err
			}
			s.Resources = string(resources)
		}
		statements = append(statements, s)
	}

	// begin transaction
	tx := m.db.Begin(ctx)
	err1 := tx.SavePermission(ctx, p)
	err2 := tx.MultiSaveStatement(ctx, statements)
	err3 := tx.Commit()
	if err1 != nil || err2 != nil || err3 != nil {
		tx.Rollback()
		log.Errorw("failed to save policy", "permission err", err1, "statement err", err2, "commit err", err3)
		return errors.New("save policy transaction failed")
	}
	return nil
}

func (m *Module) handleDeletePolicy(ctx context.Context, block *tmctypes.ResultBlock, event *permissiontypes.EventDeletePolicy) error {
	// begin transaction
	tx := m.db.Begin(ctx)
	policyIDHash := common.BigToHash(event.PolicyId.BigInt())
	err1 := tx.UpdatePermission(ctx, &models.Permission{
		PolicyID:        policyIDHash,
		Removed:         true,
		UpdateTimestamp: block.Block.Time.Unix(),
	})
	err2 := tx.RemoveStatements(ctx, policyIDHash)
	err3 := tx.Commit()
	if err1 != nil || err2 != nil || err3 != nil {
		tx.Rollback()
		log.Errorw("failed to delete policy", "permission err", err1, "statement err", err2, "commit err", err3)
		return errors.New("delete policy transaction failed")
	}
	return nil
}
