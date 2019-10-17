FROM golang:alpine

ARG APPNAME="asira_lender"

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev;
#  tzdata wget gcc libc-dev make openssl py-pip;

RUN go get -u github.com/golang/dep/cmd/dep

CMD if [ "${APPENV}" = "staging" ] || [ "${APPENV}" = "production" ] ; then \
        cp deploy/conf.yaml config.yaml \
    elif [ "${ENV}" = "dev" ] ; then \
        cp deploy/dev-config.yaml config.yaml ; \
    fi \
    && dep ensure -v \
    && go build -v -o $GOPATH/bin/"${APPNAME}" \
    && "${APPNAME}" run \
    && "${APPNAME}" migrate up \
EXPOSE 8000