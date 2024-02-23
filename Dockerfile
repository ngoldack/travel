FROM debian:bookworm AS builder

LABEL org.opencontainers.image.authors="nicolas.goldack@proton.me"

SHELL ["/bin/bash", "-c", "source ~/.bashrc"]

RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    pkg-config \
    libssl-dev \
    libzmq3-dev \
    libzmq5 \
    libzmq5-dev

RUN curl -fsSL https://moonrepo.dev/install/proto.sh | bash -s -- --yes

WORKDIR /build

COPY .prototools .
RUN proto use

RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY go.mod go.sum ./
COPY *.go ./
RUN go mod download

COPY . .

RUN task build

