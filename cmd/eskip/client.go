package main

import "github.com/zalando/skipper/innkeeper"

func createInnkeeperClient(m *medium) (*innkeeper.Client, error) {
	auth := innkeeper.CreateInnkeeperAuthentication(innkeeper.AuthOptions{InnkeeperAuthToken: m.oauthToken})

	ic, err := innkeeper.New(innkeeper.Options{
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
