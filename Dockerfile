FROM golang:1.4.1
WORKDIR /go/src/spruce
COPY . .
RUN ls
RUN go get -d -v ./...
#RUN go install -v ./...
RUN go build -v -x /go/src/spruce/spruced.go
RUN ls /go/src/spruce
EXPOSE 6998
EXPOSE 6999
CMD ["/go/src/spruce/spruced"]