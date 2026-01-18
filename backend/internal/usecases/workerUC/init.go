package workerUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"time"
)

type IWorkerUsecase interface {
	UpdateUsageStatistics(ratePerByte float64, now time.Time) error
}

type WorkerUsecase struct {
	postgres postgresRepo.IPostgresRepository
	influxdb influxdbRepo.IInfluxdbRepository
}

func NewWorkerUsecase(postgres postgresRepo.IPostgresRepository, influxdb influxdbRepo.IInfluxdbRepository) IWorkerUsecase {
	return &WorkerUsecase{
		postgres: postgres,
		influxdb: influxdb,
	}
}
