# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM  golang:1.16-stretch AS builder
ADD . /src
WORKDIR /src
RUN cd /src && echo $(ls -1 /src)
RUN go mod download
RUN go build -o /bridge e2e/celo-celo/example/main.go

# # final stage
FROM debian:stretch-slim
RUN apt-get -y update && apt-get -y upgrade && apt-get install ca-certificates wget -y

COPY --from=builder /bridge ./
RUN chmod +x ./bridge

ENTRYPOINT ["./bridge"]