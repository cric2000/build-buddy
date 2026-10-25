FROM golang:1.26.8-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /buildbuddy ./cmd/buildbuddy

FROM scratch
COPY --from=build /buildbuddy /buildbuddy
USER 65532:65532
ENV HOST=0.0.0.0 PORT=8080
EXPOSE 8080
ENTRYPOINT ["/buildbuddy"]
CMD ["serve"]
