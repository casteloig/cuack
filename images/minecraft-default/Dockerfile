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

RUN mkdir -p /cuackDocker
RUN chown -R minecraft /cuackDocker

USER minecraft

WORKDIR /home/minecraft

RUN wget https://launcher.mojang.com/v1/objects/a16d67e5807f57fc4e550299cf20226194497dc2/server.jar
COPY eula.txt .
COPY server.properties .
COPY start.sh .

RUN sudo chmod +x start.sh

EXPOSE 25565
EXPOSE 25575

CMD ["./start.sh"]