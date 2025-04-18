/*
 * Copyright 2018-2025 the original author or authors.
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

	"github.com/buildpacks/libcnb/v2"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak/v2/log"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/image-labels/v4/labels"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		logger log.Logger
		ctx    libcnb.DetectContext
	)

	it("fails without any interesting environment variables set", func() {
		Expect(labels.NewDetect(logger)(ctx)).To(Equal(libcnb.DetectResult{Pass: false}))
	})

	context("$BP_IMAGE_LABELS", func() {
		it.Before(func() {
			t.Setenv("BP_IMAGE_LABELS", "")
		})

		it("passes with $BP_IMAGE_LABELS", func() {
			Expect(labels.NewDetect(logger)(ctx)).To(Equal(libcnb.DetectResult{
				Pass: true,
				Plans: []libcnb.BuildPlan{
					{
						Provides: []libcnb.BuildPlanProvide{
							{Name: "image-labels"},
						},
						Requires: []libcnb.BuildPlanRequire{
							{Name: "image-labels"},
						},
					},
				},
			}))
		})
	})

	for k := range labels.Labels {
		context(fmt.Sprintf("$%s", k), func() {
			it.Before(func() {
				Expect(os.Setenv(k, "")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv(k)).To(Succeed())
			})

			it(fmt.Sprintf("passes with $%s", k), func() {
				Expect(labels.NewDetect(logger)(ctx)).To(Equal(libcnb.DetectResult{
					Pass: true,
					Plans: []libcnb.BuildPlan{
						{
							Provides: []libcnb.BuildPlanProvide{
								{Name: "image-labels"},
							},
							Requires: []libcnb.BuildPlanRequire{
								{Name: "image-labels"},
							},
						},
					},
				}))
			})
		})
	}
}
