#!/bin/bash

# oras login quay.io -u mangelajo

set -e

FLASHER_OCI_CONTAINER="${1:-quay.io/jumpstarter-dev/jumpstarter-flasher-test:latest}"
BUNDLE_FILES=${2:-"./test/"}

echo "Building and pushing ${FLASHER_OCI_CONTAINER}"

set -x

cd "${BUNDLE_FILES}"
MANIFESTS=
for file in *.yaml; do
	MANIFESTS="${MANIFESTS} ${file}:application/yaml "
done
DATA_FILES=
while IFS= read -r -d '' file; do
	DATA_FILES="${DATA_FILES} ${file}:application/octet-stream "
done < <(find ./data -type f -not -name .gitkeep -print0)

# shellcheck disable=SC2086  # Variables contain space-separated file arguments that must be expanded
oras push "$FLASHER_OCI_CONTAINER" \
	--artifact-type application/vnd.oci.bundle.v1 \
	${MANIFESTS} \
	${DATA_FILES}
