package redirect

import (
	"context"
	"net/http"
	"net/url"
)

type GrpcServer struct{}

func CreateGrpcServer() *GrpcServer {
	impl := &GrpcServer{}
	return impl
}

func (s *GrpcServer) GetRedirection(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	var response *GetRedirectionResponse
	response, err := CalculateRedirection(ctx, in)

	peeled, err := ParamPeeling(ctx, &GetRedirectionMessage{
		Url: response.Location,
	})
	if err != nil {
		return response, err
	}
	response = peeled

	return response, err
}

func CalculateRedirection(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
	//make business logic
	response := &GetRedirectionResponse{
		Location:       in.Url,
		HttpStatusCode: http.StatusTemporaryRedirect,
	}

	// peeled, err := ParamPeeling(ctx, in)
	// if err != nil {
	// 	return response, err
	// }
	// response = peeled

	// maybe want to make an other redirection
	// check itxxxz
	if response.Location != in.Url {
		redirectTo := &GetRedirectionMessage{
			SessionID:  in.SessionID,
			RequestID:  in.RequestID,
			Url:        response.Location,
			HttpMethod: in.HttpMethod}

		r, err := CalculateRedirection(ctx, redirectTo)
		if err != nil {
			return response, err
		}
		response = r
	}
	return response, nil
}

func ParamPeeling(ctx context.Context, in *GetRedirectionMessage) (*GetRedirectionResponse, error) {
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
		newURLStr = newURLStr + u.Host
		newURLStr = newURLStr + u.Path

		query := u.Query()
		for key := range query {
			if parampeelingKeyCheck(key, u) == true {
				delete(query, key)
			}
		}
		u.RawQuery = query.Encode()
		newURLStr = newURLStr + "?" + u.RawQuery

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
