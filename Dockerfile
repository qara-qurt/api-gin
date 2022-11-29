FROM golang:1.19-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build cmd/main.go


# FROM golang:1.19-alpine

# WORKDIR /app

# COPY --from=builder app/main .
# COPY --from=builder app/configs/config.yml .

EXPOSE 8080

CMD [ "./main" ]
