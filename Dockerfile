FROM golang

ADD . /go/src/ewa

RUN go get -v github.com/spf13/cobra/cobra
RUN go get -v github.com/spf13/viper
RUN go get -v github.com/boltdb/bolt

RUN go install ewa

#ENTRYPOINT /go/bin/ewa


