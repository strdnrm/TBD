FROM golang:1.19.3-alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

WORKDIR /final_project

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/apiserver/ 

EXPOSE 8080

CMD ["/apiserver"]