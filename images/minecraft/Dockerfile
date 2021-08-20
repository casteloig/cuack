FROM adoptopenjdk:16-jre

RUN apt-get update
RUN apt-get install -y \
        sudo \
        net-tools \
        curl \
        wget \
        unzip

RUN adduser --disabled-password --geco "" minecraft
RUN adduser minecraft sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN usermod -aG sudo minecraft

USER minecraft

WORKDIR /home/minecraft

RUN wget https://launcher.mojang.com/v1/objects/35139deedbd5182953cf1caa23835da59ca3d7cd/server.jar
COPY eula.txt .
COPY server.properties .
COPY start.sh .

RUN sudo chmod +x start.sh

EXPOSE 25565

CMD ["./start.sh"]