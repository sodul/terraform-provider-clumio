#!/bin/bash

VERSION=$1
if [ -z ${VERSION} ]; then
  echo "Usage: ./installer.sh <version>. An example of <version> would be v0.1.0."
  exit 1
fi
OS=$(uname | awk '{print tolower($0)}')
arch=$(uname -m)
case ${arch} in
  x86_64)
    ARCH="amd64"
    ;;
  i386|i686)
    ARCH="386"
    ;;
  armv8l|aarch64)
    ARCH="arm64"
    ;;
  armv7l|arm)
    ARCH="arm"
    ;;
  *)
    echo "Not supported architecture. Exiting..."
    exit 1
esac

VERSION_NUMBER=${VERSION:1}
TF_PLUGIN_DIR="${HOME}/.terraform.d/plugins/clumio.com/providers/clumio/${VERSION_NUMBER}/${OS}_${ARCH}"
if ! mkdir -p ${TF_PLUGIN_DIR}; then
   echo "Error creating directory ${TF_PLUGIN_DIR}. Exiting..."
  exit 1
fi

PROVIDER_NAME=terraform-provider-clumio_${VERSION_NUMBER}_${OS}_${ARCH}
BINARY="https://github.com/clumio-code/terraform-provider-clumio/releases/download/${VERSION}/${PROVIDER_NAME}.zip"

if ! mkdir ${PROVIDER_NAME}; then
  echo "Error creating directory ${PROVIDER_NAME}. Exiting..."
  exit 1
fi

cd ${PROVIDER_NAME}

if command -v curl &> /dev/null; then
  if ! curl -sSL ${BINARY} -o "${PROVIDER_NAME}.zip"; then
    echo "Error downloading ${BINARY}. Exiting..."
    exit 1
  fi
elif command -v wget &> /dev/null; then
  if ! wget -q ${BINARY}; then
    echo "Error downloading ${BINARY}. Exiting..."
    exit 1
  fi
else:
  echo "wget or curl is required to download the binary. Exiting..."
  exit 1
fi

if ! unzip -q ${PROVIDER_NAME}.zip; then
  echo "Error unzipping ${PROVIDER_NAME}.zip. Exiting..."
  exit 1
fi

if ! cp terraform-provider-clumio_${VERSION} ${TF_PLUGIN_DIR}; then
  echo "Error copying terraform-provider-clumio_${VERSION} to ${TF_PLUGIN_DIR}. Exiting..."
  exit 1
fi
cd ..
# cleanup downloaded files
rm -rf ${PROVIDER_NAME}*
if [ -f ${TF_PLUGIN_DIR}/terraform-provider-clumio_${VERSION} ]; then
  echo "Clumio Terraform Provider installed successfully."
fi
