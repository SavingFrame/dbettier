package messages

// LogLevel defines the severity of a log entry
type LogLevel int

const (
	LogInfo LogLevel = iota
	LogSuccess
	LogWarning
	LogError
	LogSQL
)

// AddLogMsg is a message to add a log entry to the log panel
type AddLogMsg struct {
	Message string
	Level   LogLevel
}
