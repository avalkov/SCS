FROM amd64/golang:1.22.3-bullseye AS server_builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/scs_server /app/apps/server/main.go

FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=server_builder /app/scs_server /root/scs_server

RUN chmod +x /root/scs_server

CMD ["/root/scs_server"]
