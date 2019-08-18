package redirect

import (
	"context"

	session "github.com/hobord/infra/session"
)

// GrpcServer is base struct
type GrpcServer struct {
	configState *RedirectionConfigState
}

// CreateGrpcServer make an instace of GrpcServer
func CreateGrpcServer() *GrpcServer {
	configState := &RedirectionConfigState{}
	configState.loadConfigs("")

	srv := &GrpcServer{
		configState: configState,
	}

	return srv
}

// GetRedirection is implementing RedirectService rcp function
func (s *GrpcServer) GetRedirection(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	sessionValues := &session.Values{}
	redirections := make(map[string]int32)

	request := Request{
		Url:         in.Url,
		HttpMethod:  in.HttpMethod,
		HttpHeaders: in.Headers,
		RequestID:   in.RequestID,
		SessionID:   in.SessionID,
	}

	response, err := CalculateRedirections(ctx, s.configState, request, sessionValues, redirections)
	if err != nil {
		return &response, err
	}
	r := ParamPeeling(ctx, s.configState, Request{
		Url:         response.Location,
		HttpMethod:  in.HttpMethod,
		HttpHeaders: in.Headers,
		RequestID:   in.RequestID,
		SessionID:   in.SessionID,
	})
	if r.HttpStatusCode != 200 {
		return &r, nil
	}

	return &response, err
}
