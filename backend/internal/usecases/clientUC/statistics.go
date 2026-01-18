package clientUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"net/http"
	"time"
)

func (uc *ClientUsecase) GetUsageStatistics(clientID string) (responsedto.ClientUsageStatistic, int, string, error) {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfNextMonth := time.Date(firstDayOfMonth.Year(), firstDayOfMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-24 * time.Hour)
	totalDaysInMonth := lastDayOfCurrentMonth.Day()
	totalDaysBeforeNextMonth := totalDaysInMonth - now.Day()

	thirtyDaysStatistic, err := uc.influxdb.GetTotalRequestsInLastXDaysByClient(30, clientID)
	if err != nil {
		return responsedto.ClientUsageStatistic{}, http.StatusBadRequest, "Failed to get last thrthy days usage statistic", err
	}

	monthToDateTotal, err := uc.influxdb.GetTotalByDateRangeByClient(firstDayOfMonth, firstDayOfNextMonth, clientID)
	if err != nil {
		return responsedto.ClientUsageStatistic{}, http.StatusBadRequest, "Failed to get month to date usage statistic", err
	}

	finalResponse := responsedto.ClientUsageStatistic{
		MonthToDate: models.MonthToDateStatistic{
			Month: firstDayOfMonth.Month().String(),
		},
		LastThirtyDaysStatistic: thirtyDaysStatistic,
		MetricType:              "data_transferred",
		UnitOfMeasurement:       "GB",
	}

	if monthToDateTotal > 0 {
		finalResponse.MonthToDate.MonthToDateUsage = monthToDateTotal / 1000000000
		finalResponse.MonthToDate.ForecastedEndOfMonthUsage = (monthToDateTotal / 1000000000) + float64(totalDaysBeforeNextMonth)*thirtyDaysStatistic.Average
	}

	return finalResponse, http.StatusOK, "", nil
}
