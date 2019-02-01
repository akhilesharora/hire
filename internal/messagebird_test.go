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

type MessagebirdTestSuite struct {
	suite.Suite
	repo *Server
}

func TestMessagebirdTestSuite(t *testing.T) {
	mailerTestSuite := MessagebirdTestSuite{}
	q := make(chan *Messages, 100)
	mailerTestSuite.repo= NewServer(testutils.GetConfig(), &q)
	// Start SMS worker
	go func() {
		mailerTestSuite.repo.MessagebirdWorker(q)
	}()
	suite.Run(t, &mailerTestSuite)
	time.Sleep(1 * time.Second)
}

func (s *MessagebirdTestSuite) Test01SMSRequests() {
	cl := &http.Client{}
	payload := `{"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}`
	req, err := http.NewRequest("POST", s.repo.conf.ServerAddr, bytes.NewBuffer([]byte(payload)))
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
}

// Test for not allowing GET
func (s *MessagebirdTestSuite) Test02GETMethodNotAllowed() {
	cl := &http.Client{}
	payload := `{"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}`
	req, err := http.NewRequest("GET", s.repo.conf.ServerAddr, bytes.NewBuffer([]byte(payload)))
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
	s.Assert().Equal(r.StatusCode, http.StatusMethodNotAllowed)
}

// Test for content more than 160 chars
func (s *MessagebirdTestSuite) Test03ContentLengthExceeded() {
	cl := &http.Client{}
	payload := `{"recipient":31612345678,"originator":"MessageBird","message":"This is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test message.This is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test messageThis is a test message"}`
	req, err := http.NewRequest("GET", s.repo.conf.ServerAddr, bytes.NewBuffer([]byte(payload)))
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
	s.Assert().Equal(r.StatusCode, 405)
}

func (s *MessagebirdTestSuite) Test04SendSMS() {
	msg := &Messages{
		Recipients: []string{"31612345678"},
		Originator: s.repo.conf.Originator,
		Message: "Hello",
	}
	err := s.repo.sendSms(msg, s.repo.newDefaultClient())
	s.Assert().NoError(err)
}

//{"errors":[{"code":2,"description":"Request not allowed (incorrect access_key)","parameter":"access_key"}]}