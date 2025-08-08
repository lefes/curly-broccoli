FROM golang:1.24 AS build

WORKDIR /src
COPY . .

RUN go mod download
RUN go vet ./...

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/curly-broccoli

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /
CMD ["/app"]