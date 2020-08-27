FROM golang as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -o penpal-api .


FROM ubuntu

WORKDIR /app

COPY --from=builder /build/penpal-api /app/penpal-api
COPY email.tmpl /app/email.tmpl

CMD ["/app/penpal-api"]