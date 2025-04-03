FROM ghcr.io/geonet/base-images/almalinux:8.10

LABEL org.opencontainers.image.authors="geonetci@gns.cri.nz"

RUN dnf -y --refresh upgrade && \
    dnf -y install \
    expect \
    python3.11 \
    python3.11-pip \
    python3.11-wheel \
    tree \
    wget \
    which

RUN alternatives --set python /usr/bin/python3.11

RUN getent group seiscomp >/dev/null || groupadd -g 3560 seiscomp
RUN getent passwd seiscomp >/dev/null || useradd -u 3560 -d /home/seiscomp -g seiscomp seiscomp

USER seiscomp
WORKDIR /home/seiscomp
ENV PATH="/home/seiscomp/.local/bin:$PATH"
RUN wget https://data.gempa.de/gsm/gempa-gsm.tar.gz
RUN tar -xvf gempa-gsm.tar.gz 

RUN python3 -m pip install --user cryptography humanize natsort requests tqdm bs4 gnupg

WORKDIR gsm 

RUN mkdir -p log && \
    mkdir -p packages

COPY --chown=seiscomp:seiscomp \
     --chmod=0660 \
     gsm.exp .

RUN chmod u+x ./gsm.exp

RUN ./gsm.exp

RUN ./gsm update
