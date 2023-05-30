# Using a multi stage build for a small file size
# Only need to produce an image with just a binary

# Build Stage or builder stage
# builds binary file. The AS after from FROM denotes the stage
FROM golang:1.19.6-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o bankapi main.go
# install curl
RUN apk add curl
# Install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz


# Second Stage
# Run stage
FROM alpine:3.17
WORKDIR /app
# first arg: copy binary produced from builder stage to this new image
# next is the path of the file we want to copy.
# final is arg is the target location where we want to put the final image
COPY --from=builder /app/bankapi .
# copy downloaded migrate binary from builder
# prefix original migrate file to /app since that is the work dir in the builder stage
COPY --from=builder /app/migrate.linux-amd64 ./migrate
# copy .env file
COPY app.env .
# copy start script
COPY start.sh .
# copy wait-for script
COPY wait-for.sh .
# copy migration files to current working dir under the name migrations
COPY pkg/db/migrations ./migrations

EXPOSE 5000
# this gets passed to the "$@" in start.sh
CMD [ "/app/bankapi" ]
# makes the start.sh script the executor
ENTRYPOINT [ "/app/start.sh" ]