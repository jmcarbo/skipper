package flowid

import (
	"github.com/zalando/skipper/filters"
	"github.com/zalando/skipper/filters/filtertest"
	"net/http"
	"testing"
)

const (
	testFlowId    = "FLOW-ID-FOR-TESTING"
	invalidFlowId = "[<>] (o) [<>]"
)

var (
	testFlowIdSpec           = New()
	filterConfigWithReuse    = []interface{}{ReuseParameterValue}
	filterConfigWithoutReuse = []interface{}{"dummy"}
)

func TestNewFlowIdGeneration(t *testing.T) {
	f, _ := testFlowIdSpec.CreateFilter(filterConfigWithReuse)
	fc := buildfilterContext()
	f.Request(fc)

	flowId := fc.Request().Header.Get(HeaderName)
	if flowId == "" {
		t.Errorf("flowId not generated")
	}
}

func TestFlowIdReuseExisting(t *testing.T) {
	f, _ := testFlowIdSpec.CreateFilter(filterConfigWithReuse)
	fc := buildfilterContext(HeaderName, testFlowId)
	f.Request(fc)

	flowId := fc.Request().Header.Get(HeaderName)
	if flowId != testFlowId {
		t.Errorf("Got wrong flow id. Expected '%s' got '%s'", testFlowId, flowId)
	}
}

func TestFlowIdIgnoreReuseExisting(t *testing.T) {
	f, _ := testFlowIdSpec.CreateFilter(filterConfigWithoutReuse)
	fc := buildfilterContext(HeaderName, testFlowId)
	f.Request(fc)

	flowId := fc.Request().Header.Get(HeaderName)
	if flowId == testFlowId {
		t.Errorf("Got wrong flow id. Expected a newly generated flowid but got the test flow id '%s'", flowId)
	}
}

func TestFlowIdRejectInvalidReusedFlowId(t *testing.T) {
	f, _ := testFlowIdSpec.CreateFilter(filterConfigWithReuse)
	fc := buildfilterContext(HeaderName, invalidFlowId)
	f.Request(fc)

	flowId := fc.Request().Header.Get(HeaderName)
	if flowId == invalidFlowId {
		t.Errorf("Got wrong flow id. Expected a newly generated flowid but got the test flow id '%s'", flowId)
	}
}

func TestFlowIdWithSpecificLen(t *testing.T) {
	fc := []interface{}{ReuseParameterValue, float64(42.0)}
	f, _ := testFlowIdSpec.CreateFilter(fc)
	fctx := buildfilterContext()
	f.Request(fctx)

	flowId := fctx.Request().Header.Get(HeaderName)

	l := len(flowId)
	if l != 42 {
		t.Errorf("Wrong flowId len. Expected %d, got %d", 42, l)
	}
}

func TestFlowIdWithInvalidParameters(t *testing.T) {
	fc := []interface{}{true}
	_, err := testFlowIdSpec.CreateFilter(fc)
	if err != filters.ErrInvalidFilterParameters {
		t.Errorf("Expected an invalid parameters error, got %v", err)
	}

	fc = []interface{}{"", float64(MinLength - 1)}
	_, err = testFlowIdSpec.CreateFilter(fc)
	if err != filters.ErrInvalidFilterParameters {
		t.Errorf("Expected an invalid parameters error, got %v", err)
	}
}

func buildfilterContext(headers ...string) filters.FilterContext {
	r, _ := http.NewRequest("GET", "http://example.org", nil)
	for i := 0; i < len(headers); i += 2 {
		r.Header.Set(headers[i], headers[i+1])
	}
	return &filtertest.Context{FRequest: r}
}
