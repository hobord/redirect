package redirect

import (
	"context"
	"net/http"
	"net/url"
)

// ParamPeeling is peeling the specific url query parameters
func ParamPeeling(ctx context.Context, configState *RedirectionConfigState, request Request) GetRedirectionResponse {
	response := GetRedirectionResponse{}

	u, err := url.Parse(request.URL)
	if err != nil {
		response.HttpStatusCode = 200
		return response
	}
	if u.RawQuery == "" {
		response.HttpStatusCode = 200
		return response
	}

	// Build the new url
	newURLStr := u.Scheme + "://"
	if u.User.String() != "" {
		newURLStr = newURLStr + u.User.String() + "@"
	}
	newURLStr = newURLStr + u.Host
	newURLStr = newURLStr + u.Path

	query := u.Query()
	peeled := false
	for key := range query {
		if paramPeelingKeyCheck(configState, key, u) == true {
			peeled = true
			delete(query, key)
		}
	}
	// if we not peeled anything then just return 200
	if peeled == false {
		response.HttpStatusCode = 200
		return response
	}

	// finish the new url build
	u.RawQuery = query.Encode()
	if u.RawQuery != "" {
		newURLStr = newURLStr + "?" + u.RawQuery
	}
	if u.Fragment != "" {
		newURLStr = newURLStr + "#" + u.Fragment
	}

	response = GetRedirectionResponse{
		Location:       newURLStr,
		HttpStatusCode: http.StatusTemporaryRedirect,
	}
	return response
}

// Param Peeling business logic to here
func paramPeelingKeyCheck(configState *RedirectionConfigState, param string, u *url.URL) bool {
	// TODO: Param Peeling business logic by url to here
	if host, ok := configState.ParamPeeling[u.Host]; ok {
		if keys, ok := host[u.Scheme]; ok {
			for _, key := range keys {
				if param == key {
					return true
				}
			}
		}
	}
	return false
}
