FROM debian:buster-slim

RUN apt-get update
RUN apt-get install -y \
        sudo \
        net-tools \
        curl \
        wget \
        unzip \
        ca-certificates \
        gosu \
        xz-utils

RUN adduser --disabled-password --geco "" factorio
RUN adduser factorio sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN usermod -aG sudo factorio

RUN mkdir -p /cuackDocker
RUN chown -R factorio /cuackDocker

RUN \
        curl -sSL https://www.factorio.com/get-download/1.1.39/headless/linux64 -o /tmp/factorio_headless_x64_1.1.39.tar.xz \
        && tar -xJf /tmp/factorio_headless_x64_1.1.39.tar.xz -C /opt \
        && bash -c 'mkdir -p {/factorio,/factorio/config,/factorio/mods,/factorio/saves,/factorio/scenarios}' \
        && ln -s /factorio/config /opt/factorio/config \
        && ln -s /factorio/mods /opt/factorio/mods \
        && ln -s /factorio/saves /opt/factorio/saves \
        && ln -s /factorio/scenarios /opt/factorio/scenarios 

EXPOSE 34197/udp
EXPOSE 27015/tcp

USER factorio

RUN sudo chown -R  factorio /opt/factorio/ \
        && sudo chown -R factorio /factorio \
        && sudo mkdir /mnt/cuack \
        && sudo chown -R factorio /mnt/cuack/

COPY start.sh /opt/factorio/bin/x64/start.sh
RUN sudo chmod a+x /opt/factorio/bin/x64/start.sh

CMD ["/opt/factorio/bin/x64/start.sh"]