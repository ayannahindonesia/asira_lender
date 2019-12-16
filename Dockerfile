# === Lintasarta Dockerfile ===
FROM golang:alpine  AS build-env

ARG APPNAME="asira_lender"
ARG ENV="staging"

#RUN adduser -D -g '' golang
#USER root

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev;
RUN apk --no-cache add curl
#  tzdata wget gcc libc-dev make openssl py-pip;
RUN go get -u github.com/golang/dep/cmd/dep

RUN cd $GOPATH/src/"${APPNAME}"
RUN cp deploy/conf.yaml config.yaml
RUN dep ensure -v
RUN go build -v -o "${APPNAME}-res"

RUN ls -alh $GOPATH/src/
RUN ls -alh $GOPATH/src/"${APPNAME}"
RUN ls -alh $GOPATH/src/"${APPNAME}"/vendor

FROM alpine

WORKDIR /go/src/
WORKDIR /go/src/migration/

COPY --from=build-env /go/src/asira_lender/asira_lender-res /go/src/asira_lender
COPY --from=build-env /go/src/asira_lender/deploy/conf.yaml /go/src/config.yaml
COPY --from=build-env /go/src/asira_lender/permissions.yaml /go/src/permissions.yaml
COPY --from=build-env /go/src/asira_lender/migration/00001_init_tables.sql /go/src/migration/00001_init_tables.sql
COPY --from=build-env /go/src/asira_lender/migration/image_dummy.txt /go/src/migration/image_dummy.txt
COPY --from=build-env /go/src/asira_lender/migration/migration.go /go/src/migration/migration.go

#ENTRYPOINT /app/asira_lender-res
CMD ["/go/src/asira_lender","run"]

EXPOSE 8000


# === Dev Dockerfile ===
# FROM golang:alpine

# ARG APPNAME="asira_lender"

# ADD . $GOPATH/src/"${APPNAME}"
# WORKDIR $GOPATH/src/"${APPNAME}"

# RUN apk add --update git gcc libc-dev tzdata;
# #  tzdata wget gcc libc-dev make openssl py-pip;

# ENV TZ=Asia/Jakarta

# RUN go get -u github.com/golang/dep/cmd/dep

# CMD if [ "${APPENV}" = "staging" ] || [ "${APPENV}" = "production" ] ; then \
#         cp deploy/conf.yaml config.yaml ; \
#     elif [ "${APPENV}" = "dev" ] ; then \
#         cp deploy/dev-config.yaml config.yaml ; \
#     fi \
#     && dep ensure -v \
#     && go build -v -o $GOPATH/bin/"${APPNAME}" \
#     && "${APPNAME}" run \
#     && "${APPNAME}" migrate up \
# EXPOSE 8000