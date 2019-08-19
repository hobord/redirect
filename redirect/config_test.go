package redirect

import (
	"log"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigLoad(t *testing.T) {
	// config := &RedirectionServiceConfig{}
	// config.GetConf()
	// t.Logf("Loaded: %v", config)

	// var runtime_viper = viper.New()

	// runtime_viper.SetConfigFile("../config.yml")
	// runtime_viper.SetConfigName("config")         // name of config file (without extension)
	// runtime_viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	// runtime_viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	// runtime_viper.AddConfigPath(".")              // optionally look for config in the working directory
	// runtime_viper.AddConfigPath("..")             // optionally look for config in the working directory

	// if err := runtime_viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file, %s", err)
	// }
	// err := runtime_viper.Unmarshal(config)
	// if err != nil {
	// 	log.Fatalf("unable to decode into struct, %v", err)
	// }
	// t.Logf("Loaded: %v", config)

	// cfgStore := &configstore{}
	// cfgStore.LoadConfigs("../configs")
	// t.Logf("Loaded: %v", cfgStore.Main.cfg)

}

func TestLoadConfigs(t *testing.T) {
	cfgState := &RedirectionConfigState{}
	cfgState.loadConfigs("../configs/test")
	t.Logf("Loaded: %v", cfgState)
}

func TestParampeelingConfigLoader(t *testing.T) {
	file := "../configs/test/peeling_example.yaml"
	cfgState := &RedirectionConfigState{}
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	cfgState.parampeelingConfigLoader(v)
	t.Logf("Loaded: %v", cfgState.ParamPeeling)
}

func TestRedirectionsConfigLoader(t *testing.T) {
	file := "../configs/test/redirections_example.yml"
	cfgState := &RedirectionConfigState{}
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	cfgState.redirectionsConfigLoader(v)
	t.Logf("Loaded: %v", cfgState.RedirectionHosts)
}
