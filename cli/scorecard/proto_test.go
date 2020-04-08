// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scorecard

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/forseti-security/config-validator/pkg/api/validator"
)

func jsonToInterface(jsonStr string) map[string]interface{} {
	var interfaceVar map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &interfaceVar)
	return interfaceVar
}

func TestDataTypeTransformation(t *testing.T) {
	fileContent, err := ioutil.ReadFile(testRoot + "/shared/iam_policy_audit_logs.json")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	asset := jsonToInterface(string(fileContent))
	wantedName := "//cloudresourcemanager.googleapis.com/projects/23456"

	pbAsset := &validator.Asset{}
	err = protoViaJSON(asset, pbAsset)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	t.Run("protoViaJSON - CAI asset with unknown fieldsto Proto", func(t *testing.T) {
		if pbAsset.Name != wantedName {
			t.Errorf("wanted %s pbAsset.Name, got %s", wantedName, pbAsset.Name)
		}
	})
	t.Run("interfaceViaJSON", func(t *testing.T) {
		var gotInterface interface{}
		gotInterface, err = interfaceViaJSON(pbAsset)
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		gotName := gotInterface.(map[string]interface{})["name"]
		if gotName != wantedName {
			t.Errorf("wanted %s, got %s", wantedName, gotName)
		}
	})
	t.Run("stringViaJSON", func(t *testing.T) {
		gotStr, err := stringViaJSON(pbAsset)
		wantedStr := `{"name":"//cloudresourcemanager.googleapis.com/projects/23456","assetType":"cloudresourcemanager.googleapis.com/Project","iamPolicy":{"version":1,"bindings":[{"role":"roles/owner","members":["user:user@example.com"]}],"etag":"WwAA1Aaa/BA="},"ancestors":["projects/1234","organizations/56789"]}`
		if err != nil {
			t.Fatal("unexpected error", err)
		}
		if gotStr != wantedStr {
			t.Errorf("wanted %s, got %s", wantedStr, gotStr)
		}
	})
}