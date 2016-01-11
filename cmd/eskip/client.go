package main

import (
	"errors"
	"github.com/zalando/skipper/eskipfile"
	innk "github.com/zalando/skipper/innkeeper"
    etcdc "github.com/zalando/skipper/etcd"
)

func createInnkeeperClient(m *medium) (*innk.Client, error) {
	auth := innk.CreateInnkeeperAuthentication(innk.AuthOptions{InnkeeperAuthToken: m.oauthToken})
	return innk.New(innk.Options{
		Address:        m.urls[0].String(),
		Insecure:       false,
		Authentication: auth})
}

func createClient(m *medium) (interface{}, error) {
	switch m.typ {
	case innkeeper:
		return createInnkeeperClient(m)
    case etcd:
        return etcdc.New(urlsToStrings(m.urls), m.path), nil
	case inline:
		return &inlineClient{routes: m.eskip}, nil
	case file:
		return eskipfile.Open(m.path)
	default:
		return nil, errors.New("this is not yet implemented")
	}
}
