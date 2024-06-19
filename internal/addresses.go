package internal

import (
	"net/http"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/saiset-co/sai-cyclone-indexer/logger"
)

var mu = &sync.Mutex{}

func (is *InternalService) addAddress(data interface{}) (string, int, error) {
	address, ok := data.(string)
	if !ok {
		return "address should be a string", http.StatusBadRequest, nil
	}
	_, exist := is.addresses[address]
	if exist {
		return "Already exists", http.StatusOK, nil
	}

	mu.Lock()
	is.addresses[address] = struct{}{}
	mu.Unlock()

	err := is.rewriteAddressesFile()
	if err != nil {
		mu.Lock()
		delete(is.addresses, address)
		mu.Unlock()
		return "", http.StatusInternalServerError, err
	}

	return "Added", http.StatusOK, nil
}

func (is *InternalService) deleteAddress(data interface{}) (string, int, error) {
	address, ok := data.(string)
	if !ok {
		return "address should be a string", http.StatusBadRequest, nil
	}

	_, exist := is.addresses[address]
	if !exist {
		return "deleteAddress", http.StatusOK, nil
	}

	mu.Lock()
	delete(is.addresses, address)
	mu.Unlock()

	err := is.rewriteAddressesFile()
	if err != nil {
		logger.Logger.Error("deleteAddress", zap.Error(err))
		return "deleteAddress", http.StatusInternalServerError, err
	}

	return "deleteAddress", http.StatusOK, nil
}

func (is *InternalService) loadAddresses() error {
	fileBytes, err := os.ReadFile(filePathAddresses)
	if err != nil {
		logger.Logger.Error("loadAddresses", zap.Error(err))
		return err
	}

	var addressArray []string
	err = jsoniter.Unmarshal(fileBytes, &addressArray)
	if err != nil {
		logger.Logger.Error("loadAddresses", zap.Error(err))
		return err
	}

	mu.Lock()
	for _, address := range addressArray {
		is.addresses[address] = struct{}{}
	}
	mu.Unlock()

	return nil
}

func (is *InternalService) rewriteAddressesFile() error {
	addressArray := make([]string, len(is.addresses))
	mu.Lock()
	i := 0
	for k, _ := range is.addresses {
		addressArray[i] = k
		i++
	}
	mu.Unlock()

	jsonBytes, err := jsoniter.Marshal(&addressArray)
	if err != nil {
		logger.Logger.Error("rewriteAddressesFile", zap.Error(err))
		return err
	}

	err = os.WriteFile(filePathAddresses, jsonBytes, os.ModePerm)

	return err
}
