# Changelog

## [0.6.2](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.6.1...flagd-proxy/v0.6.2) (2024-05-10)


### 🧹 Chore

* bump go deps to latest ([#1307](https://github.com/open-feature/flagd/issues/1307)) ([004ad08](https://github.com/open-feature/flagd/commit/004ad083dc01538791148d6233e453d2a3009fcd))

## [0.6.1](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.6.0...flagd-proxy/v0.6.1) (2024-04-19)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.9.0 ([#1281](https://github.com/open-feature/flagd/issues/1281)) ([3cfb052](https://github.com/open-feature/flagd/commit/3cfb0523cc857dd2019d712c621afe81c2b41398))

## [0.6.0](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.5.2...flagd-proxy/v0.6.0) (2024-04-10)


### ⚠ BREAKING CHANGES

* allow custom seed when using targetingKey override for fractional op ([#1266](https://github.com/open-feature/flagd/issues/1266))
  * This is a breaking change only to the extent that it changes the assignment of evaluated flag values.
      Previously, flagd's `fractional` op would internally concatenate any specified bucketing property with the `flag-key`.
      This improved apparent "randomness" by reducing the chances that users were assigned a bucket of the same ordinality across multiple flags.
      However, sometimes it's desireable to have such predictibility, so now **flagd will use the bucketing value as is**.
      If you are specifying a bucketing value in a `fractional` rule, and want to maintain the previous assignments, you can do this concatenation manually:
      `{ "var": "user.name" }` => `{"cat": [{ "var": "$flagd.flagKey" }, { "var": "user.name" }]}`.
      This will result in the same assignment as before.
      Please note, that if you do not specify a bucketing key at all (the shorthand version of the `fractional` op), flagd still uses a concatentation of the `flag-key` and `targetingKey` as before; this behavior has not changed.

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.8.2 ([#1255](https://github.com/open-feature/flagd/issues/1255)) ([9005089](https://github.com/open-feature/flagd/commit/9005089b3e7c8ec4c1e52b42a59c0c05983647a2))


### ✨ New Features

* allow custom seed when using targetingKey override for fractional op ([#1266](https://github.com/open-feature/flagd/issues/1266)) ([f62bc72](https://github.com/open-feature/flagd/commit/f62bc721e8ebc07e27fbe7b9ca085a8771295d65))


### 🧹 Chore

* update go deps ([#1279](https://github.com/open-feature/flagd/issues/1279)) ([219789f](https://github.com/open-feature/flagd/commit/219789fca8a929d552e4e8d1f6b6d5cd44505f43))

## [0.5.2](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.5.1...flagd-proxy/v0.5.2) (2024-03-27)


### ✨ New Features

* OFREP support for flagd  ([#1247](https://github.com/open-feature/flagd/issues/1247)) ([9d12fc2](https://github.com/open-feature/flagd/commit/9d12fc20702a86e8385564659be88f07ad36d9e5))

## [0.5.1](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.5.0...flagd-proxy/v0.5.1) (2024-03-15)


### 🐛 Bug Fixes

* update protobuff CVE-2024-24786 ([#1249](https://github.com/open-feature/flagd/issues/1249)) ([fd81c23](https://github.com/open-feature/flagd/commit/fd81c235fb4a09dfc42289ac316ac3a1d7eff58c))


### 🧹 Chore

* move packaging & isolate service implementations  ([#1234](https://github.com/open-feature/flagd/issues/1234)) ([b58fab3](https://github.com/open-feature/flagd/commit/b58fab3df030ef7e9e10eafa7a0141c05aa05bbd))

## [0.5.0](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.4.2...flagd-proxy/v0.5.0) (2024-02-20)


### ⚠ BREAKING CHANGES

* new proto (flagd.sync.v1) for sync sources ([#1214](https://github.com/open-feature/flagd/issues/1214))

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.5 ([#1198](https://github.com/open-feature/flagd/issues/1198)) ([ce38845](https://github.com/open-feature/flagd/commit/ce388458b9c8a686a7b6ff38b532c941d43d842c))


### ✨ New Features

* new proto (flagd.sync.v1) for sync sources ([#1214](https://github.com/open-feature/flagd/issues/1214)) ([544234e](https://github.com/open-feature/flagd/commit/544234ebd9f9be5f54c2865a866575a7869a56c0))


### 🧹 Chore

* **deps:** update golang docker tag to v1.22 ([#1201](https://github.com/open-feature/flagd/issues/1201)) ([d14c69e](https://github.com/open-feature/flagd/commit/d14c69e93e56d32a37b2428f1db2d4ac79563597))

## [0.4.2](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.4.1...flagd-proxy/v0.4.2) (2024-02-05)


### 🐛 Bug Fixes

* add signal handling to SyncFlags grpc ([#1176](https://github.com/open-feature/flagd/issues/1176)) ([5c8ed7c](https://github.com/open-feature/flagd/commit/5c8ed7c6dd29ffe43c1f1f0e2843683570873443))
* **deps:** update module github.com/open-feature/flagd/core to v0.7.4 ([#1119](https://github.com/open-feature/flagd/issues/1119)) ([e998e41](https://github.com/open-feature/flagd/commit/e998e41f7c6fc8007458dff08e66aa19c7b7b0e7))

## [0.4.1](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.4.0...flagd-proxy/v0.4.1) (2024-01-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.3 ([#1104](https://github.com/open-feature/flagd/issues/1104)) ([b6c00c7](https://github.com/open-feature/flagd/commit/b6c00c7615040399b60f9085a8238d417445546d))
* **deps:** update module github.com/spf13/viper to v1.18.2 ([#1069](https://github.com/open-feature/flagd/issues/1069)) ([f0d6206](https://github.com/open-feature/flagd/commit/f0d620698abbde6ef455c2dd64b02a52eac96a89))

## [0.4.0](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.3.2...flagd-proxy/v0.4.0) (2023-12-22)


### ⚠ BREAKING CHANGES

* remove deprecated flags ([#1075](https://github.com/open-feature/flagd/issues/1075))

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.2 ([#1056](https://github.com/open-feature/flagd/issues/1056)) ([81e83ea](https://github.com/open-feature/flagd/commit/81e83ea0a4aa78d853ea7700cb06bb2a0f329619))
* **deps:** update module github.com/spf13/viper to v1.18.0 ([#1060](https://github.com/open-feature/flagd/issues/1060)) ([9dfa689](https://github.com/open-feature/flagd/commit/9dfa6899ed3a25a5c34f8b0ebd152b01b1097dec))


### 🧹 Chore

* refactoring component structure ([#1044](https://github.com/open-feature/flagd/issues/1044)) ([0c7f78a](https://github.com/open-feature/flagd/commit/0c7f78a95fa4ad2a8b2afe2f6023b9c6d4fd48ed))
* remove deprecated flags ([#1075](https://github.com/open-feature/flagd/issues/1075)) ([49f6fe5](https://github.com/open-feature/flagd/commit/49f6fe5679425b31b1e1cf39a2a2e4767b2e1db9))

## [0.3.2](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.3.1...flagd-proxy/v0.3.2) (2023-12-05)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.1 ([#1037](https://github.com/open-feature/flagd/issues/1037)) ([0ed9b68](https://github.com/open-feature/flagd/commit/0ed9b68341d026681c684a726b215ff910fe2a00))

## [0.3.1](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.3.0...flagd-proxy/v0.3.1) (2023-11-28)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.0 ([#1014](https://github.com/open-feature/flagd/issues/1014)) ([deec49e](https://github.com/open-feature/flagd/commit/deec49e99ef52f62adbf278a8f58936acbb86b9d))


### 🔄 Refactoring

* Rename metrics-port to management-port ([#1012](https://github.com/open-feature/flagd/issues/1012)) ([5635e38](https://github.com/open-feature/flagd/commit/5635e38703cae835a53e9cce83d5bc42d00091e2))

## [0.3.0](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.13...flagd-proxy/v0.3.0) (2023-11-15)


### ⚠ BREAKING CHANGES

* OFO APIs were updated to version v1beta1, since they are more stable now. Resources of the alpha versions are no longer supported in flagd or flagd-proxy.

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.8 ([#1006](https://github.com/open-feature/flagd/issues/1006)) ([c9b48bd](https://github.com/open-feature/flagd/commit/c9b48bd0b617f6d3c04c8924b1d6650ba17de81a))


### ✨ New Features

* support OFO v1beta1 API ([#997](https://github.com/open-feature/flagd/issues/997)) ([bb6f5bf](https://github.com/open-feature/flagd/commit/bb6f5bf0fc382ade75d80a34d209beaa2edc459d))


### 🧹 Chore

* move e2e tests to test ([#1005](https://github.com/open-feature/flagd/issues/1005)) ([a94b639](https://github.com/open-feature/flagd/commit/a94b6399e529ca03c6034eb86ec4028d7e8c2a82))

## [0.2.13](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.12...flagd-proxy/v0.2.13) (2023-11-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.7 ([#966](https://github.com/open-feature/flagd/issues/966)) ([c038a3a](https://github.com/open-feature/flagd/commit/c038a3a3700eee82afa3e2cb2484614ec6ed566c))
* **deps:** update module github.com/spf13/cobra to v1.8.0 ([#993](https://github.com/open-feature/flagd/issues/993)) ([05c7870](https://github.com/open-feature/flagd/commit/05c7870cc7662117f85e9c6528508327ae320b83))


### 🔄 Refactoring

* migrate to connectrpc/connect-go ([#990](https://github.com/open-feature/flagd/issues/990)) ([7dd5b2b](https://github.com/open-feature/flagd/commit/7dd5b2b4c284481bcba5a9c45bd6c85ad1dc6d33))

## [0.2.12](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.11...flagd-proxy/v0.2.12) (2023-10-12)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.6 ([#916](https://github.com/open-feature/flagd/issues/916)) ([1f80e4d](https://github.com/open-feature/flagd/commit/1f80e4db9f8d1ba24884a71f2f8d552499ab5fe2))
* **deps:** update module github.com/spf13/viper to v1.17.0 ([#956](https://github.com/open-feature/flagd/issues/956)) ([31d015d](https://github.com/open-feature/flagd/commit/31d015d329ae9c1da3ec13878078371bcbf43fbf))
* **deps:** update module go.uber.org/zap to v1.26.0 ([#917](https://github.com/open-feature/flagd/issues/917)) ([e57e206](https://github.com/open-feature/flagd/commit/e57e206c937d5b11b81d46ee57b3e92cc454dd88))

## [0.2.11](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.10...flagd-proxy/v0.2.11) (2023-09-14)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.5 ([#900](https://github.com/open-feature/flagd/issues/900)) ([c2ddcbf](https://github.com/open-feature/flagd/commit/c2ddcbfe49b8507fe463c11eb2b031bbc331792a))

## [0.2.10](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.9...flagd-proxy/v0.2.10) (2023-09-08)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.4 ([#880](https://github.com/open-feature/flagd/issues/880)) ([ebb543d](https://github.com/open-feature/flagd/commit/ebb543d6eec18134e44ee7fe623fd2a336a1cf8d))
* **deps:** update opentelemetry-go monorepo ([#868](https://github.com/open-feature/flagd/issues/868)) ([d48317f](https://github.com/open-feature/flagd/commit/d48317f61d7db7ba0398dc9ab7cdd174a0b87555))


### 🧹 Chore

* upgrade to go 1.20 ([#891](https://github.com/open-feature/flagd/issues/891)) ([977167f](https://github.com/open-feature/flagd/commit/977167fb8db330b62726097616dcd691267199ad))

## [0.2.9](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.8...flagd-proxy/v0.2.9) (2023-08-30)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.3 ([#794](https://github.com/open-feature/flagd/issues/794)) ([9671964](https://github.com/open-feature/flagd/commit/96719649affeb1f8412e8b25f52d7292281d8230))


### 🧹 Chore

* **deps:** update golang docker tag to v1.21 ([#822](https://github.com/open-feature/flagd/issues/822)) ([effe29d](https://github.com/open-feature/flagd/commit/effe29d50e33e6c06ef40d7f83f1b3f0df6bd1a2))

## [0.2.8](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.7...flagd-proxy/v0.2.8) (2023-08-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.2 ([#779](https://github.com/open-feature/flagd/issues/779)) ([f34de59](https://github.com/open-feature/flagd/commit/f34de59fc8e636be043ce89758950d6ea3fe7376))
* **deps:** update module go.uber.org/zap to v1.25.0 ([#786](https://github.com/open-feature/flagd/issues/786)) ([40d0aa6](https://github.com/open-feature/flagd/commit/40d0aa66cf422db6811206d777b55396a96f330f))

## [0.2.7](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.6...flagd-proxy/v0.2.7) (2023-07-28)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.1 ([#745](https://github.com/open-feature/flagd/issues/745)) ([d290d8f](https://github.com/open-feature/flagd/commit/d290d8fda8aa84ed2db6454fdd26e60b028e3f7f))

## [0.2.6](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.5...flagd-proxy/v0.2.6) (2023-07-27)


### ✨ New Features

* **flagd-proxy:** introduce zero-downtime ([#752](https://github.com/open-feature/flagd/issues/752)) ([ed5e6e5](https://github.com/open-feature/flagd/commit/ed5e6e5f3ee0a923c33dbf1a8bf20f80adec71bd))

## [0.2.5](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.4...flagd-proxy/v0.2.5) (2023-07-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.4 ([#693](https://github.com/open-feature/flagd/issues/693)) ([33705a6](https://github.com/open-feature/flagd/commit/33705a67300ec70760ba0baeb610f5a2e931205f))
* **deps:** update module github.com/spf13/viper to v1.16.0 ([#679](https://github.com/open-feature/flagd/issues/679)) ([798a975](https://github.com/open-feature/flagd/commit/798a975bb1a47420e814b6dd439f1cece1a263e5))


### 🔄 Refactoring

* **flagd-proxy:** update build.Dockerfile with buildkit caching ([#725](https://github.com/open-feature/flagd/issues/725)) ([06f3d2e](https://github.com/open-feature/flagd/commit/06f3d2eecbcff16bcf2fdfcab33b24c9e697e849))
* remove protobuf dependency from eval package ([#701](https://github.com/open-feature/flagd/issues/701)) ([34ffafd](https://github.com/open-feature/flagd/commit/34ffafd9a777da3f11bd3bfa81565e774cc63214))

## [0.2.4](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.3...flagd-proxy/v0.2.4) (2023-06-07)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.3 ([#634](https://github.com/open-feature/flagd/issues/634)) ([1bc7e99](https://github.com/open-feature/flagd/commit/1bc7e99473bc0c7bcacfb40030562e556d3895d6))


### 🧹 Chore

* update otel dependencies ([#649](https://github.com/open-feature/flagd/issues/649)) ([2114e41](https://github.com/open-feature/flagd/commit/2114e41c38951247866c0b408e5f933282902e70))


### ✨ New Features

* telemetry improvements ([#653](https://github.com/open-feature/flagd/issues/653)) ([ea02cba](https://github.com/open-feature/flagd/commit/ea02cba24bde982d55956fe54de1e8f27226bfc6))


### 🔄 Refactoring

* introduce additional linting rules + fix discrepancies ([#616](https://github.com/open-feature/flagd/issues/616)) ([aef0b90](https://github.com/open-feature/flagd/commit/aef0b9042dcbe5b3f9a7e97960b27366fe50adfe))
* introduce isyncstore interface ([#660](https://github.com/open-feature/flagd/issues/660)) ([c0e2fa0](https://github.com/open-feature/flagd/commit/c0e2fa00736d46db98f72114a449b2e2bf998e3d))

## [0.2.3](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.2...flagd-proxy/v0.2.3) (2023-05-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.2 ([#613](https://github.com/open-feature/flagd/issues/613)) ([218f435](https://github.com/open-feature/flagd/commit/218f435f0212fa24483b2af25e184e154e575eb1))
* **deps:** update module github.com/spf13/cobra to v1.7.0 ([#587](https://github.com/open-feature/flagd/issues/587)) ([12b3477](https://github.com/open-feature/flagd/commit/12b34773a68f6ae7e7e605aebc9f7075eb819994))


### ✨ New Features

* Introduce connect traces ([#624](https://github.com/open-feature/flagd/issues/624)) ([28bac6a](https://github.com/open-feature/flagd/commit/28bac6a54aed79cb8d84a147ffea296c36f5bd51))

## [0.2.2](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.1...flagd-proxy/v0.2.2) (2023-04-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.1 ([#579](https://github.com/open-feature/flagd/issues/579)) ([58eed62](https://github.com/open-feature/flagd/commit/58eed62f5021e5c7a01a171067b725bf3ff83965))

## [0.2.1](https://github.com/open-feature/flagd/compare/flagd-proxy/v0.2.0...flagd-proxy/v0.2.1) (2023-04-12)


### ✨ New Features

* flagd OTEL collector ([#586](https://github.com/open-feature/flagd/issues/586)) ([494bec3](https://github.com/open-feature/flagd/commit/494bec33dcc1ddf0fa5cd0866f06265618408f5e))

## [0.2.0](https://github.com/open-feature/flagd/compare/flagd-proxy-v0.1.2...flagd-proxy/v0.2.0) (2023-03-30)


### ⚠ BREAKING CHANGES

* rename `kube-flagd-proxy` to `flagd-proxy` ([#576](https://github.com/open-feature/flagd/issues/576))

### ✨ New Features

* rename `kube-flagd-proxy` to `flagd-proxy` ([#576](https://github.com/open-feature/flagd/issues/576)) ([223de99](https://github.com/open-feature/flagd/commit/223de99ee3efbcd601bf75ab1f6258eeac0c426e))

## Changelog
