#!/bin/bash
set -e

export HUB=${HUB}
export DEST_REPO=${DEST_REPO}
export DEST_TRANSPORT_USER=${DEST_TRANSPORT_USER:-"docker"}
export DEST_TRANSPORT_PASSWORD=${DEST_TRANSPORT_PASSWORD:-"docker"}

echo "## Check Package Version ##################"
bash --version
skopeo --version

echo "## Login dest TRANSPORT ##################"
set -x
skopeo login -u ${DEST_TRANSPORT_USER} -p ${DEST_TRANSPORT_PASSWORD} ${DEST_REPO}

image-sync