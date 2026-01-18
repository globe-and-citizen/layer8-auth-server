package influxdbRepo

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/config"
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type IInfluxdbRepository interface {
	IsConnected(ctx context.Context) error
	GetTotalRequestsInLastXDaysByClient(days int, clientID string) (models.Statistics, error)
	GetTotalByDateRangeByClient(start time.Time, end time.Time, clientID string) (float64, error)
	GetTotalUsageStatisticsByDateRangeForEachClient(start time.Time, end time.Time) ([]models.ClientUsageStatisticsByRange, error)
}

func NewInfluxdbRepository(conf config.InfluxDB2Config) IInfluxdbRepository {
	influxdb2Client := influxdb2.NewClient(conf.Url, conf.Token)

	return &InfluxdbRepository{
		config: conf,
		client: influxdb2Client,
	}
}

type InfluxdbRepository struct {
	config config.InfluxDB2Config
	client influxdb2.Client
}

// IsConnected returns nil if InfluxDB responds to a trivial query, else an error.
func (r *InfluxdbRepository) IsConnected(ctx context.Context) error {
	queryAPI := r.client.QueryAPI(r.config.Org)

	// lightweight flux query; uses the configured bucket and a short time range
	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -1m)
	|> limit(n:1)`, r.config.Bucket)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return err
	}
	// ensure we consume and check for query errors
	defer result.Close()

	for result.Next() {
		// got at least one record -> healthy
		return nil
	}
	if result.Err() != nil {
		return result.Err()
	}
	// no records but no error -> server responded; treat as healthy
	return nil
}

func (r *InfluxdbRepository) GetTotalRequestsInLastXDaysByClient(days int, clientID string) (models.Statistics, error) {
	result := make([]models.UsageStatisticPerDate, 0)

	queryAPI := r.client.QueryAPI(r.config.Org)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -%dd)
	|> filter(fn: (r) => r["_measurement"] == "total_byte_transferred")
	|> filter(fn: (r) => r["_field"] == "counter")
	|> filter(fn: (r) => r["client_id"] == "%s")
	|> group(columns: ["client_id"])
	|> aggregateWindow(every: 1d, fn: sum, createEmpty: true)
	|> yield(name: "sum")`, r.config.Bucket, days, clientID)

	rawDataFromInflux, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return models.Statistics{}, err
	}

	var totalRequest float64
	for rawDataFromInflux.Next() {
		rawDataPointer := rawDataFromInflux.Record()
		unparsedTotal := rawDataPointer.ValueByKey("_value")
		decimalValueTotal, err := strconv.ParseFloat(fmt.Sprint(unparsedTotal), 64)
		if err != nil {
			decimalValueTotal = 0
		}

		var totalForThisPeriod float64
		if decimalValueTotal > 0 {
			totalRequest += decimalValueTotal / 1000000000
			totalForThisPeriod = decimalValueTotal / 1000000000
		}

		at := rawDataPointer.ValueByKey("_time").(time.Time)
		result = append(result, models.UsageStatisticPerDate{
			Date:  at.Format("Mon, 02 Jan 2006"),
			Total: totalForThisPeriod,
		})
	}

	var averageRequest float64
	if totalRequest > 0 {
		averageRequest = totalRequest / float64(len(result))
	}

	return models.Statistics{
		Total:   totalRequest,
		Average: averageRequest,
		Details: result,
	}, nil
}

func (r *InfluxdbRepository) GetTotalByDateRangeByClient(start time.Time, end time.Time, clientID string) (float64, error) {
	queryAPI := r.client.QueryAPI(r.config.Org)

	query := fmt.Sprintf(`
	from(bucket: "%s")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "total_byte_transferred")
	|> filter(fn: (r) => r["client_id"] == "%s")
	|> filter(fn: (r) => r["_field"] == "counter")
	|> group(columns: ["client_id"])
	|> sum()`, r.config.Bucket, start.Format(time.RFC3339), end.Format(time.RFC3339), clientID)

	rawDataFromInflux, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 0, err
	}

	// TODO: assert that rawDataFromInflux contains only one record
	var decimalValueTotal float64
	for rawDataFromInflux.Next() {
		rawDataPointer := rawDataFromInflux.Record()
		unparsedTotal := rawDataPointer.ValueByKey("_value")
		decimalValueTotal, err = strconv.ParseFloat(fmt.Sprint(unparsedTotal), 64)
		if err != nil {
			decimalValueTotal = 0
		}
	}

	return decimalValueTotal, err
}

func (r *InfluxdbRepository) GetTotalUsageStatisticsByDateRangeForEachClient(start time.Time, end time.Time) ([]models.ClientUsageStatisticsByRange, error) {
	queryAPI := r.client.QueryAPI(r.config.Org)

	query := fmt.Sprintf(`
	from(bucket: "%s")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r["_measurement"] == "total_byte_transferred")
	|> filter(fn: (r) => r["_field"] == "counter")
	|> group(columns: ["client_id"])
	|> sum()`, r.config.Bucket, start.Format(time.RFC3339), end.Format(time.RFC3339))

	queryResult, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	response := make([]models.ClientUsageStatisticsByRange, 0)
	for queryResult.Next() {
		record := queryResult.Record()

		clientId := fmt.Sprint(record.ValueByKey("client_id"))
		rawTotalBytes := record.ValueByKey("_value")

		totalBytes, err := strconv.ParseFloat(fmt.Sprint(rawTotalBytes), 64)
		if err != nil {
			return nil, err
		}

		if totalBytes == 0 {
			continue
		}

		response = append(response, models.ClientUsageStatisticsByRange{
			ClientId:   clientId,
			TotalBytes: totalBytes,
		})
	}

	return response, nil
}
