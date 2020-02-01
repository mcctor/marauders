package users

import "github.com/mcctor/marauders/http"

const (
	ResultsPerPage = 10
	FirstPage      = 1
)

var (
	Href = http.ServerAddr + "/v1/users/"
)
