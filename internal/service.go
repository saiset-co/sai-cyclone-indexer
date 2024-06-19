package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/saiset-co/sai-cyclone-indexer/logger"
	saiService "github.com/saiset-co/sai-service/service"
)

const (
	filePathAddresses   = "./addresses.json"
	filePathLatestBlock = "./latest_handled_block.data"
)

type InternalService struct {
	Context      *saiService.Context
	currentBlock int
	client       http.Client
	addresses    map[string]struct{}
	notifier     Notifier
	storage      Storage
}

func (is *InternalService) Init() {
	is.client = http.Client{
		Timeout: 5 * time.Second,
	}
	is.addresses = make(map[string]struct{})

	is.notifier = NewNotifier(
		cast.ToString(is.Context.GetConfig("notifier.sender_id", "")),
		cast.ToString(is.Context.GetConfig("notifier.token", "")),
		cast.ToString(is.Context.GetConfig("notifier.url", "")),
	)

	is.storage = NewStorage(
		cast.ToString(is.Context.GetConfig("storage.collection", "")),
		cast.ToString(is.Context.GetConfig("storage.token", "")),
		cast.ToString(is.Context.GetConfig("storage.url", "")),
	)

	err := is.loadAddresses()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Logger.Error("loadAddresses", zap.Error(err))
	}

	fileBytes, err := os.ReadFile(filePathLatestBlock)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Logger.Error("can't read "+filePathLatestBlock, zap.Error(err))
		}
	} else {
		latestHandledBlock, err := strconv.Atoi(string(fileBytes))
		if err != nil {
			logger.Logger.Error("strconv.Atoi", zap.Error(err))
		}

		is.currentBlock = latestHandledBlock
	}

	startBlock := cast.ToInt(is.Context.GetConfig("start_block", 0))
	if is.currentBlock < startBlock {
		is.currentBlock = startBlock
	}
}

func (is *InternalService) Process() {
	sleepDuration := cast.ToDuration(is.Context.GetConfig("sleep_duration", 2))

	for {
		select {
		case <-is.Context.Context.Done():
			logger.Logger.Debug("saiCosmosIndexer loop is done")
			return
		default:
			if len(is.addresses) == 0 {
				time.Sleep(time.Second * sleepDuration)
				continue
			}

			latestBlockHeight, err := is.getLatestBlock()
			if err != nil {
				logger.Logger.Error("getLatestBlock", zap.Error(err))
				time.Sleep(time.Second * sleepDuration)
				continue
			}

			if is.currentBlock >= latestBlockHeight {
				time.Sleep(time.Second * sleepDuration)
				continue
			}

			err = is.handleBlockTxs()
			if err != nil {
				logger.Logger.Error("handleBlockTxs", zap.Error(err))
				time.Sleep(time.Second * sleepDuration)
				continue
			}

			is.currentBlock += 1
		}
	}
}

func (is *InternalService) handleBlockTxs() error {
	blockTxs, err := is.getBlockTxs()
	if err != nil {
		logger.Logger.Error("handleBlockTxs", zap.Error(err))
		return err
	}

	logger.Logger.Debug("handleBlockTxs", zap.Any("blockTxs", blockTxs))

	var txArray []interface{}

	for _, txRes := range blockTxs.Transactions {
		var isReceiver = false
		_, isSender := is.addresses[txRes.Sender]

		CBytes, err := json.Marshal(txRes.Exec.VmResponse.C)
		if err != nil {
			logger.Logger.Error("handleBlockTxs", zap.Error(err))
		}

		for addr, _ := range is.addresses {
			if strings.Contains(string(CBytes), addr) {
				isReceiver = true
				break
			}

			if _, ok := txRes.Exec.VmResponse.D[addr]; ok {
				isReceiver = true
				break
			}

			if _, ok := txRes.Exec.VmResponse.R[addr]; ok {
				isReceiver = true
				break
			}

			if _, ok := txRes.Exec.VmResponse.T[addr]; ok {
				isReceiver = true
				break
			}

			if _, ok := txRes.Exec.VmResponse.V[addr]; ok {
				isReceiver = true
				break
			}
		}

		if !isReceiver && !isSender {
			continue
		}

		txArray = append(txArray, txRes)

		go is.sendTxNotification(txRes)
	}

	err = is.sendTxsToStorage(txArray)
	if err != nil {
		logger.Logger.Error("handleBlockTxs", zap.Error(err))
		return err
	}

	err = is.rewriteLastHandledBlock(is.currentBlock)

	return err
}

func (is *InternalService) sendTxsToStorage(txs []interface{}) error {
	err := is.storage.Save(txs)
	if err != nil {
		logger.Logger.Error("is.notifier.SendTx", zap.Error(err))
	}

	return err
}

func (is *InternalService) sendTxNotification(tx interface{}) {
	err := is.notifier.Notify(tx)
	if err != nil {
		logger.Logger.Error("is.notifier.SendTx", zap.Error(err))
	}
}
