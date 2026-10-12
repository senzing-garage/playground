# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_SENZINGSDK_RUNTIME=senzing/senzingsdk-runtime:4.2.4@sha256:ca9f5e5d2fb4e7bc0e0b65f919e8777f37f3313b5ca8f9aa2d6abdcf9b18d9d7
ARG IMAGE_BUILDER=golang:1.27.1-bookworm@sha256:648f440f42a0958804efb24df176f806f9d353b41f1c0627f666428e40310f6b
ARG IMAGE_FPM=dockter/fpm:1.1.0@sha256:a92ac598d35f1a7a4a659e26bd0e1bd25f317aafdcd4be8bf2795314c421d89b
ARG IMAGE_FINAL=alpine@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# -----------------------------------------------------------------------------
# Stage: senzingsdk_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_SENZINGSDK_RUNTIME} AS senzingsdk_runtime

# -----------------------------------------------------------------------------
# Stage: builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_BUILDER} AS builder
ENV REFRESHED_AT=2026-10-12
LABEL Name="senzing/go-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Build arguments.

ARG PROGRAM_NAME="unknown"
ARG BUILD_VERSION=0.0.0
ARG BUILD_ITERATION=0
ARG GO_PACKAGE_NAME="unknown"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/${GO_PACKAGE_NAME}

# Copy files from prior stage.

COPY --from=senzingsdk_runtime  "/opt/senzing/er/lib/"   "/opt/senzing/er/lib/"
COPY --from=senzingsdk_runtime  "/opt/senzing/er/sdk/c/" "/opt/senzing/er/sdk/c/"

# Build go program.

WORKDIR ${GOPATH}/src/${GO_PACKAGE_NAME}
RUN make linux/amd64

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/${GO_PACKAGE_NAME}/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: fpm
#  - Reference: https://github.com/jordansissel/fpm/blob/master/Dockerfile
#  - FPM: https://fpm.readthedocs.io/en/latest/cli-reference.html
# -----------------------------------------------------------------------------

FROM ${IMAGE_FPM} AS fpm
ENV REFRESHED_AT=2026-10-12
LABEL Name="senzing/fpm-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Use arguments from prior stage.

ARG PROGRAM_NAME
ARG BUILD_VERSION
ARG BUILD_ITERATION
ARG GO_PACKAGE_NAME

# Copy files from prior stage.

COPY --from=builder "/output/linux-amd64/*"    "/output/linux-amd64/"

# Create Linux RPM package.

RUN fpm \
      --input-type dir \
      --output-type rpm \
      --name ${PROGRAM_NAME} \
      --package /output/${PROGRAM_NAME}-${BUILD_VERSION}.rpm \
      --version ${BUILD_VERSION} \
      --iteration ${BUILD_ITERATION} \
      /output/linux-amd64/=/usr/bin

# Create Linux DEB package.

RUN fpm \
      --deb-no-default-config-files \
      --input-type dir \
      --iteration ${BUILD_ITERATION} \
      --name ${PROGRAM_NAME} \
      --output-type deb \
      --package /output/${PROGRAM_NAME}-${BUILD_VERSION}.deb \
      --version ${BUILD_VERSION} \
      /output/linux-amd64/=/usr/bin

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS final
ENV REFRESHED_AT=2026-10-12
LABEL Name="senzing/final-stage" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"
HEALTHCHECK CMD ["/app/healthcheck.sh"]

# Copy files from repository.

COPY ./rootfs /

# Use arguments from prior stage.

ARG PROGRAM_NAME

# Copy files from prior step.

COPY --from=fpm "/output/*"                           "/output/"
COPY --from=fpm "/output/linux-amd64/${PROGRAM_NAME}" "/output/linux-amd64/${PROGRAM_NAME}"

USER 1001
CMD ["/bin/bash"]
