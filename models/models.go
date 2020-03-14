package models

import (
	"time"
)

type Model struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Subscription struct {
	Model
	ChatID      int64
	AddressArID int64
}

type AddressArCheckStatus string

const (
	AddressStatusNoWork AddressArCheckStatus = "nowork"
	AddressStatusWork                        = "work"
	AddressStatusInit                        = "init"
)

type AddressAr struct {
	Model
	CheckStatus    AddressArCheckStatus
	ServiceMessage string
	TakenAt        time.Time
	CheckedAt      time.Time
	Subscriptions  []Subscription `gorm:"foreignkey:AddressArID"`
}
