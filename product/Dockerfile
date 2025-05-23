FROM golang:alpine 

RUN mkdir /app

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o product main.go
EXPOSE 8081
CMD [ "/app/product" ]