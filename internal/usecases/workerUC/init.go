package workerUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/repositories/ethRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"math/big"
	"time"
)

type IWorkerUsecase interface {
	UpdateUsageBalance(ratePerByte *big.Int, now time.Time) error
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
