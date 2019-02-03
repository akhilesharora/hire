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
	s.cl = NewDefaultClient(s.conf.AccessKey)
	return s
}

// Example: {"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}
type Message struct {
	Recipient int64 `json:"recipient"`
	Originator string `json:"originator"`
	Message string `json:"message,omitempty"`
}

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

	var msg Message
	err = json.Unmarshal(data, &msg)
	if err!= nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithError(err).Error("Failed to parse the message")
		return
	}

	if msg.Originator == "" {
		msg.Originator = s.conf.Originator
	}

	err = isValidMessage(&msg)
	if err!=nil{
		w.Write([]byte(err.Error()))
		log.WithError(err).Error("Failed to parse the message")
		return
	}

	s.q <- &Messages{
		Recipients: []string{strconv.FormatInt(msg.Recipient,10)},
		Originator: msg.Originator,
		Message: msg.Message,
	}
	w.WriteHeader(http.StatusOK)
}

func isValidMessage(msg *Message) error {
	if !orRgx.MatchString(msg.Originator) {
		return &CustomError{Code:InvalidOriginator, Msg: ErrContentLimitExceeded, Description:"Content limit exceeded or Invalid format", Parameter: "originator"}
	}
	if  len(msg.Message)>MaxContentLength {
		return &CustomError{Code:InvalidMessageBody, Msg: ErrContentLimitExceeded, Description:"Content limit exceeded", Parameter: "message"}
	}
	return nil
}

func (s *Server) MessagebirdWorker(q <-chan *Messages) {
	log.Println("Initializing worker")
	tick := time.Tick(1 * time.Second)
	// Forever
	for {
		select {
		case <-tick:
			msg := <-q
			err := s.cl.sendSms(msg)
			if err != nil {
				fmt.Println("Error:", err)
				log.WithError(err).WithFields(log.Fields{"Message": msg})
			}
		}
	}
}

type Messages struct {
	Recipients []string
	Originator string
	Message string
}