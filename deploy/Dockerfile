FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /telemetryd ./cmd/telemetryd
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /telemetryd /telemetryd
EXPOSE 8086
ENTRYPOINT ["/telemetryd"]
