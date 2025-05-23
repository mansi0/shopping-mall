FROM golang:alpine 

RUN mkdir /app

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o cart main.go
EXPOSE 8080
CMD [ "/app/cart" ]