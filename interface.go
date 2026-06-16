package main

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RouteConfig struct {
	To        string     `json:"to"`
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`
}
