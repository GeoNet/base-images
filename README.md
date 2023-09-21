<!-- generate TOC using `go run sigs.k8s.io/mdtoc@latest --inplace README.md` -->
<!-- toc -->
- [base-images](#base-images)
- [Purpose](#purpose)
- [images](#images)
- [Features](#features)
- [Layout](#layout)
- [Usage](#usage)
  - [Sync an image](#sync-an-image)
    - [Mutably tagged images and auto-updating](#mutably-tagged-images-and-auto-updating)
  - [Build with apko](#build-with-apko)
  - [Build with apko and Melange](#build-with-apko-and-melange)
  - [Build with Docker](#build-with-docker)
- [Tips](#tips)
  - [What is an image ref?](#what-is-an-image-ref)
  - [Obtaining a digest](#obtaining-a-digest)
  - [Tagging built images](#tagging-built-images)
  - [Verifying images](#verifying-images)
  - [See the tree of attached signatures and SBOMs](#see-the-tree-of-attached-signatures-and-sboms)
  - [View the SBOM in the attestation](#view-the-sbom-in-the-attestation)
- [Tooling](#tooling)
- [Patterns for discussion](#patterns-for-discussion)
- [Good reads](#good-reads)
<!-- /toc -->

# base-images

> pipelines to sync, build and scan containers that are vendored or built for GeoNet

# Purpose

- sync to vendor images
- build container images
- run security scans again each destination image in configuration
- separate implementation from the data
  - implementation can be swapped out underneath, if need be and maintain the function

The repo is mostly concerned with base images or images used in build processes.

# images

list of vendored base images

| Image                                              | Description                                                                    |
|----------------------------------------------------|--------------------------------------------------------------------------------|
| ghcr.io/geonet/base-images/alpine                  | standard Alpine base image                                                     |
| ghcr.io/geonet/base-images/datadog/agent           | datadog agent                                                                  |
| ghcr.io/geonet/base-images/debian                  | standard Debian slim image                                                     |
| ghcr.io/geonet/base-images/git                     | Alpine with git installed                                                      |
| ghcr.io/geonet/base-images/git-ssh                 | Alpine with git and ssh installed                                              |
| ghcr.io/geonet/base-images/go                      | a minimal image for building Go applications                                   |
| ghcr.io/geonet/base-images/hadolint/hadolint       | a Dockerfile linter                                                            |
| ghcr.io/geonet/base-images/hashicorp/terraform     | Terraform image                                                                |
| ghcr.io/geonet/base-images/mkdocs_plus             | a custom build image for docs                                                  |
| ghcr.io/geonet/base-images/nginx                   | minimal nginx image                                                            |
| ghcr.io/geonet/base-images/node                    | Alpine based node image                                                        |
| ghcr.io/geonet/base-images/python                  | Alpine and Debian base Python images                                           |
| ghcr.io/geonet/base-images/siderolabs-conform      | conform policy for pull request commits                                        |
| ghcr.io/geonet/base-images/static                  | scratch with timezone and ca-certificate data pre-installed. Great for Go apps |
| ghcr.io/geonet/base-images/texlive                 | a latex image                                                                  |
| ghcr.io/geonet/base-images/ubi8                    | Red Hat Universal Base Image 8                                                 |
| ghcr.io/geonet/base-images/ubi8-minimal            | Red Hat Universal Base Image 8 minimal                                         |
| ghcr.io/geonet/base-images/centos                  | Centos 7 and stream variants available as tags eg :centos7 :stream8            |
| ghcr.io/geonet/base-images/chart-centos7           | Centos 7 base image updated to add extract packages needed for the chart app   |
| ghcr.io/geonet/base-images/curl                    | A minimal image only containing curl                                           |
| ghcr.io/geonet/base-images/owasp/zap2docker-stable | An image to run OWASP's Zed Attack Proxy security web scanner                  |
| ghcr.io/geonet/base-images/alpine-xslt             | Alpine with libxslt for XML parsing                                            |
| ghcr.io/geonet/base-images/alpine-gdal             | Alpine with gdal for raster and vector geospatial data formats                 |
| ghcr.io/geonet/base-images/fedora                  | fedora for build tasks                                                         |
| ghcr.io/geonet/base-images/fedora-coreos           | fedora coreos for edge devices                                                 |
| ghcr.io/geonet/base-images/alpine-iputils          | includes tools like ping                                                       |
| ghcr.io/geonet/base-images/shellcheck              | shellcheck bash scripts                                                                               |


for tags, check [config.yaml](./config.yaml).

# Features

- declarative image sync management, given source and destination
- declarative apko build management
  - with Melange (Alpine APK building) integration
  - generated SBOMs
  - container image signing
- automatic security scanning for each image synced and built
- multiple build modes
  - apko
  - apko + melange
  - docker (fallback)
- automatic trigger of builds, sync and scan every week

# Layout

the structure of the repo is as follows:

- `config.yaml`: define configuration about runtime
- `images/NAME/{images.yaml|Dockerfile,*}`: images to build configurations
- `.github/workflows/{sync,scan,build}.yml`: lifecycle actions

# Usage

## Sync an image

Images are synced by specifying source and destinations like this

```yaml
sync:
  - source: docker.io/alpine:latest
    destination: ghcr.io/somecoolorg/images/alpine:latest
```

Images will only be synced if the digest of the source and destination don't match or if the `sync[].always` key is set to `true`

### Mutably tagged images and auto-updating

Some upstream images have mutable tags which eventually cause some images to become unresolveable.
To mitigate this, a scheduled workflow is created run to automatically produce a PR to fix the resolution.

The workflow also updates an existing branch and PR if one already exists and hasn't been merged.

see: [.github/workflows/update-image-digests.yml](./.github/workflows/update-image-digests.yml)

## Build with apko

Images for building are specified with the source being an apko formatted YAML file and an image destination

```yaml
build:
  - source: ./images/acoolthing/image.yaml
    destination: ghcr.io/somecoolorg/images/acoolthing:latest
```

with `./images/acoolthing/image.yaml` being something like this

```yaml
contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/v3.17/main
    - https://dl-cdn.alpinelinux.org/alpine/v3.17/community
  packages:
    - alpine-base
    - busybox
    - ca-certificates-bundle

entrypoint:
  command: /bin/sh -c

archs:
- x86_64
- aarch64
```

## Build with apko and Melange

Specify a source, destination and as many melangeConfigs

```yaml
build:
  - source: ./images/coolthing/image.yaml
    destination: ghcr.io/somecoolorg/images/coolthing:latest
    melangeConfigs:
      - ./images/coolthing/pkg-hello.yaml
```

with each melange config being a path to the package to build.
An entry in the apko YAML must be set `contents.repositories[last index]` to `@local /github/workspace/packages`, then packages can be installed with `NAME@local`; like this apko configuration

```yaml
contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/v3.17/main
    - https://dl-cdn.alpinelinux.org/alpine/v3.17/community
    - '@local /github/workspace/packages'
  packages:
    - alpine-base
    - busybox
    - ca-certificates-bundle
    - hello@local

entrypoint:
  command: /bin/sh -c

archs:
- x86_64
```

## Build with Docker

Only use this if you have to. It is better to build with apko.

Images can be built with Docker like this

```yaml
build:
  - source: ./images/oldschool/Dockerfile
    destination: ghcr.io/somecoolorg/images/oldschool:latest
```

note that the context of the build runs in the folder which the source is in.

# Tips

## What is an image ref?

IMAGE_REF refers to an image reference in the format of either

- ghcr.io/somecoolorg/images/hello:latest (just a tag)
- ghcr.io/somecoolorg/images/hello@sha256:a61743b19423a01827ba68a1ec81a09d04f84bc69848303905ecbc73824fb88b (just a digest)
- ghcr.io/somecoolorg/images/hello:latest@sha256:a61743b19423a01827ba68a1ec81a09d04f84bc69848303905ecbc73824fb88b (a tag and a digest)

## Obtaining a digest

the digest can be obtained through

```shell
crane digest IMAGE_REF
```

to use it, it can be included when an image ref is formatted like the examples above in _What is an image ref?_.

## Tagging built images

determine the digest of the image, for example given the image is `ghcr.io/somecoolorg/images/hello:latest`

``` shell
crane digest ghcr.io/somecoolorg/images/hello:latest
```

the digest might be `sha256:a61743b19423a01827ba68a1ec81a09d04f84bc69848303905ecbc73824fb88b`.

add a new sync image to the sync key like

```yaml
  - source: ghcr.io/somecoolorg/images/hello@sha256:a61743b19423a01827ba68a1ec81a09d04f84bc69848303905ecbc73824fb88b
    destination: ghcr.io/somecoolorg/images/hello:some-cool-tag
```

## Verifying images

Images are signed using `cosign` with [keyless](https://docs.sigstore.dev/cosign/keyless/), this means that

- signatures depend on the build or key infrastructure
- there's no worry of key rotation

Images are able to be verified through

```shell
cosign verify \
  --certificate-identity-regexp 'https://github.com/GeoNet/base-images/.github/workflows/(sync|build).yml@refs/heads/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  IMAGE_REF
```

please note: the _certificate identity_ field will vary between

- <https://github.com/GeoNet/base-images/.github/workflows/sync.yml@refs/heads/main>
- <https://github.com/GeoNet/base-images/.github/workflows/build.yml@refs/heads/main>

## See the tree of attached signatures and SBOMs

produces a nice and readible tree of signatures, attestations and SBOMs related to the IMAGE_REF

```shell
cosign tree IMAGE_REF
```

## View the SBOM in the attestation

verify the attestation (raw base64 encoded JSON)

```shell
cosign verify-attestation \
  --certificate-identity-regexp 'https://github.com/GeoNet/base-images/.github/workflows/(sync|build).yml@refs/heads/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  IMAGE_REF
```

since SBOMs are the predicate of a signed attestation instead of just uploaded, it requires an extra layer to retrieve their content

```shell
cosign verify-attestation IMAGE_REF --certificate-identity-regexp 'https://github.com/GeoNet/base-images/.github/workflows/(sync|build).yml@refs/heads/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com | jq -r .payload | base64 -d | jq -r .predicate.Data
```

using [sigs.k8s.io/bom](https://sigs.k8s.io/bom), the SBOM attestation can be visualised

```shell
🐚 cosign verify-attestation IMAGE_REF --certificate-identity-regexp 'https://github.com/GeoNet/base-images/.github/workflows/(sync|build).yml@refs/heads/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com | jq -r .payload | base64 -d | jq -r .predicate.Data | bom document outline -

Verification for ghcr.io/GeoNet/base-images/hello:latest --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - Existence of the claims in the transparency log was verified offline
  - The signatures were verified against the specified public key
&{Version:SPDX-2.3 DataLicense:CC0-1.0 ID:SPDXRef-DOCUMENT Name:sbom-sha256:e1097a462313187e2a61ba229179efafa8e62bd6166dc94698c9dcffdff05c9c Namespace:https://spdx.org/spdxdocs/apko/ Creator:{Person: Organization:Chainguard, Inc Tool:[apko (v0.7.3-45-g39b16b5)]} Created:1970-01-01 00:00:00 +0000 UTC LicenseListVersion:3.16 Packages:map[SPDXRef-Package-sha256-63bd196f274790f29c3956c9df4a8525aacec14c1c2ecb684cf7462c53f0471f:0xc000314000] Files:map[] ExternalDocRefs:[]}
PACKAGE:  <nil>
               _
 ___ _ __   __| |_  __
/ __| '_ \ / _` \ \/ /
\__ \ |_) | (_| |>  <
|___/ .__/ \__,_/_/\_\
    |_|

 📂 SPDX Document sbom-sha256:e1097a462313187e2a61ba229179efafa8e62bd6166dc94698c9dcffdff05c9c
  │
  │ 📦 DESCRIBES 1 Packages
  │
  ├ sha256:63bd196f274790f29c3956c9df4a8525aacec14c1c2ecb684cf7462c53f0471f
  │  │ 🔗 2 Relationships
  │  ├ CONTAINS PACKAGE sha256:e1097a462313187e2a61ba229179efafa8e62bd6166dc94698c9dcffdff05c9c@unknown
  │  │  │ 🔗 4 Relationships
  │  │  ├ CONTAINS PACKAGE alpine-baselayout-data@3.4.0-r0
  │  │  ├ CONTAINS PACKAGE ca-certificates-bundle@20220614-r4
  │  │  ├ CONTAINS PACKAGE musl@1.2.3-r4
  │  │  └ CONTAINS PACKAGE hello@2.12-r0
  │  │  │  │ 🔗 3 Relationships
  │  │  │  ├ CONTAINS FILE /usr/bin/hello (/usr/bin/hello)
  │  │  │  ├ CONTAINS FILE /usr/share/info/hello.info (/usr/share/info/hello.info)
  │  │  │  └ CONTAINS FILE /usr/share/man/man1/hello.1 (/usr/share/man/man1/hello.1)
  │  │  │
  │  │
  │  └ GENERATED_FROM PACKAGE github.com/GeoNet/base-images@6d8476f50ccdd7864eb0d46b2036f44516a55336
  │
  └ 📄 DESCRIBES 0 Files
```

# Tooling

| Name    | Description                                                                   | Links                                                              | Related/Alternatives                     |
|---------|-------------------------------------------------------------------------------|--------------------------------------------------------------------|------------------------------------------|
| crane   | an officially supported cli container registry tool from Google               | <https://github.com/google/go-containerregistry/tree/main/cmd/crane> | skopeo                                   |
| yq      | a cli YAML parser                                                             | <https://github.com/mikefarah/yq>                                    | jq                                       |
| cosign  | a cli container image artifact signing utility by Sigstore (Linux Foundation) | <https://github.com/sigstore/cosign>                                 | ...                                      |
| melange | a cli Alpine APK package declarative builder supported by Chainguard          | <https://github.com/chainguard-dev/melange>                          | ...                                      |
| apko    | a cli tool for declaratively building Alpine based container images           | <https://github.com/chainguard-dev/apko>                             | ko (<https://ko.build> - Linux Foundation) |
| docker  | a container ecosystem, primarily for development                              | <https://docker.io>                                                  | podman                                   |
| trivy   | a container image scanner                                                     | <https://github.com/aquasecurity/trivy>                              | clair                                    |
| syft    | a cli tool to generate sboms based on container images and filesystems        | <https://github.com/anchore/syft>                                    |                                          |

# Patterns for discussion

> semi-related topics

- require signed container images for use in production
  - <https://docs.sigstore.dev/cosign/sign>
  - <https://kyverno.io/policies/other/verify_image>
  - <https://docs.sigstore.dev/policy-controller/overview>
- use distroless container images as base images
  - don't ship package managers into production
  - <https://github.com/chainguard-images>
- run containers as non-root users
  - limit workload privilege
  - <https://github.com/chainguard-dev/apko/blob/main/docs/apko_file.md#accounts-top-level-element>
  - set `pod.spec.containers.securityContext.runAsUser`; <https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-pod>
  - `docker run --user 10000 ...` <https://docs.docker.com/engine/reference/commandline/run>
  - <https://docs.aws.amazon.com/config/latest/developerguide/ecs-task-definition-nonroot-user.html>
- use a static tag with digest instead of latest
  - ensure that you are using the version of the image that you expect
  - don't use e.g: `alpine:latest`, instead use an immutable image reference like `alpine:3.17.3@sha256:b6ca290b6b4cdcca5b3db3ffa338ee0285c11744b4a6abaa9627746ee3291d8d` to ensure an expected version is always used
  - version full digests can be manually resolved for consumption using `crane digest alpine:3.17.3` (for example)
  - when using systems like Knative, these digests are automatically resolved
- deploy containers with a read-only root filesystem
  - <https://docs.aws.amazon.com/config/latest/developerguide/ecs-containers-readonly-access.html>
  - `docker run --read-only ...`; <https://docs.docker.com/engine/reference/commandline/run>
  - set `pod.spec.containers.securityContext.readOnlyRootFilesystem` to `true`; <https://kubernetes.io/docs/tasks/configure-pod-container/security-context>
- build Go apps with [ko](https://ko.build)
  - <https://ko.build/advanced/migrating-from-dockerfile>

# Good reads

- <https://docs.sigstore.dev/history>
