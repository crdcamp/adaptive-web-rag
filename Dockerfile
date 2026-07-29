FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# RUN pip install -r requirements.txt We'll figure this part out later

COPY . .
RUN go build -o adaptive-web-rag .


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app .

EXPOSE 8082

CMD ["./adaptive-web-rag"]
