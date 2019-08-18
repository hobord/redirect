package redirect

import (
	"context"
	"net/http"
	"net/url"

	session "github.com/hobord/infra/session"
)

type Request struct {
	Url         string
	HttpMethod  string
	HttpHeaders map[string]string
	RequestID   string
	SessionID   string
}

// CalculateRedirections make recursive the redirections, with infinitive redirect loop detection.
func CalculateRedirections(ctx context.Context, configState *RedirectionConfigState, request Request, sessionValues *session.Values, redirections map[string]int32) (GetRedirectionResponse, error) {
	//make business logic
	response := GetRedirectionResponse{
		Location:       request.Url,
		HttpStatusCode: http.StatusOK,
	}

	// Apply all rules
	newRedirection, err := applyRedirectionRules(ctx, configState, request, sessionValues)
	if err != nil {
		return response, nil // TODO: it is ok?
	}

	if newRedirection.HttpStatusCode == http.StatusOK {
		return response, nil
	}

	// infinitive redirect loop detection
	if httpStatusCode, ok := redirections[newRedirection.Location]; ok {
		if httpStatusCode == newRedirection.HttpStatusCode {
			return response, nil
		}
	}
	redirections[newRedirection.Location] = newRedirection.HttpStatusCode
	response = newRedirection

	// We have changes, lets make a new loop
	redirectTo := Request{
		SessionID:   request.SessionID,
		RequestID:   request.RequestID,
		Url:         response.Location,
		HttpHeaders: request.HttpHeaders,
		HttpMethod:  request.HttpMethod}

	r, err := CalculateRedirections(ctx, configState, redirectTo, sessionValues, redirections)
	if err != nil {
		return response, err
	}
	response = r

	return response, nil
}

// applyRedirectionRules is apply the redirection rules
func applyRedirectionRules(ctx context.Context, configState *RedirectionConfigState, request Request, sessionValues *session.Values) (GetRedirectionResponse, error) {
	response := GetRedirectionResponse{
		Location:       request.Url,
		HttpStatusCode: http.StatusOK,
	}

	u, err := url.Parse(request.Url)
	if err != nil {
		return response, err
	}

	// TODO: make businesslogic to here
	if u.Host == "index.hu" {
		u.Host = "444.hu"
	} else if u.Host == "444.hu" {
		u.Host = "888.hu"
	}
	response.Location = u.String()
	response.HttpStatusCode = http.StatusTemporaryRedirect

	if host, ok := configState.RedirectionHosts[u.Host]; ok {
		if rules, ok := host[u.Scheme]; ok {
			for _, rule := range rules {
				switch rule.Type {
				case "Hash":
					if redirectTo, found := rule.TargetsByURL[u.Host]; found {
						response.Location = redirectTo.Target
						if redirectTo.HTTPStatusCode > 0 {
							response.HttpStatusCode = redirectTo.HTTPStatusCode
						} else {
							response.HttpStatusCode = rule.HTTPStatusCode
						}
						goto End
					}
				case "Regexp":
					if rule.Regexp.MatchString(u.String()) {
						response.Location = rule.Target
						response.HttpStatusCode = rule.HTTPStatusCode
						goto End
					}
				case "CustomLogic":
					switch rule.LogicName {
					case "condition":
						// r = functionName(ctx, config, request, sessionValues)
						// if r.HttpStatusCode != 200 {
						// 	response = r
						//   	goto End
						// }
					}
				}
			}
		}
	}
End:
	// END of businesslogic

	return response, nil
}
