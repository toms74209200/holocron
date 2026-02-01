FROM golang:1.25-bookworm AS builder

WORKDIR /server

COPY ./server/go.mod go.mod
COPY ./server/go.sum go.sum
COPY ./server/main.go main.go
COPY ./database/ /database
COPY ./server/oapi-codegen.yaml oapi-codegen.yaml
COPY ./spec/openapi.yml /spec/openapi.yml
COPY ./server/internal/ internal/
COPY ./server/Makefile Makefile

RUN go mod download
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
    && make generate \
    && CGO_ENABLED=0 go build -o server

FROM gcr.io/distroless/static-debian12:nonroot

EXPOSE 8080

COPY --from=builder --chown=nonroot:nonroot /server/server /
CMD [ "/server" ]