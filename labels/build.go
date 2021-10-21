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
	"sort"
	"strings"

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
		words, err := ParseLabels(s)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to parse %s\n%w", s, err)
		}

		keys := []string{}
		for k := range words {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			result.Labels = append(result.Labels, libcnb.Label{Key: key, Value: words[key]})
		}
	}

	return result, nil
}

// ReadToNext rune in string consuming the character
//
// It returns a string of read characters, a string of the remaining characters
// and the specific rune that was read.
func ReadToNext(buf string, chars string) (string, string, rune) {
	i := strings.IndexAny(buf, chars)

	if i > 0 && i < len(buf) {
		return buf[0:i], buf[i+1:], rune(buf[i])
	}
	if i > 0 && i == len(buf) {
		return buf[0 : i-1], "", rune(buf[i])
	}
	if i == 0 && i < len(buf) {
		return "", buf[i+1:], rune(buf[i])
	}
	return buf, "", rune(0)
}

// ReadKey from the string
//
// A key is either all words before the equals sign, or a single or
// double quoted word group, again before an equals sign.
//
// If key is quoted and there are characters between the ending quote
// and the equals sign, those are discarded.
//
// It returns the key, a string with the remainder of the characters
// or an error.
func ReadKey(buf string) (string, string, error) {
	item := ""
	needClosingQuote := false
	key, rest, ch := ReadToNext(buf, `"'=`)
	for {
		// unescape, escaped single and double quotes
		if (ch == '"' || ch == '\'') && strings.HasSuffix(key, `\`) {
			key = fmt.Sprintf("%s%c", strings.TrimSuffix(key, `\`), ch)
			ch = '\\'
		}

		item += key

		// end quote
		if (ch == '"' || ch == '\'') && needClosingQuote {
			invalid, rest, _ := ReadToNext(rest, `=`)
			if len(strings.TrimSpace(invalid)) > 0 {
				return "", rest, fmt.Errorf("unable to have characters after a trailing quote")
			}
			return item, rest, nil
		}

		// beginning quote
		if (ch == '"' || ch == '\'') && !needClosingQuote {
			needClosingQuote = true
		}

		// end of key or no more data
		if (ch == '=' || rest == "") && needClosingQuote {
			return "", rest, fmt.Errorf("unable to find a closing quote")
		} else if ch == '=' || rest == "" {
			return item, rest, nil
		}

		key, rest, ch = ReadToNext(rest, `"'=`)
	}
}

// ReadValue from the string
//
// A value is either all words up to the first space character, or a
// single or double quoted word group, again before the first space.
//
// If value is quoted and there are characters between the ending quote
// and the equals sign, those are discarded.
//
// It returns the value, a string with the remainder of the characters
// or an error.
func ReadValue(buf string) (string, string, error) {
	item := ""
	needClosingQuote := false
	value, rest, ch := ReadToNext(buf, `"' `)
	for {
		// unescape, escaped single and double quotes
		if (ch == '"' || ch == '\'') && strings.HasSuffix(value, `\`) {
			value = fmt.Sprintf("%s%c", strings.TrimSuffix(value, `\`), ch)
			ch = '\\'
		}

		item += value

		// end quote
		if (ch == '"' || ch == '\'') && needClosingQuote {
			invalid, rest, _ := ReadToNext(rest, ` `)
			if len(strings.TrimSpace(invalid)) > 0 {
				return "", rest, fmt.Errorf("unable to have characters after a trailing quote")
			}
			return item, rest, nil
		}

		// beginning quote
		if (ch == '"' || ch == '\'') && !needClosingQuote {
			needClosingQuote = true
		}

		// successful end of value
		if ch == ' ' && !needClosingQuote {
			return item, rest, nil
		}

		// embedded space
		if ch == ' ' && needClosingQuote {
			item += " "
		}

		// end of data buffer
		if rest == "" && needClosingQuote {
			return "", rest, fmt.Errorf("unable to find a closing quote")
		} else if rest == "" {
			return item, rest, nil
		}

		value, rest, ch = ReadToNext(rest, `"' `)
	}
}

func ParseLabels(rest string) (map[string]string, error) {
	m := make(map[string]string)

	var (
		key, val string
		err      error
		pos      int
	)

	for {
		before := len(rest)
		key, rest, err = ReadKey(rest)
		pos += before - len(rest)

		if err != nil {
			return nil, fmt.Errorf("unable to read key ending at char %d\n%w", pos-1, err)
		}

		if len(key) == 0 {
			return nil, fmt.Errorf("unable to have empty key ending at char %d", pos-1)
		}

		before = len(rest)
		val, rest, err = ReadValue(rest)
		pos += before - len(rest)

		if err != nil {
			if len(rest) > 0 {
				pos -= 1
			}
			return nil, fmt.Errorf("unable to read value ending at char %d\n%w", pos, err)
		}

		m[key] = val

		if rest == "" {
			return m, nil
		}
	}
}
