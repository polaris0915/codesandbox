package settings

import (
	"testing"
)

func TestInit(t *testing.T) {

	Init()

	t.Logf("RemoteConfig: %+v\n", RemoteConfig)
}
