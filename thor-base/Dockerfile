FROM golang:1.18-alpine AS base

# Set up dependencies
ENV PACKAGES git openssh-client build-base

#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories


# Install dependencies
RUN apk add --update $PACKAGES

# Add source files
RUN mkdir -p ./thor
COPY ./ ./thor/

#RUN go env -w GOPROXY="https://goproxy.cn,direct"


FROM base AS build

RUN  cd thor && make thor && make sender && make query && cp ./bin/thor /usr/bin/ && cp ./bin/sender /usr/bin/ && cp ./bin/query /usr/bin

FROM alpine

WORKDIR /root

COPY  --from=build /usr/bin/thor /usr/bin/thor
COPY  --from=build /usr/bin/sender /usr/bin/sender
COPY  --from=build /usr/bin/query /usr/bin/query

# Add entrypoint script
COPY ./deploy/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod u+x /usr/local/bin/entrypoint.sh

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]

