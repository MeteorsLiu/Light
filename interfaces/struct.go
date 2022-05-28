package interfaces

type UploadPayload struct {
	IP    string
	Rates string
}

type Config struct {
	RemoteControllerAddr string `json:"remoteAddr"`
	DeviceToken          string `json:"deviceToken"`
	ServerToken          string `json:"serverToken"`
	DeviceUUID           string `json:"deviceUUID"`
}
