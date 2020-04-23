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
	"os"

	"github.com/buildpacks/libcnb"
)

type Detect struct{}

func (Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	var pass bool

	for k, _ := range Labels {
		_, ok := os.LookupEnv(k)
		pass = pass || ok
	}

	_, ok := os.LookupEnv("BP_IMAGE_LABELS")
	pass = pass || ok

	if !pass {
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
