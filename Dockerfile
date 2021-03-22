FROM golang:1.16.2

WORKDIR /source
ADD . .

RUN make build

# Second stage - minimal image
FROM alpine

RUN apk update && apk add --no-cache git
COPY --from=0 /source/server /app/server

ENTRYPOINT ["/app/server"]