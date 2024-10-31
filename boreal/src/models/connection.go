package models

import (
	"net/http"
)

type Connection struct {
	Request  *http.Request
	Response http.ResponseWriter
}
