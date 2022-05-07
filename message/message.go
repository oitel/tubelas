package message

type Message struct {
	ID        uint64 `json:"id" db:"id"`
	Timestamp int64  `json:"ts" db:"ts"`
	Text      string `json:"text" db:"text"`
}
