package types

import "time"

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SuccessResponseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type NotificationSuccess struct {
	Status string `json:"status"` // всегда "success"
}

type CountDaysData struct {
	CompanyId       string    `json:"company_id"`
	SubscriptionEnd time.Time `json:"subscription_end"`
	DaysLeft        int       `json:"days_left"`
	Expired         bool      `json:"expired"`
}
