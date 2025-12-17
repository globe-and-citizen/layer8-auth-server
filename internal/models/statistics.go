package models

type MonthToDateStatistic struct {
	Month                     string  `json:"month"`
	MonthToDateUsage          float64 `json:"month_to_date_usage"`
	ForecastedEndOfMonthUsage float64 `json:"forecasted_end_of_month_usage"`
}

type Statistics struct {
	Total            float64                 `json:"total"`
	Average          float64                 `json:"average"`
	StatisticDetails []UsageStatisticPerDate `json:"details"`
}

type UsageStatisticPerDate struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type ClientUsageStatisticsByRange struct {
	ClientId   string
	TotalBytes float64
}
