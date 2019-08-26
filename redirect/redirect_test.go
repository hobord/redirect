package redirect

import (
	"context"
	"regexp"
	"testing"

	"github.com/hobord/infra/session"
)

func TestApplyRedirectionRules(t *testing.T) {
	cases := []struct {
		id                   int
		urlstr               string
		wantedLocation       string
		wantedHTTPStatusCode int32
		desc                 string
	}{
		{
			1,
			"http://site.com/path",
			"http://newsite.hu/path",
			307,
			"Hash redirection test",
		},
		{
			2,
			"http://site.com/path/",
			"http://newsite.hu/path",
			307,
			"Hash redirection test",
		},
	}
	for _, testcase := range cases {
		ctx := context.Context(context.Background())
		configState := &RedirectionConfigState{}
		configState.RedirectionHosts = make(map[string]redirectionRulesByProtcols)
		configState.RedirectionHosts["site.com"] = make(map[string][]RedirectionRule)

		rule1 := RedirectionRule{
			Type:         "Hash",
			TargetsByURL: make(map[string]redirectionTarget),
		}
		rule1.TargetsByURL["http://site.com/path"] = redirectionTarget{
			Target:         "http://newsite.hu/path",
			HTTPStatusCode: 307,
		}

		r, err := regexp.Compile("http(s?):\\/\\/site.com\\/(.*)")
		if err != nil {
			t.Error(err) // TODO: errorlog
		}
		rule2 := RedirectionRule{
			Type:           "Regexp",
			Regexp:         r,
			HTTPStatusCode: 307,
			Target:         "http://newsite.hu/path",
		}
		configState.RedirectionHosts["site.com"]["http"] = []RedirectionRule{rule1, rule2}

		request := Request{
			URL: testcase.urlstr,
		}
		sessionValues := &session.Values{}

		// redirections := make(map[string]int32)
		result, err := applyRedirectionRules(ctx, configState, request, sessionValues)
		if err != nil {
			t.Errorf("Error in caed id: %v, %v", testcase.id, err)
		}
		if result.HttpStatusCode != testcase.wantedHTTPStatusCode {
			t.Errorf("Error with id: %v, wrong status code result (wanted: %v, result %v)", testcase.id, testcase.wantedHTTPStatusCode, result.HttpStatusCode)
		}
		if result.HttpStatusCode != 200 && result.Location != testcase.wantedLocation {
			t.Errorf("Error with id: %v, wrong location restult: %v", testcase.id, result.Location)
		}
		t.Log(result)
	}
}

/*
func TestGetRedirection(t *testing.T) {
	cases := []struct {
		id                   int
		urlstr               string
		wantedLocation       string
		wantedHTTPStatusCode int32
	}{
		{
			1,
			"http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok",
			"http://888.hu/path/subpath/?foo=bar&other=ok",
			307,
		},
	}
	for _, testcase := range cases {
		grpServer := &GrpcServer{}
		ctx := context.Context(context.Background())
		msg := &GetRedirectionMessage{
			Url: testcase.urlstr,
		}
		result, err := grpServer.GetRedirection(ctx, msg)
		if err != nil {
			t.Errorf("Error in caed id: %v, %v", testcase.id, err)
		}
		if result.HttpStatusCode != testcase.wantedHTTPStatusCode {
			t.Errorf("Error with id: %v, wrong status code result (wanted: %v, result %v)", testcase.id, testcase.wantedHTTPStatusCode, result.HttpStatusCode)
		}
		if result.HttpStatusCode != 200 && result.Location != testcase.wantedLocation {
			t.Errorf("Error with id: %v, wrong location restult: %v", testcase.id, result.Location)
		}
	}
}

func TestCalculateRedirection(t *testing.T) {

	cases := []struct {
		id                   int
		urlstr               string
		wantedLocation       string
		wantedHTTPStatusCode int32
	}{
		{
			1,
			"http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok",
			"http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok",
			200,
		},
	}
	for _, testcase := range cases {
		ctx := context.Context(context.Background())
		config := &RedirectionServiceConfig{}
		request := Request{
			Url: testcase.urlstr,
		}
		sessionValues := &session.Values{}
		redirections := make(map[string]int32)
		result, err := CalculateRedirections(ctx, config, request, sessionValues, redirections)
		if err != nil {
			t.Errorf("Error in caed id: %v, %v", testcase.id, err)
		}
		if result.HttpStatusCode != testcase.wantedHTTPStatusCode {
			t.Errorf("Error with id: %v, wrong status code result (wanted: %v, result %v)", testcase.id, testcase.wantedHTTPStatusCode, result.HttpStatusCode)
		}
		if result.HttpStatusCode != 200 && result.Location != testcase.wantedLocation {
			t.Errorf("Error with id: %v, wrong location restult: %v", testcase.id, result.Location)
		}
	}
}
*/
