# Build stage
FROM golang:latest as build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . /src

RUN apt-get update && \
  apt-get upgrade -y ca-certificates

RUN go test .

RUN STATIC="-extldflags '-static'" \
    STATICENV="CGO_ENABLED=0 GOOS=linux GOARCH=amd64" \
    go build -o /src/kalamar-plugin-externalresources .

# Main stage
FROM busybox:glibc

EXPOSE 5722

WORKDIR /

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /src/kalamar-plugin-externalresources /kalamar-plugin-externalresources
COPY --from=build /src/templates /templates

ENTRYPOINT [ "./kalamar-plugin-externalresources" ]

LABEL maintainer="korap@ids-mannheim.de"
LABEL description="Docker Image for Kalamar-Plugin-ExternalResources, a frontend plugin to link texts to external resources"
LABEL repository="https://github.com/KorAP/..."

# docker build -f Dockerfile -t korap/kalamar-plugin-externalresources .
# docker run --rm --network host -v ${PWD}/db/:/db/:z -v ${PWD}/.env:/.env korap/kalamar-plugin-externalresources