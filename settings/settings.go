package settings

import (
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

var (
	ContainersConfig containersConfig
	LoggerConfig     loggerConfig
	MysqlConfig      mysqlConfig
	RemoteConfig     remoteConfig
	EnvConfig        envConfig
)

func Init() {
	vp := viper.New()
	_, file, _, _ := runtime.Caller(0)
	configPath, _ := filepath.Abs(filepath.Dir(file) + "/application.dev.yaml")
	vp.SetConfigFile(configPath)
	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}
	vp.UnmarshalKey("EnvConfig", &EnvConfig)
	vp.UnmarshalKey("ContainersConfig", &ContainersConfig)
	vp.UnmarshalKey("LoggerConfig", &LoggerConfig)
	vp.UnmarshalKey("MysqlConfig", &MysqlConfig)
	//vp.UnmarshalKey("JwtConfig", &JwtConfig)
	vp.UnmarshalKey("RemoteConfig", &RemoteConfig)
}

type envConfig struct {
	RootPath string `json:"rootPath" mapstructure:"RootPath"`
}

type cpp struct {
	ImageName        string   `json:"imageName" mapstructure:"ImageName"`
	HostMountDir     string   `json:"hostMountDir" mapstructure:"HostMountDir"`
	DockerMountDir   string   `json:"dockerMountDir" mapstructure:"DockerMountDir"`
	ContainerName    string   `json:"containerName" mapstructure:"ContainerName"`
	ContainerId      string   `json:"containerId" mapstructure:"ContainerId"`
	InitContainerCmd []string `json:"initContainerCmd" mapstructure:"InitContainerCmd"`
}

type containersConfig struct {
	Cpp cpp `json:"cpp" mapstructure:"Cpp"`
}

type loggerConfig struct {
	LogPath string `json:"logPath" mapstructure:"LogPath"`
}

type mysqlConfig struct {
	MysqlDBName     string `json:"mysqlDBName" mapstructure:"MysqlDBName"`
	MysqlDBPassword string `json:"mysqlDBPassword" mapstructure:"MysqlDBPassword"`
}

type jwtConfig struct {
	Key      string `json:"key" mapstructure:"Key"`
	NeedAuth bool   `json:"needAuth" mapstructure:"NeedAuth"`
}

type remoteConfig struct {
	Url       string    `json:"url" mapstructure:"Url"`
	Method    string    `json:"method" mapstructure:"Method"`
	JwtConfig jwtConfig `json:"jwtConfig" mapstructure:"JwtConfig"`
}
