FROM golang:alpine
WORKDIR /apiCurrency
COPY ./ ./
RUN go build -o currencyApp main.go
CMD ["./currencyApp"]