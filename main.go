package main

import (
	"github.com/saiset-co/sai-cyclone-indexer/internal"
	"github.com/saiset-co/sai-cyclone-indexer/logger"
	saiService "github.com/saiset-co/sai-service/service"
)

func main() {
	svc := saiService.NewService("saiCycloneIndexer")
	is := internal.InternalService{Context: svc.Context}

	svc.RegisterConfig("config.yml")

	logger.Logger = svc.Logger

	is.Init()

	svc.RegisterTasks([]func(){
		is.Process,
	})

	svc.RegisterHandlers(
		is.NewHandler(),
	)

	svc.Start()
}
