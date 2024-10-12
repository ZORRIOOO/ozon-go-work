package utils

import (
	"homework/loms/internal/infra/kafka/types"
	orderModel "homework/loms/internal/model/order"
	"time"
)

func RepackPayload(id orderModel.Id, eventStatus orderModel.Status) types.MessagePayload {
	return types.MessagePayload{
		OrderId:  id,
		Event:    eventStatus,
		Datetime: time.Now().UTC().String(),
	}
}
