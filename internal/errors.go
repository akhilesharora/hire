package internal

import (
	"encoding/json"
	"log"
	"strconv"
)

const (
	ErrContentLimitExceeded ="CONTENT_LIMIT_EXCEEDED"
	InvalidFormat = "INVALID_FORMAT"
)

type Code uint32

const (
	OK                 Code = 0
	InvalidOriginator       = 1
	InvalidMessageBody		= 2
	_maxCode                = 3
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case InvalidOriginator:
		return "CONTENT_LIMIT_EXCEEDED_OR_INVALID_FORMAT"
	case InvalidMessageBody:
		return "CONTENT_LIMIT_EXCEEDED"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

type CustomError struct {
	Code Code
	Msg string
	Description  string
	Parameter string
}

func (e *CustomError) Error() string {
	b, err := json.Marshal(e)
	if err!= nil{
		log.Println(err)
		return ""
	}
	return string(b)
}