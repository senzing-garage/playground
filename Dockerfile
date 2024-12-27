# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.23.4-bookworm
ARG IMAGE_FINAL=senzing/senzingsdk-runtime-beta:latest

# -----------------------------------------------------------------------------
# Stage: senzingsdk_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS senzingsdk_runtime

# -----------------------------------------------------------------------------
# Stage: builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_BUILDER} AS builder
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/go-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Run as "root" for system installation.

USER root

# Install packages via apt-get.

RUN apt-get update \
 && apt-get -y install \
        libsqlite3-dev \
        python3 \
        python3-dev \
        python3-pip \
        python3-venv \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Create and activate virtual environment.

RUN python3 -m venv /app/venv
ENV PATH="/app/venv/bin:$PATH"

# Install packages via PIP.

COPY requirements.txt .
RUN pip3 install --upgrade pip \
 && pip3 install -r requirements.txt \
 && rm requirements.txt

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/playground

# Copy files from prior stage.

COPY --from=senzingsdk_runtime  "/opt/senzing/er/lib/"   "/opt/senzing/er/lib/"
COPY --from=senzingsdk_runtime  "/opt/senzing/er/sdk/c/" "/opt/senzing/er/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/

# Build go program.

WORKDIR ${GOPATH}/src/playground
RUN make build-with-libsqlite3

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/playground/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS final
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/playground" \
      Maintainer="support@senzing.com" \
      Version="0.4.6"


ARG BUILD_USER="senzing"
ARG BUILD_UID="1001"
ARG BUILD_GID="101"

HEALTHCHECK CMD ["/app/healthcheck.sh"]
USER root

# Install packages via apt-get.

RUN export STAT_TMP=$(stat --format=%a /tmp) \
 && chmod 777 /tmp \
 && apt-get update \
 && apt-get -y install \
        gnupg2 \
        jq \
        libodbc1 \
        libsqlite3-dev \
        postgresql-client \
        supervisor \
        unixodbc \
 && chmod ${STAT_TMP} /tmp \
 && rm -rf /var/lib/apt/lists/*

# Install Java-11.

RUN mkdir -p /etc/apt/keyrings \
 && wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public > /etc/apt/keyrings/adoptium.asc

RUN echo "deb [signed-by=/etc/apt/keyrings/adoptium.asc] https://packages.adoptium.net/artifactory/deb $(awk -F= '/^VERSION_CODENAME/{print$2}' /etc/os-release) main" >> /etc/apt/sources.list

RUN export STAT_TMP=$(stat --format=%a /tmp) \
 && chmod 777 /tmp \
 && apt-get update \
 && apt-get -y install \
        curl \
        python3-venv \
        temurin-11-jdk \
 && chmod ${STAT_TMP} /tmp \
 && rm -rf /var/lib/apt/lists/*
 
# Install go.

RUN wget -O /tmp/go1.linux-amd64.tar.gz https://go.dev/dl/go1.23.3.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf /tmp/go1.linux-amd64.tar.gz 

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/playground /app/playground
COPY --from=builder /app/venv /app/venv

# Prepare jupyter lab environment.

RUN chmod --recursive 777 /app /examples /tmp 

# Create ${BUILD_USER} user.

RUN useradd --no-log-init --create-home --shell /bin/bash --uid "${BUILD_UID}" --no-user-group "${BUILD_USER}"

# Run as non-root container

USER ${BUILD_USER}
ENV HOME=/home/${BUILD_USER}
WORKDIR ${HOME}

# Activate virtual environment.

ENV VIRTUAL_ENV=/app/venv
ENV PATH="/app/venv/bin:/examples/python:${PATH}"

# Install Jupyter Go Kernel.

ENV GOROOT=/usr/local/go
ENV GOPATH=${HOME}/go
ENV PATH=$PATH:${GOROOT}/bin:${GOPATH}/bin
RUN <<EOF
  echo "NB_USER=${BUILD_USER}" >> .profile
  echo "export PATH=${PATH}" >> .profile
  echo "export GOPATH=${GOPATH}" >> .profile
  echo "export GOROOT=${GOROOT}" >> .profile
EOF

RUN go install github.com/janpfeifer/gonb@latest \
 && go install golang.org/x/tools/cmd/goimports@latest \
 && go install golang.org/x/tools/gopls@latest \
 && gonb --install

# Install go packages for SDK
 
WORKDIR /examples/go
 
RUN go get -u ./... \
 && go get -t -u ./... \
 && go mod tidy

WORKDIR ${HOME}

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/
ENV SENZING_API_SERVER_ALLOWED_ORIGINS='*'
ENV SENZING_API_SERVER_BIND_ADDR='all'
ENV SENZING_API_SERVER_ENABLE_ADMIN='true'
ENV SENZING_API_SERVER_PORT='8250'
ENV SENZING_API_SERVER_SKIP_ENGINE_PRIMING='true'
ENV SENZING_API_SERVER_SKIP_STARTUP_PERF='true'
ENV SENZING_DATA_MART_SQLITE_DATABASE_FILE=/tmp/datamart
ENV SENZING_ENGINE_CONFIGURATION_JSON='{"PIPELINE": {"CONFIGPATH": "/etc/opt/senzing", "LICENSESTRINGBASE64": "", "RESOURCEPATH": "/opt/senzing/er/resources", "SUPPORTPATH": "/opt/senzing/data"}, "SQL": {"CONNECTION": "sqlite3://na:na@nowhere/IN_MEMORY_DB?mode=memory&cache=shared"}}'
ENV SENZING_TOOLS_ENABLE_ALL=true

# Runtime execution.

WORKDIR /app
CMD ["/usr/bin/supervisord", "--nodaemon"]
