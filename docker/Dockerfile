FROM golang:1.9.7-stretch

RUN curl -Lo /tmp/glide-install.sh https://glide.sh/get \
 && chmod a+x /tmp/glide-install.sh \
 && /tmp/glide-install.sh \
 && rm -rf /tmp/glide-install.sh

WORKDIR /go/src/github.com/n3integration/terraform-godaddy

CMD glide install \
 && go build -v .
