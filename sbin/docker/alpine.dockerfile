FROM golang:latest

ARG PROJECT_NAME=go-money

# install dep
RUN go get -u github.com/golang/dep/cmd/dep

# install dependencies
ADD Gopkg.toml Gopkg.lock /go/src/$PROJECT_NAME/
RUN cd /go/src/$PROJECT_NAME && dep ensure -vendor-only

# copy configuration
ADD ./config /etc/$PROJECT_NAME

# add source code
ADD . /go/src/$PROJECT_NAME/
WORKDIR /go/src/$PROJECT_NAME/

EXPOSE 8080
ENTRYPOINT ["go"]
