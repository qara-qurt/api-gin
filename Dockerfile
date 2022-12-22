FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .


# # install psql
# RUN apt-get update
# RUN apt-get -y install postgresql-client

# # make wait-for-postgres.sh executable
# RUN chmod +x wait-for-postgres.sh

RUN go build cmd/main.go


EXPOSE 8080

CMD [ "./main" ]