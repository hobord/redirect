package redirect

import (
	"context"
	"fmt"

	session "github.com/hobord/infra/session"
)

// GrpcServer is base struct
type GrpcServer struct {
	configState *RedirectionConfigState
}

// CreateGrpcServer make an instace of GrpcServer
func CreateGrpcServer() *GrpcServer {
	configState := &RedirectionConfigState{}
	configState.loadConfigs("configs")

	srv := &GrpcServer{
		configState: configState,
	}

	return srv
}

// GetRedirection is implementing RedirectService rcp function
func (s *GrpcServer) GetRedirection(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	fmt.Printf("Get redirection: %v", in)
	sessionValues := &session.Values{}
	redirections := make(map[string]int32)

	request := Request{
		URL:         in.Url,
		HTTPMethod:  in.HttpMethod,
		HTTPHeaders: in.Headers,
		RequestID:   in.RequestID,
		SessionID:   in.SessionID,
	}

	response, err := CalculateRedirections(ctx, s.configState, request, sessionValues, redirections)
	if err != nil {
		return &response, err
	}
	r := ParamPeeling(ctx, s.configState, Request{
		URL:         response.Location,
		HTTPMethod:  in.HttpMethod,
		HTTPHeaders: in.Headers,
		RequestID:   in.RequestID,
		SessionID:   in.SessionID,
	})
	if r.HttpStatusCode != 200 {
		return &r, nil
	}
	fmt.Printf("Response: %v", response)
	return &response, err
}
