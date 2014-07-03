package main

import (
	"errors"
	"github.com/soh335/go-twitterstream"
)

type Agent struct {
	twitterStreamClient *twitterstream.Client
	name                string
	params              map[string]string
	stopChan            <-chan struct{}
	byteChan            chan<- []byte
}

func NewAgent(consumerKey, consumerSecret, token, tokenSecret string, stopChan <-chan struct{}, byteChan chan<- []byte, name string, params map[string]string) (*Agent, error) {
	agent := &Agent{}

	{
		agent.twitterStreamClient = &twitterstream.Client{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
			Token:          token,
			TokenSecret:    tokenSecret,
		}
	}

	{
		agent.stopChan = stopChan
		agent.byteChan = byteChan
		agent.name = name
		agent.params = params
	}

	return agent, nil
}

func (a *Agent) Work() error {
	Info("agent ", a.name, " start")
	//TODO: reconnecting to stream
	err := a.Start()
	if err != nil {
		Warn("agent ", a.name, " receive error ", err)
	}
	return err
}

func (a *Agent) Start() error {

	var conn *twitterstream.Connection
	var err error

	switch a.name {
	case "Userstream":
		// https://dev.twitter.com/docs/api/1.1/get/user
		conn, err = a.twitterStreamClient.Userstream("POST", a.params)
	case "Filter":
		// https://dev.twitter.com/docs/api/1.1/post/statuses/filter
		conn, err = a.twitterStreamClient.Filter("GET", a.params)
	default:
		return errors.New("not support name:" + a.name)
	}
	if err != nil {
		return err
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	errorChan := make(chan error, 1)

	go func() {
		for {
			line, err := conn.Next()
			if err != nil {
				errorChan <- err
				return
			}
			a.byteChan <- line
		}
	}()

	go func() {
		for {
			select {
			case <-a.stopChan:
				// conn.Next() return err
				conn.Close()
			case err := <-errorChan:
				doneChan <- err
				return
			}
		}
	}()

	return <-doneChan
}
