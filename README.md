# `gcr.io/paketo-buildpacks/image-labels`
The Paketo Buildpack for Image Labels is a Cloud Native Buildpack that enables configuration of labels on the created image.

This buildpack allows for the configuration of both [OCI-specified][o] labels with short environment variable names, as well as arbitrary labels using a space-delimited syntax in a single environment variable.

[o]: https://github.com/opencontainers/image-spec/blob/master/annotations.md#pre-defined-annotation-keys

## Behavior
This buildpack will participate if any of the following conditions are met

* `$BP_IMAGE_LABELS` is set
* `$BP_OCI_AUTHORS` is set
* `$BP_OCI_CREATED` is set
* `$BP_OCI_DESCRIPTION` is set
* `$BP_OCI_DOCUMENTATION` is set
* `$BP_OCI_LICENSES` is set
* `$BP_OCI_REF_NAME` is set
* `$BP_OCI_REVISION` is set
* `$BP_OCI_SOURCE` is set
* `$BP_OCI_TITLE` is set
* `$BP_OCI_URL` is set
* `$BP_OCI_VENDOR` is set
* `$BP_OCI_VERSION` is set

The buildpack will do the following:

* If `$BP_IMAGE_LABELS` is set, it will split the value first along ` `, then along `=`, respecting quotes and set each of the pairs as image labels
* If `$BP_OCI_AUTHORS`  is set, it will set the value as the `org.opencontainers.image.authors` image label
* If `$BP_OCI_CREATED`  is set, it will set the value as the `org.opencontainers.image.created` image label
* If `$BP_OCI_DESCRIPTION`  is set, it will set the value as the `org.opencontainers.image.description` image lable
* If `$BP_OCI_DOCUMENTATION`  is set, it will set the value as the `org.opencontainers.image.documentation` image label
* If `$BP_OCI_LICENSES`  is set, it will set the value as the `org.opencontainers.image.licenses` image label
* If `$BP_OCI_REF_NAME`  is set, it will set the value as the `org.opencontainers.image.ref.name` image label
* If `$BP_OCI_REVISION`  is set, it will set the value as the `org.opencontainers.image.revision` image label
* If `$BP_OCI_SOURCE`  is set, it will set the value as the `org.opencontainers.image.source` image label
* If `$BP_OCI_TITLE`  is set, it will set the value as the `org.opencontainers.image.title` image label
* If `$BP_OCI_URL`  is set, it will set the value as the `org.opencontainers.image.url` image label
* If `$BP_OCI_VENDOR`  is set, it will set the value as the `org.opencontainers.image.vendor` image label
* If `$BP_OCI_VERSION`  is set, it will set the value as the `org.opencontainers.image.version` image label

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_IMAGE_LABELS` | A collection of space-delimited key-value pairs (e.g. `alpha=bravo charlie="delta echo"`) to be set as image labels.  Values containing spaces can be quoted.
| `$BP_OCI_AUTHORS` | The value for the `org.opencontainers.image.authors` image label
| `$BP_OCI_CREATED` | The value for the `org.opencontainers.image.created` image label
| `$BP_OCI_DESCRIPTION` | The value for the `org.opencontainers.image.description` image label
| `$BP_OCI_DOCUMENTATION` | The value for the `org.opencontainers.image.documentation` image label
| `$BP_OCI_LICENSES` | The value for the `org.opencontainers.image.licenses` image label
| `$BP_OCI_REF_NAME` | The value for the `org.opencontainers.image.ref.name` image label
| `$BP_OCI_REVISION` | The value for the `org.opencontainers.image.revision` image label
| `$BP_OCI_SOURCE` | The value for the `org.opencontainers.image.source` image label
| `$BP_OCI_TITLE` | The value for the `org.opencontainers.image.title` image label
| `$BP_OCI_URL` | The value for the `org.opencontainers.image.url` image label
| `$BP_OCI_VENDOR` | The value for the `org.opencontainers.image.vendor` image label
| `$BP_OCI_VERSION` | The value for the `org.opencontainers.image.version` image label


## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

