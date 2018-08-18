SMTP API
========

This provides a web service that sends email via SMTP.

Authentication and encryption are not included.

# Install

## Go Get

```bash
go get github.com/magahet/smtp-api
```

# Run

```bash
smtp-api <smtp server> <optional port>
```

# Usage

## Curl

### Form data

```bash
curl -d "from=from@example.com&subject=testing&text=sometext&to=to@example.com" -X POST http://servername/message
```

### JSON data

```bash
curl -d '{"from": "from@example.com", "subject": "testing", "text": "some text", "to": ["to@example.com"]}' -H "Content-Type: application/json" -X POST  http://servername/message
```