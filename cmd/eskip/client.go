package main

import (
	"errors"
	"github.com/zalando/skipper/eskipfile"
	innkc "github.com/zalando/skipper/innkeeper"
)

func createInnkeeperClient(m *medium) (*innkc.Client, error) {
	auth := innkc.CreateInnkeeperAuthentication(innkc.AuthOptions{InnkeeperAuthToken: m.oauthToken})
	return innkc.New(innkc.Options{
		Address:        m.urls[0].String(),
		Insecure:       false,
		Authentication: auth})
}

func createClient(m *medium) (interface{}, error) {
	switch m.typ {
	case innkeeper:
		return createInnkeeperClient(m)
	case inline:
		return &inlineClient{routes: m.eskip}, nil
	case file:
		return eskipfile.Open(m.path)
	default:
		return nil, errors.New("this is not yet implemented")
	}
}
