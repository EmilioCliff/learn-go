FROM grafana/alloy:latest

USER root

COPY ./config.alloy /etc/alloy/config.alloy

EXPOSE 4317

# ENTRYPOINT ["/bin/sh", "-c", "exec /bin/alloy run /etc/alloy/config.alloy"]