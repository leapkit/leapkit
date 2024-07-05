FROM golang:1.22-alpine as builder
RUN apk --update add build-base

WORKDIR /src/app
ADD go.mod .
RUN go mod download

ADD . .

# Building TailwindCSS with tailo
RUN go run github.com/paganotoni/tailo/cmd/build@a4899cd

# Building the migrate command
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/migrate ./cmd/migrate

# Building the app
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app

FROM alpine
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /bin/

# Copying binaries
COPY --from=builder /src/app/bin/app .
COPY --from=builder /src/app/bin/migrate .

CMD migrate && app
