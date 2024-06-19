package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/saiset-co/sai-cyclone-indexer/internal/model"
	"github.com/saiset-co/sai-cyclone-indexer/logger"
)

func (is *InternalService) getLatestBlock() (int, error) {
	query := model.Json{"number": ""}
	lb := model.LatestBlock{}

	err := is.sendQuery("/api/gw/v1/block", query, &lb)
	if err != nil {
		logger.Logger.Error("getLatestBlock", zap.Error(err))
	}

	blockHeight, err := strconv.Atoi(lb.Number)
	if err != nil {
		logger.Logger.Error("sendQuery", zap.Error(err))
		return 0, err
	}

	logger.Logger.Debug("getLatestBlock", zap.Any("blockHeight", blockHeight))

	return blockHeight, nil
}

func (is *InternalService) getBlockTxs() (model.TxResponse, error) {
	var txs model.TxResponse
	query := model.Json{"number": strconv.Itoa(is.currentBlock)}

	logger.Logger.Debug("getLatestBlock", zap.Any("query", query))

	err := is.sendQuery("/api/gw/v1/block/transactions", query, &txs)
	if err != nil {
		logger.Logger.Error("getLatestBlock", zap.Error(err))
	}

	return txs, nil
}

func (is *InternalService) rewriteLastHandledBlock(blockHeight int) error {
	return os.WriteFile(filePathLatestBlock, []byte(strconv.Itoa(blockHeight)), os.ModePerm)
}

func (is *InternalService) sendQuery(url string, data interface{}, response interface{}) error {
	node := cast.ToString(is.Context.GetConfig("node_address", ""))

	requestBody, err := json.Marshal(data)
	if err != nil {
		logger.Logger.Error("sendQuery", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", node+url, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Logger.Error("sendQuery", zap.Error(err))
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Logger.Error("sendQuery", zap.Error(err))
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		logger.Logger.Error("sendQuery", zap.Error(err))
		return err
	}

	return nil
}
