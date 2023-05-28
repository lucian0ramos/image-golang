package test

import (
	"testing"

	"github.com/lucian0ramos/image-golang/src/models"
)

func TestValidateCheckManagerErrorParamsOK(t *testing.T) {
	dataOK := models.InputManageErrors{
		DeployId: 222,
		AppName:  "Test",
		Type:     "ci/release",
		Version:  "3.1.0",
		Values:   []byte(`{"key":"value"}`),
	}

	got := models.CheckManageErrorsParams(dataOK)

	if got != nil {
		t.Error("TestValidateCheckManagerErrorParamsOK failed")
	}

}

func TestValidateCheckManagerErrorParamsKO(t *testing.T) {
	dataOK := models.InputManageErrors{
		DeployId: -13,
		AppName:  "Test1",
		Type:     "ci/release2",
		Version:  "23",
		Values:   []byte(`{"key":"value"}`),
	}

	got := models.CheckManageErrorsParams(dataOK)

	if got == nil {
		t.Error("TestValidateCheckManagerErrorParamsKO failed")
	}

}
