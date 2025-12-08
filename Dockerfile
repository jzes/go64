# To Compile Nintendo 64 programs using n64go and EmbeddedGo toolchain

FROM --platform=linux/amd64 ubuntu:24.04

# Dependências básicas
RUN apt update && apt install -y git curl build-essential
# Install go 
RUN curl -LO https://go.dev/dl/go1.24.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz

ENV GOBIN="/usr/local/go/bin"
ENV PATH="/usr/local/go/bin:${PATH}"

# Install n64go and EmbeddedGo toolchain

RUN go install github.com/embeddedgo/dl/go1.24.4-embedded@latest
RUN go install github.com/clktmr/n64/tools/n64go@v0.1.2

# Install gopls for better Go language support
RUN go install golang.org/x/tools/gopls@latest

# Download EmbeddedGo toolchain
RUN go1.24.4-embedded download

ENV GOENV="go.env"
WORKDIR /go64

