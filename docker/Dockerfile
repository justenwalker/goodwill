FROM golang:1.23

ARG PROTOC_VERSION=26.1
ARG PROTOC_GEN_GO_VERISON=1.33.0
ARG PROTOC_GEN_GRPC_VERISON=1.3.0

RUN apt-get update && apt-get install -y unzip

# By default Intel chipset (x86_64) is assumed but if the host device is an Apple
# silicon (arm) chipset based then a relevant (aarch_64) release file is used.

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERISON}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GRPC_VERISON}

RUN PLATFORM=$(go env GOARCH) \
    export ZIP=x86_64 && \
    if [ ${PLATFORM} = "arm64" ]; then export ZIP=aarch_64; fi && \
    wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-${ZIP}.zip && \
    unzip -o protoc-${PROTOC_VERSION}-linux-${ZIP}.zip -d /usr/local bin/protoc && \
    unzip -o protoc-${PROTOC_VERSION}-linux-${ZIP}.zip -d /usr/local 'include/*'