FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build 

#EXPOSE 7000

CMD ["go", "run", "./main.go"]
#ENTRYPOINT ["./main"]
