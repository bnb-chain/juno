package group

import (
	"context"
	"errors"

	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventCreateGroup       = proto.MessageName(&storagetypes.EventCreateGroup{})
	EventDeleteGroup       = proto.MessageName(&storagetypes.EventDeleteGroup{})
	EventLeaveGroup        = proto.MessageName(&storagetypes.EventLeaveGroup{})
	EventUpdateGroupMember = proto.MessageName(&storagetypes.EventUpdateGroupMember{})
)

var GroupEvents = map[string]bool{
	EventCreateGroup:       true,
	EventDeleteGroup:       true,
	EventLeaveGroup:        true,
	EventUpdateGroupMember: true,
}

func (m *Module) ExtractEvent(data interface{}, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) {
	return
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, _ common.Hash, event sdk.Event) error {
	if !GroupEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventCreateGroup:
		createGroup, ok := typedEvent.(*storagetypes.EventCreateGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateGroup", "event", typedEvent)
			return errors.New("create group event assert error")
		}
		return m.handleCreateGroup(ctx, block, createGroup)
	case EventUpdateGroupMember:
		updateGroupMember, ok := typedEvent.(*storagetypes.EventUpdateGroupMember)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateGroupMember", "event", typedEvent)
			return errors.New("update group member event assert error")
		}
		return m.handleUpdateGroupMember(ctx, block, updateGroupMember)

	case EventDeleteGroup:
		deleteGroup, ok := typedEvent.(*storagetypes.EventDeleteGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventDeleteGroup", "event", typedEvent)
			return errors.New("delete group event assert error")
		}
		return m.handleDeleteGroup(ctx, block, deleteGroup)
	case EventLeaveGroup:
		leaveGroup, ok := typedEvent.(*storagetypes.EventLeaveGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventLeaveGroup", "event", typedEvent)
			return errors.New("leave group event assert error")
		}
		return m.handleLeaveGroup(ctx, block, leaveGroup)
	}
	return nil
}

func (m *Module) handleCreateGroup(ctx context.Context, block *tmctypes.ResultBlock, createGroup *storagetypes.EventCreateGroup) error {

	var membersToAddList []*models.Group

	//create group first
	groupItem := &models.Group{
		Owner:      common.HexToAddress(createGroup.Owner),
		GroupID:    common.BigToHash(createGroup.GroupId.BigInt()),
		GroupName:  createGroup.GroupName,
		SourceType: createGroup.SourceType.String(),
		AccountID:  common.HexToAddress("0"),
		Extra:      createGroup.Extra,

		CreateAt:   block.Block.Height,
		CreateTime: block.Block.Time.UTC().Unix(),
		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    false,
	}
	membersToAddList = append(membersToAddList, groupItem)

	for _, member := range createGroup.Members {
		groupMemberItem := &models.Group{
			Owner:      common.HexToAddress(createGroup.Owner),
			GroupID:    common.BigToHash(createGroup.GroupId.BigInt()),
			GroupName:  createGroup.GroupName,
			SourceType: createGroup.SourceType.String(),
			AccountID:  common.HexToAddress(member),
			Extra:      createGroup.Extra,

			CreateAt:   block.Block.Height,
			CreateTime: block.Block.Time.UTC().Unix(),
			UpdateAt:   block.Block.Height,
			UpdateTime: block.Block.Time.UTC().Unix(),
			Removed:    false,
		}
		membersToAddList = append(membersToAddList, groupMemberItem)
	}

	return m.db.CreateGroup(ctx, membersToAddList)
}

func (m *Module) handleDeleteGroup(ctx context.Context, block *tmctypes.ResultBlock, deleteGroup *storagetypes.EventDeleteGroup) error {
	group := &models.Group{
		Owner:     common.HexToAddress(deleteGroup.Owner),
		GroupID:   common.BigToHash(deleteGroup.GroupId.BigInt()),
		GroupName: deleteGroup.GroupName,

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    true,
	}

	//update group item
	groupItem := &models.Group{
		GroupID:   common.BigToHash(deleteGroup.GroupId.BigInt()),
		AccountID: common.HexToAddress("0"),

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    true,
	}
	m.db.UpdateGroup(ctx, groupItem)

	return m.db.DeleteGroup(ctx, group)
}

func (m *Module) handleLeaveGroup(ctx context.Context, block *tmctypes.ResultBlock, leaveGroup *storagetypes.EventLeaveGroup) error {
	group := &models.Group{
		Owner:     common.HexToAddress(leaveGroup.Owner),
		GroupID:   common.BigToHash(leaveGroup.GroupId.BigInt()),
		GroupName: leaveGroup.GroupName,
		AccountID: common.HexToAddress(leaveGroup.MemberAddress),

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    true,
	}

	//update group item
	groupItem := &models.Group{
		GroupID:   common.BigToHash(leaveGroup.GroupId.BigInt()),
		AccountID: common.HexToAddress("0"),

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    false,
	}
	m.db.UpdateGroup(ctx, groupItem)

	return m.db.UpdateGroup(ctx, group)
}

func (m *Module) handleUpdateGroupMember(ctx context.Context, block *tmctypes.ResultBlock, updateGroupMember *storagetypes.EventUpdateGroupMember) error {

	membersToAdd := updateGroupMember.MembersToAdd
	membersToDelete := updateGroupMember.MembersToDelete

	var membersToAddList []*models.Group

	if len(membersToAdd) > 0 {
		for _, memberToAdd := range membersToAdd {
			groupItem := &models.Group{
				Owner:     common.HexToAddress(updateGroupMember.Owner),
				GroupID:   common.BigToHash(updateGroupMember.GroupId.BigInt()),
				GroupName: updateGroupMember.GroupName,
				AccountID: common.HexToAddress(memberToAdd),
				Operator:  common.HexToAddress(updateGroupMember.Operator),

				CreateAt:   block.Block.Height,
				CreateTime: block.Block.Time.UTC().Unix(),
				UpdateAt:   block.Block.Height,
				UpdateTime: block.Block.Time.UTC().Unix(),
				Removed:    false,
			}
			membersToAddList = append(membersToAddList, groupItem)
		}
		m.db.CreateGroup(ctx, membersToAddList)
	}

	for _, memberToDelete := range membersToDelete {
		groupItem := &models.Group{
			Owner:     common.HexToAddress(updateGroupMember.Owner),
			GroupID:   common.BigToHash(updateGroupMember.GroupId.BigInt()),
			GroupName: updateGroupMember.GroupName,
			AccountID: common.HexToAddress(memberToDelete),
			Operator:  common.HexToAddress(updateGroupMember.Operator),

			UpdateAt:   block.Block.Height,
			UpdateTime: block.Block.Time.UTC().Unix(),
			Removed:    true,
		}
		m.db.UpdateGroup(ctx, groupItem)
	}

	//update group item
	groupItem := &models.Group{
		GroupID:   common.BigToHash(updateGroupMember.GroupId.BigInt()),
		AccountID: common.HexToAddress("0"),

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    false,
	}
	m.db.UpdateGroup(ctx, groupItem)

	return nil
}
