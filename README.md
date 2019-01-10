# dockertags

[![Docker Repository on Quay](https://quay.io/repository/wantedly/dockertags/status "Docker Repository on Quay")](https://quay.io/repository/wantedly/dockertags)

`dockertags` retrieves and print available Docker image tag list from image repository.

## How to build and install

```bash
$ go get -d github.com/wantedly/dockertags
$ cd $GOPATH/src/github.com/wantedly/dockertags
$ make
$ make install
```

## Usage

```bash
$ dockertags
Usage:
  dockertags IMAGENAME ... [IMAGENAME]

$ dockertags -h
Usage of dockertags:
  -i	show image path in output
  -r string
    	regexp tags (default ".+")

```

Provide Docker image name, then available tags of given image will be shown.

```bash
$ dockertags quay.io/wantedly/dockertags
latest
master

$ dockertags mysql/mysql-server
latest
8.0
8.0.0
5.7
5.7.9
...
5.5.42
```


```bash
$ dockertags -i alpine
alpine:latest
alpine:edge
alpine:3.8
alpine:3.7
alpine:3.6
alpine:3.5
alpine:3.4
alpine:3.3
alpine:3.2
alpine:3.1
alpine:2.7
alpine:2.6

$ dockertags k8s.gcr.io/pause gcr.io/google-containers/busybox
k8s.gcr.io/pause:0.8.0
k8s.gcr.io/pause:1.0
k8s.gcr.io/pause:2.0
k8s.gcr.io/pause:3.0
k8s.gcr.io/pause:3.1
k8s.gcr.io/pause:go
k8s.gcr.io/pause:latest
k8s.gcr.io/pause:test
k8s.gcr.io/pause:test2
gcr.io/google-containers/busybox:1.24
gcr.io/google-containers/busybox:1.27
gcr.io/google-containers/busybox:1.27.2
gcr.io/google-containers/busybox:latest

$ dockertags -r="^(v?)(\d+).(\d+)$|^latest$" k8s.gcr.io/pause gcr.io/google-containers/busybox
k8s.gcr.io/pause:1.0
k8s.gcr.io/pause:2.0
k8s.gcr.io/pause:3.0
k8s.gcr.io/pause:3.1
k8s.gcr.io/pause:latest
gcr.io/google-containers/busybox:1.24
gcr.io/google-containers/busybox:1.27
gcr.io/google-containers/busybox:latest
```

## Supported image registry

- [Quay.io](https://quay.io)
  - To retrieve tags of private image, pass [API access token](http://docs.quay.io/api/) to `QUAYIO_TOKEN` environement variable.
- [Docker Hub](https://hub.docker.com)
  - To retrieve tags of private image, set `DOCKER_USERNAME` and `DOCKER_PASSWORD` for [Docker Hub](hub.docker.com). If they're not set, only tags for public images will be fetched.
- [Google Cloud Registry](https://gcr.io)
  - Private images not tested, but you can try use `GCRIO_TOKEN` variable. I hope it works :)

## License
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
