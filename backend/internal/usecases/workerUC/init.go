package workerUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/ethRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"time"
)

type IWorkerUsecase interface {
	UpdateUsageStatistics(ratePerByte float64, now time.Time) error
	ListenToEthereumEvents()
}

type WorkerUsecase struct {
	ctx      context.Context
	postgres postgresRepo.IPostgresRepository
	influxdb influxdbRepo.IInfluxdbRepository
	ethereum ethRepo.IEthereumRepository
}

func NewWorkerUsecase(
	ctx context.Context,
	postgres postgresRepo.IPostgresRepository,
	influxdb influxdbRepo.IInfluxdbRepository,
	ethereum ethRepo.IEthereumRepository,
) IWorkerUsecase {
	return &WorkerUsecase{
		ctx:      ctx,
		postgres: postgres,
		influxdb: influxdb,
		ethereum: ethereum,
	}
}
