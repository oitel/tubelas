package message

type Message struct {
	ID        uint64 `json:"id"`
	Timestamp int64  `json:"ts"`
	Text      string `json:"text"`
}
