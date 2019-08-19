package redirect

import (
	"context"
	"net/url"
	"testing"
)

func TestParampeelingKeyCheck(t *testing.T) {
	configState := &RedirectionConfigState{}
	configState.ParamPeeling = make(map[string]paramPeelingByProtocols)
	configState.ParamPeeling["site.com"] = make(map[string][]string)
	configState.ParamPeeling["site.com"]["http"] = []string{"key1", "key2"}
	cases := []struct {
		id     int
		key    string
		urlstr string
		want   bool
	}{
		{
			1, "key1", "http://site.com/?key1=123", true,
		},
	}
	for _, testcase := range cases {
		u, err := url.Parse(testcase.urlstr)
		if err != nil {
			t.Errorf("Error in convert the case url %v", err)
		}
		result := paramPeelingKeyCheck(configState, testcase.key, u)
		if result != testcase.want {
			t.Errorf("Error with id: %v", testcase.id)
		}
	}
}

func TestParamPeeling(t *testing.T) {
	cases := []struct {
		id                   int
		url                  string
		wantedLocation       string
		wantedHTTPStatusCode int32
	}{
		{
			1,
			"http://site.com/path/subpath/?foo=bar&key1=xyz&other=ok",
			"http://site.com/path/subpath/?foo=bar&other=ok",
			307,
		},
		{
			1,
			"http://site.com/path/subpath/?key1=xyz&other=ok",
			"http://site.com/path/subpath/?other=ok",
			307,
		},
		{
			1,
			"http://site.com/path/subpath/?other=ok&key1=xyz",
			"http://site.com/path/subpath/?other=ok",
			307,
		},
		{
			1,
			"http://site.com/path/subpath/?key1=xyz",
			"http://site.com/path/subpath/",
			307,
		},
		{
			2,
			"http://site.com/path/subpath/?foo=bar&other=ok",
			"http://site.com/path/subpath/?foo=bar&other=ok",
			200,
		},
		{
			3,
			"http://site.com/path/subpath/",
			"http://site.com/path/subpath/",
			200,
		},
		{
			4,
			"http://site.com/path/subpath/#dsdas",
			"http://site.com/path/subpath/#dsdas",
			200,
		},
		{
			5,
			"http://site.com/path/subpath/?foo=bar&torem=xyz&other=ok",
			"http://site.com/path/subpath/?foo=bar&torem=xyz&other=ok",
			200,
		},
	}

	for _, testcase := range cases {
		configState := &RedirectionConfigState{}
		configState.ParamPeeling = make(map[string]paramPeelingByProtocols)
		configState.ParamPeeling["site.com"] = make(map[string][]string)
		configState.ParamPeeling["site.com"]["http"] = []string{"key1", "key2"}

		ctx := context.Context(context.Background())
		request := Request{
			URL: testcase.url,
		}
		result := ParamPeeling(ctx, configState, request)
		if result.HttpStatusCode != testcase.wantedHTTPStatusCode {
			t.Errorf("Error with id: %v, wrong status code result (wanted: %v, result %v)", testcase.id, testcase.wantedHTTPStatusCode, result.HttpStatusCode)
		}

		if result.HttpStatusCode != 200 && result.Location != testcase.wantedLocation {
			t.Errorf("Error with id: %v, wrong location restult: %v", testcase.id, result.Location)
		}
	}
}
