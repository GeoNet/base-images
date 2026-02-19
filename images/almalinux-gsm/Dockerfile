FROM ghcr.io/geonet/base-images/almalinux:10.1

LABEL org.opencontainers.image.authors="geonetci@gns.cri.nz"
LABEL org.opencontainers.image.source="https://github.com/GeoNet/base-images"
LABEL org.opencontainers.image.ref.name="ghcr.io/geonet/base-images/almalinux-gsm:10.1"
LABEL org.opencontainers.image.title="AlmaLinux-GSM"
LABEL org.opencontainers.image.description="Gempa GSM on AlmaLinux 10"

RUN dnf -y --refresh upgrade && \
    dnf -y install \
    expect \
    gnupg2 \
    python3 \
    python3-pip \
    python3-wheel \
    tree \
    which

RUN getent group seiscomp >/dev/null || groupadd -g 3560 seiscomp
RUN getent passwd seiscomp >/dev/null || useradd -u 3560 -d /home/seiscomp -g seiscomp seiscomp

USER root
WORKDIR /home/seiscomp

RUN curl -O https://data.gempa.de/gsm/gempa-gsm.tar.gz
RUN tar -xvf gempa-gsm.tar.gz
RUN chown -vR 3560:3560 /home/seiscomp/gsm

USER seiscomp
WORKDIR /home/seiscomp
ENV PATH="/home/seiscomp/.local/bin:$PATH"

RUN mkdir -p /home/seiscomp/data

RUN mkdir -p /home/seiscomp/gsm/log && \
    mkdir -p /home/seiscomp/gsm/packages

RUN python3 -m pip install --user cryptography humanize natsort requests tqdm bs4 gnupg

WORKDIR /home/seiscomp/gsm

COPY --chown=seiscomp:seiscomp \
     --chmod=0660 \
     gsm.conf .

RUN --mount=type=secret,id=GSM_PASSWORD,uid=3560,gid=3560,target=/home/seiscomp/gsm/secret \
    export GSM_PASSWORD=$(cat /home/seiscomp/gsm/secret) && \
    ./gsm --password "${GSM_PASSWORD}" migrate

RUN ls -la /home/seiscomp

RUN --mount=type=secret,id=GSM_PASSWORD,uid=3560,gid=3560,target=/home/seiscomp/gsm/secret \
    export GSM_PASSWORD=$(cat /home/seiscomp/gsm/secret) && \
    ./gsm --password "${GSM_PASSWORD}" update
    
RUN ls -la /home/seiscomp

USER root

RUN --mount=type=secret,id=GSM_PASSWORD,target=/home/seiscomp/gsm/secret \
    export GSM_PASSWORD=$(cat /home/seiscomp/gsm/secret) && \
    ./gsm --password "${GSM_PASSWORD}" install seiscomp

RUN ls -la /home/seiscomp
