FROM golang:1.17.5-buster AS go_build

WORKDIR /kafu

COPY go.mod go.sum /kafu/
RUN go mod download

COPY . /kafu
RUN make

FROM debian:buster-slim

WORKDIR /kafu
RUN apt update && apt install -y ca-certificates
COPY ./entry-point.sh ./entry-point.sh
COPY --from=go_build /kafu/bin/ ./bin/

RUN chmod 755 ./entry-point.sh
ENTRYPOINT [ "./entry-point.sh" ]
