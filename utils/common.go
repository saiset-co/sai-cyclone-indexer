package utils

import (
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/saiset-co/sai-cyclone-indexer/logger"
)

func SaiQuerySender(body io.Reader, address, token string) ([]byte, error) {
	const failedResponseStatus = "NOK"

	type responseWrapper struct {
		Status string              `json:"Status"`
		Error  string              `json:"Error"`
		Result jsoniter.RawMessage `json:"result"`
		Count  int                 `json:"count"`
	}

	req, err := http.NewRequest(http.MethodPost, address, body)
	if err != nil {
		logger.Logger.Error("SaiQuerySender", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Logger.Error("SaiQuerySender", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Logger.Error("SaiQuerySender", zap.Error(err))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logger.Logger.Error("SaiQuerySender", zap.Error(err))
		return nil, fmt.Errorf("%s", resBytes)
	}

	result := responseWrapper{}
	err = jsoniter.Unmarshal(resBytes, &result)
	if err != nil {
		logger.Logger.Error("SaiQuerySender", zap.Error(err))
		return nil, err
	}

	if result.Status == failedResponseStatus {
		err = fmt.Errorf(result.Error)
		logger.Logger.Error("SaiQuerySender", zap.Error(err))
		return nil, err
	}

	return resBytes, nil
}
