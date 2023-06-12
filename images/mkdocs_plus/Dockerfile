# TODO tidy
FROM ghcr.io/geonet/base-images/python:3.11.3-bullseye

ENV DRAW_IO_RELEASE="15.8.7"

# hadolint ignore=DL3008
RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
   libx11-xcb1 libatk-bridge2.0 libgtk-3-0 libnss3 libxss1 libasound2\
   libgbm1 \
   graphviz xdg-utils libsecret-1-0 nodejs \
   build-essential sudo \
   python3-pip \
   gconf2 gconf-service libnotify4 xvfb xauth wget \
  && apt-get purge -y --auto-remove apt-transport-https \
  && apt-get clean \
  && rm -rf /tmp/* /var/lib/apt/lists/* /var/cache/apt/archives/*

RUN wget -q -O /tmp/draw-io.deb \
	https://github.com/jgraph/drawio-desktop/releases/download/v${DRAW_IO_RELEASE}/drawio-amd64-${DRAW_IO_RELEASE}.deb \
	&& dpkg -i /tmp/draw-io.deb \
	&& rm -rf /tmp/*

COPY requirements.txt /src/docs/
COPY xvfb_drawio /src/docs/
WORKDIR /src/docs

RUN pip3 install --no-cache-dir -r ./requirements.txt
RUN mv /usr/bin/drawio /usr/bin/drawiobin

# Currently, the electron application does not support
# running without X, run it with our wrapper with xvfb
COPY ./xvfb_drawio /usr/bin/drawio
