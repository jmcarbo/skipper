package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateSelect(t *testing.T) {
	for _, item := range []struct {
		msg        string
		media      []*medium
		inputType  mediaType
		outputType mediaType
		err        error
	}{{
		"no media",
		[]*medium{},
		0,
		0,
		missingInput,
	}, {
		"only output",
		[]*medium{{typ: innkeeper}},
		0,
		0,
		missingInput,
	}, {
		"missing output",
		[]*medium{{typ: file}},
		0,
		0,
		missingOutput,
	}, {
		"too many inputs",
		[]*medium{{typ: innkeeper}, {typ: inline}, {typ: file}},
		0,
		0,
		tooManyInputs,
	}, {
		"output first",
		[]*medium{{typ: innkeeper}, {typ: inline}},
		inline,
		innkeeper,
		nil,
	}, {
		"output last",
		[]*medium{{typ: inline}, {typ: innkeeper}},
		inline,
		innkeeper,
		nil,
	}, {
		"output not innkeeper",
		[]*medium{{typ: inline}, {typ: etcd}},
		0,
		0,
		missingOutput,
	}, {
		"invalid input",
		[]*medium{{typ: innkeeper}, {typ: etcd}},
		0,
		0,
		invalidInputType,
	}} {
		if input, output, err := validateSelectInsert(&args{media: item.media}); err != item.err {
			t.Error("validate and select media failed, error case", item.msg, err, item.err)
		} else if err == nil && input.typ != item.inputType {
			t.Error("validate and select media failed, input type", item.msg, input.typ, item.inputType)
		} else if err == nil && output.typ != item.outputType {
			t.Error("validate and select media failed, output type", item.msg, output.typ, item.outputType)
		}
	}
}

func TestNewInsert(t *testing.T) {
	urls, err := stringsToUrls("http://innkeeper.example.org:8080")
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range []struct {
		msg   string
		media []*medium
		err   bool
	}{{
		"invalid media",
		[]*medium{{typ: etcd}, {typ: inline}},
		true,
	}, {
		"input client fails",
		[]*medium{{typ: innkeeper}, {typ: file, path: "/"}},
		true,
	}, {
		"all fine and dandy",
		[]*medium{{typ: innkeeper, urls: urls}, {typ: inline}},
		false,
	}} {
		if c, err := newInsert(&args{media: item.media}); err == nil && item.err || err != nil && !item.err {
			t.Error("new insert failed, error case", item.msg, err, item.err)
		} else if err == nil && (c.(*insertCommand).loader == nil || c.(*insertCommand).inserter == nil) {
			t.Error("new insert failed, loader and/or inserter missing", item.msg)
		}
	}
}

func TestExecuteInsert(t *testing.T) {
	for _, item := range []struct {
		msg        string
		args       *args
		serverFail bool
		err        bool
	}{{
		"load all fail",
		&args{
			media: []*medium{{
				typ:   inline,
				eskip: "invalid eskip",
			}, {
				typ:        innkeeper,
				oauthToken: "test-token",
			}},
		},
		false,
		true,
	}, {
		"insert fail",
		&args{
			media: []*medium{{
				typ:   inline,
				eskip: "Any() -> <shunt>",
			}, {
				typ:        innkeeper,
				oauthToken: "test-token",
			}},
		},
		true,
		true,
	}, {
		"insert succeed",
		&args{
			media: []*medium{{
				typ:   inline,
				eskip: "Any() -> <shunt>; Path(\"/some\") -> \"https://www.example.org\"",
			}, {
				typ:        innkeeper,
				oauthToken: "test-token",
			}},
		},
		false,
		true,
	}} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if item.serverFail {
				w.WriteHeader(http.StatusBadRequest)
			}
		}))

		urls, err := stringsToUrls(s.URL)
		if err != nil {
			t.Error(err)
			return
		}

		for _, m := range item.args.media {
			m.urls = urls
		}

		cmd, err := newInsert(item.args)
		if err != nil {
			t.Error(err)
			return
		}

		err = cmd.execute()
		if err == nil && item.err || err != nil && !item.err {
			t.Error("insert execute failed, error case", item.msg, err, item.err)
		}
	}
}
