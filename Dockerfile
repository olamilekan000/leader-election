FROM golang:1.20.2 as builder

WORKDIR /app

# copy modules manifests
COPY go.mod go.sum ./

# Copy only the Go source files
COPY *.go ./


# build
RUN go build -o lease-based-le

FROM gcr.io/distroless/base:nonroot AS deployable

EXPOSE 8881

COPY --from=builder  /app/lease-based-le .

ENTRYPOINT ["./lease-based-le"]
