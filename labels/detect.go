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

package labels

import (
	"fmt"

	"github.com/paketo-buildpacks/libpak/v2/log"

	"github.com/buildpacks/libcnb/v2"
	"github.com/paketo-buildpacks/libpak/v2"
)

func NewDetect(l log.Logger) libcnb.DetectFunc {
	return func(context libcnb.DetectContext) (libcnb.DetectResult, error) {
		md, err := libpak.NewBuildModuleMetadata(context.Buildpack.Metadata)
		if err != nil {
			return libcnb.DetectResult{}, fmt.Errorf("unable to create build module metadata\n%w", err)
		}

		cr, err := libpak.NewConfigurationResolver(md)
		if err != nil {
			return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
		}

		var pass bool

		for k := range Labels {
			_, ok := cr.Resolve(k)
			pass = pass || ok
		}

		_, ok := cr.Resolve("BP_IMAGE_LABELS")
		pass = pass || ok

		if !pass {
			l.Body("SKIPPED: No supported environment variables were set")
			return libcnb.DetectResult{Pass: false}, nil
		}

		return libcnb.DetectResult{
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
		}, nil
	}
}
