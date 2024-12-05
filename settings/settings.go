package settings

import (
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

var (
	CreateDockerConfig createDockerConfig
	LoggerConfig       loggerConfig
)

func Init() {
	vp := viper.New()
	_, file, _, _ := runtime.Caller(0)
	configPath, _ := filepath.Abs(filepath.Dir(file) + "/application.dev.yaml")
	vp.SetConfigFile(configPath)
	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}
	vp.UnmarshalKey("CreateDockerConfig", &CreateDockerConfig)
	vp.UnmarshalKey("LoggerConfig", &LoggerConfig)

	//fmt.Printf("%+v\n", CreateDockerConfig)
}

type createDockerConfig struct {
	ImageName        string   `json:"imageName" mapstructure:"ImageName"`
	HostMountDir     string   `json:"hostMountDir" mapstructure:"HostMountDir"`
	DockerMountDir   string   `json:"dockerMountDir" mapstructure:"DockerMountDir"`
	ContainerName    string   `json:"containerName" mapstructure:"ContainerName"`
	InitContainerCmd []string `json:"initContainerCmd" mapstructure:"InitContainerCmd"`
}

type loggerConfig struct {
	LogPath string `json:"logPath" mapstructure:"LogPath"`
}
