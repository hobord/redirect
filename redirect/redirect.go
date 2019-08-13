package redirect

import (
	"context"
	"net/url"
)

type GrpcServer struct{}

func CreateGrpcServer() *GrpcServer {
	impl := &GrpcServer{}
	return impl
}

func (s *GrpcServer) GetRedirection(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	var response *GetRedirectionResponse
	response, err := calculateRedirection(ctx, in)
	return response, err
}

func calculateRedirection(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	//make business logic
	reponse := &GetRedirectionResponse{}

	// maybe want to make an other redirection
	redirectTo := &GetRedirectionMessage{
		SessionID:  in.SessionID,
		RequestID:  in.RequestID,
		Url:        reponse.Location,
		HttpMethod: in.HttpMethod}
	// check it
	r, err := calculateRedirection(ctx, redirectTo)
	if err != nil {
		return reponse, err
	}
	if r.Location != in.Url {
		return r, nil
	}

	return reponse, nil
}

func paramPeeling(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	//make business logic
	reponse := &GetRedirectionResponse{}
	// fake peeling logic

	u, err := url.Parse(in.Url)
	if err != nil {
		return reponse, err
	}
	newURLStr := in.Url
	if u.RawQuery != "" {
		newURLStr = u.Scheme + "://"
		if u.User.String() != "" {
			newURLStr = newURLStr + u.User.String() + "@"
		}
		newURLStr = newURLStr + u.Host + "/"
		newURLStr = newURLStr + u.Path

		query := u.Query()
		for key := range query {
			if parampeelingKeyCheck(key, u) == true {
				delete(query, key)
			}
		}
		u.RawQuery = query.Encode()

		if u.Fragment != "" {
			newURLStr = newURLStr + "#" + u.Fragment
		}
	}
	reponse.Location = newURLStr
	return reponse, nil
}

func parampeelingKeyCheck(key string, u *url.URL) bool {
	if key == "toremove" {
		return true
	}
	return false
}
