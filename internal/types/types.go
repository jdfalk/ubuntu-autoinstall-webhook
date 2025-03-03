package types

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
