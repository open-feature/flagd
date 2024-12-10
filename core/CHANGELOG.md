# Changelog

## [0.10.5](https://github.com/open-feature/flagd/compare/core/v0.10.4...core/v0.10.5) (2024-12-10)


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.31.2 ([#1430](https://github.com/open-feature/flagd/issues/1430)) ([0df8622](https://github.com/open-feature/flagd/commit/0df862215563545c33f518ab7a5ad42a19bf6adb))
* **deps:** update kubernetes packages to v0.31.3 ([#1454](https://github.com/open-feature/flagd/issues/1454)) ([f56d7b0](https://github.com/open-feature/flagd/commit/f56d7b043c2d80ae4fe27e996c05a7cc1c2c1b28))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.35.2-20240906125204-0a6a901b42e8.1 ([#1451](https://github.com/open-feature/flagd/issues/1451)) ([8c6d91d](https://github.com/open-feature/flagd/commit/8c6d91d538d226b10cb954c23409902e9d245cda))
* **deps:** update module github.com/diegoholiveira/jsonlogic/v3 to v3.6.0 ([#1460](https://github.com/open-feature/flagd/issues/1460)) ([dbc1da4](https://github.com/open-feature/flagd/commit/dbc1da4ba984c06972b57cf990d1d31c4b8323df))
* **deps:** update module github.com/fsnotify/fsnotify to v1.8.0 ([#1438](https://github.com/open-feature/flagd/issues/1438)) ([949c73b](https://github.com/open-feature/flagd/commit/949c73bd6ebadb30cfa3b7573b43d722f8d2a93d))
* **deps:** update module github.com/stretchr/testify to v1.10.0 ([#1455](https://github.com/open-feature/flagd/issues/1455)) ([8c843df](https://github.com/open-feature/flagd/commit/8c843df7714b1f2d120c5cac8e40c7220cc0c05b))
* **deps:** update module golang.org/x/crypto to v0.29.0 ([#1443](https://github.com/open-feature/flagd/issues/1443)) ([db96dd5](https://github.com/open-feature/flagd/commit/db96dd57b9de032fc4d15931bf907a7ed962f81b))
* **deps:** update module golang.org/x/crypto to v0.30.0 ([#1457](https://github.com/open-feature/flagd/issues/1457)) ([dbdaa19](https://github.com/open-feature/flagd/commit/dbdaa199f0667f16d2a3b91867535ce93e63373c))
* **deps:** update module golang.org/x/mod to v0.22.0 ([#1444](https://github.com/open-feature/flagd/issues/1444)) ([ed064e1](https://github.com/open-feature/flagd/commit/ed064e134fb3a5edb0ec2d976f136af7e94d7f6d))
* **deps:** update module google.golang.org/grpc to v1.68.0 ([#1442](https://github.com/open-feature/flagd/issues/1442)) ([cd27d09](https://github.com/open-feature/flagd/commit/cd27d098e6d8d8b0f681ef42d26dba1ebac67d12))
* **deps:** update module google.golang.org/grpc to v1.68.1 ([#1456](https://github.com/open-feature/flagd/issues/1456)) ([0b6e2a1](https://github.com/open-feature/flagd/commit/0b6e2a1cd64910226d348c921b08a6de8013ac90))
* **deps:** update opentelemetry-go monorepo ([#1447](https://github.com/open-feature/flagd/issues/1447)) ([68b5794](https://github.com/open-feature/flagd/commit/68b5794180da84af9adc1f2cd80f929489969c1c))


### ✨ New Features

* add context-value flag ([#1448](https://github.com/open-feature/flagd/issues/1448)) ([7ca092e](https://github.com/open-feature/flagd/commit/7ca092e478c937eca0c91357394499763545dc1c))

## [0.10.4](https://github.com/open-feature/flagd/compare/core/v0.10.3...core/v0.10.4) (2024-10-28)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.35.1-20240906125204-0a6a901b42e8.1 ([#1420](https://github.com/open-feature/flagd/issues/1420)) ([1f06d5a](https://github.com/open-feature/flagd/commit/1f06d5a1837ea2b753974e96c2a1154d6cb3e582))
* **deps:** update module github.com/prometheus/client_golang to v1.20.5 ([#1425](https://github.com/open-feature/flagd/issues/1425)) ([583ba89](https://github.com/open-feature/flagd/commit/583ba894f2de794b36b6a1cc3bfceb9c46dc9d96))
* **deps:** update module go.uber.org/mock to v0.5.0 ([#1427](https://github.com/open-feature/flagd/issues/1427)) ([0c6fd7f](https://github.com/open-feature/flagd/commit/0c6fd7fa688db992d4e58a202889cbfea07eebf6))
* **deps:** update module gocloud.dev to v0.40.0 ([#1422](https://github.com/open-feature/flagd/issues/1422)) ([e0e4709](https://github.com/open-feature/flagd/commit/e0e4709243d8301bcbb0aaaa309be66944c1d9ed))
* **deps:** update module golang.org/x/crypto to v0.28.0 ([#1416](https://github.com/open-feature/flagd/issues/1416)) ([fb272da](https://github.com/open-feature/flagd/commit/fb272da56e0eba12245309899888c18920b9a200))
* **deps:** update module google.golang.org/grpc to v1.67.1 ([#1415](https://github.com/open-feature/flagd/issues/1415)) ([85a3a6b](https://github.com/open-feature/flagd/commit/85a3a6b46233fcc7cf71a0292b46c82ac8e66d7b))


### ✨ New Features

* added custom grpc resolver ([#1424](https://github.com/open-feature/flagd/issues/1424)) ([e5007e2](https://github.com/open-feature/flagd/commit/e5007e2bcb6f049a3c54e09331065bb9abe215be))
* support azure blob sync ([#1428](https://github.com/open-feature/flagd/issues/1428)) ([5c39cfe](https://github.com/open-feature/flagd/commit/5c39cfe30a3dead4f6db2c6f9ee4c12193cd479b))

## [0.10.3](https://github.com/open-feature/flagd/compare/core/v0.10.2...core/v0.10.3) (2024-09-23)


### 🐛 Bug Fixes

* **deps:** update kubernetes package and controller runtime, fix proto lint ([#1290](https://github.com/open-feature/flagd/issues/1290)) ([94860d6](https://github.com/open-feature/flagd/commit/94860d6ceabe9eb7c1e5dd8ea139a796710d6d8b))
* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20240906125204-0a6a901b42e8.1 ([#1400](https://github.com/open-feature/flagd/issues/1400)) ([954d972](https://github.com/open-feature/flagd/commit/954d97238210f90b650493ae76277d4a8d80788a))
* **deps:** update module connectrpc.com/connect to v1.17.0 ([#1408](https://github.com/open-feature/flagd/issues/1408)) ([e7eb691](https://github.com/open-feature/flagd/commit/e7eb691094dfbf02e37d79c41f60f556415e7640))
* **deps:** update module github.com/prometheus/client_golang to v1.20.3 ([#1384](https://github.com/open-feature/flagd/issues/1384)) ([8fd16b2](https://github.com/open-feature/flagd/commit/8fd16b23b1fa8517128af36b3068ca18ebbad6c3))
* **deps:** update module github.com/prometheus/client_golang to v1.20.4 ([#1406](https://github.com/open-feature/flagd/issues/1406)) ([a0a6426](https://github.com/open-feature/flagd/commit/a0a64269b08251317676075fdea7bc65bea8a8dc))
* **deps:** update module gocloud.dev to v0.39.0 ([#1404](https://github.com/open-feature/flagd/issues/1404)) ([a3184d6](https://github.com/open-feature/flagd/commit/a3184d68413749808709baac47df3bf7400f9cdc))
* **deps:** update module golang.org/x/crypto to v0.27.0 ([#1396](https://github.com/open-feature/flagd/issues/1396)) ([f9a7d10](https://github.com/open-feature/flagd/commit/f9a7d10590d3191ea8eba0dbb340fa94d07026a4))
* **deps:** update module golang.org/x/mod to v0.21.0 ([#1397](https://github.com/open-feature/flagd/issues/1397)) ([1507e19](https://github.com/open-feature/flagd/commit/1507e19e9304bcebfbbe4376f45e9f2e82135fd2))
* **deps:** update module google.golang.org/grpc to v1.66.0 ([#1393](https://github.com/open-feature/flagd/issues/1393)) ([c96e9d7](https://github.com/open-feature/flagd/commit/c96e9d764aa51caf00fbde07cdc7d2de55b98b9e))
* **deps:** update module google.golang.org/grpc to v1.66.1 ([#1402](https://github.com/open-feature/flagd/issues/1402)) ([50c9cd3](https://github.com/open-feature/flagd/commit/50c9cd3ada2f470a22374392a5a152a487636645))
* **deps:** update module google.golang.org/grpc to v1.66.2 ([#1405](https://github.com/open-feature/flagd/issues/1405)) ([69ec28f](https://github.com/open-feature/flagd/commit/69ec28fceb597bdaad63b184943b66ccdb4af0b7))
* **deps:** update module google.golang.org/grpc to v1.67.0 ([#1407](https://github.com/open-feature/flagd/issues/1407)) ([1ad6480](https://github.com/open-feature/flagd/commit/1ad6480a0f37c4677e53065ef455f615b26b1f17))
* **deps:** update opentelemetry-go monorepo ([#1387](https://github.com/open-feature/flagd/issues/1387)) ([22aef5b](https://github.com/open-feature/flagd/commit/22aef5bbf030c619e48fbe22a16d83e071b11902))
* **deps:** update opentelemetry-go monorepo ([#1403](https://github.com/open-feature/flagd/issues/1403)) ([fc4cd3e](https://github.com/open-feature/flagd/commit/fc4cd3e547f4826ea0bb8cc1bb2304807932b4e6))
* remove dep cycle with certreloader ([#1410](https://github.com/open-feature/flagd/issues/1410)) ([5244f6f](https://github.com/open-feature/flagd/commit/5244f6f6c94f310fd80c7ab84942103cc8c18a39))


### ✨ New Features

* add mTLS support to otel exporter ([#1389](https://github.com/open-feature/flagd/issues/1389)) ([8737f53](https://github.com/open-feature/flagd/commit/8737f53444016b114ee4ae52eead0b835af0e200))

## [0.10.2](https://github.com/open-feature/flagd/compare/core/v0.10.1...core/v0.10.2) (2024-08-22)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20240215170432-1e611e2999cc.1 ([#1372](https://github.com/open-feature/flagd/issues/1372)) ([ae24595](https://github.com/open-feature/flagd/commit/ae2459504f7eccafebccec83fa1f72b08f41a978))
* **deps:** update module connectrpc.com/otelconnect to v0.7.1 ([#1367](https://github.com/open-feature/flagd/issues/1367)) ([184915b](https://github.com/open-feature/flagd/commit/184915b31726729e8ed2f7999f338bf4ed684809))
* **deps:** update module github.com/open-feature/open-feature-operator/apis to v0.2.44 ([#1368](https://github.com/open-feature/flagd/issues/1368)) ([0c68726](https://github.com/open-feature/flagd/commit/0c68726bed1cdae07f1b90447818ebbc9dc45caf))
* **deps:** update module golang.org/x/crypto to v0.26.0 ([#1379](https://github.com/open-feature/flagd/issues/1379)) ([05f6658](https://github.com/open-feature/flagd/commit/05f6658e3dc72182adbff9197c8980641af8c53f))
* **deps:** update module golang.org/x/mod to v0.20.0 ([#1377](https://github.com/open-feature/flagd/issues/1377)) ([797d7a4](https://github.com/open-feature/flagd/commit/797d7a4bbafc73e6882e5998df500ae4fe98fbbc))
* **deps:** update module golang.org/x/sync to v0.8.0 ([#1378](https://github.com/open-feature/flagd/issues/1378)) ([4804c17](https://github.com/open-feature/flagd/commit/4804c17a67ea9761079ecade34ccb3446643050b))


### ✨ New Features

* add 'watcher' interface to file sync ([#1365](https://github.com/open-feature/flagd/issues/1365)) ([61fff43](https://github.com/open-feature/flagd/commit/61fff43e288daac88efb127ada20276c01ed5928))
* added new grpc sync config option to allow setting max receive message size. ([#1358](https://github.com/open-feature/flagd/issues/1358)) ([bed077b](https://github.com/open-feature/flagd/commit/bed077bac9da3b6e3bd45ca54046e40a595fcba6))
* Support blob type sources and GCS as an example of such source. ([#1366](https://github.com/open-feature/flagd/issues/1366)) ([21f2c9a](https://github.com/open-feature/flagd/commit/21f2c9a5d64cbfe2fc841080850a2c582e8f4ba6))


### 🧹 Chore

* **deps:** update dependency go to v1.22.6 ([#1297](https://github.com/open-feature/flagd/issues/1297)) ([50b92c1](https://github.com/open-feature/flagd/commit/50b92c17cfd872d3e6b95fef3b3d96444e563715))

## [0.10.1](https://github.com/open-feature/flagd/compare/core/v0.10.0...core/v0.10.1) (2024-07-08)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.4.0-20240215170432-1e611e2999cc.2 ([#1342](https://github.com/open-feature/flagd/issues/1342)) ([efdd921](https://github.com/open-feature/flagd/commit/efdd92139903b89ac986a62ff2cf4f5cfef91cde))
* **deps:** update module golang.org/x/crypto to v0.25.0 ([#1351](https://github.com/open-feature/flagd/issues/1351)) ([450cbc8](https://github.com/open-feature/flagd/commit/450cbc84ca55eef3fccc768003e358a8e589668e))
* **deps:** update module golang.org/x/mod to v0.19.0 ([#1349](https://github.com/open-feature/flagd/issues/1349)) ([6ee89b4](https://github.com/open-feature/flagd/commit/6ee89b44ca4aca8f6236603fc3f969e814907bd6))
* **deps:** update module google.golang.org/grpc to v1.65.0 ([#1346](https://github.com/open-feature/flagd/issues/1346)) ([72a6b87](https://github.com/open-feature/flagd/commit/72a6b876e880ff0b43440d9b63710c7a87536988))
* **deps:** update opentelemetry-go monorepo ([#1347](https://github.com/open-feature/flagd/issues/1347)) ([37fb3cd](https://github.com/open-feature/flagd/commit/37fb3cd81d5436e9d8cd3ea490a3951ae5794130))

## [0.10.0](https://github.com/open-feature/flagd/compare/core/v0.9.3...core/v0.10.0) (2024-06-27)


### ⚠ BREAKING CHANGES

* support emitting errors from the bulk evaluator ([#1338](https://github.com/open-feature/flagd/issues/1338))

### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.4.0-20240215170432-1e611e2999cc.1 ([#1333](https://github.com/open-feature/flagd/issues/1333)) ([494062f](https://github.com/open-feature/flagd/commit/494062fed891fab0fb659352142dbbc97c8f1492))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.34.2-20240215170432-1e611e2999cc.2 ([#1330](https://github.com/open-feature/flagd/issues/1330)) ([32291ad](https://github.com/open-feature/flagd/commit/32291ad93d25d79299a7a02381df70e2719c4fbc))
* **deps:** update module connectrpc.com/connect to v1.16.2 ([#1289](https://github.com/open-feature/flagd/issues/1289)) ([8bacb7c](https://github.com/open-feature/flagd/commit/8bacb7c59c17956dda3cf9d2a7bc6f139885a656))
* **deps:** update module github.com/open-feature/open-feature-operator/apis to v0.2.43 ([#1331](https://github.com/open-feature/flagd/issues/1331)) ([fecd769](https://github.com/open-feature/flagd/commit/fecd769e5f2c4aa7bf1a0fb13f0543e7b8045af8))
* **deps:** update module golang.org/x/crypto to v0.24.0 ([#1335](https://github.com/open-feature/flagd/issues/1335)) ([2a31a17](https://github.com/open-feature/flagd/commit/2a31a1740303991412e0169e50a064823cce0560))
* **deps:** update module golang.org/x/mod to v0.18.0 ([#1336](https://github.com/open-feature/flagd/issues/1336)) ([5fa83f7](https://github.com/open-feature/flagd/commit/5fa83f7ab266320d8d8f1388a6c2f2cac922275a))
* **deps:** update opentelemetry-go monorepo ([#1314](https://github.com/open-feature/flagd/issues/1314)) ([e9f1a7a](https://github.com/open-feature/flagd/commit/e9f1a7a04828f36691e694375b3c665140bc7dee))
* readable error messages  ([#1325](https://github.com/open-feature/flagd/issues/1325)) ([7ff33ef](https://github.com/open-feature/flagd/commit/7ff33effcc47e31c5b7fdc33385d8128db2163fc))


### ✨ New Features

* add mandatory flags property in bulk response ([#1339](https://github.com/open-feature/flagd/issues/1339)) ([b20266e](https://github.com/open-feature/flagd/commit/b20266ed5e0c16bf14769f300297f7f2b0ab2dcd))
* support emitting errors from the bulk evaluator ([#1338](https://github.com/open-feature/flagd/issues/1338)) ([b9c099c](https://github.com/open-feature/flagd/commit/b9c099cb7fa002a509a82c81b467f5e784c27e82))
* support relative weighting for fractional evaluation ([#1313](https://github.com/open-feature/flagd/issues/1313)) ([f82c094](https://github.com/open-feature/flagd/commit/f82c094f5c47f99ef37b7392bfd39cec3ec7ba51))

## [0.9.3](https://github.com/open-feature/flagd/compare/core/v0.9.2...core/v0.9.3) (2024-06-06)


### 🐛 Bug Fixes

* fixes store merge when selector is used ([#1322](https://github.com/open-feature/flagd/issues/1322)) ([ed5025d](https://github.com/open-feature/flagd/commit/ed5025d8f28fdf92a5c7dceaec7fd6df7f979e3b))


### 🧹 Chore

* adapt telemetry setup error handling ([#1315](https://github.com/open-feature/flagd/issues/1315)) ([20bcb78](https://github.com/open-feature/flagd/commit/20bcb78d11dbb16aab2b14d5869bb990a0f7bca5))

## [0.9.2](https://github.com/open-feature/flagd/compare/core/v0.9.1...core/v0.9.2) (2024-05-10)


### ✨ New Features

* improve error log and add flag disabled handling for ofrep ([#1306](https://github.com/open-feature/flagd/issues/1306)) ([39ae4fe](https://github.com/open-feature/flagd/commit/39ae4fe11380af5c6e23c4aaae45b5ec17cf32d6))


### 🧹 Chore

* bump go deps to latest ([#1307](https://github.com/open-feature/flagd/issues/1307)) ([004ad08](https://github.com/open-feature/flagd/commit/004ad083dc01538791148d6233e453d2a3009fcd))

## [0.9.1](https://github.com/open-feature/flagd/compare/core/v0.9.0...core/v0.9.1) (2024-04-19)


### 🐛 Bug Fixes

* missing/nil custom variables in fractional operator ([#1295](https://github.com/open-feature/flagd/issues/1295)) ([418c5cd](https://github.com/open-feature/flagd/commit/418c5cd7c07fad61674a751872a1256b5062799c))


### ✨ New Features

* move json logic operator registration to resolver ([#1291](https://github.com/open-feature/flagd/issues/1291)) ([b473457](https://github.com/open-feature/flagd/commit/b473457ddff28789fee1eeb6704491b6aa3525e3))

## [0.9.0](https://github.com/open-feature/flagd/compare/core/v0.8.2...core/v0.9.0) (2024-04-10)


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

### ✨ New Features

* allow custom seed when using targetingKey override for fractional op ([#1266](https://github.com/open-feature/flagd/issues/1266)) ([f62bc72](https://github.com/open-feature/flagd/commit/f62bc721e8ebc07e27fbe7b9ca085a8771295d65))


### 🧹 Chore

* refactor evaluation core ([#1259](https://github.com/open-feature/flagd/issues/1259)) ([0e6604c](https://github.com/open-feature/flagd/commit/0e6604cd038dc13d7d40e622523320bf03efbcd0))
* update go deps ([#1279](https://github.com/open-feature/flagd/issues/1279)) ([219789f](https://github.com/open-feature/flagd/commit/219789fca8a929d552e4e8d1f6b6d5cd44505f43))
* wire evaluation ctx to store methods ([#1273](https://github.com/open-feature/flagd/issues/1273)) ([0075932](https://github.com/open-feature/flagd/commit/00759322594f309ca9236156f296805a09f5f9fe))

## [0.8.2](https://github.com/open-feature/flagd/compare/core/v0.8.1...core/v0.8.2) (2024-03-27)


### ✨ New Features

* OFREP support for flagd  ([#1247](https://github.com/open-feature/flagd/issues/1247)) ([9d12fc2](https://github.com/open-feature/flagd/commit/9d12fc20702a86e8385564659be88f07ad36d9e5))

## [0.8.1](https://github.com/open-feature/flagd/compare/core/v0.8.0...core/v0.8.1) (2024-03-15)


### 🐛 Bug Fixes

* occasional panic when watched YAML files change ([#1246](https://github.com/open-feature/flagd/issues/1246)) ([6249d12](https://github.com/open-feature/flagd/commit/6249d12ec452073ed881b6e5faf716332c7f132a))
* update protobuff CVE-2024-24786 ([#1249](https://github.com/open-feature/flagd/issues/1249)) ([fd81c23](https://github.com/open-feature/flagd/commit/fd81c235fb4a09dfc42289ac316ac3a1d7eff58c))


### 🧹 Chore

* move packaging & isolate service implementations  ([#1234](https://github.com/open-feature/flagd/issues/1234)) ([b58fab3](https://github.com/open-feature/flagd/commit/b58fab3df030ef7e9e10eafa7a0141c05aa05bbd))

## [0.8.0](https://github.com/open-feature/flagd/compare/core/v0.7.5...core/v0.8.0) (2024-02-20)


### ⚠ BREAKING CHANGES

* new proto (flagd.sync.v1) for sync sources ([#1214](https://github.com/open-feature/flagd/issues/1214))

### 🐛 Bug Fixes

* **deps:** update github.com/open-feature/flagd-schemas digest to 8c72c14 ([#1212](https://github.com/open-feature/flagd/issues/1212)) ([4add9fd](https://github.com/open-feature/flagd/commit/4add9fd1c47e3e3dea818a6b262273fadb7edb81))
* **deps:** update kubernetes packages to v0.29.2 ([#1213](https://github.com/open-feature/flagd/issues/1213)) ([b0c805f](https://github.com/open-feature/flagd/commit/b0c805f7f58979f927e60c22c94c0448af459c7d))
* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.15.0-20240215170432-1e611e2999cc.1 ([#1219](https://github.com/open-feature/flagd/issues/1219)) ([4c4f08a](https://github.com/open-feature/flagd/commit/4c4f08afabf7d646973768500db028ab0b4c7d68))
* **deps:** update module golang.org/x/crypto to v0.19.0 ([#1203](https://github.com/open-feature/flagd/issues/1203)) ([f0ff317](https://github.com/open-feature/flagd/commit/f0ff3177f67c832d62694cdf44b766344da5483f))
* **deps:** update module golang.org/x/mod to v0.15.0 ([#1202](https://github.com/open-feature/flagd/issues/1202)) ([6ca8e6d](https://github.com/open-feature/flagd/commit/6ca8e6d33f6646698605fb4b5b99f8a3ee1ddbed))
* **deps:** update module golang.org/x/net to v0.21.0 ([#1204](https://github.com/open-feature/flagd/issues/1204)) ([bccf365](https://github.com/open-feature/flagd/commit/bccf365fa2e5f443208ec70b1244bdb4f07ced04))
* **deps:** update module google.golang.org/grpc to v1.61.1 ([#1210](https://github.com/open-feature/flagd/issues/1210)) ([10cc63e](https://github.com/open-feature/flagd/commit/10cc63e7992b4ae8d861f5296afcc78417e645cb))
* **deps:** update opentelemetry-go monorepo ([#1199](https://github.com/open-feature/flagd/issues/1199)) ([422ebaa](https://github.com/open-feature/flagd/commit/422ebaa30b8bb0246bcdf8c0cc2be0a5870eb9e9))


### ✨ New Features

* new proto (flagd.sync.v1) for sync sources ([#1214](https://github.com/open-feature/flagd/issues/1214)) ([544234e](https://github.com/open-feature/flagd/commit/544234ebd9f9be5f54c2865a866575a7869a56c0))

## [0.7.5](https://github.com/open-feature/flagd/compare/core/v0.7.4...core/v0.7.5) (2024-02-05)


### 🐛 Bug Fixes

* add signal handling to SyncFlags grpc ([#1176](https://github.com/open-feature/flagd/issues/1176)) ([5c8ed7c](https://github.com/open-feature/flagd/commit/5c8ed7c6dd29ffe43c1f1f0e2843683570873443))
* **deps:** update kubernetes packages to v0.29.1 ([#1156](https://github.com/open-feature/flagd/issues/1156)) ([899e6b5](https://github.com/open-feature/flagd/commit/899e6b505abe63b7b599858a0d55e6a69be08993))
* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.14.0-20231031123731-ac2ec0f39838.1 ([#1170](https://github.com/open-feature/flagd/issues/1170)) ([8b3c8d6](https://github.com/open-feature/flagd/commit/8b3c8d6c87cbebd6f324771cb8aa1f6990a74cf4))
* **deps:** update module golang.org/x/crypto to v0.18.0 ([#1138](https://github.com/open-feature/flagd/issues/1138)) ([53569d9](https://github.com/open-feature/flagd/commit/53569d9cd88de1073a7e49b1a835adee4b0e8ef2))
* **deps:** update module golang.org/x/net to v0.20.0 ([#1139](https://github.com/open-feature/flagd/issues/1139)) ([fdb1d0c](https://github.com/open-feature/flagd/commit/fdb1d0c909373e23a8c8fca435ef77205526f730))
* **deps:** update module google.golang.org/grpc to v1.61.0 ([#1164](https://github.com/open-feature/flagd/issues/1164)) ([11ccecd](https://github.com/open-feature/flagd/commit/11ccecd5ac2e23991ac76ee079630f527559db1d))
* **deps:** update opentelemetry-go monorepo ([#1155](https://github.com/open-feature/flagd/issues/1155)) ([436fefe](https://github.com/open-feature/flagd/commit/436fefedf67afda1cbf97ff12f7c3071cb833d9a))


### ✨ New Features

* add targeting validation ([#1146](https://github.com/open-feature/flagd/issues/1146)) ([b727dd0](https://github.com/open-feature/flagd/commit/b727dd00c561c27e54dbe4bafff0ca2d82487a42))
* **core:** support any auth scheme in HTTP-sync auth header ([#1152](https://github.com/open-feature/flagd/issues/1152)) ([df65966](https://github.com/open-feature/flagd/commit/df6596634e1ca960592d5e4825985aeee781a25d))


### 🧹 Chore

* update test dependencies, fix for otel api change and update renovate configuration ([#1188](https://github.com/open-feature/flagd/issues/1188)) ([3270346](https://github.com/open-feature/flagd/commit/32703464d5c637dbb06c2c857070e9f038977c01))


### 📚 Documentation

* update schemas ([#1158](https://github.com/open-feature/flagd/issues/1158)) ([396c618](https://github.com/open-feature/flagd/commit/396c618bac13c5c8eb2aadc29bb126a83fec1b56))

## [0.7.4](https://github.com/open-feature/flagd/compare/core/v0.7.3...core/v0.7.4) (2024-01-04)


### 🐛 Bug Fixes

* add custom marshalling options ([#1117](https://github.com/open-feature/flagd/issues/1117)) ([e8e49de](https://github.com/open-feature/flagd/commit/e8e49de909fba6791a26e73a59d5b7ab54284499))
* **deps:** update kubernetes packages to v0.29.0 ([#1082](https://github.com/open-feature/flagd/issues/1082)) ([751a79a](https://github.com/open-feature/flagd/commit/751a79a16571bea943aa1297c9828d1149d0c319))
* **deps:** update module connectrpc.com/connect to v1.14.0 ([#1108](https://github.com/open-feature/flagd/issues/1108)) ([0a41aca](https://github.com/open-feature/flagd/commit/0a41acae082e65b2f0b1e1f36304556f77eba72a))
* **deps:** update module github.com/prometheus/client_golang to v1.18.0 ([#1110](https://github.com/open-feature/flagd/issues/1110)) ([745bbb0](https://github.com/open-feature/flagd/commit/745bbb079af50f8fdb870e2234edd70c7a7a52f9))
* **deps:** update module golang.org/x/crypto to v0.17.0 [security] ([#1090](https://github.com/open-feature/flagd/issues/1090)) ([26681de](https://github.com/open-feature/flagd/commit/26681ded7389ca4ecb3d9612aedad8c700796dc9))
* **deps:** update module google.golang.org/protobuf to v1.32.0 ([#1106](https://github.com/open-feature/flagd/issues/1106)) ([e0d3b34](https://github.com/open-feature/flagd/commit/e0d3b34aa2fad5435a296dda896c9d430563b6ea))

## [0.7.3](https://github.com/open-feature/flagd/compare/core/v0.7.2...core/v0.7.3) (2023-12-22)


### 🐛 Bug Fixes

* **deps:** update golang.org/x/exp digest to 6522937 ([#1032](https://github.com/open-feature/flagd/issues/1032)) ([78b23d2](https://github.com/open-feature/flagd/commit/78b23d25fcd2ea49675fc96963565afbdf0acf25))
* **deps:** update module connectrpc.com/connect to v1.13.0 ([#1070](https://github.com/open-feature/flagd/issues/1070)) ([63f86ea](https://github.com/open-feature/flagd/commit/63f86ea699ff1f25639101a36e4fc25c76c94c1e))
* **deps:** update module github.com/diegoholiveira/jsonlogic/v3 to v3.4.0 ([#1068](https://github.com/open-feature/flagd/issues/1068)) ([5c5d5ab](https://github.com/open-feature/flagd/commit/5c5d5abc38540277deb4e11f41f79ff49273d659))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.5.2 ([#1059](https://github.com/open-feature/flagd/issues/1059)) ([cefea3e](https://github.com/open-feature/flagd/commit/cefea3ee035726940ee5ac51e8f7fc92cc0eeac3))
* **deps:** update module google.golang.org/grpc to v1.60.0 ([#1074](https://github.com/open-feature/flagd/issues/1074)) ([bf3e9d8](https://github.com/open-feature/flagd/commit/bf3e9d82b40afdc89edc519f40bc0f67af48f1bf))
* **deps:** update module google.golang.org/grpc to v1.60.1 ([#1092](https://github.com/open-feature/flagd/issues/1092)) ([5bf1368](https://github.com/open-feature/flagd/commit/5bf1368e68012937b57bf27168646e784524ae9b))
* make sure sync builder is initialized to avoid nil pointer access ([#1076](https://github.com/open-feature/flagd/issues/1076)) ([ebcd616](https://github.com/open-feature/flagd/commit/ebcd616e0df1ab56fd11bcf4f53fa7cf13f5e1d6))


### ✨ New Features

* support new flagd.evaluation and flagd.sync schemas ([#1083](https://github.com/open-feature/flagd/issues/1083)) ([e9728aa](https://github.com/open-feature/flagd/commit/e9728aae8352e77e6564a88e37d13f87021526ef))


### 🧹 Chore

* refactoring component structure ([#1044](https://github.com/open-feature/flagd/issues/1044)) ([0c7f78a](https://github.com/open-feature/flagd/commit/0c7f78a95fa4ad2a8b2afe2f6023b9c6d4fd48ed))
* renaming of evaluation components ([#1064](https://github.com/open-feature/flagd/issues/1064)) ([d39f31d](https://github.com/open-feature/flagd/commit/d39f31d65b57a9b033a4976e5cea8fdab716769d))
* use client-go library for retrieving FeatureFlag CRs ([#1077](https://github.com/open-feature/flagd/issues/1077)) ([c86dff0](https://github.com/open-feature/flagd/commit/c86dff0bd058ca65f4ccd34c0bc5c0a4ad5b11a6))

## [0.7.2](https://github.com/open-feature/flagd/compare/core/v0.7.1...core/v0.7.2) (2023-12-05)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/open-feature-operator to v0.5.0 ([#1039](https://github.com/open-feature/flagd/issues/1039)) ([eb128d9](https://github.com/open-feature/flagd/commit/eb128d97ecb8b6916f9c5a32c21c698207f82be5))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.5.1 ([#1046](https://github.com/open-feature/flagd/issues/1046)) ([0321935](https://github.com/open-feature/flagd/commit/0321935992ef50c2fc62bd659d8858c25c88efa7))
* **deps:** update module golang.org/x/crypto to v0.16.0 ([#1033](https://github.com/open-feature/flagd/issues/1033)) ([b79aaf2](https://github.com/open-feature/flagd/commit/b79aaf2c2ca5fe9c43c1460d6cab349d5d68b7e1))
* **deps:** update module golang.org/x/net to v0.19.0 ([#1034](https://github.com/open-feature/flagd/issues/1034)) ([c6426b2](https://github.com/open-feature/flagd/commit/c6426b20fa241066fc844a2a000812015790cf12))
* various edge cases in targeting ([#1041](https://github.com/open-feature/flagd/issues/1041)) ([ca38c16](https://github.com/open-feature/flagd/commit/ca38c165c606ffa4517872f170375dde37f4ac5b))

## [0.7.1](https://github.com/open-feature/flagd/compare/core/v0.7.0...core/v0.7.1) (2023-11-28)


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.28.4 ([#1016](https://github.com/open-feature/flagd/issues/1016)) ([ae470e3](https://github.com/open-feature/flagd/commit/ae470e37f3368c81484f0c54366ccf059ea2cea6))
* **deps:** update opentelemetry-go monorepo ([#1019](https://github.com/open-feature/flagd/issues/1019)) ([23ae555](https://github.com/open-feature/flagd/commit/23ae555ad73128dff46e911fccfef76306f5c550))


### 🔄 Refactoring

* Rename metrics-port to management-port ([#1012](https://github.com/open-feature/flagd/issues/1012)) ([5635e38](https://github.com/open-feature/flagd/commit/5635e38703cae835a53e9cce83d5bc42d00091e2))

## [0.7.0](https://github.com/open-feature/flagd/compare/core/v0.6.8...core/v0.7.0) (2023-11-15)


### ⚠ BREAKING CHANGES

* OFO APIs were updated to version v1beta1, since they are more stable now. Resources of the alpha versions are no longer supported in flagd or flagd-proxy.

### ✨ New Features

* support OFO v1beta1 API ([#997](https://github.com/open-feature/flagd/issues/997)) ([bb6f5bf](https://github.com/open-feature/flagd/commit/bb6f5bf0fc382ade75d80a34d209beaa2edc459d))

## [0.6.8](https://github.com/open-feature/flagd/compare/core/v0.6.7...core/v0.6.8) (2023-11-13)


### 🐛 Bug Fixes

* **deps:** update golang.org/x/exp digest to 9a3e603 ([#929](https://github.com/open-feature/flagd/issues/929)) ([f8db930](https://github.com/open-feature/flagd/commit/f8db930da22f52ffddd4533bdbaf8188f89250d0))
* **deps:** update kubernetes packages to v0.28.3 ([#974](https://github.com/open-feature/flagd/issues/974)) ([d7d205f](https://github.com/open-feature/flagd/commit/d7d205f457e46de3385610ee27db6dfd41323cd1))
* **deps:** update module github.com/diegoholiveira/jsonlogic/v3 to v3.3.1 ([#971](https://github.com/open-feature/flagd/issues/971)) ([f1a40b8](https://github.com/open-feature/flagd/commit/f1a40b862e7be69b542ab65a645c29c56b8fe307))
* **deps:** update module github.com/diegoholiveira/jsonlogic/v3 to v3.3.2 ([#975](https://github.com/open-feature/flagd/issues/975)) ([b53c14a](https://github.com/open-feature/flagd/commit/b53c14afabd6b3f2403d9e2ea2bd26c8f1cfe8e4))
* **deps:** update module github.com/fsnotify/fsnotify to v1.7.0 ([#981](https://github.com/open-feature/flagd/issues/981)) ([727b9d2](https://github.com/open-feature/flagd/commit/727b9d2046ccd998d2ab721df10d88b84445494f))
* **deps:** update module golang.org/x/mod to v0.14.0 ([#991](https://github.com/open-feature/flagd/issues/991)) ([87bc12d](https://github.com/open-feature/flagd/commit/87bc12dd968e9fcf524eb0b69462222f29c6ba34))
* **deps:** update module golang.org/x/net to v0.17.0 [security] ([#963](https://github.com/open-feature/flagd/issues/963)) ([7f54bd1](https://github.com/open-feature/flagd/commit/7f54bd1fb3fdbb7dff3a7b097f804ce843bb6e3a))
* **deps:** update module golang.org/x/net to v0.18.0 ([#1000](https://github.com/open-feature/flagd/issues/1000)) ([e9347cc](https://github.com/open-feature/flagd/commit/e9347cc88a4174c984432571710588ac57665cb7))
* **deps:** update module golang.org/x/sync to v0.5.0 ([#992](https://github.com/open-feature/flagd/issues/992)) ([bd24536](https://github.com/open-feature/flagd/commit/bd24536722f3c3c99b5946de0b62dba603992ba7))
* **deps:** update module google.golang.org/grpc to v1.59.0 ([#972](https://github.com/open-feature/flagd/issues/972)) ([7d0f1f2](https://github.com/open-feature/flagd/commit/7d0f1f225ebfe83681c83cf38bf5c253a7b0ebef))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.3 ([#976](https://github.com/open-feature/flagd/issues/976)) ([b33c9c9](https://github.com/open-feature/flagd/commit/b33c9c97cf21eb317e4b0b1f452571d85de145b1))
* **deps:** update opentelemetry-go monorepo ([#1001](https://github.com/open-feature/flagd/issues/1001)) ([9798aeb](https://github.com/open-feature/flagd/commit/9798aeb248764400128048b65ba27baba71b07e6))


### 🔄 Refactoring

* migrate to connectrpc/connect-go ([#990](https://github.com/open-feature/flagd/issues/990)) ([7dd5b2b](https://github.com/open-feature/flagd/commit/7dd5b2b4c284481bcba5a9c45bd6c85ad1dc6d33))

## [0.6.7](https://github.com/open-feature/flagd/compare/core/v0.6.6...core/v0.6.7) (2023-10-12)


### 🐛 Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.17.0 ([#939](https://github.com/open-feature/flagd/issues/939)) ([9065cba](https://github.com/open-feature/flagd/commit/9065cba599ac07225b613c5acd9403ab24462078))
* **deps:** update module github.com/rs/cors to v1.10.1 ([#946](https://github.com/open-feature/flagd/issues/946)) ([1c39862](https://github.com/open-feature/flagd/commit/1c39862297746e66189ece87892b6e4694294fb6))
* **deps:** update module go.uber.org/zap to v1.26.0 ([#917](https://github.com/open-feature/flagd/issues/917)) ([e57e206](https://github.com/open-feature/flagd/commit/e57e206c937d5b11b81d46ee57b3e92cc454dd88))
* **deps:** update module golang.org/x/mod to v0.13.0 ([#952](https://github.com/open-feature/flagd/issues/952)) ([be61450](https://github.com/open-feature/flagd/commit/be61450ec3fab3c294c9813df193c98e374900aa))
* **deps:** update module golang.org/x/sync to v0.4.0 ([#949](https://github.com/open-feature/flagd/issues/949)) ([faa24a6](https://github.com/open-feature/flagd/commit/faa24a6c6f34330364e1ba0c3a847943f4e55150))
* **deps:** update module google.golang.org/grpc to v1.58.1 ([#915](https://github.com/open-feature/flagd/issues/915)) ([06d95de](https://github.com/open-feature/flagd/commit/06d95ded9b69c9c598d08f8a6ef73ec598a817af))
* **deps:** update module google.golang.org/grpc to v1.58.2 ([#928](https://github.com/open-feature/flagd/issues/928)) ([90f1878](https://github.com/open-feature/flagd/commit/90f1878ae482ea4a684615f733392a38301de68d))
* **deps:** update module google.golang.org/grpc to v1.58.3 ([#960](https://github.com/open-feature/flagd/issues/960)) ([fee1558](https://github.com/open-feature/flagd/commit/fee1558da4f5418499fb09fe356d16f008423eb7))
* **deps:** update opentelemetry-go monorepo ([#943](https://github.com/open-feature/flagd/issues/943)) ([e7cee41](https://github.com/open-feature/flagd/commit/e7cee41630e5684999b3689dbf1ab66234c65f6e))
* erroneous warning about prop overwrite ([#924](https://github.com/open-feature/flagd/issues/924)) ([673b76a](https://github.com/open-feature/flagd/commit/673b76aeff5c27e8e031e59da5f0ed1871d3f749))


### ✨ New Features

* add `$flagd.timestamp` to json evaluator ([#958](https://github.com/open-feature/flagd/issues/958)) ([a1b04e7](https://github.com/open-feature/flagd/commit/a1b04e778df4d2f6beb881183a82161018151479))

## [0.6.6](https://github.com/open-feature/flagd/compare/core/v0.6.5...core/v0.6.6) (2023-09-14)


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.28.2 ([#911](https://github.com/open-feature/flagd/issues/911)) ([2eda6ab](https://github.com/open-feature/flagd/commit/2eda6ab5e528f12a9ce6b6818e08abb0d783b23d))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.2 ([#907](https://github.com/open-feature/flagd/issues/907)) ([9976851](https://github.com/open-feature/flagd/commit/9976851d792ff3eb5fde18f19e397738eb7cacaf))
* **deps:** update opentelemetry-go monorepo ([#906](https://github.com/open-feature/flagd/issues/906)) ([5a41226](https://github.com/open-feature/flagd/commit/5a4122658039aafcf080fcc6655c2a679622ed69))
* use 32bit murmur calculation (64 is not stable) ([#913](https://github.com/open-feature/flagd/issues/913)) ([db8dca4](https://github.com/open-feature/flagd/commit/db8dca421cb0dba2968d47e5cc162d81401298db))

## [0.6.5](https://github.com/open-feature/flagd/compare/core/v0.6.4...core/v0.6.5) (2023-09-08)


### 🐛 Bug Fixes

* **deps:** update module github.com/rs/cors to v1.10.0 ([#893](https://github.com/open-feature/flagd/issues/893)) ([fe61fbe](https://github.com/open-feature/flagd/commit/fe61fbe47a4e58562cbcb1c5201281fae1adafaf))
* **deps:** update module golang.org/x/crypto to v0.13.0 ([#888](https://github.com/open-feature/flagd/issues/888)) ([1a9037a](https://github.com/open-feature/flagd/commit/1a9037a5b058e44fa844392d0110696b032eff6e))
* **deps:** update module golang.org/x/net to v0.15.0 ([#889](https://github.com/open-feature/flagd/issues/889)) ([233d976](https://github.com/open-feature/flagd/commit/233d97694826d0e018be19a78259188802aba37f))
* **deps:** update module google.golang.org/grpc to v1.58.0 ([#896](https://github.com/open-feature/flagd/issues/896)) ([853b76d](https://github.com/open-feature/flagd/commit/853b76dfa3babfebd8bdbcd3e0913380f077b8ab))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.1 ([#882](https://github.com/open-feature/flagd/issues/882)) ([ca3d85a](https://github.com/open-feature/flagd/commit/ca3d85a51c0ed1c1def54d7304d4b9fe69622662))
* **deps:** update opentelemetry-go monorepo ([#868](https://github.com/open-feature/flagd/issues/868)) ([d48317f](https://github.com/open-feature/flagd/commit/d48317f61d7db7ba0398dc9ab7cdd174a0b87555))


### 🧹 Chore

* upgrade to go 1.20 ([#891](https://github.com/open-feature/flagd/issues/891)) ([977167f](https://github.com/open-feature/flagd/commit/977167fb8db330b62726097616dcd691267199ad))

## [0.6.4](https://github.com/open-feature/flagd/compare/core/v0.6.3...core/v0.6.4) (2023-08-30)


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.28.0 ([#841](https://github.com/open-feature/flagd/issues/841)) ([cc195e1](https://github.com/open-feature/flagd/commit/cc195e1dde052d583656d5e5b49caec50f832365))
* **deps:** update kubernetes packages to v0.28.1 ([#860](https://github.com/open-feature/flagd/issues/860)) ([f3237c2](https://github.com/open-feature/flagd/commit/f3237c2d324fbb15fd5f7fe337a0601af3b537bb))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.36 ([#799](https://github.com/open-feature/flagd/issues/799)) ([fa4da4b](https://github.com/open-feature/flagd/commit/fa4da4b0115e9fb40ab038b996e1e32b9f6a47ab))
* **deps:** update module golang.org/x/crypto to v0.12.0 ([#797](https://github.com/open-feature/flagd/issues/797)) ([edae3fd](https://github.com/open-feature/flagd/commit/edae3fd466c0be62a0256c268e85cb337c9536f2))
* **deps:** update module golang.org/x/net to v0.14.0 ([#798](https://github.com/open-feature/flagd/issues/798)) ([92c2f26](https://github.com/open-feature/flagd/commit/92c2f2676163688130737b34a115374cb5631247))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.15.1 ([#795](https://github.com/open-feature/flagd/issues/795)) ([13d62fd](https://github.com/open-feature/flagd/commit/13d62fd0fc4749f19dba0a18e1fda46a723380c5))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.0 ([#856](https://github.com/open-feature/flagd/issues/856)) ([88d832a](https://github.com/open-feature/flagd/commit/88d832a9d49a4bc1d6156849a59227ecab07f96e))


### ✨ New Features

* add flag key to hash in fractional evaluation ([#847](https://github.com/open-feature/flagd/issues/847)) ([ca6a35f](https://github.com/open-feature/flagd/commit/ca6a35fd72462177f45a116e9009fc30b3588b83))
* add gRPC healthchecks ([#863](https://github.com/open-feature/flagd/issues/863)) ([da30b7b](https://github.com/open-feature/flagd/commit/da30b7babffd8487c992fa41519787c8d78ebdba))
* support nested props in fractional evaluator ([#869](https://github.com/open-feature/flagd/issues/869)) ([50ff739](https://github.com/open-feature/flagd/commit/50ff739178fb732e38a220bb6a071260af1f2469))


### 🧹 Chore

* deprecate fractionalEvaluation for fractional ([#873](https://github.com/open-feature/flagd/issues/873)) ([243fef9](https://github.com/open-feature/flagd/commit/243fef9e1f0ed00ccf5d9a389e10d9ad6a197fb1))
* replace xxh3 with murmur3 in bucket algorithm ([#846](https://github.com/open-feature/flagd/issues/846)) ([c3c9e4e](https://github.com/open-feature/flagd/commit/c3c9e4e40aeae7e75b1b9ab13bb9a40264be84e5))

## [0.6.3](https://github.com/open-feature/flagd/compare/core/v0.6.2...core/v0.6.3) (2023-08-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/diegoholiveira/jsonlogic/v3 to v3.3.0 ([#785](https://github.com/open-feature/flagd/issues/785)) ([ee9c54b](https://github.com/open-feature/flagd/commit/ee9c54b6b5cd51b947aae1ff6309ffae07ce89eb))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.35 ([#783](https://github.com/open-feature/flagd/issues/783)) ([9ff0b5b](https://github.com/open-feature/flagd/commit/9ff0b5b1bd3bb95581eab83f944aa60e179b207a))
* **deps:** update module go.uber.org/zap to v1.25.0 ([#786](https://github.com/open-feature/flagd/issues/786)) ([40d0aa6](https://github.com/open-feature/flagd/commit/40d0aa66cf422db6811206d777b55396a96f330f))
* **deps:** update module golang.org/x/net to v0.13.0 ([#784](https://github.com/open-feature/flagd/issues/784)) ([f57d023](https://github.com/open-feature/flagd/commit/f57d023174d9cc74b7d8260055f82b84a2bdcc52))
* metric descriptions match the otel spec ([#789](https://github.com/open-feature/flagd/issues/789)) ([34befcd](https://github.com/open-feature/flagd/commit/34befcdfedc5f0479cb0ae77fe148849c341d33e))


### ✨ New Features

* add new configuration "sync-interval" which controls the HTTP polling interval ([#404](https://github.com/open-feature/flagd/issues/404)) ([ace62c7](https://github.com/open-feature/flagd/commit/ace62c7a6ab2b5b5d26642286deb6db406391d8f))
* include falsy json fields ([#792](https://github.com/open-feature/flagd/issues/792)) ([37d91a0](https://github.com/open-feature/flagd/commit/37d91a09836f07e07b12acd13850ea5c7c9252cd))

## [0.6.2](https://github.com/open-feature/flagd/compare/core/v0.6.1...core/v0.6.2) (2023-07-28)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.3.0-20230720212818-3675556880a1.1 ([#747](https://github.com/open-feature/flagd/issues/747)) ([fb17bc6](https://github.com/open-feature/flagd/commit/fb17bc6a5c715f507b2838c150dc8a2f139a38fb))
* **deps:** update module golang.org/x/net to v0.12.0 ([#734](https://github.com/open-feature/flagd/issues/734)) ([777b28b](https://github.com/open-feature/flagd/commit/777b28b1d512245b0046d11197f6dfa341b317d2))


### ✨ New Features

* grpc selector as scope ([#761](https://github.com/open-feature/flagd/issues/761)) ([7246e6d](https://github.com/open-feature/flagd/commit/7246e6dce648c6445f90d71fc172bbab209d9928))

## [0.6.1](https://github.com/open-feature/flagd/compare/core/v0.6.0...core/v0.6.1) (2023-07-27)


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.27.4 ([#756](https://github.com/open-feature/flagd/issues/756)) ([dcc10f3](https://github.com/open-feature/flagd/commit/dcc10f33f5fd9a8936241725ea811b90b4f136be))
* **deps:** update module github.com/bufbuild/connect-go to v1.10.0 ([#771](https://github.com/open-feature/flagd/issues/771)) ([c74103f](https://github.com/open-feature/flagd/commit/c74103faec068f14c87ad3ec227f5b802dbfac43))
* **deps:** update module google.golang.org/grpc to v1.57.0 ([#773](https://github.com/open-feature/flagd/issues/773)) ([be8bf04](https://github.com/open-feature/flagd/commit/be8bf045093d89099eead2cccb86a5a7275e25d5))


### ✨ New Features

* **flagd-proxy:** introduce zero-downtime ([#752](https://github.com/open-feature/flagd/issues/752)) ([ed5e6e5](https://github.com/open-feature/flagd/commit/ed5e6e5f3ee0a923c33dbf1a8bf20f80adec71bd))
* **flagd:** custom error handling for OTel errors ([#769](https://github.com/open-feature/flagd/issues/769)) ([bda1a92](https://github.com/open-feature/flagd/commit/bda1a92785c4348fe306a1d259b7bea91bd01c41))

## [0.6.0](https://github.com/open-feature/flagd/compare/core/v0.5.4...core/v0.6.0) (2023-07-13)


### ⚠ BREAKING CHANGES

* rename metrics and service ([#730](https://github.com/open-feature/flagd/issues/730))

### 🔄 Refactoring

* remove protobuf dependency from eval package ([#701](https://github.com/open-feature/flagd/issues/701)) ([34ffafd](https://github.com/open-feature/flagd/commit/34ffafd9a777da3f11bd3bfa81565e774cc63214))


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.27.3 ([#708](https://github.com/open-feature/flagd/issues/708)) ([5bf3a69](https://github.com/open-feature/flagd/commit/5bf3a69aa4bf95ce77ad08491bcce420620525d3))
* **deps:** update module github.com/bufbuild/connect-go to v1.9.0 ([#722](https://github.com/open-feature/flagd/issues/722)) ([75223e2](https://github.com/open-feature/flagd/commit/75223e2fc01c4dcd0291b46a0d50b8815b31654c))
* **deps:** update module github.com/bufbuild/connect-opentelemetry-go to v0.4.0 ([#739](https://github.com/open-feature/flagd/issues/739)) ([713e2a9](https://github.com/open-feature/flagd/commit/713e2a9834546963615046de1b6125e7fa6bf20d))
* **deps:** update module github.com/prometheus/client_golang to v1.16.0 ([#709](https://github.com/open-feature/flagd/issues/709)) ([b8bedd2](https://github.com/open-feature/flagd/commit/b8bedd2b895026eace8204ae4ffcff771f7e8e97))
* **deps:** update module golang.org/x/crypto to v0.10.0 ([#647](https://github.com/open-feature/flagd/issues/647)) ([7f1d7e6](https://github.com/open-feature/flagd/commit/7f1d7e66669b88b2c56b32f9cdd9be354ebcfc8e))
* **deps:** update module golang.org/x/mod to v0.11.0 ([#705](https://github.com/open-feature/flagd/issues/705)) ([42813be](https://github.com/open-feature/flagd/commit/42813bef092ba7fffed0dd94166bfd01ea8a7582))
* **deps:** update module golang.org/x/mod to v0.12.0 ([#729](https://github.com/open-feature/flagd/issues/729)) ([7b109c7](https://github.com/open-feature/flagd/commit/7b109c705aceb652ac2675bd0ffe82420983798b))
* **deps:** update module golang.org/x/net to v0.11.0 ([#706](https://github.com/open-feature/flagd/issues/706)) ([27d893f](https://github.com/open-feature/flagd/commit/27d893fe78417f7b8418003edc401ab5a6c21fb9))
* **deps:** update module golang.org/x/sync to v0.3.0 ([#707](https://github.com/open-feature/flagd/issues/707)) ([7852efb](https://github.com/open-feature/flagd/commit/7852efb84e9f071b2b482b1968d799888b6882dc))
* **deps:** update module google.golang.org/grpc to v1.56.1 ([#710](https://github.com/open-feature/flagd/issues/710)) ([8f16573](https://github.com/open-feature/flagd/commit/8f165739aee8f28800e200b357203e88a3fd5938))
* **deps:** update module google.golang.org/grpc to v1.56.2 ([#738](https://github.com/open-feature/flagd/issues/738)) ([521cc30](https://github.com/open-feature/flagd/commit/521cc30cde1971be000ec10d93f6d70b9b2260ee))
* **deps:** update module google.golang.org/protobuf to v1.31.0 ([#720](https://github.com/open-feature/flagd/issues/720)) ([247239e](https://github.com/open-feature/flagd/commit/247239e76b9de1a619aad9e957ed8b44ae534b77))
* **deps:** update opentelemetry-go monorepo ([#648](https://github.com/open-feature/flagd/issues/648)) ([c12dad8](https://github.com/open-feature/flagd/commit/c12dad89a8e761154f57739ded594b2783a14f8a))


### ✨ New Features

* **flagD:** support zero downtime during upgrades ([#731](https://github.com/open-feature/flagd/issues/731)) ([7df8d39](https://github.com/open-feature/flagd/commit/7df8d3994b75991b5e49a65728ef5e4b24a85dde))
* rename metrics and service ([#730](https://github.com/open-feature/flagd/issues/730)) ([09c0198](https://github.com/open-feature/flagd/commit/09c0198f76a200b1b6a1f48e9c94ec0547283ca2))

## [0.5.4](https://github.com/open-feature/flagd/compare/core/v0.5.3...core/v0.5.4) (2023-06-07)


### ✨ New Features

* add `sem_ver` jsonLogic evaluator ([#675](https://github.com/open-feature/flagd/issues/675)) ([a8d8ab6](https://github.com/open-feature/flagd/commit/a8d8ab6b4495457a40a2c32b8bd5be48b1fd6941))
* add `starts_with` and `ends_with` json evaluators ([#658](https://github.com/open-feature/flagd/issues/658)) ([f932b8f](https://github.com/open-feature/flagd/commit/f932b8f4c834a5ebe27ebb860c26fdea8da20598))
* telemetry improvements ([#653](https://github.com/open-feature/flagd/issues/653)) ([ea02cba](https://github.com/open-feature/flagd/commit/ea02cba24bde982d55956fe54de1e8f27226bfc6))


### 🐛 Bug Fixes

* **deps:** update module github.com/bufbuild/connect-go to v1.8.0 ([#683](https://github.com/open-feature/flagd/issues/683)) ([13bb13d](https://github.com/open-feature/flagd/commit/13bb13daa11068481ba97f3432ae08de78392a91))
* **deps:** update module github.com/bufbuild/connect-opentelemetry-go to v0.3.0 ([#669](https://github.com/open-feature/flagd/issues/669)) ([e899435](https://github.com/open-feature/flagd/commit/e899435c29c32264ea2477436e69ce92c7775ee9))
* **deps:** update module github.com/prometheus/client_golang to v1.15.1 ([#636](https://github.com/open-feature/flagd/issues/636)) ([b22279d](https://github.com/open-feature/flagd/commit/b22279df469dc78f9d3e5bc4a59ab6baf539a8ae))
* **deps:** update module github.com/stretchr/testify to v1.8.3 ([#662](https://github.com/open-feature/flagd/issues/662)) ([2e06d58](https://github.com/open-feature/flagd/commit/2e06d582ee9c8abfd57f8945d91261eab6cf9854))
* **deps:** update module github.com/stretchr/testify to v1.8.4 ([#678](https://github.com/open-feature/flagd/issues/678)) ([ca8c9d6](https://github.com/open-feature/flagd/commit/ca8c9d66a0c6b21129c4c36a3c10dcf3be869ee7))
* **deps:** update module golang.org/x/mod to v0.10.0 ([#682](https://github.com/open-feature/flagd/issues/682)) ([16199ce](https://github.com/open-feature/flagd/commit/16199ceac9ebbae68dafbd6c21239f64f8c32511))
* **deps:** update module golang.org/x/net to v0.10.0 ([#644](https://github.com/open-feature/flagd/issues/644)) ([ccd9d35](https://github.com/open-feature/flagd/commit/ccd9d351df153039a124064f30e5829610773f27))
* **deps:** update module golang.org/x/sync to v0.2.0 ([#638](https://github.com/open-feature/flagd/issues/638)) ([7f4a7db](https://github.com/open-feature/flagd/commit/7f4a7db8139294a21b3415710c143f182d93264a))
* **deps:** update module google.golang.org/grpc to v1.55.0 ([#640](https://github.com/open-feature/flagd/issues/640)) ([c0d7328](https://github.com/open-feature/flagd/commit/c0d732866262240e340fe10f8ac0f6ff2a5c4f8c))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.15.0 ([#665](https://github.com/open-feature/flagd/issues/665)) ([9490ed6](https://github.com/open-feature/flagd/commit/9490ed62e2fc589af8ae7ee26bfd559797a1f83c))
* fix connect error code handling for disabled flags ([#670](https://github.com/open-feature/flagd/issues/670)) ([86a8012](https://github.com/open-feature/flagd/commit/86a8012efcfeb3e967657f6143c143b457d64ca2))
* remove disabled flags from bulk evaluation ([#672](https://github.com/open-feature/flagd/issues/672)) ([d2ce988](https://github.com/open-feature/flagd/commit/d2ce98838edf63b88ee9fb5ae6f8d534e1112e7e))


### 🔄 Refactoring

* introduce additional linting rules + fix discrepancies ([#616](https://github.com/open-feature/flagd/issues/616)) ([aef0b90](https://github.com/open-feature/flagd/commit/aef0b9042dcbe5b3f9a7e97960b27366fe50adfe))
* introduce isyncstore interface ([#660](https://github.com/open-feature/flagd/issues/660)) ([c0e2fa0](https://github.com/open-feature/flagd/commit/c0e2fa00736d46db98f72114a449b2e2bf998e3d))


### 🧹 Chore

* refactor json logic evaluator to pass custom operators as options ([#691](https://github.com/open-feature/flagd/issues/691)) ([1c9bff9](https://github.com/open-feature/flagd/commit/1c9bff9a523037c3654b592dc08c193aa3295e9e))
* update otel dependencies ([#649](https://github.com/open-feature/flagd/issues/649)) ([2114e41](https://github.com/open-feature/flagd/commit/2114e41c38951247866c0b408e5f933282902e70))

## [0.5.3](https://github.com/open-feature/flagd/compare/core/v0.5.2...core/v0.5.3) (2023-05-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/bufbuild/connect-go to v1.6.0 ([#585](https://github.com/open-feature/flagd/issues/585)) ([8f2f467](https://github.com/open-feature/flagd/commit/8f2f467af52a3686196a821eec61954d89d3f71d))
* **deps:** update module github.com/bufbuild/connect-go to v1.7.0 ([#625](https://github.com/open-feature/flagd/issues/625)) ([1b24fc9](https://github.com/open-feature/flagd/commit/1b24fc923a405b337634009831ef0b9792953ce5))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.34 ([#604](https://github.com/open-feature/flagd/issues/604)) ([3e6a84b](https://github.com/open-feature/flagd/commit/3e6a84b455a330f541784b346d3b6199f8b423f7))
* **deps:** update module github.com/prometheus/client_golang to v1.15.0 ([#608](https://github.com/open-feature/flagd/issues/608)) ([0597a8f](https://github.com/open-feature/flagd/commit/0597a8f23d0914b26f06b1335b49fe7c18ecb4f9))
* **deps:** update module github.com/rs/cors to v1.9.0 ([#609](https://github.com/open-feature/flagd/issues/609)) ([97066c1](https://github.com/open-feature/flagd/commit/97066c107d14777eaad8a05b3bb051639af3179c))
* **deps:** update module github.com/rs/xid to v1.5.0 ([#614](https://github.com/open-feature/flagd/issues/614)) ([e3dfbc6](https://github.com/open-feature/flagd/commit/e3dfbc6753bddc9e2f5c4a2633a7010123cf2f97))
* **deps:** update module golang.org/x/crypto to v0.8.0 ([#595](https://github.com/open-feature/flagd/issues/595)) ([36016d7](https://github.com/open-feature/flagd/commit/36016d7940fa772c01dd61b071b2c9ec753cfb75))


### ✨ New Features

* Introduce connect traces ([#624](https://github.com/open-feature/flagd/issues/624)) ([28bac6a](https://github.com/open-feature/flagd/commit/28bac6a54aed79cb8d84a147ffea296c36f5bd51))


### 🧹 Chore

* add instructions for windows and fix failing unit tests ([#632](https://github.com/open-feature/flagd/issues/632)) ([6999d67](https://github.com/open-feature/flagd/commit/6999d6722581ab8e2e14bfd4b2d0341fe5216684))

## [0.5.2](https://github.com/open-feature/flagd/compare/core/v0.5.1...core/v0.5.2) (2023-04-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.32 [security] ([#606](https://github.com/open-feature/flagd/issues/606)) ([6f721af](https://github.com/open-feature/flagd/commit/6f721af379fcb0f1a74410637a313477148ef863))
* eventing configuration setup ([#605](https://github.com/open-feature/flagd/issues/605)) ([edfbe51](https://github.com/open-feature/flagd/commit/edfbe5191651f25da991b507a3feedcbbe3c66f1))


### ✨ New Features

* introduce metrics for failed evaluations ([#584](https://github.com/open-feature/flagd/issues/584)) ([77664cd](https://github.com/open-feature/flagd/commit/77664cdf53a868f56ca040bdfe3f4930cd9a8fb4))
* otel traces for flag evaluation ([#598](https://github.com/open-feature/flagd/issues/598)) ([1757035](https://github.com/open-feature/flagd/commit/175703548f88469f25d749e320ee48030c9f9074))

## [0.5.1](https://github.com/open-feature/flagd/compare/core/v0.5.0...core/v0.5.1) (2023-04-12)


### 🧹 Chore

* move startServer functions into errGroups ([#566](https://github.com/open-feature/flagd/issues/566)) ([0223c23](https://github.com/open-feature/flagd/commit/0223c23fbd72322cf6ecbe6736968b3e6c6132bb))


### ✨ New Features

* flagd OTEL collector ([#586](https://github.com/open-feature/flagd/issues/586)) ([494bec3](https://github.com/open-feature/flagd/commit/494bec33dcc1ddf0fa5cd0866f06265618408f5e))


### 🔄 Refactoring

* remove connect-go from flagd-proxy and replace with grpc ([#589](https://github.com/open-feature/flagd/issues/589)) ([425de9a](https://github.com/open-feature/flagd/commit/425de9a1c2d1574779b905ac6debb9edfc156b15))


### 🐛 Bug Fixes

* flagd-proxy locking bug fix ([#592](https://github.com/open-feature/flagd/issues/592)) ([b166122](https://github.com/open-feature/flagd/commit/b1661225c912ee11ba4749f7ef157a0335e8781f))

## [0.5.0](https://github.com/open-feature/flagd/compare/core/v0.4.5...core/v0.5.0) (2023-03-30)


### ⚠ BREAKING CHANGES

* rename `kube-flagd-proxy` to `flagd-proxy` ([#576](https://github.com/open-feature/flagd/issues/576))
* unify sources configuration handling ([#560](https://github.com/open-feature/flagd/issues/560))

### 🧹 Chore

* move credential builder for grpc sync into seperate component ([#536](https://github.com/open-feature/flagd/issues/536)) ([7314fee](https://github.com/open-feature/flagd/commit/7314feea8c7bc90aac0528a9e1be0759a7a60c15))
* refactor configuration handling for startup ([#551](https://github.com/open-feature/flagd/issues/551)) ([8dfbde5](https://github.com/open-feature/flagd/commit/8dfbde5bbffd16fb66797a750d15f0226edf54a7))
* refactor middleware setup in server ([#554](https://github.com/open-feature/flagd/issues/554)) ([01016c7](https://github.com/open-feature/flagd/commit/01016c7df7c5f653cdadf151539432f692f36251))
* refactor service configuration objects ([#545](https://github.com/open-feature/flagd/issues/545)) ([c7b29ed](https://github.com/open-feature/flagd/commit/c7b29edcfe9dab61eaa585011a690d47829601b6)), closes [#524](https://github.com/open-feature/flagd/issues/524)
* unify sources configuration handling ([#560](https://github.com/open-feature/flagd/issues/560)) ([7f4888a](https://github.com/open-feature/flagd/commit/7f4888a1676e49acecf328685623566ea057ffcf))


### 🐛 Bug Fixes

* **deps:** update module google.golang.org/grpc to v1.54.0 ([#548](https://github.com/open-feature/flagd/issues/548)) ([99ba5ec](https://github.com/open-feature/flagd/commit/99ba5ece76d98124c108bc6280bee03a5c0cd25d))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.14.6 ([#572](https://github.com/open-feature/flagd/issues/572)) ([bed9458](https://github.com/open-feature/flagd/commit/bed94584a30bb6752284ece152cb114a102cbe8a))
* fixing silent lint failures ([#550](https://github.com/open-feature/flagd/issues/550)) ([30c8022](https://github.com/open-feature/flagd/commit/30c8022e891d1d278c096dd2438137ced7552678))
* nil pointer fix + export constructors ([#555](https://github.com/open-feature/flagd/issues/555)) ([78adb81](https://github.com/open-feature/flagd/commit/78adb81f4eb7a5b7fdb7075fa0bf8afa6d03dc72))


### ✨ New Features

* expose Impression metric ([#556](https://github.com/open-feature/flagd/issues/556)) ([77e0a33](https://github.com/open-feature/flagd/commit/77e0a33be24dcd0b6e239e5ed709167167c14171))
* Introduce kube-proxy-metrics ([#558](https://github.com/open-feature/flagd/issues/558)) ([ad0baeb](https://github.com/open-feature/flagd/commit/ad0baeb08fa67c94356d6a3f298283373bd5211b))
* rename `kube-flagd-proxy` to `flagd-proxy` ([#576](https://github.com/open-feature/flagd/issues/576)) ([223de99](https://github.com/open-feature/flagd/commit/223de99ee3efbcd601bf75ab1f6258eeac0c426e))
* refactor core module into multiple packages ([#530](https://github.com/open-feature/flagd/issues/530)) ([9d68d0b](https://github.com/open-feature/flagd/commit/9d68d0b45815facdf6079ffcd7864f720ccb8475))

## [0.4.5](https://github.com/open-feature/flagd/compare/core/v0.4.4...core/v0.4.5) (2023-03-20)


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.26.3 ([#533](https://github.com/open-feature/flagd/issues/533)) ([6ddd5b2](https://github.com/open-feature/flagd/commit/6ddd5b29806f3101cf122bfc4196ba7d0ef4c025))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.31 ([#527](https://github.com/open-feature/flagd/issues/527)) ([7928130](https://github.com/open-feature/flagd/commit/7928130b10906b10f4501630f16a71bdd8e4e365))
* **deps:** update module google.golang.org/protobuf to v1.29.1 [security] ([#504](https://github.com/open-feature/flagd/issues/504)) ([59db0bb](https://github.com/open-feature/flagd/commit/59db0bba43a9c002378fdced2fcf4729d619e28b))
* **deps:** update module google.golang.org/protobuf to v1.30.0 ([#499](https://github.com/open-feature/flagd/issues/499)) ([f650338](https://github.com/open-feature/flagd/commit/f650338e01e721a9d24e2ed6f6fe585dbb6beb42))


### ✨ New Features

* grpc connection options to flagd configuration options ([#532](https://github.com/open-feature/flagd/issues/532)) ([aa74951](https://github.com/open-feature/flagd/commit/aa74951f43b662ff2df53e68d347fc10e6d23bb8))
* Introduce flagd kube proxy ([#495](https://github.com/open-feature/flagd/issues/495)) ([440864c](https://github.com/open-feature/flagd/commit/440864ce87174618321c9d5146221490d8f07b24))

## [0.4.4](https://github.com/open-feature/flagd/compare/core-v0.4.3...core/v0.4.4) (2023-03-10)


### ✨ New Features

* Restructure for monorepo setup ([#486](https://github.com/open-feature/flagd/issues/486)) ([ed2993c](https://github.com/open-feature/flagd/commit/ed2993cd67b8a46db3beb6bb8a360e1aa20349da))
