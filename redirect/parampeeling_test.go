package redirect

import (
	"context"
	"net/url"
	"testing"
)

func TestParampeelingKeyCheck(t *testing.T) {
	cases := []struct {
		id     int
		key    string
		urlstr string
		want   bool
	}{
		{
			1, "toremove", "NOTIMPLEMENTED", true,
		},
		{
			2, "notremove", "NOTIMPLEMENTED", false,
		},
		{
			3, "", "NOTIMPLEMENTED", false,
		},
	}
	for _, testcase := range cases {
		u, err := url.Parse(testcase.urlstr)
		if err != nil {
			t.Errorf("Error in convert the case url %v", err)
		}
		result := parampeelingKeyCheck(testcase.key, u)
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
		wantedHttpStatusCode int32
	}{
		{
			1,
			"http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok",
			"http://index.hu/path/subpath/?foo=bar&other=ok",
			307,
		},
		{
			1,
			"http://index.hu/path/subpath/?toremove=xyz&other=ok",
			"http://index.hu/path/subpath/?other=ok",
			307,
		},
		{
			1,
			"http://index.hu/path/subpath/?other=ok&toremove=xyz",
			"http://index.hu/path/subpath/?other=ok",
			307,
		},
		{
			1,
			"http://index.hu/path/subpath/?toremove=xyz",
			"http://index.hu/path/subpath/",
			307,
		},
		{
			2,
			"http://index.hu/path/subpath/?foo=bar&other=ok",
			"http://index.hu/path/subpath/?foo=bar&other=ok",
			200,
		},
		{
			3,
			"http://index.hu/path/subpath/",
			"http://index.hu/path/subpath/",
			200,
		},
		{
			4,
			"http://index.hu/path/subpath/#dsdas",
			"http://index.hu/path/subpath/#dsdas",
			200,
		},
		{
			5,
			"http://index.hu/path/subpath/?foo=bar&torem=xyz&other=ok",
			"http://index.hu/path/subpath/?foo=bar&torem=xyz&other=ok",
			200,
		},
	}

	for _, testcase := range cases {
		ctx := context.Context(context.Background())
		request := Request{
			Url: testcase.url,
		}
		result := ParamPeeling(ctx, request)
		if result.HttpStatusCode != testcase.wantedHttpStatusCode {
			t.Errorf("Error with id: %v, wrong status code result (wanted: %v, result %v)", testcase.id, testcase.wantedHttpStatusCode, result.HttpStatusCode)
		}

		if result.HttpStatusCode != 200 && result.Location != testcase.wantedLocation {
			t.Errorf("Error with id: %v, wrong location restult: %v", testcase.id, result.Location)
		}
	}
}
