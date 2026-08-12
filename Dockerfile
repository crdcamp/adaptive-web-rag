# Build the go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the code and build
COPY . .
RUN go build -o adaptive-web-rag .


# Setup Python
FROM python:3.12-slim

WORKDIR /app

# Set up virtual environment
ENV VIRTUAL_ENV=/opt/venv
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# Copy prompt directory
COPY prompts/ ./prompts

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
RUN python -m playwright install --with-deps chromium

# Copy the built Go executable from the builder stage
COPY --from=builder /app/adaptive-web-rag .

# Copy Python script and directory into the container
COPY crawl_data/ ./crawl_data
COPY crawl.py .

EXPOSE 8082

CMD ["./adaptive-web-rag"]
