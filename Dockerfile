FROM jennal/goplay:latest
MAINTAINER jennal <jennalcn@gmail.com>
LABEL maintainer "jennalcn@gmail.com"
ADD . /go/src/github.com/jennal/goplay-master
RUN go install github.com/jennal/goplay-master
ENTRYPOINT ["/go/bin/goplay-master"]
EXPOSE 6812