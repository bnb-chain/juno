package payment

import (
	"context"
	"errors"

	paymenttypes "github.com/bnb-chain/greenfield/x/payment/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventPaymentAccountUpdate = proto.MessageName(&paymenttypes.EventPaymentAccountUpdate{})
	EventStreamRecordUpdate   = proto.MessageName(&paymenttypes.EventStreamRecordUpdate{})
)

var paymentEvents = map[string]bool{
	EventPaymentAccountUpdate: true,
	EventStreamRecordUpdate:   true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, _ common.Hash, event sdk.Event) error {
	if !paymentEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventPaymentAccountUpdate:
		paymentAccountUpdate, ok := typedEvent.(*paymenttypes.EventPaymentAccountUpdate)
		if !ok {
			log.Errorw("type assert error", "type", "EventPaymentAccountUpdate", "event", typedEvent)
			return errors.New("update payment account event assert error")
		}
		return m.handlePaymentAccountUpdate(ctx, block, paymentAccountUpdate)
	case EventStreamRecordUpdate:
		streamRecordUpdate, ok := typedEvent.(*paymenttypes.EventStreamRecordUpdate)
		if !ok {
			log.Errorw("type assert error", "type", "EventStreamRecordUpdate", "event", typedEvent)
			return errors.New("update stream record event assert error")
		}
		return m.handleEventStreamRecordUpdate(ctx, streamRecordUpdate)
	}

	return nil
}

func (m *Module) handlePaymentAccountUpdate(ctx context.Context, block *tmctypes.ResultBlock, paymentAccountUpdate *paymenttypes.EventPaymentAccountUpdate) error {
	paymentAccount := &models.PaymentAccount{
		Addr:       common.HexToAddress(paymentAccountUpdate.Addr),
		Owner:      common.HexToAddress(paymentAccountUpdate.Owner),
		Refundable: paymentAccountUpdate.Refundable,
		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().Unix(),
	}

	return m.db.SavePaymentAccount(ctx, paymentAccount)
}

func (m *Module) handleEventStreamRecordUpdate(ctx context.Context, streamRecordUpdate *paymenttypes.EventStreamRecordUpdate) error {
	streamRecord := &models.StreamRecord{
		Account:         common.HexToAddress(streamRecordUpdate.Account),
		CrudTimestamp:   streamRecordUpdate.CrudTimestamp,
		NetflowRate:     (*common.Big)(streamRecordUpdate.NetflowRate.BigInt()),
		StaticBalance:   (*common.Big)(streamRecordUpdate.StaticBalance.BigInt()),
		BufferBalance:   (*common.Big)(streamRecordUpdate.BufferBalance.BigInt()),
		LockBalance:     (*common.Big)(streamRecordUpdate.LockBalance.BigInt()),
		Status:          streamRecordUpdate.Status.String(),
		SettleTimestamp: streamRecordUpdate.SettleTimestamp,
	}

	outflows, err := jsoniter.Marshal(streamRecordUpdate.OutFlows)
	if err != nil {
		return errors.New("marshal stream record outflows failed")
	}

	streamRecord.OutFlows = outflows

	return m.db.SaveStreamRecord(ctx, streamRecord)
}
