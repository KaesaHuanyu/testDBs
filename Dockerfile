FROM daocloud.io/golang
MAINTAINER kaesa.li@daocloud.io
WORKDIR /gopath/app
ENV GOPATH /gopath/app
RUN apt-get update
ADD . /gopath/app
RUN go install testdbs
RUN rm -fr /gopath/app/src
EXPOSE 2333
