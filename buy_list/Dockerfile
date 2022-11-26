FROM golang:1.19.3-alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# RUN mkdir /buy_list_bot
WORKDIR /buy_list

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main    

EXPOSE 8080

CMD ["/main"]