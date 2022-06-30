FROM golang:latest
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV PORT 8000
EXPOSE 6969
RUN go build
CMD ["./app"]