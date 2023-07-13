package payment

import (
	"context"
	"errors"

	paymenttypes "github.com/bnb-chain/greenfield/x/payment/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventPaymentAccountUpdate = proto.MessageName(&paymenttypes.EventPaymentAccountUpdate{})
	EventStreamRecordUpdate   = proto.MessageName(&paymenttypes.EventStreamRecordUpdate{})
)

var PaymentEvents = map[string]bool{
	EventPaymentAccountUpdate: true,
	EventStreamRecordUpdate:   true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, _ common.Hash, event sdk.Event) error {
	_, err := m.Handle(ctx, block, common.Hash{}, event, true)
	return err
}

func (m *Module) Handle(ctx context.Context, block *tmctypes.ResultBlock, _ common.Hash, event sdk.Event, OperateDB bool) (interface{}, error) {
	log.Infof("Handle")
	if !PaymentEvents[event.Type] {
		log.Infof("event type: %s", event.Type)
		return nil, nil
	}
	log.Infof("Handle1")
	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return nil, err
	}
	log.Infof("Handle2")
	switch event.Type {
	case EventPaymentAccountUpdate:
		log.Infof("Aaaaaaaaaa")
		paymentAccountUpdate, ok := typedEvent.(*paymenttypes.EventPaymentAccountUpdate)
		if !ok {
			log.Errorw("type assert error", "type", "EventPaymentAccountUpdate", "event", typedEvent)
			return nil, errors.New("update payment account event assert error")
		}
		data := m.handlePaymentAccountUpdate(ctx, block, paymentAccountUpdate)
		if !OperateDB {
			return data, nil
		}
		return m.db.SavePaymentAccount(ctx, data), nil
	case EventStreamRecordUpdate:
		log.Infof("bbbbbbbbbbbb")
		streamRecordUpdate, ok := typedEvent.(*paymenttypes.EventStreamRecordUpdate)
		if !ok {
			log.Errorw("type assert error", "type", "EventStreamRecordUpdate", "event", typedEvent)
			return nil, errors.New("update stream record event assert error")
		}
		data := m.handleEventStreamRecordUpdate(ctx, streamRecordUpdate)
		if !OperateDB {
			log.Infof("data:%v", data)
			return data, nil
		}
		return nil, m.db.SaveStreamRecord(ctx, data)
	}
	log.Infof("Handle3")
	return nil, nil
}

func (m *Module) ExtractEvent(ctx context.Context, block *tmctypes.ResultBlock, _ common.Hash, event sdk.Event) (interface{}, error) {
	return m.Handle(ctx, block, common.Hash{}, event, false)
}

func (m *Module) handlePaymentAccountUpdate(ctx context.Context, block *tmctypes.ResultBlock, paymentAccountUpdate *paymenttypes.EventPaymentAccountUpdate) *models.PaymentAccount {
	return &models.PaymentAccount{
		Addr:       common.HexToAddress(paymentAccountUpdate.Addr),
		Owner:      common.HexToAddress(paymentAccountUpdate.Owner),
		Refundable: paymentAccountUpdate.Refundable,
		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
	}
}

func (m *Module) handleEventStreamRecordUpdate(ctx context.Context, streamRecordUpdate *paymenttypes.EventStreamRecordUpdate) *models.StreamRecord {
	return &models.StreamRecord{
		Account:         common.HexToAddress(streamRecordUpdate.Account),
		CrudTimestamp:   streamRecordUpdate.CrudTimestamp,
		NetflowRate:     (*common.Big)(streamRecordUpdate.NetflowRate.BigInt()),
		StaticBalance:   (*common.Big)(streamRecordUpdate.StaticBalance.BigInt()),
		BufferBalance:   (*common.Big)(streamRecordUpdate.BufferBalance.BigInt()),
		LockBalance:     (*common.Big)(streamRecordUpdate.LockBalance.BigInt()),
		Status:          streamRecordUpdate.Status.String(),
		SettleTimestamp: streamRecordUpdate.SettleTimestamp,
	}
}
