FROM amd64/golang:1.22.3-bullseye AS client_builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/scs_client /app/apps/client/main.go

FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=client_builder /app/scs_client /root/scs_client
COPY ./commands /root/commands

RUN chmod +x /root/scs_client

ENTRYPOINT ["/root/scs_client", "/root/commands/commands1.txt", "/root/commands/commands2.txt", "/root/commands/commands3.txt"]
