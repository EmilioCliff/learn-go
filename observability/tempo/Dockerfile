FROM grafana/tempo:latest

COPY config.yml /etc/tempo/config.yml

USER root

EXPOSE 3200
EXPOSE 4319

CMD ["--config.file=/etc/tempo/config.yml"]