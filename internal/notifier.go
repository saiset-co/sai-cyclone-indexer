package internal

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/saiset-co/sai-cyclone-indexer/logger"
	"github.com/saiset-co/sai-cyclone-indexer/utils"
)

type Notifier interface {
	Notify(data interface{}) error
}

type notifier struct {
	senderID string
	address  string
	token    string
}

type notificationData struct {
	From string      `json:"from"`
	Tx   interface{} `json:"tx"`
}

type notificationRequest struct {
	Method string           `json:"method"`
	Data   notificationData `json:"data"`
}

func NewNotifier(senderID, token, address string) Notifier {
	return &notifier{
		senderID: senderID,
		token:    token,
		address:  address,
	}
}

func (n *notifier) Notify(tx interface{}) error {
	req := notificationRequest{
		Method: "notify",
		Data: notificationData{
			From: n.senderID,
			Tx:   tx,
		},
	}

	payload, err := jsoniter.Marshal(&req)
	if err != nil {
		logger.Logger.Error("Notify", zap.Error(err))
		return err
	}

	_, err = utils.SaiQuerySender(bytes.NewReader(payload), n.address, n.token)

	return err
}
