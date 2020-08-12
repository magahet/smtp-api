FROM golang as builder

WORKDIR /tmp/smtp-api

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -o smpt-api .


FROM ubuntu

WORKDIR /app

COPY --from=builder tmp/smtp-api/smtp-api /app/smtp-api
COPY run.sh /app/run.sh

CMD ["/app/run.sh"]