package types

// Event represents a webhook event.
type Event struct {
	Origin      string  `json:"origin"`
	Timestamp   float64 `json:"timestamp"`
	EventType   string  `json:"event_type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Result      string  `json:"result,omitempty"` // Some events may have a "result"
	Files       []File  `json:"files,omitempty"`  // Optional files attribute
	SourceIP    string  `json:"source_ip"`        // IP of the client
	Status      string  `json:"status,omitempty"`
	Progress    int     `json:"progress,omitempty"`
	Message     string  `json:"message,omitempty"`
}

// File represents an optional file entry in the event payload.
type File struct {
	Content  string `json:"content"`  // Base64 encoded content
	Path     string `json:"path"`     // File path
	Encoding string `json:"encoding"` // Encoding format (e.g., base64)
}

// DBConfig holds database connection configurations.
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Migration holds migration information.
type Migration struct {
	ID   int
	Name string
	SQL  string
}

// LoggerConfig represents logger configuration options.
type LoggerConfig struct {
	Level  string
	Format string
}

// HTTPServerConfig holds server configuration details.
type HTTPServerConfig struct {
	Address string
	Port    int
}

// User represents a user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
