FROM golang:1.22-alpine as builder
RUN apk --update add build-base

WORKDIR /src/app
ADD go.mod .
RUN go mod download

ADD . .
# Installing kit
RUN go install github.com/leapkit/leapkit/kit@v0.0.7

# Building TailwindCSS with tailo
RUN go run github.com/paganotoni/tailo/cmd/build@a4899cd

# Building the app
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app

FROM alpine
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /bin/

# Copying binaries
COPY --from=builder /src/app/bin/app .
COPY --from=builder /src/app/bin/db .

# Copying migrations folder
COPY --from=builder /src/todox/internal/migrations migrations

CMD kit db migrate --migrations.folder=migrations && app
