package response

import "net/http"

type Status struct {
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}

var StatusOk = Status{
	StatusCode: http.StatusOK,
	StatusText: "Status: OK",
}
var StatusNotFound = Status{
	StatusCode: http.StatusNotFound,
	StatusText: "Status: Not Found",
}
var StatusInternalServerError = Status{
	StatusCode: http.StatusInternalServerError,
	StatusText: "Status: Internal Server error",
}
var StatusCreated = Status{
	StatusCode: http.StatusCreated,
	StatusText: "Status: Created",
}
var StatusBadRequest = Status{
	StatusCode: http.StatusBadRequest,
	StatusText: "Status: Bad Request",
}
var StatusUnauthorized = Status{
	StatusCode: http.StatusUnauthorized,
	StatusText: "Status: Unauthorized",
}
var StatusForbidden = Status{
	StatusCode: http.StatusForbidden,
	StatusText: "Status: Forbidden",
}
var StatusConflict = Status{
	StatusCode: http.StatusConflict,
	StatusText: "Status: Conflict",
}

func (r Status) Equals(other Status) bool {
	return r.StatusCode == other.StatusCode
}
