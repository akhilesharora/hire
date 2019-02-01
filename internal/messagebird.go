package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	sdk "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
	"github.com/messagebird/internal/config"
	log "github.com/sirupsen/logrus"
)

const MaxContentLength = 160

type Server struct {
	conf *config.ServerConfig
	q chan *Messages
}

func NewServer(conf *config.ServerConfig, q *chan *Messages) *Server {
	return &Server{
		conf: conf,
		q: *q,
	}
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

	if r.ContentLength > MaxContentLength {
		log.WithError(err).Error(http.ErrContentLength)
		w.WriteHeader(405)
		w.Write([]byte(http.ErrContentLength.Error()))
		return
	}
	var msg Message
	err = json.Unmarshal(data, &msg)
	if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithError(err).Error("Failed to parse the message")
		return
	}

	if msg.Originator == "" {
		msg.Originator = s.conf.Originator
	}

	//@TODO- Validation of phone number
	s.q <- &Messages{
		Recipients: []string{strconv.FormatInt(msg.Recipient,10)},
		Originator: msg.Originator,
		Message: msg.Message,
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) newDefaultClient() *sdk.Client {
	return sdk.New(s.conf.AccessKey)
}
func (s *Server) MessagebirdWorker(q <-chan *Messages) {
	log.Println("Initializing worker")
	tick := time.Tick(1 * time.Second)

	// Forever
	for {
		select {
		case <-tick:
			fmt.Println("Here", time.Now().UTC())
			case msg:= <-q:
			client := s.newDefaultClient()
			if client == nil {
				log.Println("Messagebird client died")
				continue
			}
			err := sendSms(msg, client)
			if err!= nil{
				fmt.Println("Error:", err)
				log.WithError(err).WithFields(log.Fields{"Message": msg, "Client": client})
			}
			continue
		}
	}
}

type Messages struct {
	Recipients []string
	Originator string
	Message string
}

func sendSms(msg *Messages, client *sdk.Client) error {
	message, err := sms.Create(client, msg.Originator, msg.Recipients, msg.Message,nil)
	if err!= nil {
		return err
	}
	log.Println("Messagebird Message:", message)
	return nil
}