FROM ghcr.io/geonet/base-images/debian:bullseye as builder
RUN apt-get update
RUN apt-get -y install python3-pip
COPY pyproject.toml /src/docs/pyproject.toml
RUN pip install '/src/docs'
WORKDIR /src/docs
