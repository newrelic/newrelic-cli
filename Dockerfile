FROM ubuntu:22.04

RUN apt update && apt install -y python3 ca-certificates curl && rm -rf /var/lib/apt/lists/*
ADD https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/master/files/docker/systemctl3.py /usr/bin/systemctl
RUN chmod +x /usr/bin/systemctl

CMD ["/usr/bin/systemctl"]
