package hive

import "encoding/json"

// MessageLog use to format json for Logchan.Message
type MessageLog struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Mode    string `json:"mode"`
}

// NewMessageLog creates a new instance of MessageLog
func NewMessageLog(w *Worker, c string, isRead bool) *MessageLog {
	m := &MessageLog{
		ID:      w.ID,
		Content: c,
		Mode:    "w",
	}

	if isRead {
		m.Mode = "r"
	}

	return m
}

// JSON converts *MessageLog to json
func (m *MessageLog) JSON() []byte {
	b, _ := json.Marshal(m)
	return b
}
