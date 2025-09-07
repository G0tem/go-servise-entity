package dto

import (
	"encoding/json"
)

type NotifyMessage struct {
	// uuid in string representation
	UserId string `json:"user_id"`
	// milliseconds since zero unix time (1970-01-01)
	NotifyTimestamp int64 `json:"notify_timestamp"`
	// message payload
	MessageParams interface{} `json:"params"`
	// uuid in string representation
	TargetUserId string `json:"target_user_id"`
}

func (notifyMessage NotifyMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(notifyMessage)
}
