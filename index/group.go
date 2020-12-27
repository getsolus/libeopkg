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
	"github.com/getsolus/libeopkg/shared"
)

// A Group as seen through the eyes of XML
type Group struct {
	// ID of this group, i.e. "multimedia"
	Name string
	// Translated short name
	LocalName shared.LocalisedFields
	// Display icon for this Group
	Icon string
}
