package redirect

type metadataYaml struct {
	Name string `yaml:"name"`
}

// ParamPeelingConfigYaml config in yaml
type ParamPeelingConfigYaml struct {
	APIVersion string           `yaml:"apiVersion"` // RedirectionService/v1
	Kind       string           `yaml:"kind"`       // ParamPeelingConfig
	Metadata   metadataYaml     `yaml:"metadata"`
	Spec       paramPeelingSpec `yaml:"spec,omitempty"`
}

type paramPeelingSpec struct {
	Hosts     []string `yaml:"hosts"`
	Protocols []string `yaml:"protocols"`
	Params    []string `yaml:"params"`
}

//RedirectionsConfigYaml config in yaml
type RedirectionsConfigYaml struct {
	APIVersion string                `yaml:"apiVersion"` // RedirectionService/v1
	Kind       string                `yaml:"kind"`       // RedirectionsConfig
	Metadata   metadataYaml          `yaml:"metadata"`
	Spec       configRedirectionSpec `yaml:"spec,omitempty"`
}

type configRedirectionSpec struct {
	Hosts     []string `yaml:"hosts"`
	Protocols []string `yaml:"protocols"`
	Rules     []struct {
		Type            string                  `yaml:"type"`
		FileURL         string                  `yaml:"fileUrl,omitempty"`
		TargetsByURL    []redirectionTargetYaml `yaml:"targetsByURL,omitempty"`
		RegexExpression string                  `yaml:"expression,omitempty"`
		LogicName       string                  `yaml:"logicName,omitempty"`
		Target          string                  `yaml:"target,omitempty"`
		HTTPStatusCode  int32                   `yaml:"httpStatusCode,omitempty"`
	} `yaml:"rules,omitempty"`
}
type redirectionTargetYaml struct {
	Src            string `yaml:"src"`
	Target         string `yaml:"target"`
	HTTPStatusCode int32  `yaml:"httpStatusCode,omitempty"`
}
