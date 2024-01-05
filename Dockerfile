FROM golang:1.21.0 AS build

ENV GOOS=linux \
  GOARCH=amd64 \
  CGO_ENABLED=0

WORKDIR /srv/app/pkg
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. .
RUN go build -o /srv/app/app cmd/main.go

# run stage
FROM golang:1.21.0-alpine as run
ENV GIN_MODE release

ARG PORT
EXPOSE ${PORT}

WORKDIR /srv
RUN mkdir -p /srv
COPY --from=build /srv/app/app /srv/app

CMD ["/srv/app"]