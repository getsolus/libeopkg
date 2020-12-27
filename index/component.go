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

const (
	// DefaultMaintainerName is the catch-all name for Solus maintainers
	DefaultMaintainerName = "Solus Team"
	// DefaultMaintainerEmail is the catch-all email for Solus maintainers
	DefaultMaintainerEmail = "copyright@getsol.us"
)

// A Component as seen through the eyes of XML
type Component struct {
	// ID of this component, i.e. "system.base"
	Name string
	// Translated short name
	LocalName shared.LocalisedFields
	// Translated summary
	Summary shared.LocalisedFields
	// Translated description
	Description shared.LocalisedFields
	// Which group this component belongs to
	Group string
	// Maintainer for this component
	Maintainer struct {
		// Name of the component maintainer
		Name string
		// Contact e-mail address of component maintainer
		Email string // Contact e-mail address of component maintainer
	}
}
