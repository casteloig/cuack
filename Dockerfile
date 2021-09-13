FROM golang:1.17.1-buster

RUN adduser --disabled-password --geco "" cuack
RUN adduser cuack sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN usermod -aG sudo cuack

WORKDIR /go/src/cuack
USER cuack

COPY bin/cuack-ctl ./cuack-ctl

ENTRYPOINT ["/go/src/cuack/cuack-ctl"]
