package redirect

import "regexp"

//RedirectionConfigState config state
type RedirectionConfigState struct {
	RedirectionHosts map[string]redirectionRulesByProtcols // haskey by hostname: www.site.com
	ParamPeeling     map[string]paramPeelingByProtocols
}

type redirectionRulesByProtcols map[string][]RedirectionRule // haskeys http / https

// RedirectionRule redirection rule in state
type RedirectionRule struct {
	Type           string
	LogicName      string
	Regexp         *regexp.Regexp
	TargetsByURL   map[string]redirectionTarget
	Target         string
	HTTPStatusCode int32
}

type redirectionTarget struct {
	Target         string
	HTTPStatusCode int32
}

type paramPeelingByProtocols map[string][]string // haskeys http / https
