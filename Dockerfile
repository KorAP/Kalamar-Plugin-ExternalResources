# Build stage
FROM golang:latest as build

RUN apt-get update && \
  apt-get upgrade -y ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . /src

RUN CGO_ENABLED=0 go test .

# Build static
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v \
    -ldflags "-extldflags '-static' -s -w" \
    --trimpath \
    -o /src/external-big .

FROM gruebel/upx:latest as upx

COPY --from=build /src/external-big /external-big

# Compress the binary and copy it to final image
RUN upx --best --lzma -o /external /external-big

# Main stage
FROM scratch AS final

WORKDIR /

EXPOSE 5722

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /src/templates /templates
COPY --from=build /src/i18n      /i18n
COPY --from=upx   /external      /external

ENTRYPOINT [ "/external" ]

LABEL maintainer="korap@ids-mannheim.de"
LABEL description="Docker Image for Kalamar-Plugin-ExternalResources, a frontend plugin to link texts to external resources"
LABEL repository="https://github.com/KorAP/Kalamar-Plugin-ExternalResources"

# docker build -f Dockerfile -t korap/kalamar-plugin-externalresources:latest .
# docker run --rm --network host -v ${PWD}/db/:/db/:z -v ${PWD}/.env:/.env korap/kalamar-plugin-externalresources:latest