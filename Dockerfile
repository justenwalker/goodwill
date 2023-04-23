ARG GO_VER="1.20.3"
ARG PROTOC_VER="22.2"
ARG PROTOC_GEN_GO_VER="1.30.0"
ARG PROTOC_GEN_GRPC_VER="1.54.0"
ARG JDK_VER="20"
ARG MAVEN_VER="3.9.1"
FROM ubuntu:20.04 as build
ARG GO_VER
ARG PROTOC_VER
ARG PROTOC_GEN_GO_VER
ARG PROTOC_GEN_GRPC_VER
ARG JDK_VER
ARG MAVEN_VER
RUN apt update && apt install -y curl git xzdec wget unzip build-essential
RUN curl -sLo /tmp/jdk.tgz "https://download.java.net/java/GA/jdk${JDK_VER}/bdc68b4b9cbc4ebcb30745c85038d91d/36/GPL/openjdk-${JDK_VER}_linux-aarch64_bin.tar.gz"
RUN curl -sLo /tmp/protoc.zip "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-aarch_64.zip"
RUN curl -sLo /tmp/go.tgz "https://golang.org/dl/go${GO_VER}.linux-amd64.tar.gz"
RUN curl -sLo /tmp/maven.tgz "https://dlcdn.apache.org/maven/maven-3/${MAVEN_VER}/binaries/apache-maven-${MAVEN_VER}-bin.tar.gz"
RUN curl -sLo /tmp/signify.txz "https://github.com/aperezdc/signify/releases/download/v31/signify-31.tar.xz"
RUN tar -xzf /tmp/go.tgz -C /opt
RUN tar -xzf /tmp/jdk.tgz -C /opt
RUN mkdir -p /opt/maven && tar -xzf /tmp/maven.tgz -C /opt/maven --strip-components=1
RUN mkdir -p /opt/signify && tar -xJf /tmp/signify.txz -C /opt/signify --strip-components=1
RUN mkdir -p /opt/protoc && unzip -d /opt/protoc /tmp/protoc.zip
ENV JAVA_HOME=/opt/jdk-${JDK_VER}
ENV PATH=${JAVA_HOME}/bin:/opt/go/bin:/opt/protoc/bin:$PATH
ENV GOPATH=/opt/gopath
ENV GOROOT=/opt/go
RUN mkdir -p $GOPATH
RUN git clone https://github.com/magefile/mage && \
    cd mage && \
    go run bootstrap.go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VER}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GRPC_VER}
RUN cd /opt/signify && make BUNDLED_LIBBSD_VERIFY_GPG=0 BUNDLED_LIBBSD=1

FROM ubuntu:20.04
ARG GO_VER
ARG PROTOC_VER
ARG PROTOC_GEN_GO_VER
ARG PROTOC_GEN_GRPC_VER
ARG JDK_VER
ARG MAVEN_VER
COPY --from=build /opt/go /usr/local/go
COPY --from=build /opt/protoc /usr/local/
COPY --from=build /opt/jdk-${JDK_VER} /opt/jdk-${JDK_VER}
COPY --from=build /opt/maven /opt/maven
COPY --from=build /opt/gopath/bin /usr/local/bin
COPY --from=build /opt/signify/signify /usr/local/bin/signify
COPY scripts/entrypoint.sh /entrypoint.sh
ENV JAVA_HOME=/opt/jdk-${JDK_VER}
ENV PATH=${JAVA_HOME}/bin:/opt/maven/bin:/usr/local/go/bin:/usr/local/bin:$PATH
RUN apt update && apt install -y curl git unzip ca-certificates && \
    rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]