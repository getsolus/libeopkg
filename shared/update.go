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

package shared

// An Update forms part of a package's history, describing the version, release,
// etc, for each release of the package.
type Update struct {
	Release int    `xml:"release,attr"`
	Type    string `xml:"type,attr,omitempty"`
	Date    string
	Version string
	Comment struct {
		Value string `xml:",cdata"`
	}
	Name struct {
		Value string `xml:",cdata"`
	}
	Email string
}
