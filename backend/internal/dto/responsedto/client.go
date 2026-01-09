package responsedto

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"globe-and-citizen/layer8/auth-server/backend/pkg/scram"
)

type ClientRegisterPrecheck struct {
	scram.ServerRegisterFirstMessage `json:",inline"`
}

type ClientLoginPrecheck struct {
	scram.ServerLoginFirstMessage `json:",inline"`
}

type ClientLogin struct {
	scram.ServerLoginFinalMessage `json:",inline"`
	Token                         string `json:"token"`
}

type ClientProfile struct {
	ID              string `json:"id"`
	Secret          string `json:"secret"`
	Name            string `json:"name"`
	RedirectURI     string `json:"redirect_uri"`
	BackendURI      string `json:"backend_uri"`
	NTorCertificate string `json:"ntor_certificate"`
}

type ClientUsageStatistic struct {
	MetricType              string                      `json:"metric_type"`
	UnitOfMeasurement       string                      `json:"unit_of_measurement"`
	MonthToDate             models.MonthToDateStatistic `json:"month_to_date"`
	LastThirtyDaysStatistic models.Statistics           `json:"last_thirty_days_statistic"`
}

type ClientGetUnpaidAmount struct {
	UnpaidAmount int `json:"unpaid_amount"`
}

type ClientGetNTorCertificate struct {
	ClientID    string `json:"client_id"`
	Certificate string `json:"cert"`
}
