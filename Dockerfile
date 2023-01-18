FROM golang:1.17.3-alpine AS build
ARG VERSION_BIN=0.1.1
ARG VERSION=0.1.0
ARG BIN_NAME=mininote
RUN apk add gcc musl-dev
WORKDIR /src
COPY ./src .
RUN go mod vendor
RUN CGO_ENABLED=0 go build -ldflags="-X main.Version=${VERSION_BIN}" -o "${BIN_NAME}" ./
RUN CGO_ENABLED=0 go build -ldflags="-X main.Version=${VERSION}" -o "migrate" ./cmd/migrate/main.go

FROM scratch AS bin
ARG BIN_NAME=mininote
COPY --from=build /src/"${BIN_NAME}" /"${BIN_NAME}"
COPY --from=build /src/migrate /migrate
COPY --from=build /src/migrations /migrations
COPY --from=build /src/.env /.env

ENTRYPOINT [ "/mininote" ]

EXPOSE 9077
