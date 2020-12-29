//
// Copyright © 2017-2020 Solus Project
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
//

package index

import (
	"os"
	"testing"
)

// Our test files, known to produce a valid delta package
const (
	index = "../testdata/eopkg-index.xml"
)

func TestSave(t *testing.T) {
	i, err := Load(index)
	if err != nil {
		t.Fatalf("Should have loaded successfully: %s", err)
	}
	os.Mkdir("TESTING", 0755)
	if err = i.Save("TESTING"); err != nil {
		t.Errorf("Should have saved successfully: %s", err)
	}
	//os.RemoveAll("TESTING")
}
