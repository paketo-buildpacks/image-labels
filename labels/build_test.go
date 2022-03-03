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

package labels_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/image-labels/v4/labels"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx   libcnb.BuildContext
		build labels.Build
	)

	assertMap := func(args string, expected map[string]string, err string) {
		m, e := labels.ParseLabels(args)

		if err != "" {
			Expect(e).To(MatchError(ContainSubstring(err)))
			return
		}

		Expect(e).ToNot(HaveOccurred())

		for k, v := range expected {
			Expect(m).To(HaveKeyWithValue(k, v))
		}
	}

	context("$BP_IMAGE_LABELS", func() {

		it.Before(func() {
			Expect(os.Setenv("BP_IMAGE_LABELS", `alpha=bravo charlie="delta echo" foxtrot='golf hotel'`)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_IMAGE_LABELS")).To(Succeed())
		})

		it("sets image labels", func() {
			Expect(build.Build(ctx)).To(Equal(libcnb.BuildResult{
				Labels: []libcnb.Label{
					{Key: "alpha", Value: "bravo"},
					{Key: "charlie", Value: "delta echo"},
					{Key: "foxtrot", Value: "golf hotel"},
				},
			}))
		})
	})

	for k, v := range labels.Labels {
		context(fmt.Sprintf("$%s", k), func() {

			it.Before(func() {
				Expect(os.Setenv(k, "test-value")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv(k)).To(Succeed())
			})

			it(fmt.Sprintf("passes with $%s", k), func() {
				Expect(build.Build(ctx)).To(Equal(libcnb.BuildResult{
					Labels: []libcnb.Label{
						{Key: v, Value: "test-value"},
					},
				}))
			})
		})
	}

	context("Parser Implementation", func() {
		context("ReadToNext", func() {
			it("reads to a double quote", func() {
				read, left, ch := labels.ReadToNext(`my " string`, `"`)
				Expect(ch).To(Equal('"'))
				Expect(read).To(Equal("my "))
				Expect(left).To(Equal(" string"))
			})

			it("reads to a double quote at beginning", func() {
				read, left, _ := labels.ReadToNext(`"my string`, `"`)
				Expect(read).To(Equal(""))
				Expect(left).To(Equal("my string"))
			})

			it("reads to a double quote at end", func() {
				read, left, _ := labels.ReadToNext(`my string"`, `"`)
				Expect(read).To(Equal("my string"))
				Expect(left).To(Equal(""))
			})

			it("reads to a double quote at end", func() {
				read, left, _ := labels.ReadToNext(`my" string"`, `"`)
				Expect(read).To(Equal("my"))
				Expect(left).To(Equal(" string\""))

				read, left, _ = labels.ReadToNext(left, `"`)
				Expect(read).To(Equal(" string"))
				Expect(left).To(Equal(""))
			})

			it("reads to one of the next characters", func() {
				read, left, _ := labels.ReadToNext(`my string"`, `"s`)
				Expect(read).To(Equal("my "))
				Expect(left).To(Equal(`tring"`))
			})
		})

		context("ReadKey", func() {
			it("Reads a simple key", func() {
				item, left, err := labels.ReadKey(`key=after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("key"))
				Expect(left).To(Equal("after"))

				item, left, err = labels.ReadKey(`key =after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("key "))
				Expect(left).To(Equal("after"))
			})

			it("Reads an key with spaces", func() {
				item, left, err := labels.ReadKey(`dont care=after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("dont care"))
				Expect(left).To(Equal("after"))

				item, left, err = labels.ReadKey(`dont care =after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("dont care "))
				Expect(left).To(Equal("after"))
			})

			it("Reads a quoted key", func() {
				item, left, err := labels.ReadKey(`"also works"=after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("also works"))
				Expect(left).To(Equal("after"))

				item, left, err = labels.ReadKey(`"also works" =after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("also works"))
				Expect(left).To(Equal("after"))
			})

			it("Reads a key with embedded quote", func() {
				item, left, err := labels.ReadKey(`"also\" works"=after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal(`also" works`))
				Expect(left).To(Equal("after"))

				item, left, err = labels.ReadKey(`'also\' works' =after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal(`also' works`))
				Expect(left).To(Equal("after"))
			})

			it("Fails cause of missing quote", func() {
				_, _, err := labels.ReadKey(`'foo bar=`)
				Expect(err).To(MatchError("unable to find a closing quote"))
			})
		})

		context("ReadValue", func() {
			it("Reads a simple value", func() {
				item, left, err := labels.ReadValue(`after another=val`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("after"))
				Expect(left).To(Equal("another=val"))
			})

			it("Reads a quoted value", func() {
				item, left, err := labels.ReadValue(`'foo bar' another=val`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("foo bar"))
				Expect(left).To(Equal("another=val"))
			})

			it("Reads a key with embedded quote", func() {
				item, left, err := labels.ReadValue(`"also\" works" after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal(`also" works`))
				Expect(left).To(Equal("after"))

				item, left, err = labels.ReadValue(`'also\' works' after`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal(`also' works`))
				Expect(left).To(Equal("after"))
			})

			it("Reads value at end", func() {
				item, left, err := labels.ReadValue(`foobar`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("foobar"))
				Expect(left).To(Equal(""))

				item, left, err = labels.ReadValue(`'foo bar'`)
				Expect(err).ToNot(HaveOccurred())
				Expect(item).To(Equal("foo bar"))
				Expect(left).To(Equal(""))
			})

			it("Fails cause of missing quote", func() {
				_, _, err := labels.ReadValue(`'foo bar`)
				Expect(err).To(MatchError("unable to find a closing quote"))
			})
		})
	})

	context("Parses labels", func() {
		it("parses simple label", func() {
			assertMap("key=value", map[string]string{"key": "value"}, "")
		})

		it("parses simple label with quotes", func() {
			assertMap(`key="value"`, map[string]string{"key": "value"}, "")
			assertMap(`key='value'`, map[string]string{"key": "value"}, "")
			assertMap(`"key"='value'`, map[string]string{"key": "value"}, "")
			assertMap(`'key'='value'`, map[string]string{"key": "value"}, "")
		})

		it("parses embedded quotes", func() {
			assertMap(`key='val\'ue'`, map[string]string{"key": "val'ue"}, "")
			assertMap(`key="val\"ue"`, map[string]string{"key": "val\"ue"}, "")
		})

		it("parses two simple labels", func() {
			assertMap(`key='value' foo="bar"`, map[string]string{"key": "value", "foo": "bar"}, "")
		})

		it("doesn't require quotes", func() {
			assertMap(`some-label=(example)value`, map[string]string{"some-label": "(example)value"}, "")
		})

		it("fails missing quote", func() {
			assertMap(`"bad=val`, nil, "unable to read key ending at char 4")
		})

		it("parses unbalanced labels", func() {
			assertMap(`key='value`, nil, "unable to read value ending at char 10")
			assertMap(`key=value"`, nil, "unable to read value ending at char 10")
		})

		it("parses a complex label", func() {
			assertMap(`some-label=(example)value some-label-2=""hi there" test='hi'`,
				nil, "unable to read value ending at char 43\nunable to have characters after a trailing quote")
		})

		it("parses with embedded equal signs", func() {
			assertMap(`foo=bar=baz`,
				map[string]string{"foo": "bar=baz"}, "")
			assertMap(`foo="bar=baz"`,
				map[string]string{"foo": "bar=baz"}, "")
		})

		it("fails on an empty key", func() {
			assertMap(`""=bar`,
				nil, "unable to have empty key ending at char 2")
		})

		it("fails if characters after a quote", func() {
			assertMap(`"foo"junk=bar`,
				nil, "unable to read key ending at char 9\nunable to have characters after a trailing quote")
			assertMap(`foo="bar"junk`,
				nil, "unable to read value ending at char 13\nunable to have characters after a trailing quote")
		})

		it("parses complicated statement", func() {
			assertMap(`some-label=(example)value some-label-2=""hi t""here="" test=\'hi\'`,
				nil, "unable to read value ending at char 43\nunable to have characters after a trailing quote")
		})

		it("parses complicated statement", func() {
			assertMap(`some-label=(example)value some-label-2="hi t"here="" test=\'hi\'`,
				nil, "unable to read value ending at char 52\nunable to have characters after a trailing quote")
		})
	})

}
