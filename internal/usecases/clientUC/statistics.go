package clientUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"time"
)

func (uc *ClientUsecase) GetUsageStatistics(
	ctx context.Context,
	clientID string,
) (*responsedto.ClientUsageStatistic, *ucerror.UCError) {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfNextMonth := time.Date(firstDayOfMonth.Year(), firstDayOfMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-24 * time.Hour)
	totalDaysInMonth := lastDayOfCurrentMonth.Day()
	totalDaysBeforeNextMonth := totalDaysInMonth - now.Day()

	thirtyDaysStatistic, err := uc.influxdb.GetTotalRequestsInLastXDaysByClient(ctx, 30, clientID)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to get last thirty days usage statistic: %w", err), consts.ErrInternalServer)
	}

	monthToDateTotal, err := uc.influxdb.GetTotalByDateRangeByClient(ctx, firstDayOfMonth, firstDayOfNextMonth, clientID)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to get month to date usage statistic: %w", err), consts.ErrInternalServer)
	}

	finalResponse := responsedto.ClientUsageStatistic{
		MonthToDate: models.MonthToDateStatistic{
			Month: firstDayOfMonth.Month().String(),
		},
		LastThirtyDaysStatistic: *thirtyDaysStatistic,
		MetricType:              "data_transferred",
		UnitOfMeasurement:       "GB",
	}

	if monthToDateTotal > 0 {
		finalResponse.MonthToDate.MonthToDateUsage = monthToDateTotal / 1000000000
		finalResponse.MonthToDate.ForecastedEndOfMonthUsage = (monthToDateTotal / 1000000000) + float64(totalDaysBeforeNextMonth)*thirtyDaysStatistic.Average
	}

	return &finalResponse, nil
}
