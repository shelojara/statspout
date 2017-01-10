FROM golang

ADD . /go

RUN go install github.com/mijara/statspout

EXPOSE 8080

ENTRYPOINT ["statspout"]

