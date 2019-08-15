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
func CalculateRedirections(ctx context.Context, request Request, sessionValues *session.Values, redirections map[string]int32) (GetRedirectionResponse, error) {
	//make business logic
	response := GetRedirectionResponse{
		Location:       request.Url,
		HttpStatusCode: http.StatusOK,
	}

	// Apply all rules
	newRedirection, err := applyRedirectionRules(ctx, request, sessionValues)
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

	r, err := CalculateRedirections(ctx, redirectTo, sessionValues, redirections)
	if err != nil {
		return response, err
	}
	response = r

	return response, nil
}

// applyRedirectionRules is apply the redirection rules
func applyRedirectionRules(ctx context.Context, request Request, sessionValues *session.Values) (GetRedirectionResponse, error) {
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

	// TODO: TEST IT
	if rules, ok := hostsRules[u.Host]; ok {
		for _, rule := range rules {
			switch rule.Type {
			case RuleHashTable:
				if rule.DefaultHTTPStatusCode != 0 {
					response.HttpStatusCode = rule.DefaultHTTPStatusCode
				}
				if hashRule, ok := rule.HasmapRules[request.Url]; ok {
					response.Location = hashRule.Target
					if hashRule.HTTPStatusCode != 0 {
						response.HttpStatusCode = hashRule.HTTPStatusCode
					}
				}
			case RuleRegexp:
			case CustomLogic:
			}
		}
	}

	// END of businesslogic

	return response, nil
}

type RuleType int

const (
	RuleRegexp RuleType = iota
	RuleHashTable
	CustomLogic
)

type HashRule struct {
	Target         string
	HTTPStatusCode int32
}
type HashRules map[string]HashRule

type Rule struct {
	Type                  RuleType
	Methods               []string
	HTTPHeaders           []string
	Expression            string
	LogicName             string
	DefaultHTTPStatusCode int32
	FilePath              string
	HasmapRules           HashRules
}
type OrderedRules []Rule
type HostsRules map[string]OrderedRules

var hostsRules HostsRules

func init() {
	hostsRules = make(HostsRules)
}
