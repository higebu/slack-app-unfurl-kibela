FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /work/app
COPY . .
RUN go build -o app .

# runtime image
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /work/app /bin

CMD /bin/app
