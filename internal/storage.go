package internal

import (
	"go.uber.org/zap"

	"github.com/saiset-co/sai-cyclone-indexer/logger"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
)

type Storage interface {
	Save(data []interface{}) error
}

type storage struct {
	saiStorage adapter.SaiStorage
	collection string
}

func NewStorage(collection, token, address string) Storage {
	return &storage{
		collection: collection,
		saiStorage: adapter.SaiStorage{
			Url:   address,
			Token: token,
		},
	}
}

func (s *storage) Save(txs []interface{}) error {
	storageRequest := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: s.collection,
			Documents:  txs,
		},
	}

	_, err := s.saiStorage.Send(storageRequest)
	if err != nil {
		logger.Logger.Error("Save", zap.Error(err))
		return err
	}

	return nil
}
