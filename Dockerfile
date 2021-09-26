FROM alpine:3.14.2
EXPOSE 80

RUN apk add --no-cache git make musl-dev go

ENV GO111MODULE=on
ENV GOFLAGS="-mod=vendor"
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR $GOPATH

RUN go get -u github.com/jeffory/s3-policy-tester
CMD ["s3-policy-tester", "run"]