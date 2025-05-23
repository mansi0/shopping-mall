FROM golang:alpine 

RUN mkdir /app

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o customer main.go
EXPOSE 8082
CMD [ "/app/customer" ]

# FROM alpine

# RUN mkdir /app

# WORKDIR /app

# COPY --from=builder /app/customer /app

# EXPOSE 8082

# CMD [ "/app/customer" ]
