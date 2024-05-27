FROM alpine:latest as builder

ARG NGINX_VERSION=1.23.2
ARG NGINX_RTMP_VERSION=1.2.2


RUN	apk update		&&	\
    apk add				\
    git			\
    gcc			\
    binutils		\
    gmp			\
    isl26			\
    libgomp			\
    libatomic		\
    libgcc			\
    openssl			\
    pkgconf			\
    pkgconfig		\
    mpc1			\
    libstdc++		\
    ca-certificates		\
    libssh2			\
    curl			\
    expat			\
    pcre			\
    musl-dev		\
    libc-dev		\
    pcre-dev		\
    zlib-dev		\
    openssl-dev		\
    curl			\
    make


RUN	cd /tmp/									&&	\
    curl --remote-name http://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz			&&	\
    git clone https://github.com/arut/nginx-rtmp-module.git -b v${NGINX_RTMP_VERSION}

RUN	cd /tmp										&&	\
    tar xzf nginx-${NGINX_VERSION}.tar.gz						&&	\
    cd nginx-${NGINX_VERSION}							&&	\
    ./configure										\
    --sbin-path=/usr/local/sbin/nginx \
    --conf-path=/etc/nginx/nginx.conf \
    --error-log-path=/var/log/nginx/error.log \
    --pid-path=/var/run/nginx/nginx.pid \
    --lock-path=/var/lock/nginx/nginx.lock \
    --http-log-path=/var/log/nginx/access.log \
    --http-client-body-temp-path=/tmp/nginx-client-body \
    --with-http_ssl_module \
    --with-threads \
    --with-ipv6 \
    --add-module=../nginx-rtmp-module					&&	\
    make										&&	\
    make install

FROM alpine:latest
LABEL org.opencontainers.image.authors="geonetci@gns.cri.nz"
RUN apk update		&& \
    apk add			   \
    openssl		   \
    libstdc++	   \
    ca-certificates	   \
    pcre

COPY --from=0 /usr/local/sbin/nginx /usr/local/sbin/nginx
COPY --from=0 /tmp/nginx-rtmp-module/stat.xsl /etc/nginx/conf/stat.xsl

ADD run.sh /

EXPOSE 1935
EXPOSE 80

CMD /run.sh
