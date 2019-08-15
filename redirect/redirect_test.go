package redirect

import (
	"context"
	"testing"

	"github.com/hobord/infra/session"
)

func TestGetRedirection(t *testing.T) {
	msg := &GetRedirectionMessage{
		Url: "http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok#bookmark",
	}
	grpServer := &GrpcServer{}
	ctx := context.Context(context.Background())
	result, err := grpServer.GetRedirection(ctx, msg)
	if err != nil {
		t.Errorf("I Got %v", err)
	}
	t.Logf("Result: %v", result.Location)
}

func TestCalculateRedirection(t *testing.T) {
	ctx := context.Context(context.Background())
	request := Request{
		Url: "http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok#bookmark",
	}
	sessionValues := &session.Values{}
	redirections := make(map[string]int32)

	result, err := CalculateRedirections(ctx, request, sessionValues, redirections)
	if err != nil {
		t.Errorf("I Got %v", err)
	}
	t.Logf("Result: %v", result.Location)
}
