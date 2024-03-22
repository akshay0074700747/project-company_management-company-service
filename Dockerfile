FROM golang:1.21.5-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o company-service

FROM busybox:latest

WORKDIR /company-service

COPY --from=build /app/cmd/company-service .

COPY --from=build /app/cmd/.env .

EXPOSE 50003

CMD [ "./company-service" ]