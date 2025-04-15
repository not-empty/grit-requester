package requesterV2

import (
	"testing"
)

func TestSetGetConf(t *testing.T) {
	conf := StaticConfig{}
	conf.Set("test", MSAuthConf{
		Token:   "test",
		Secret:  "test",
		Context: "test",
		BaseUrl: "test",
	})

	conf.Set("test2", MSAuthConf{
		Token:   "test2",
		Secret:  "test2",
		Context: "test2",
		BaseUrl: "test2",
	})

	returnedConf, err := conf.Get("test")

	if err != nil {
		t.Error("Failed to get conf", err.Error())
	}

	if returnedConf.Context != "test" {
		t.Errorf("Get returned wrong conf expected test and returned %s", returnedConf.Context)
	}
}

func TestGetInvalidParameter(t *testing.T) {
	conf := StaticConfig{}
	conf.Set("test", MSAuthConf{
		Token:   "test",
		Secret:  "test",
		Context: "test",
		BaseUrl: "test",
	})

	returnedConf, err := conf.Get("test2")

	if err == nil || returnedConf.Context != "" {
		t.Errorf("Expected an error but got nil")
	}
}

func TestTryGetConfInEmptyMap(t *testing.T) {
	conf := StaticConfig{}

	returnedConf, err := conf.Get("test")

	if err == nil || returnedConf.Context != "" {
		t.Errorf("Expected an error but got nil")
	}
}
