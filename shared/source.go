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

// Source provides the information relating to the source package within
// each binary package.
// This source identifies one or more packages coming from the same origin,
// i.e they have the same *source name*.
type Source struct {
	Name     string
	Homepage string `xml:"Homepage,omitempty"`
	Packager Packager
}
