package dto

import "encoding/json"

type EntityMessage struct {
	// uuid in string representation
	UserId string `json:"user_id"`
	// milliseconds since zero unix time (1970-01-01)
	Timestamp int64 `json:"timestamp"`
	// message payload
	MessageParams interface{} `json:"params"`
}

func (entityMessage EntityMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(entityMessage)
}
