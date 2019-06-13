FROM golang
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR $GOPATH/src/github.com/HDIOES/hundredToOneBackend
COPY Gopkg.toml Gopkg.lock ./
COPY . ./
RUN dep ensure
RUN go install github.com/HDIOES/hundredToOneBackend
RUN cp dbconfig.json $GOPATH/bin/
RUN cp -r migrations/ $GOPATH/bin/
WORKDIR $GOPATH/bin
ENTRYPOINT ["./hundredToOneBackend"]