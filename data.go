package main

import (
	"github.com/mattn/go-scan"
)

//TODO: refactoring

type Tweet struct {
	data interface{}
}

func (t *Tweet) IsTweet() bool {
	//TODO: smart
	_, hasText := t.data.(map[string]interface{})["text"]
	_, hasUser := t.data.(map[string]interface{})["user"]
	return hasText && hasUser
}

func (t *Tweet) UserScreenName() (string, error) {
	var screenName string
	if err := scan.ScanTree(t.data, "/user/screen_name", &screenName); err != nil {
		return "", err
	}
	return screenName, nil
}

func (t *Tweet) Text() (string, error) {
	var text string
	if err := scan.ScanTree(t.data, "/text", &text); err != nil {
		return "", err
	}
	return text, nil
}

type User struct {
	data interface{}
}

func (u *User) ScreenName() (string, error) {
	var screenName string
	if err := scan.ScanTree(u.data, "/screen_name", &screenName); err != nil {
		return "", err
	}
	return screenName, nil
}

type List struct {
	data interface{}
}

func (l *List) Name() (string, error) {
	var name string
	if err := scan.ScanTree(l.data, "/name", &name); err != nil {
		return "", err
	}
	return name, nil
}
