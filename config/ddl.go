package config

type appJsonc struct {
	Netdev netdev `json:"netdev"`
	Smtp   smtp   `json:"smtp"`
}

type netdev struct {
	Period   int64  `json:"period"`
	Receiver string `json:"receiver"`
}

type smtp struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`

	SslOnConnect bool `json:"sslOnConnect"`
}
