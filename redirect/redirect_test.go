package redirect

import (
	"context"
	// "net/url"
	"testing"
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
	msg := &GetRedirectionMessage{
		Url: "http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok#bookmark",
	}
	ctx := context.Context(context.Background())
	result, err := CalculateRedirection(ctx, msg)
	if err != nil {
		t.Errorf("I Got %v", err)
	}
	t.Logf("Result: %v", result.Location)
}
