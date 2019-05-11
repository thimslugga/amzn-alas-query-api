FROM golang:alpine as build

RUN apk --update add git upx ca-certificates

# Setup work env
RUN mkdir /app /tmp/gocode
ADD . /app/
WORKDIR /app

# Required envs for GO
ENV GOPATH=/tmp/gocode
ENV GOOS=linux
ENV GOARCH=amd64

# Disable CGO so we can use a scratch container
ENV CGO_ENABLED=0

RUN go get -d ./...
RUN go build -o /app/alas-query-api .
RUN upx /app/alas-query-api


# Use a scratch container so nothing but the app is present
FROM scratch

# To ensure the aws-sdk can use a shared credentials file (should only be used for testing)
ENV HOME /

# Copy ca-certificates from build for aws cert verification
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from build
COPY --from=build /app/alas-query-api /app/alas-query-api

EXPOSE 8443

# Start the server
CMD ["/app/alas-query-api"]
