package group

import (
	"context"
	"errors"

	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

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

var groupEvents = map[string]bool{
	EventCreateGroup:       true,
	EventDeleteGroup:       true,
	EventLeaveGroup:        true,
	EventUpdateGroupMember: true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, _ common.Hash, event sdk.Event) error {
	if !groupEvents[event.Type] {
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
	for _, member := range createGroup.Members {
		groupItem := &models.Group{
			Owner:      common.HexToAddress(createGroup.OwnerAddress),
			GroupID:    common.BigToHash(createGroup.GroupId.BigInt()),
			GroupName:  createGroup.GroupName,
			SourceType: createGroup.SourceType.String(),
			AccountID:  common.HexToHash(member),

			CreateAt:   block.Block.Height,
			CreateTime: block.Block.Time.UTC().Unix(),
			UpdateAt:   block.Block.Height,
			UpdateTime: block.Block.Time.UTC().Unix(),
			Removed:    false,
		}
		membersToAddList = append(membersToAddList, groupItem)
	}

	return m.db.CreateGroup(ctx, membersToAddList)
}

func (m *Module) handleDeleteGroup(ctx context.Context, block *tmctypes.ResultBlock, deleteGroup *storagetypes.EventDeleteGroup) error {
	group := &models.Group{
		Owner:     common.HexToAddress(deleteGroup.OwnerAddress),
		GroupID:   common.BigToHash(deleteGroup.GroupId.BigInt()),
		GroupName: deleteGroup.GroupName,

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    true,
	}
	return m.db.DeleteGroup(ctx, group)
}

func (m *Module) handleLeaveGroup(ctx context.Context, block *tmctypes.ResultBlock, leaveGroup *storagetypes.EventLeaveGroup) error {
	group := &models.Group{
		Owner:     common.HexToAddress(leaveGroup.OwnerAddress),
		GroupID:   common.BigToHash(leaveGroup.GroupId.BigInt()),
		GroupName: leaveGroup.GroupName,
		AccountID: common.HexToHash(leaveGroup.MemberAddress),

		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
		Removed:    true,
	}
	return m.db.UpdateGroup(ctx, group)
}

func (m *Module) handleUpdateGroupMember(ctx context.Context, block *tmctypes.ResultBlock, updateGroupMember *storagetypes.EventUpdateGroupMember) error {

	membersToAdd := updateGroupMember.MembersToAdd
	membersToDelete := updateGroupMember.MembersToDelete

	var membersToAddList []*models.Group
	for _, memberToAdd := range membersToAdd {
		groupItem := &models.Group{
			Owner:           common.HexToAddress(updateGroupMember.OwnerAddress),
			GroupID:         common.BigToHash(updateGroupMember.GroupId.BigInt()),
			GroupName:       updateGroupMember.GroupName,
			AccountID:       common.HexToHash(memberToAdd),
			OperatorAddress: common.HexToAddress(updateGroupMember.OperatorAddress),

			UpdateAt:   block.Block.Height,
			UpdateTime: block.Block.Time.UTC().Unix(),
			Removed:    false,
		}
		membersToAddList = append(membersToAddList, groupItem)
	}

	m.db.CreateGroup(ctx, membersToAddList)

	for _, memberToDelete := range membersToDelete {
		groupItem := &models.Group{
			Owner:           common.HexToAddress(updateGroupMember.OwnerAddress),
			GroupID:         common.BigToHash(updateGroupMember.GroupId.BigInt()),
			GroupName:       updateGroupMember.GroupName,
			AccountID:       common.HexToHash(memberToDelete),
			OperatorAddress: common.HexToAddress(updateGroupMember.OperatorAddress),

			UpdateAt:   block.Block.Height,
			UpdateTime: block.Block.Time.UTC().Unix(),
			Removed:    true,
		}
		m.db.UpdateGroup(ctx, groupItem)
	}

	return nil
}
