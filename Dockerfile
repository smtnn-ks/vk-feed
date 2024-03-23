FROM golang:1.22-alpine3.19 as build
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN ./init-db.sh
RUN go build -o /bin/app ./cmd

FROM scratch 
COPY --from=build /bin/app /bin/app
COPY --from=build /usr/src/app/.env .env
CMD ["/bin/app"]

