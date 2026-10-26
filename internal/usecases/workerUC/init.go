package workerUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/repositories/ethRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"time"
)

type IWorkerUsecase interface {
	UpdateUsageBalance(now time.Time) error
	ListenToEthereumEvents()
}

type WorkerUsecase struct {
	config   config.WorkerConfig
	logger   log.ILogger
	ctx      context.Context
	postgres postgresRepo.IPostgresRepository
	influxdb influxdbRepo.IInfluxdbRepository
	ethereum ethRepo.IEthereumRepository
}

func NewWorkerUsecase(
	config config.WorkerConfig,
	logger log.ILogger,
	ctx context.Context,
	postgres postgresRepo.IPostgresRepository,
	influxdb influxdbRepo.IInfluxdbRepository,
	ethereum ethRepo.IEthereumRepository,
) IWorkerUsecase {
	return &WorkerUsecase{
		config:   config,
		logger:   logger,
		ctx:      ctx,
		postgres: postgres,
		influxdb: influxdb,
		ethereum: ethereum,
	}
}
