FROM golang:1.23.1
WORKDIR /tmp

RUN apt-get update && apt-get upgrade -y
# install requirements
RUN apt-get install -y clang gcc-aarch64-linux-gnu gcc-mingw-w64-x86-64 xz-utils gcc-arm-linux-gnueabi
# install node
RUN apt-get install -y npm


WORKDIR /opt

## ADD ZIG to use as cross compiler for windows arm64, only needed for CGO
# download from: https://ziglang.org/download/
# zig is only needed to compule with CGO
#RUN wget https://ziglang.org/builds/zig-linux-x86_64-0.14.0-dev.2506+32354d119.tar.xz
#RUN  tar -xJf zig-linux-x86_64-0.14.0-dev.2506+32354d119.tar.xz && \
#     rm zig-linux-x86_64-0.14.0-dev.2506+32354d119.tar.xz && \
#     mv zig-linux-x86_64-0.14.0-dev.2506+32354d119 zig
#
#ENV PATH="/opt/zig:${PATH}"

# install some utilities
RUN apt-get install -y joe bash-completion

# install golangci lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
## install go-licence-detectors
RUN go install go.elastic.co/go-licence-detector@latest

## install goreleaser oss
RUN wget https://github.com/goreleaser/goreleaser/releases/download/v2.5.0/goreleaser_2.5.0_amd64.deb && \
    dpkg -i goreleaser_2.5.0_amd64.deb


WORKDIR /project