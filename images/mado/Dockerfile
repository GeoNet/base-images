FROM ghcr.io/geonet/base-images/fedora:41 AS build

RUN dnf update -yq && dnf install -yq cargo

RUN cargo install --tag v0.2.2 --git https://github.com/akiomik/mado.git mado

FROM ghcr.io/geonet/base-images/static:latest AS runner

COPY --from=build /root/.cargo/bin/mado /usr/bin/mado

CMD [ "mado" ]
