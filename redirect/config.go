package redirect

import (
	"encoding/csv"
	fmt "fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/viper"
)

var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)

type RedirectionServiceConfig struct {
	APIVersion string       `yaml:"apiVersion"` // RedirectionService/v1
	Kind       string       `yaml:"kind"`       // RedirectionServiceConfig
	Metadata   confMetadata `yaml:"metadata"`
	Spec       struct {
		ParamPeeling []struct {
			Host       string                 `yaml:"host"`
			Protocols  []string               `yaml:"protocols,omitempty"`
			ConfigName string                 `yaml:"configName,omitempty"`
			Spec       configParamPeelingSpec `yaml:"spec,omitempty"`
		} `yaml:"paramPeeling,omitempty"`
		Redirections []struct {
			Host       string                `yaml:"host"`
			Protocols  []string              `yaml:"protocols,omitempty"`
			ConfigName string                `yaml:"configName,omitempty"`
			Spec       configRedirectionSpec `yaml:"spec,omitempty"`
		} `yaml:"redirections,omitempty"`
		RedirectionHostHasmap map[string]struct {
			Protocols []string
			Spec      configRedirectionSpec
		}
	} `yaml:"spec"`
}
type confMetadata struct {
	Name string `yaml:"name"`
}

type ParamPeelingConfig struct {
	APIVersion string                 `yaml:"apiVersion"` // RedirectionService/v1
	Kind       string                 `yaml:"kind"`       // ParampeelingConfig
	Metadata   confMetadata           `yaml:"metadata"`
	Spec       configParamPeelingSpec `yaml:"spec,omitempty"`
}

type configParamPeelingSpec struct {
	Params []string `yaml:"params"`
}

type RedirectionsConfig struct {
	APIVersion string                `yaml:"apiVersion"` // RedirectionService/v1
	Kind       string                `yaml:"kind"`       // RedirectionsConfig
	Metadata   confMetadata          `yaml:"metadata"`
	Spec       configRedirectionSpec `yaml:"spec,omitempty"`
}

type configRedirectionSpec struct {
	Rules []struct {
		Type           string   `yaml:"type"`
		HTTPMethods    []string `yaml:"httpMethods,omitempty"`
		FileURL        string   `yaml:"fileUrl,omitempty"`
		HTTPStatusCode int32    `yaml:"httpStatusCode,omitempty"`
		HashMap        map[string]struct {
			Target         string `yaml:target,omitempty`
			HTTPStatusCode int32  `yaml:httpStatusCode,omitempty`
		} `yaml: "hasmap,omitempty"`
		LogicName  string         `yaml:"logicName,omitempty"`
		Expression string         `yaml:"expression,omitempty"`
		Regexp     *regexp.Regexp `yaml:"-"`
		Target     string         `yaml:target,omitempty`
	} `yaml:"rules,omitempty"`
}

// func (c *RedirectionServiceConfig) GetConf() *RedirectionServiceConfig {

// 	yamlFile, err := ioutil.ReadFile("../config.yml")
// 	if err != nil {
// 		log.Printf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, c)
// 	if err != nil {
// 		log.Fatalf("Unmarshal: %v", err)
// 	}

// 	return c
// }

// func (c *ParamPeelingConfig) GetConf(fileName string) *ParamPeelingConfig {
// 	yamlFile, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		log.Printf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, c)
// 	if err != nil {
// 		log.Fatalf("Unmarshal: %v", err)
// 	}

// 	return c
// }

// func (c *RedirectionsConfig) GetConf(fileName string) *RedirectionsConfig {
// 	yamlFile, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		log.Printf("yamlFile.Get err   #%v ", err)
// 	}
// 	err = yaml.Unmarshal(yamlFile, c)
// 	if err != nil {
// 		log.Fatalf("Unmarshal: %v", err)
// 	}

// 	return c
// }

type configstore struct {
	Main struct {
		viper *viper.Viper
		cfg   *RedirectionServiceConfig
	}
	ParamPeelingConfigs map[string]struct {
		viper *viper.Viper
		cfg   *ParamPeelingConfig
	}
	RedirectionsConfigs map[string]struct {
		viper *viper.Viper
		cfg   *RedirectionsConfig
	}
}

func (cst *configstore) LoadConfigs(root string) {
	cst.ParamPeelingConfigs = make(map[string]struct {
		viper *viper.Viper
		cfg   *ParamPeelingConfig
	})
	cst.RedirectionsConfigs = make(map[string]struct {
		viper *viper.Viper
		cfg   *RedirectionsConfig
	})

	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(file)

		v := viper.New()
		v.SetConfigFile(file)
		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		// var cfg interface{}

		cfgKind := v.GetString("kind")
		cfgName := v.GetString("metadata.name")
		switch cfgKind {
		case "RedirectionServiceConfig":
			cfg := &RedirectionServiceConfig{}
			err = v.Unmarshal(&cfg)
			if err != nil {
				log.Fatalf("unable to decode into struct, %v", err)
			}
			cst.Main.viper = v
			cst.Main.cfg = cfg
		case "ParampeelingConfig":
			cfg := &ParamPeelingConfig{}
			err = v.Unmarshal(&cfg)
			if err != nil {
				log.Fatalf("unable to decode into struct, %v", err)
			}
			cst.ParamPeelingConfigs[cfgName] = struct {
				viper *viper.Viper
				cfg   *ParamPeelingConfig
			}{
				viper: v,
				cfg:   cfg,
			}
		case "RedirectionsConfig":
			cfg := &RedirectionsConfig{}
			err = v.Unmarshal(&cfg)
			if err != nil {
				log.Fatalf("unable to decode into struct, %v", err)
			}
			cst.RedirectionsConfigs[cfgName] = struct {
				viper *viper.Viper
				cfg   *RedirectionsConfig
			}{
				viper: v,
				cfg:   cfg,
			}
		}

	}

	cst.ParseConfigs()
}

func (cst *configstore) ParseConfigs() {
	for i, c := range cst.Main.cfg.Spec.ParamPeeling {
		if cfg, ok := cst.ParamPeelingConfigs[c.ConfigName]; ok {
			cst.Main.cfg.Spec.ParamPeeling[i].Spec = cfg.cfg.Spec
		}
	}

	for i, c := range cst.Main.cfg.Spec.Redirections {
		if cfg, ok := cst.RedirectionsConfigs[c.ConfigName]; ok {
			cst.Main.cfg.Spec.Redirections[i].Spec = cfg.cfg.Spec
		}
	}

	cst.Main.cfg.Spec.RedirectionHostHasmap = make(map[string]struct {
		Protocols []string
		Spec      configRedirectionSpec
	})
	for _, c := range cst.Main.cfg.Spec.Redirections {
		for i, rule := range c.Spec.Rules {
			switch rule.Type {
			case "Regex":
				r, err := regexp.Compile(rule.Expression)
				if err != nil {
					panic("hibas regex")
					// TODO remove this rule and continue
				}
				c.Spec.Rules[i].Regexp = r
			case "Hash":
				if rule.FileURL != "" {
					//TODO: Load csv file into hasmap
					// Open CSV file
					f, err := os.Open(rule.FileURL)
					if err != nil {
						panic(err)
					}
					defer f.Close()
					lines, err := csv.NewReader(f).ReadAll()
					if err != nil {
						panic(err)
					}
					// c.Spec.Rules[i].HashMap = make
					// Loop through lines & turn into object
					for _, line := range lines {
						i, err := strconv.ParseInt(line[2], 10, 32)
						if err != nil {
							panic(err)
						}
						c.Spec.Rules[i].HashMap[line[0]] = struct {
							Target         string `yaml:target,omitempty`
							HTTPStatusCode int32  `yaml:httpStatusCode,omitempty`
						}{
							Target:         line[1],
							HTTPStatusCode: int32(i),
						}
					}
				}
			}

		}
		cst.Main.cfg.Spec.RedirectionHostHasmap[c.Host] = struct {
			Protocols []string
			Spec      configRedirectionSpec
		}{
			Protocols: c.Protocols,
			Spec:      c.Spec,
		}
	}
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
	Type                  RuleType `json:"name" yaml:"name"`
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
