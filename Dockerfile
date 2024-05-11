FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /agent ./cmd/agent

FROM postgres:9.6 AS build-release-stage

COPY --from=build-stage /agent /agent

RUN echo "" > /etc/apt/sources.list.d/pgdg.list
RUN echo 'deb http://archive.debian.org/debian/ stretch main contrib non-free' > /etc/apt/sources.list
RUN apt-get update
RUN apt-get install --yes --force-yes iptables
RUN debconf-set-selections <<EOF
iptables-persistent iptables-persistent/autosave_v4 boolean true
iptables-persistent iptables-persistent/autosave_v6 boolean true
EOF

RUN apt-get install --yes --force-yes iptables-persistent
ENV PG_MAX_WAL_SENDERS 8
ENV PG_WAL_KEEP_SEGMENTS 8

COPY ./docker/scripts/setup-replication.sh /docker-entrypoint-initdb.d/
COPY ./docker/scripts/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint-initdb.d/setup-replication.sh /docker-entrypoint.sh