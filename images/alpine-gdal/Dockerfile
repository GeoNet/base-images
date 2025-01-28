# Use the official Alpine Linux image as base
FROM alpine:3.21
LABEL org.opencontainers.image.authors="geonetci@gns.cri.nz"
RUN apk update && apk add --no-cache ca-certificates-bundle tzdata gdal gdal-driver-HDF5 gdal-driver-PNG
