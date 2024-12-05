package settings

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	// Expected configuration values
	expectedConfig := createDockerConfig{
		ImageName:      "ubuntu:20.04",
		HostMountDir:   "/Users/alin-youlinlin/Desktop/polaris-all_projects/codesandbox/user_code",
		DockerMountDir: "/user_code",
		ContainerName:  "cpp-sandbox",
		InitContainerCmd: []string{
			"bash",
			"-c",
			"apt-get update && apt-get install -y build-essential g++ vim",
		},
	}

	// Run the Init function to load the configuration
	Init()

	// Compare the loaded configuration with the expected values
	if !reflect.DeepEqual(CreateDockerConfig, expectedConfig) {
		t.Errorf("Loaded CreateDockerConfig does not match the expected configuration.\nExpected: %+v\nGot: %+v", expectedConfig, CreateDockerConfig)
	}
}
