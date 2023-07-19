# FlagD Proxy Zero downtime test

## How to run

Clone this repository and run the following command:

```shell
FLAGD_PROXY_IMG=your-flagd-image FLAGD_PROXY_IMG_ZD=your-flagd-second-image ZD_CLIENT_IMG=your-zd-client-image make run-zd-test
```

This will create a flagd-proxy and a job in `flagd-zd-test` namespace,
where the test will be run.

Please be aware, you need to build your custom image for the zd-client
and two images for flagD first.

To build your images using [ko](https://github.com/ko-build/ko),
you need to login to your repository, where the images will be pushed:

```shell
ko login your_repository_server -u username -p password
```

Afterwards, use this command to build flagd-proxy or zd-client:

```shell
KO_DOCKER_REPO=your_repository_server ko build . --bare --tags your-tag
```
