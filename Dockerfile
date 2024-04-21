FROM ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y \
    curl && \
    rm -rf /var/lib/apt/lists/*

RUN curl -O https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz && \
    rm go1.21.3.linux-amd64.tar.gz

ENV SPAMBOT_APITOKEN=7160401681:AAG0rdibqQ3dftK7mgiVGlN1U2vMTty4iKI \
    PATH=$PATH:/usr/local/go/bin \
    GO111MODULE=on \
    GOPROXY=https://proxy.golang.org

WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src ./

RUN go build -o /app/bot .

CMD ["/app/bot"]

#ENTRYPOINT ["/bin/bash"]
