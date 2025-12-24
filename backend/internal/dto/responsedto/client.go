package responsedto

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
)

type ClientRegisterPrecheck struct {
	Salt           string `json:"salt"`
	IterationCount int    `json:"iterationCount"`
}

type ClientLoginPrecheck struct {
	Salt      string `json:"salt"`
	IterCount int    `json:"iter_count"`
	Nonce     string `json:"nonce"`
}

type ClientLogin struct {
	ServerSignature string `json:"server_signature"`
	Token           string `json:"token"`
}

type ClientProfile struct {
	ID              string `json:"id"`
	Secret          string `json:"secret"`
	Name            string `json:"name"`
	RedirectURI     string `json:"redirect_uri"`
	BackendURI      string `json:"backend_uri"`
	X509Certificate string `json:"x509_certificate"`
}

type ClientUsageStatistic struct {
	MetricType              string                      `json:"metric_type"`
	UnitOfMeasurement       string                      `json:"unit_of_measurement"`
	MonthToDate             models.MonthToDateStatistic `json:"month_to_date"`
	LastThirtyDaysStatistic models.Statistics           `json:"last_thirty_days_statistic"`
}

type ClientUnpaidAmount struct {
	UnpaidAmount int `json:"unpaid_amount"`
}
