FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o avito-shop-service cmd/main.go

FROM alpine
WORKDIR /etc/avito-shop-service
ENV PATH=/etc/avito-shop-service:${PATH}
COPY --from=build /src/avito-shop-service .

ENTRYPOINT ["avito-shop-service"]