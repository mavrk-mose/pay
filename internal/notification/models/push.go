package models

import (
	"github.com/lib/pq"
	"time"
)

type Platform string

const (
	PlatformIOS     = "ios"
	PlatformAndroid = "android"
)

type Push struct {
	ID        int64          `db:"id" json:"-"`
	DeviceID  string         `db:"device_id" json:"device_id"`
	Platform  Platform       `db:"platform" json:"platform"`
	PushToken string         `db:"push_token" json:"push_token"`
	Addresses pq.StringArray `db:"addresses" json:"addresses"`
	CreatedAt time.Time      `db:"created_at" json:"-"`
	UpdatedAt time.Time      `db:"updated_at" json:"-"`
}

type PushItem struct {
	Address   string
	UserToken string
	Amount    float64
	Direction Direction
	Platform  Platform
}
