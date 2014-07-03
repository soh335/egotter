package main

import (
	"encoding/json"
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"regexp"
)

var NotChainError = errors.New("not chain error")

type Matchers struct {
	bytChan  chan []byte
	matchers []Matcher
	stopChan <-chan struct{}
	handle   func(msg string)
}

func NewMatchers(config *Config, stopChan <-chan struct{}, handle func(msg string)) (*Matchers, error) {

	anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.Twitter.Token, config.Twitter.TokenSecret)
	user, err := api.GetSelf(nil)
	if err != nil {
		return nil, err
	}

	screenName := user.ScreenName

	keywords := make([]*regexp.Regexp, len(config.Keywords))
	for i, keyword := range config.Keywords {
		keywords[i] = regexp.MustCompile(keyword)
	}

	//TODO: configuration for selecting, ordering matcher
	matchers := []Matcher{
		&FilterScreenNameMatch{
			screenNames: config.FilterScreenNames,
		},
		NewDuplicateMatch(30),
		&EventMatch{
			events:     config.Events,
			screenName: screenName,
		},
		&ReTweetMatch{
			screenName: screenName,
		},
		&MentionMatch{
			screenName: screenName,
		},
		&KeywordMatch{
			keywords:   keywords,
			screenName: screenName,
		},
	}

	r := &Matchers{
		bytChan:  make(chan []byte, 10),
		matchers: matchers,
		stopChan: stopChan,
		handle:   handle,
	}

	return r, nil
}

func (r *Matchers) Work() error {

	for {
		select {
		case byt := <-r.bytChan:
			var data map[string]interface{}
			if err := json.Unmarshal(byt, &data); err != nil {
				return err
			}
			for _, matcher := range r.matchers {
				msg, err := matcher.Match(data)

				if err == NotChainError {
					break
				} else if err != nil {
					return err
				}

				if msg != "" {
					go r.handle(msg)
					break
				}
			}
		case <-r.stopChan:
			return errors.New("receive stop notification")
		}
	}
}
