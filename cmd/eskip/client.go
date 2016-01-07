package main

import (
	"errors"
	innkc "github.com/zalando/skipper/innkeeper"
)

func createInnkeeperClient(m *medium) (*innkc.Client, error) {
	auth := innkc.CreateInnkeeperAuthentication(innkc.AuthOptions{InnkeeperAuthToken: m.oauthToken})

	ic, err := innkc.New(innkc.Options{
		Address:        m.urls[0].String(),
		Insecure:       false,
		Authentication: auth})

	if err != nil {
		return nil, err
	}
	return ic, nil
}

func createClient(m *medium) (interface{}, error) {
	switch m.typ {
	case innkeeper:
		return createInnkeeperClient(m)
	case inline:
		return &inlineClient{routes: m.eskip}, nil
	default:
		return nil, errors.New("this is not yet implemented")
	}
}
