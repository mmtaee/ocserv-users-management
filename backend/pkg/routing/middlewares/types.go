package middlewares

type Unauthorized struct {
	Error string `json:"error"`
}

type PermissionDenied struct {
	Error string `json:"error"`
}
