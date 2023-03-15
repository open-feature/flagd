# Snap

flagD can be released on the snapstore as a snap package.
The homepage for the snap is found [here](https://snapcraft.io/flagd/)

## Login

`snapcraft login`

## Build

`snapcraft`

Run this command from `snap` directory.

## Release

```shell
snapcraft upload flagd_<VERSION>_amd64.snap --release=candidate
```

## Promotion

```shell
snapcraft promote flagd --from-channel=candidate --to-channel=stable
```
