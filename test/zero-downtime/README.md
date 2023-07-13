# FlagD Zero downtime test

## How to run

Clone this repository and run the following command to deploy a standalone flagD:

```shell
IMG=your-flagd-image make deploy-dev-env
```

This will create a flagd deployment `flagd-dev` namespace.

To run the test, execute:

```shell
IMG=your-flagd-image IMG_ZD=your-flagd-image2 make run-zd-test
```

Please be aware, you need to build your two custom images with different tags for flagD first.

To build your images using Docker execute:

```shell
docker build . -t image-name:tag -f flagd/build.Dockerfile
```
