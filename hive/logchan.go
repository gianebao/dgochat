package hive

// Logchan represents the log messages created by the hive
type Logchan struct {
	Fatal   chan string
	Info    chan string
	Warning chan string
	Error   chan string
	Message chan string
}

// NewLogchan creates a new Logchan instance
func NewLogchan() *Logchan {
	return &Logchan{
		Fatal:   make(chan string),
		Info:    make(chan string),
		Warning: make(chan string),
		Error:   make(chan string),
		Message: make(chan string),
	}
}
