/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package labels

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.BuildResult{}

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	for k, v := range Labels {
		if s, ok := cr.Resolve(k); ok {
			result.Labels = append(result.Labels, libcnb.Label{Key: v, Value: s})
		}
	}

	if s, ok := cr.Resolve("BP_IMAGE_LABELS"); ok {
		err, words := parseLabels(s)
		if err != nil {
			return libcnb.BuildResult{}, err
		}
		for name, val := range words {
			result.Labels = append(result.Labels, libcnb.Label{Key: name, Value: val})
		}
	}

	return result, nil
}

func parseLabels(args string) (error, map[string]string) {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(args, f)

	// create and fill the map
	m := make(map[string]string)
	for _, item := range items {
		item = strings.ReplaceAll(item, "\"", "")
		item = strings.ReplaceAll(item, "'", "")
		x := strings.Split(item, "=")
		if len(x) == 1 { // Missing value
			return fmt.Errorf("unable to parse labels: %s", item), nil
		}

		m[x[0]] = x[1]
	}
	return nil, m
}
