FROM golang:1.20-buster AS builder

WORKDIR "/go/src/app"

COPY . .

RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o /strmr cmd/main.go

FROM debian:stable-slim as prod

COPY --from=builder "/strmr" .

RUN apt update -y && apt install -y libc6

RUN echo 'DPkg::Post-Invoke {"/bin/rm -f /var/cache/apt/archives/*.deb || true";};' | tee /etc/apt/apt.conf.d/clean

COPY configs /configs

CMD ./strmr
