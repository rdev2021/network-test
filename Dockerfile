FROM golang:alpine3.19 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY pkg/ ./pkg/
COPY ui/ ./ui/
RUN go build -o /network-test

FROM gcr.io/distroless/base-debian11:nonroot-arm64 AS build-release-stage
WORKDIR /
COPY --from=build-stage /network-test /network-test
COPY --from=build-stage /app/ui/ /ui/
EXPOSE 9091
USER nonroot:nonroot
ENTRYPOINT ["/network-test"]