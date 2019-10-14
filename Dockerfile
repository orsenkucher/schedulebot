FROM golang:1.12 as build
WORKDIR /go/src/app
COPY . .
ENV GO111MODULE on
RUN go build -v -o /app main.go

FROM gcr.io/distroless/base
COPY --from=build /app /app
CMD ["/app"]