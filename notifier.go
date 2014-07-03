package main

import (
	"github.com/soh335/go-imkayaccom"
)

type ImKayacComNotifier struct {
	client *imkayaccom.Client
}

func NewImKayacComNotifier(user, password, secret string) *ImKayacComNotifier {
	var client *imkayaccom.Client

	if password != "" {
		client = imkayaccom.NewPasswordClient(user, password)
	} else if secret != "" {
		client = imkayaccom.NewSecretClient(user, secret)
	} else {
		client = imkayaccom.NewNoPasswordClient(user)
	}

	return &ImKayacComNotifier{client}
}

//TODO: handler support
func (i *ImKayacComNotifier) Notify(msg string) {
	Info("notify:", msg)
	i.client.Post(msg, "")
}
