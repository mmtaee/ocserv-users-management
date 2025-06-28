package models

type IPBan struct {
	IP       string `json:"IP"`
	Since    string `json:"Since"`
	SinceAlt string `json:"_Since"` // maps to "_Since" in JSON
	Score    int    `json:"Score"`
}

type Iroute struct {
	ID       int      `json:"ID"`
	Username string   `json:"Username"`
	VHost    string   `json:"vhost"`
	Device   string   `json:"Device"`
	IP       string   `json:"IP"`
	IRoutes  []string `json:"iRoutes"`
}

type OnlineUserSession struct {
	Username    string `json:"Username"`
	Group       string `json:"Groupname"`
	AverageRX   string `json:"Average RX"`
	AverageTX   string `json:"Average TX"`
	ConnectedAt string `json:"_Connected at"`
}
