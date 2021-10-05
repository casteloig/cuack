FROM golang:1.17.1-buster

RUN adduser --home /home/cuack --disabled-password --geco "" cuack
RUN adduser cuack sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN usermod -aG sudo cuack

WORKDIR /go/src/cuack
USER cuack

RUN mkdir /home/cuack/.config
COPY bin/cuack-ctl ./cuack-ctl

#ENTRYPOINT ["/go/src/cuack/cuack-ctl"]
ENTRYPOINT ["/bin/bash"]
