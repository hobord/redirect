package redirect

import (
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// config := &RedirectionServiceConfig{}
	// config.GetConf()
	// t.Logf("Loaded: %v", config)

	// var runtime_viper = viper.New()

	// runtime_viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001","/config/hugo.yml")
	// runtime_viper.AddSecureRemoteProvider("etcd","http://127.0.0.1:4001","/config/hugo.yaml","/etc/secrets/mykeyring.gpg")

	// runtime_viper.SetConfigType("yml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	// err := runtime_viper.ReadRemoteConfig()

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

	cfgStore := &configstore{}
	cfgStore.LoadConfigs("../configs")
	t.Logf("Loaded: %v", cfgStore.Main.cfg)

}
