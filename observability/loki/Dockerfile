FROM grafana/loki:latest

COPY config.yml /etc/loki/config.yaml

USER root

EXPOSE 3100

CMD ["-config.file=/etc/loki/config.yaml"]