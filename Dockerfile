FROM golang:1.13
COPY . .
RUN go build -o web-backend

FROM ubuntu:20:04
COPY --from=0 /web-backend .
CMD ["./web-backend"]
