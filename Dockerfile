# Using a multi stage build for a small file size
# Only need to produce an image with just a binary

# Build Stage or builder stage
# builds binary file. The AS after from FROM denotes the stage
FROM golang:1.19.6-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o bankgo main.go


# Second Stage
# Run stage
FROM alpine:3.17
WORKDIR /app
# first arg: copy binary produced from builder stage to this new image
# next is the path of the file we want to copy.
# final is arg is the target location where we want to put the final image
COPY --from=builder /app/bankgo .
# copy .env file
COPY app.env .

EXPOSE 5000
CMD [ "/app/bankgo" ]
