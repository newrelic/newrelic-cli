FROM ubuntu:22.04

# Run apt but don't delete the cache as we might use it in dependant Dockerfiles
RUN apt update && apt install -y python3 ca-certificates curl
ADD https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/master/files/docker/systemctl3.py /usr/bin/systemctl
RUN chmod +x /usr/bin/systemctl

CMD ["/usr/bin/systemctl"]
