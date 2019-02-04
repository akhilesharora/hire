package internal

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/akhilesharora/hire/internal/testutils"
)

const urlScheme  = "http"

func (s *MessagebirdTestSuite) Test01DoSendSMSRequest() {
	cl := &http.Client{}
	payload := `{"recipient":31620286093,"originator":"MessageBird","message":"This is a test message."}`
	req, err := http.NewRequest("POST", s.restEndpoint, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := cl.Do(req)
	if err != nil {
		log.Println(err)
	}
	body := r.Body
	_, err = ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	defer body.Close()
	s.Assert().NoError(err)
	s.Assert().NotNil(r)
	s.Assert().Equal(http.StatusOK, r.StatusCode)
}

// Test for request method allowed
func (s *MessagebirdTestSuite) Test02GETMethodNotAllowed() {
	cl := &http.Client{}
	payload := `{"recipient":31620286093,"originator":"MessageBird","message":"This is a test message."}`
	req, err := http.NewRequest("GET", s.restEndpoint, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := cl.Do(req)
	if err != nil {
		log.Println(err)
	}
	body := r.Body
	_, err = ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	defer body.Close()
	s.Assert().Equal(http.StatusMethodNotAllowed,r.StatusCode)
}

// Test for content more than 160 chars
func (s *MessagebirdTestSuite) Test03ContentLengthExceeded() {
	cl := &http.Client{}
	payload := `{"recipient":31620286093,"originator":"MessageBird","message":"This is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test message.This is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test message"}`
	req, err := http.NewRequest("POST", s.restEndpoint, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := cl.Do(req)
	if err != nil {
		log.Println(err)
	}
	body := r.Body
	_, err = ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	defer body.Close()
	s.Assert().Equal(200,r.StatusCode)
}

func (s *MessagebirdTestSuite) Test04SendSMS() {
	msg := &Messages{
		Recipients: []string{"31620286093"},
		Originator: "MessageBird",
		Message: "This is a test message",
	}
	err := s.server.cl.sendSms(msg)
	s.Assert().NoError(err)
}

func (s *MessagebirdTestSuite) Test04InvalidOriginatorValidations() {
	cl := &http.Client{}
	payload := `{"recipient":31620286093,"originator":"MessageBird1234","message":"This is a test message"}`
	req, err := http.NewRequest("POST", s.restEndpoint, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := cl.Do(req)
	if err != nil {
		log.Println(err)
	}
	body := r.Body
	b , err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	defer body.Close()
	s.Assert().Equal(200,r.StatusCode)
	s.Assert().EqualError(&CustomError{Code:InvalidOriginator, Msg: InvalidFormat, Description:"Invalid format or Content limit exceeded more than 11 chars", Parameter: "originator"}, string(b))

}
func (s *MessagebirdTestSuite) Test04BadRequestInvalidMessageValidation() {
	cl := &http.Client{}
	payload := `{"recipient":"31620286093"","originator":"MessageBird1234","message":"This is a test message"}`
	req, err := http.NewRequest("POST", s.restEndpoint, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := cl.Do(req)
	if err != nil {
		log.Println(err)
	}
	body := r.Body
	_ , err = ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
	}
	defer body.Close()
	s.Assert().Equal(http.StatusBadRequest,r.StatusCode)
}

func TestMessagebirdTestSuite(t *testing.T) {
	m := MessagebirdTestSuite{}
	q := make(chan *Messages, 1000)
	m.server = NewServer(testutils.GetConfig(), &q)
	// Start an instance of test server
	http.HandleFunc("/", m.server.Handler )
	go testutils.StartTestServerInstance()
	// Start SMS worker
	go func(q <-chan *Messages) {
		m.server.MessagebirdWorker(q)
	}(q)
	m.restEndpoint = makeFQD(urlScheme, m.server.conf.ServerAddr)
	suite.Run(t, &m)
	time.Sleep(1 * time.Second)
}

type MessagebirdTestSuite struct {
	suite.Suite
	server *Server
	restEndpoint string
}
// Append Protocol with server address
func makeFQD(urlScheme, url string) string {
	return urlScheme+"://"+url
}

