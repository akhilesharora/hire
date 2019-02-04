package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/akhilesharora/hire/internal/config"
)

const MaxContentLength = 160
var orRgx = regexp.MustCompile("^[a-zA-Z0-9]{1,11}$")

type Server struct {
	conf *config.ServerConfig
	q chan *Messages
	cl SmsClient
}

func NewServer(conf *config.ServerConfig, q *chan *Messages) *Server {
	s:=  &Server{
		conf: conf,
		q: *q,
	}
	s.cl = newDefaultClient(s.conf.AccessKey)
	return s
}

type Message struct {
	Recipient int64 `json:"recipient"`
	Originator string `json:"originator"`
	Message string `json:"message,omitempty"`
}

// HTTP handler interface for incoming sms requests
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method!= "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err!= nil {
		log.WithError(err).Error("Failed to parse the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validation on content length
	if r.ContentLength > MaxContentLength {
		e := &CustomError{Code:InvalidMessageBody, Msg: ErrContentLimitExceeded, Description:"Content limit exceeded"}
		w.Write([]byte(e.Error()))
		log.WithError(e)
		return
	}

	var msg Message
	err = json.Unmarshal(data, &msg)
	if err!= nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithError(err).Error("Failed to parse the message")
		return
	}

	if !isValidMessage(&msg) {
		e := &CustomError{Code:InvalidOriginator, Msg: InvalidFormat, Description:"Invalid format or Content limit exceeded more than 11 chars", Parameter: "originator"}
		w.Write([]byte(e.Error()))
		log.WithError(e)
		return
	}

	// Push incoming sms messages to queue
	s.q <- &Messages{
		Recipients: []string{strconv.FormatInt(msg.Recipient,10)},
		Originator: msg.Originator,
		Message: msg.Message,
	}
	w.WriteHeader(http.StatusOK)
}

// Check if Message is valid
func isValidMessage(msg *Message) bool {
	if !orRgx.MatchString(msg.Originator) || len(msg.Originator) == 0 {
		return false
	}
	return true
}

// Worker for processing SMS requests
func (s *Server) MessagebirdWorker(q <-chan *Messages) {
	log.Println("Initializing worker")
	tick := time.Tick(1 * time.Second)
	// Forever
	for {
		select {
		// Rate limit request with 1 req/s
		case <-tick:
			msg := <-q
			err := s.cl.sendSms(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

type Messages struct {
	Recipients []string
	Originator string
	Message string
}

