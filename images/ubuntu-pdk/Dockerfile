FROM ubuntu:latest
RUN apt-get update && apt-get -yq install wget build-essential
RUN . /etc/lsb-release \
    && wget https://apt.puppet.com/puppet-tools-release-${DISTRIB_CODENAME}.deb \
    && dpkg -i puppet-tools-release-${DISTRIB_CODENAME}.deb
RUN apt-get -q update
RUN apt-get -y install pdk
WORKDIR /src
VOLUME [ "/src" ]
