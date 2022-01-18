FROM  golang:1.17 as builder
RUN apt-get update && apt-get install -y pkg-config protobuf-compiler libprotobuf-dev
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs  | bash -s -- -y
ENV PATH=$PATH:/root/.cargo/bin
WORKDIR /code
COPY . .
RUN make


FROM debian:buster-slim
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update \
	&& apt-get install -y \
	ca-certificates \
	tzdata \
	&& apt-get clean autoclean \
	&& apt-get autoremove --yes \
	&& rm -rf /var/lib/{apt,dpkg,cache,log}

COPY --from=builder /code/bin/linux/influxd /usr/bin/influxd
COPY docker/influxd/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["influxd"]
EXPOSE 8086
