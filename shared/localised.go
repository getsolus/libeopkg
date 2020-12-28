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

import (
	"strings"
)

// LocalisedField is used in various parts of the eopkg metadata to provide
// a field value with an xml:lang attribute describing the language
type LocalisedField struct {
	Value string `xml:",cdata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

// FixMissingLocalLanguage should be used on a set of LocalisedField to restore
// the missing "en" that is required in the very first field set.
func (fields *LocalisedFields) FixMissingLocalLanguage() {
	if fields == nil {
		return
	}
	field := &(*fields)[0]
	if field.Lang == "" {
		field.Lang = "en"
	}
}

// LocalisedFields is a list of more than one translation of the same field
type LocalisedFields []LocalisedField

// Clean removes unnecessary whitespace from each field and then fixes any missing language attributes
func (fields *LocalisedFields) Clean() {
	for i, field := range *fields {
		(*fields)[i].Value = strings.TrimSpace(field.Value)
	}
	fields.FixMissingLocalLanguage()
}
