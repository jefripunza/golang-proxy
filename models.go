package main

import (
	"time"
)

type ProxyRoute struct {
	ID                      uint      `gorm:"primaryKey" json:"id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	Domain                  string    `gorm:"uniqueIndex" json:"domain"`
	SchemaType              string    `json:"schema_type"`
	TargetURL               string    `json:"target_url"`
	DynamicResolveURL       string    `json:"dynamic_resolve_url"`
	UseBasicAuth            bool      `json:"use_basic_auth"`
	BasicAuthUsername       string    `json:"basic_auth_username"`
	BasicAuthPassword       string    `json:"basic_auth_password"`
	UseValidationMiddleware bool      `json:"use_validation_middleware"`
	ValidationMiddlewareURL string    `json:"validation_middleware_url"`
	SSLActive               bool      `json:"ssl_active"`
	LogPathPrefix           string    `json:"log_path_prefix"`
}

type ProxyLog struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Timestamp       time.Time `gorm:"index" json:"timestamp"`
	Domain          string    `gorm:"index" json:"domain"`
	Path            string    `json:"path"`
	Method          string    `json:"method"`
	StatusCode      int       `json:"status_code"`
	ResponseTimeMs  int64     `json:"response_time_ms"`
	SourceIP        string    `json:"source_ip"`
	ErrorMessage    string    `json:"error_message"`
	RequestHeaders  string    `json:"request_headers"`
	ResponseHeaders string    `json:"response_headers"`
	RequestBody     string    `json:"request_body"`
	ResponseBody    string    `json:"response_body"`
}

type BlockedIP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IPAddress string    `gorm:"uniqueIndex" json:"ip_address"`
	Reason    string    `json:"reason"`
}
