FROM golang:1.5
COPY . /go/src/skipper
WORKDIR /go/src/skipper
RUN go-wrapper download
RUN go-wrapper install
RUN go get ./... && go install ./...
