package main

import (
	"github.com/zalando/skipper/eskip"
	// "github.com/zalando/skipper/eskipfile"
	// etcdclient "github.com/zalando/skipper/etcd"
	"io"
	"io/ioutil"
	// "os"
)

type readClient interface {
	LoadAndParseAll() ([]*eskip.RouteInfo, error)
}

type stdinReader struct {
	reader io.Reader
}

type inlineClient struct {
	routes string
}

type idsReader struct {
	ids []string
}

// func createReadClient(m *medium) (readClient, error) {
// 	// no output, no client
// 	if m == nil {
// 		return nil, nil
// 	}
//
// 	switch m.typ {
// 	case innkeeper:
// 		return createInnkeeperClient(m)
//
// 	case etcd:
// 		return etcdclient.New(urlsToStrings(m.urls), m.path), nil
//
// 	case stdin:
// 		return &stdinReader{reader: os.Stdin}, nil
//
// 	case file:
// 		return eskipfile.Open(m.path)
//
// 	case inline:
// 		return &inlineClient{routes: m.eskip}, nil
//
// 	case inlineIds:
// 		return &idsReader{ids: m.ids}, nil
//
// 	default:
// 		return nil, invalidInputType
// 	}
// }

func (r *stdinReader) LoadAndParseAll() ([]*eskip.RouteInfo, error) {
	// this pretty much disables continuous piping,
	// but since the reset command first upserts all
	// and deletes the diff only after, it may not
	// even be consistent to do continous piping.
	// May change in the future.
	doc, err := ioutil.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}

	routes, err := eskip.Parse(string(doc))

	if err != nil {
		return nil, err
	}

	return routesToRouteInfos(routes), nil
}

func (ic *inlineClient) LoadAll() ([]*eskip.Route, error) {
	return eskip.Parse(ic.routes)
}

func (r *idsReader) LoadAndParseAll() ([]*eskip.RouteInfo, error) {
	routeInfos := make([]*eskip.RouteInfo, len(r.ids))
	for i, id := range r.ids {
		routeInfos[i] = &eskip.RouteInfo{eskip.Route{Id: id}, nil}
	}

	return routeInfos, nil
}

func routesToRouteInfos(routes []*eskip.Route) (routeInfos []*eskip.RouteInfo) {
	for _, route := range routes {
		routeInfos = append(routeInfos, &eskip.RouteInfo{*route, nil})
	}
	return
}
