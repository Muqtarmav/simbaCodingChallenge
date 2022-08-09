FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build github.com/Muqtarmav/simbaCodingChallenge/Dockerfile

EXPOSE 8080

#CMD [ "/main.go" ]
