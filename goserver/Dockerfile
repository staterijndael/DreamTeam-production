FROM veselisk/dt-deps AS builder

WORKDIR /app
COPY . .

RUN go install dt/rpc/services/errors/error_helper
RUN go install dt/notification/notifier
RUN go generate ./... > gen.log
RUN wire

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/dt

FROM scratch
WORKDIR /app
COPY --from=builder /app /app

EXPOSE 9003

CMD ["/app/dt"]
