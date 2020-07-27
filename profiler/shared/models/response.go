package models

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
)

// ResponseTemplate contain json format for lucent-api response
type ResponseTemplate struct {
	RC       string      `json:"rc"`
	Message  interface{} `json:"message,omitempty"`
	Messages []string    `json:"messages,omitempty"`
	Payload  interface{} `json:"payload,omitempty"`
	Sign     string      `json:"sign,omitempty"`
}

// for error response
type ResponseError struct {
	RC      string `json:"rc"`
	Message string `json:"message"`
}

// component related with error response
type ResponseCode struct {
	RC         string `json:"rc"`
	Message    string `json:"message"`
	HttpStatus int    `json:"http_status"`
}

func ReturnError(c echo.Context, errorCode string) error {
	errCode := GetErrorMessage(errorCode)
	return c.JSON(errCode.HttpStatus, ResponseError{RC: errCode.RC, Message: errCode.Message})
}

//getErrorMessage is collection of listed error, responseCode list can be seen in notion.
func GetErrorMessage(rc string) ResponseCode {
	var (
		response ResponseCode
	)

	switch rc {
	default:
		response = ResponseCode{
			RC:         rc,
			Message:    "bad request",
			HttpStatus: http.StatusBadRequest,
		}
	case "02":
		response = ResponseCode{
			RC:         rc,
			Message:    "maintenance",
			HttpStatus: http.StatusServiceUnavailable,
		}
	case "50A":
		response = ResponseCode{
			RC:         rc,
			Message:    "request timeout",
			HttpStatus: http.StatusRequestTimeout,
		}
	case "50B":
		response = ResponseCode{
			RC:         rc,
			Message:    "request limit exceed",
			HttpStatus: http.StatusTooManyRequests,
		}
	case "50C":
		response = ResponseCode{
			RC:         rc,
			Message:    "empty body",
			HttpStatus: http.StatusBadRequest,
		}
	case "50D":
		response = ResponseCode{
			RC:         rc,
			Message:    "something when wrong",
			HttpStatus: http.StatusInternalServerError,
		}
	case "50E":
		response = ResponseCode{
			RC:         rc,
			Message:    "stack error",
			HttpStatus: http.StatusInternalServerError,
		}
	case "51A":
		response = ResponseCode{
			RC:         rc,
			Message:    "missing authentication token",
			HttpStatus: http.StatusUnauthorized,
		}
	case "51B":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid credential",
			HttpStatus: http.StatusUnauthorized,
		}
	case "51C":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid signature",
			HttpStatus: http.StatusUnauthorized,
		}
	case "51D":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid or expired token",
			HttpStatus: http.StatusUnauthorized,
		}
	case "51E":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid roles",
			HttpStatus: http.StatusForbidden,
		}
	case "52A":
		response = ResponseCode{
			RC:         rc,
			Message:    "inactive data",
			HttpStatus: http.StatusBadRequest,
		}
	case "52B":
		response = ResponseCode{
			RC:         rc,
			Message:    "data not found or unavailable",
			HttpStatus: http.StatusBadRequest,
		}
	case "52C":
		response = ResponseCode{
			RC:         rc,
			Message:    "request validation mismatch",
			HttpStatus: http.StatusUnprocessableEntity,
		}
	case "52D":
		response = ResponseCode{
			RC:         rc,
			Message:    "duplicate data",
			HttpStatus: http.StatusBadRequest,
		}
	case "52E":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid or expired otp",
			HttpStatus: http.StatusBadRequest,
		}
	case "52F":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid amount",
			HttpStatus: http.StatusBadRequest,
		}
	case "52G":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid file type",
			HttpStatus: http.StatusBadRequest,
		}
	case "52H":
		response = ResponseCode{
			RC:         rc,
			Message:    "invalid csv column header",
			HttpStatus: http.StatusBadRequest,
		}
	case "52L":
		response = ResponseCode{
			RC:         rc,
			Message:    "file too large",
			HttpStatus: http.StatusBadRequest,
		}
	// InternalError
	case "53A":
		response = ResponseCode{
			RC:         rc,
			Message:    "fail update data",
			HttpStatus: http.StatusInternalServerError,
		}
	case "53S":
		response = ResponseCode{
			RC:         rc,
			Message:    "something went wrong",
			HttpStatus: http.StatusInternalServerError,
		}
	}

	return response
}

// CurrentTZTime
// Custom type use for converting date from Mongo to current timezone defined in config
// This type implement "MarshalJSON", so when marshalling to JSON will automatically convert to current TZ

type CurrentTZTime struct {
	time.Time
}

func NowTZTime() time.Time {
	var tz string
	tz = os.Getenv("LUCENT_TZ")
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	loc, _ := time.LoadLocation(tz)
	return time.Now().In(loc)
}

func (t CurrentTZTime) MarshalJSON() ([]byte, error) {
	var tz string
	tz = os.Getenv("LUCENT_TZ")
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	// current time ISO = RFC3339
	stamp := fmt.Sprintf("\"%s\"", t.Time.In(loc).Format("2006-01-02T15:04:05Z07:00"))
	return []byte(stamp), nil
}

// ResponseValidationError return response when validation error occured
func ResponseValidationError(c echo.Context, v interface{}, m map[string]string) error {
	messages := make([]string, 0, len(m))
	for _, value := range m {
		messages = append(messages, value)
	}
	sort.Strings(messages)

	response := ResponseTemplate{
		RC:       "52C",
		Message:  "request validation mismatch",
		Messages: messages,
	}

	return c.JSON(http.StatusUnprocessableEntity, response)

	// return echo.NewHTTPError(http.StatusUnprocessableEntity, response)

}
