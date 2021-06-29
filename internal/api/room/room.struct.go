package room

type chat struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
	Sender      string `json:"sender"`
	Timestamp   int64  `json:"timestamp"`
}

type notice struct {
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
	Timestamp   int64  `json:"timestamp"`
}

type exception struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type room struct {
	RoomId      string `json:"roomId"`
	Host        string `json:"host"`
	Ongoing     bool   `json:"ongoing"`
	MemberCount int32  `json:"memberCount"`
}

type message struct {
	Type    string `json:"type"`    // chat | command
	Message string `json:"message"` // string | (command: start)
}
