#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

GH_REF=$(git rev-parse HEAD)
PUSH_TO_PROJECT="europe-west4-docker.pkg.dev/prj-zen-c-artifact-reg-5bhv/garm"

VERSION=$(git describe --tags --match='v[0-9]*' --always)
AZURE_REF=v0.1.0
OPENSTACK_REF=v0.1.0
LXD_REF=v0.1.0
INCUS_REF=v0.1.0
AWS_REF=v0.1.0
GCP_REF=v0.1.0
EQUINIX_REF=v0.1.0
K8S_REF=v0.3.2
docker buildx build \
  --provenance=false \
  --platform linux/amd64 \
  --label "org.opencontainers.image.source=https://github.com/cloudbase/garm/tree/${GH_REF}" \
  --label "org.opencontainers.image.description=GARM ${GH_REF}" \
  --label "org.opencontainers.image.licenses=Apache 2.0" \
  --build-arg="GARM_REF=${GH_REF}" \
  --build-arg="AZURE_REF=${AZURE_REF}" \
  --build-arg="OPENSTACK_REF=${OPENSTACK_REF}" \
  --build-arg="LXD_REF=${LXD_REF}" \
  --build-arg="INCUS_REF=${INCUS_REF}" \
  --build-arg="AWS_REF=${AWS_REF}" \
  --build-arg="GCP_REF=${GCP_REF}" \
  --build-arg="EQUINIX_REF=${EQUINIX_REF}" \
  --build-arg="K8S_REF=${K8S_REF}" \
  -t ${PUSH_TO_PROJECT}/garm-release:"${GH_REF}" \
  -t ${PUSH_TO_PROJECT}/garm-release:"main-${GH_REF}" \
  --push .
