FROM golang
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get -v github.com/rubenv/sql-migrate/...
WORKDIR $GOPATH/src/github.com/HDIOES/hundredToOneBackend
COPY Gopkg.toml Gopkg.lock ./
COPY . ./
RUN dep ensure
RUN go install github.com/HDIOES/hundredToOneBackend
RUN sql-migrate up -env test
RUN cp dbconfig.json $GOPATH/bin/
WORKDIR $GOPATH/bin
ENTRYPOINT ["./hundredToOneBackend"]