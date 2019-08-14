package redirect

import (
	"context"
	"log"
	// "net/url"
	"testing"
)

func TestCalculateRedirection(t *testing.T) {
	msg := &GetRedirectionMessage{
		SessionID: "in.SessionID",
		RequestID: "in.RequestID",
		Url:       "http://index.hu?foo=bar&toremove=xyz",
	}
	ctx := context.Context(context.Background())
	result, err := CalculateRedirection(ctx, msg)
	if err != nil {
		t.Errorf("I Got %v", err)
	}
	log.Printf("result %v", result)
}
