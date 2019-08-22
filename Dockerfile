# build stage
FROM golang:alpine AS build-env

# need git to run 'go get'
RUN apk add --no-cache git

WORKDIR /go/src
COPY . .
RUN go get -d -v ./...
RUN go build -v -o server

# final stage
FROM alpine

# create a non-root user
RUN addgroup -g 1001 -S echo && adduser -u 1001 -S echo -G echo
RUN mkdir /app && chown -R echo:echo /app
USER 1001

WORKDIR /app
COPY --from=build-env /go/src/server /app/
COPY --from=build-env /go/src/public /app/public
ENTRYPOINT ./server
EXPOSE 8081