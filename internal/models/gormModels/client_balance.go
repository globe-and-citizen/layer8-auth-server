package gormModels

import (
	"math/big"
	"time"
)

type AccountStatus string

const (
	// AccountOwing - when balance is positive
	AccountOwing AccountStatus = "owing"
	// AccountZeroed - when balance is equal to 0
	AccountZeroed AccountStatus = "zeroed"
	// AccountOverpaid - when balance is negative
	AccountOverpaid AccountStatus = "overpaid"
)

func (AccountStatus) GetStatus(balance *big.Int) AccountStatus {
	switch balance.Cmp(big.NewInt(0)) {
	case -1:
		return AccountOverpaid
	case 1:
		return AccountOwing
	default:
		return AccountZeroed
	}
}

type ClientBalance struct {
	ID                 uint          `gorm:"primaryKey; autoIncrement; not null"`
	ClientID           string        `gorm:"column:client_id; unique; not null"`
	BalanceWei         string        `gorm:"column:balance_wei" type:"numeric(78,0)"`
	Status             AccountStatus `gorm:"column:status" default:"zeroed"`
	LastUsageUpdatedAt time.Time     `gorm:"column:last_usage_updated_at"`
}

func (ClientBalance) TableName() string {
	return "client_balance"
}

type ClientPaymentReceipt struct {
	ID            uint      `gorm:"primaryKey; autoIncrement; not null"`
	ClientID      string    `gorm:"column:client_id; not null"`
	PaidAmountWei string    `gorm:"column:paid_amount_wei" type:"numeric(78,0)"`
	PaidAt        time.Time `gorm:"column:paid_at"`
	TxID          string    `gorm:"column:tx_id"`
}

func (ClientPaymentReceipt) TableName() string {
	return "client_payment_receipt"
}
