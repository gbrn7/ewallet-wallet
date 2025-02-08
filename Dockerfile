FROM golang:1.23

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env .

RUN go build -o ewallet-wallet

RUN chmod +x ewallet-wallet

EXPOSE 8081

CMD [ "./ewallet-wallet" ]