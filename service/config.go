package gomoney

// AppConfig ...
type AppConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Host string `json:"host"`
}
