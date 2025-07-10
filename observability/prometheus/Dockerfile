FROM prom/prometheus:latest

COPY prometheus.yml /etc/prometheus/prometheus.yml

USER root

EXPOSE 9090

CMD ["--config.file=/etc/prometheus/prometheus.yml", "--storage.tsdb.path=/prometheus", "--web.enable-remote-write-receiver", "--enable-feature=native-histograms"]