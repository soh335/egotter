package main

import (
	"errors"
	"fmt"
	"github.com/mattn/go-scan"
	"regexp"
	"sync"
)

type Matcher interface {
	Match(data map[string]interface{}) (string, error)
}

type FilterScreenNameMatch struct {
	screenNames []string
}

//TODO: filter event from screennames
func (s *FilterScreenNameMatch) Match(data map[string]interface{}) (string, error) {
	tweet := &Tweet{data}
	if !tweet.IsTweet() {
		return "", nil
	}

	tweetScreenName, err := tweet.UserScreenName()
	if err != nil {
		return "", err
	}

	for _, screenName := range s.screenNames {
		if screenName == tweetScreenName {
			return "", NotChainError
		}
	}

	return "", nil
}

type DuplicateMatch struct {
	m    sync.Mutex
	list []string
}

func NewDuplicateMatch(size int) *DuplicateMatch {
	d := &DuplicateMatch{
		list: make([]string, size),
	}
	return d
}

func (d *DuplicateMatch) Match(data map[string]interface{}) (string, error) {
	id, ok := data["id_str"].(string)
	if !ok {
		return "", nil
	}

	d.m.Lock()
	defer d.m.Unlock()

	for _, item := range d.list {
		if item == id {
			Info("duplicated ", id)
			return "", NotChainError
		}
	}

	copy(d.list[1:], d.list[0:cap(d.list)-1])
	d.list[0] = id

	return "", nil
}

type ReTweetMatch struct {
	screenName string
}

func (r *ReTweetMatch) Match(data map[string]interface{}) (string, error) {
	tweet := &Tweet{data}
	if !tweet.IsTweet() {
		return "", nil
	}
	retweetedStatus, hasRetweetedStatus := data["retweeted_status"]
	if !hasRetweetedStatus {
		return "", nil
	}
	reTweet := &Tweet{retweetedStatus}
	reTweetScreenName, err := reTweet.UserScreenName()
	if err != nil {
		return "", err
	}
	if reTweetScreenName != r.screenName {
		return "", nil
	}
	text, err := tweet.Text()
	if err != nil {
		return "", err
	}
	screenName, err := tweet.UserScreenName()
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("retweet: user:%s tweet:%s", screenName, text)
	return message, nil
}

type MentionMatch struct {
	screenName string
}

func (m *MentionMatch) Match(data map[string]interface{}) (string, error) {
	tweet := &Tweet{data}
	if !tweet.IsTweet() {
		return "", nil
	}
	var userMentions []map[string]interface{}
	if err := scan.ScanTree(data, "/entities/user_mentions", &userMentions); err != nil {
		return "", err
	}

	for _, userMention := range userMentions {
		if userMention["screen_name"].(string) == m.screenName {
			text, err := tweet.Text()
			if err != nil {
				return "", err
			}
			screenName, err := tweet.UserScreenName()
			if err != nil {
				return "", err
			}
			message := fmt.Sprintf("mention: user:%s tweet:%s", screenName, text)
			return message, nil
		}
	}

	return "", nil
}

type KeywordMatch struct {
	keywords   []*regexp.Regexp
	screenName string
}

func (k *KeywordMatch) Match(data map[string]interface{}) (string, error) {
	tweet := &Tweet{data}
	if !tweet.IsTweet() {
		return "", nil
	}
	for _, keyword := range k.keywords {
		text, err := tweet.Text()
		if err != nil {
			return "", err
		}
		if keyword.MatchString(text) {
			screenName, err := tweet.UserScreenName()
			if err != nil {
				return "", err
			}
			if screenName == k.screenName {
				return "", err
			}
			message := fmt.Sprintf("keyword: user:%s tweet:%s", screenName, text)
			return message, nil
		}
	}

	return "", nil
}

type EventMatch struct {
	events     []string
	screenName string
}

func (e *EventMatch) Match(data map[string]interface{}) (string, error) {
	event, hasEvent := data["event"]
	if !hasEvent {
		return "", nil
	}
	event = event.(string)

	subscribeEvent := false
	for _, _event := range e.events {
		if event == _event {
			subscribeEvent = true
			break
		}
	}

	if !subscribeEvent {
		return "", nil
	}

	switch event {
	case "favorite", "unfavorite":
		user := &User{data["source"]}
		sourceScreenName, err := user.ScreenName()
		if err != nil {
			return "", err
		}

		if sourceScreenName == e.screenName {
			return "", nil
		}

		tweet := &Tweet{data["target_object"]}
		text, err := tweet.Text()
		if err != nil {
			return "", err
		}
		message := fmt.Sprintf("user:%s %s tweet:%s", sourceScreenName, event, text)
		return message, nil
	case "follow", "unfollow", "block", "unblock":
		sourceUser := &User{data["source"]}
		sourceScreenName, err := sourceUser.ScreenName()
		if err != nil {
			return "", err
		}

		if sourceScreenName == e.screenName {
			return "", nil
		}

		targetUser := &User{data["target"]}
		targetScreenName, err := targetUser.ScreenName()
		if err != nil {
			return "", err
		}
		message := fmt.Sprintf("user:%s %s user:%s", sourceScreenName, event, targetScreenName)
		return message, nil
	case "list_member_added", "list_member_removed", "list_user_subscribed", "list_user_unsubscribed":
		sourceUser := &User{data["source"]}
		sourceScreenName, err := sourceUser.ScreenName()
		if err != nil {
			return "", err
		}

		if sourceScreenName == e.screenName {
			return "", nil
		}

		targetUser := &User{data["target"]}
		targetScreenName, err := targetUser.ScreenName()
		if err != nil {
			return "", err
		}
		list := &List{data["target_object"]}
		listName, err := list.Name()
		if err != nil {
			return "", err
		}
		message := fmt.Sprintf("user:%s %s list:%s of user:%s", sourceScreenName, event, listName, targetScreenName)
		return message, nil
	case "list_created", "list_destroyed", "list_updated":
		return "", nil
	case "user_update":
		return "", nil
	case "access_revoked":
		return "", nil
	default:
		return "", errors.New("not support event " + event.(string))
	}
	panic("not reach")
}
