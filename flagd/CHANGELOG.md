# Changelog

## [0.13.2](https://github.com/open-feature/flagd/compare/flagd/v0.13.1...flagd/v0.13.2) (2026-01-09)


### 🐛 Bug Fixes

* **security:** update module github.com/open-feature/flagd/core to v0.13.1 [security] ([#1846](https://github.com/open-feature/flagd/issues/1846)) ([a8a57ad](https://github.com/open-feature/flagd/commit/a8a57adb1d49640bfc14bf02cf77961652f7516a))

## [0.13.1](https://github.com/open-feature/flagd/compare/flagd/v0.13.0...flagd/v0.13.1) (2025-12-27)


### 🐛 Bug Fixes

* **security:** update go for various stdlib CVEs ([#1840](https://github.com/open-feature/flagd/issues/1840)) ([6dcb36d](https://github.com/open-feature/flagd/commit/6dcb36d2d6b55b7fe0b6107ac9a25baf302c5cdc))

## [0.13.0](https://github.com/open-feature/flagd/compare/flagd/v0.12.9...flagd/v0.13.0) (2025-12-23)


### 🐛 Bug Fixes

* fixing sync return format missing flag layer, adding full e2e suite ([#1827](https://github.com/open-feature/flagd/issues/1827)) ([570693d](https://github.com/open-feature/flagd/commit/570693d9e7b3200c0865e7ebb3b467ccfc38bb88))
* **security:** update module github.com/go-viper/mapstructure/v2 to v2.4.0 [security] ([#1784](https://github.com/open-feature/flagd/issues/1784)) ([037e30b](https://github.com/open-feature/flagd/commit/037e30b2f87897499580b25497049b88da7e386c))
* **security:** update module golang.org/x/crypto to v0.45.0 [security] ([#1826](https://github.com/open-feature/flagd/issues/1826)) ([7e0762b](https://github.com/open-feature/flagd/commit/7e0762b921ea70bed7915bcaab50e450e0a51158))


### ✨ New Features

* add support for http-based ofrep metrics ([#1803](https://github.com/open-feature/flagd/issues/1803)) ([fcd19b3](https://github.com/open-feature/flagd/commit/fcd19b310b319c9837b41c19d6f082ac750cb81d))
* cleanup evaluator interface ([#1793](https://github.com/open-feature/flagd/issues/1793)) ([aa504f7](https://github.com/open-feature/flagd/commit/aa504f7077093746f886248a4766d9ae5587bf3d))
* enable parsing of array flag configurations for flagd ([#1797](https://github.com/open-feature/flagd/issues/1797)) ([97c6ffa](https://github.com/open-feature/flagd/commit/97c6ffaf2b51765ccd6aaec38c2902ed2ac8f5f3))
* multi-project support via selectors and flagSetId namespacing ([#1702](https://github.com/open-feature/flagd/issues/1702)) ([f9ce46f](https://github.com/open-feature/flagd/commit/f9ce46f1032e7cb423e0e5c75a7c02f91ab5a88f))
* normalize selector in sync (use header as in OFREP and RPC)  ([#1815](https://github.com/open-feature/flagd/issues/1815)) ([c1f06cb](https://github.com/open-feature/flagd/commit/c1f06cb00fc5425d6ee73d6cfca314e9711913df))


### 🧹 Chore

* **refactor:** use memdb for flag storage ([#1697](https://github.com/open-feature/flagd/issues/1697)) ([5c5c1cf](https://github.com/open-feature/flagd/commit/5c5c1cfe84890c4cdd74c9b82504fd2632965221))


### 🔄 Refactoring

* store cleanup ([#1705](https://github.com/open-feature/flagd/issues/1705)) ([bcff8d7](https://github.com/open-feature/flagd/commit/bcff8d757b6d0ca69bccee26ba41880bdf2b5040))

## [0.12.9](https://github.com/open-feature/flagd/compare/flagd/v0.12.8...flagd/v0.12.9) (2025-07-28)


### ✨ New Features

* Add toggle for disabling getMetadata request ([#1693](https://github.com/open-feature/flagd/issues/1693)) ([e8fd680](https://github.com/open-feature/flagd/commit/e8fd6808608caa7ff7e54792fe97d88e7c294486))

## [0.12.8](https://github.com/open-feature/flagd/compare/flagd/v0.12.7...flagd/v0.12.8) (2025-07-21)


### 🐛 Bug Fixes

* update to latest otel semconv ([#1668](https://github.com/open-feature/flagd/issues/1668)) ([81855d7](https://github.com/open-feature/flagd/commit/81855d76f94a09251a19a05f830cc1d11ab6b566))


### 🧹 Chore

* **deps:** update module github.com/open-feature/flagd/core to v0.11.8 ([#1685](https://github.com/open-feature/flagd/issues/1685)) ([c07ffba](https://github.com/open-feature/flagd/commit/c07ffba55426d538224d8564be5f35339d2258d0))

## [0.12.7](https://github.com/open-feature/flagd/compare/flagd/v0.12.6...flagd/v0.12.7) (2025-07-15)


### 🧹 Chore

* **deps:** update module github.com/open-feature/flagd/core to v0.11.6 ([#1683](https://github.com/open-feature/flagd/issues/1683)) ([b6da282](https://github.com/open-feature/flagd/commit/b6da282f8a98082ba3733593d501d14842cbd97f))

## [0.12.6](https://github.com/open-feature/flagd/compare/flagd/v0.12.5...flagd/v0.12.6) (2025-07-10)


### 🐛 Bug Fixes

* **security:** update module github.com/go-viper/mapstructure/v2 to v2.3.0 [security] ([#1667](https://github.com/open-feature/flagd/issues/1667)) ([caa0ed0](https://github.com/open-feature/flagd/commit/caa0ed04eb9d5d01136deb71b8fcc4da72aa1910))


### ✨ New Features

* add sync_context to SyncFlags ([#1642](https://github.com/open-feature/flagd/issues/1642)) ([07a45d9](https://github.com/open-feature/flagd/commit/07a45d9b2275584fa92ff33cbe5e5c7d7864db38))

## [0.12.5](https://github.com/open-feature/flagd/compare/flagd/v0.12.4...flagd/v0.12.5) (2025-06-13)


### ✨ New Features

* add server-side deadline to sync service ([#1638](https://github.com/open-feature/flagd/issues/1638)) ([b70fa06](https://github.com/open-feature/flagd/commit/b70fa06b66e1fe8a28728441a7ccd28c6fe6a0c6))
* updating context using headers ([#1641](https://github.com/open-feature/flagd/issues/1641)) ([ba34815](https://github.com/open-feature/flagd/commit/ba348152b6e7b6bd7473bb11846aac7db316c88e))

## [0.12.4](https://github.com/open-feature/flagd/compare/flagd/v0.12.3...flagd/v0.12.4) (2025-05-28)


### 🐛 Bug Fixes

* Bump OpenTelemetry instrumentation dependencies ([#1616](https://github.com/open-feature/flagd/issues/1616)) ([11db29d](https://github.com/open-feature/flagd/commit/11db29dc3ab0639c3a9129c9bf449e0bd8bfc45c))


### ✨ New Features

* add traces to ofrep endpoint ([#1593](https://github.com/open-feature/flagd/issues/1593)) ([a5d43bc](https://github.com/open-feature/flagd/commit/a5d43bc1de9ca2757899678436a3e04d833da196))


### 🧹 Chore

* **deps:** update dependency go to v1.24.1 ([#1559](https://github.com/open-feature/flagd/issues/1559)) ([cd46044](https://github.com/open-feature/flagd/commit/cd4604471bba0a1df67bf87653a38df3caf9d20f))
* **security:** upgrade dependency versions ([#1632](https://github.com/open-feature/flagd/issues/1632)) ([761d870](https://github.com/open-feature/flagd/commit/761d870a3c563b8eb1b83ee543b41316c98a1d48))

## [0.12.3](https://github.com/open-feature/flagd/compare/flagd/v0.12.2...flagd/v0.12.3) (2025-03-25)


### 🐛 Bug Fixes

* add support for unix socket connection in sync service ([#1518](https://github.com/open-feature/flagd/issues/1518)) ([#1560](https://github.com/open-feature/flagd/issues/1560)) ([e2203a1](https://github.com/open-feature/flagd/commit/e2203a13ca4debd5b43c94a9b00acd2d40152bb2))
* **deps:** update module github.com/open-feature/flagd/core to v0.11.2 ([#1570](https://github.com/open-feature/flagd/issues/1570)) ([e151b1f](https://github.com/open-feature/flagd/commit/e151b1f97524a568e361103bf7a388f2598e5861))
* **deps:** update module github.com/prometheus/client_golang to v1.21.1 ([#1576](https://github.com/open-feature/flagd/issues/1576)) ([cd95193](https://github.com/open-feature/flagd/commit/cd95193f71fd465ffd1b177fa492aa84d8414a87))
* **deps:** update module google.golang.org/grpc to v1.71.0 ([#1578](https://github.com/open-feature/flagd/issues/1578)) ([5c2c64f](https://github.com/open-feature/flagd/commit/5c2c64f878b8603dd37cbfd79b0e1588e4b5a3c6))
* incorrect metadata returned per source ([#1599](https://github.com/open-feature/flagd/issues/1599)) ([b333e11](https://github.com/open-feature/flagd/commit/b333e11ecfe54f72c44ee61b3dcb1f2a487c94d4))


### ✨ New Features

* accept version numbers which are not strings ([#1589](https://github.com/open-feature/flagd/issues/1589)) ([6a13796](https://github.com/open-feature/flagd/commit/6a137967a258e799cbac9e3bb3927a07412c2a7b))

## [0.12.2](https://github.com/open-feature/flagd/compare/flagd/v0.12.1...flagd/v0.12.2) (2025-02-21)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.36.5-20250127221518-be6d1143b690.1 ([#1549](https://github.com/open-feature/flagd/issues/1549)) ([d3eb44e](https://github.com/open-feature/flagd/commit/d3eb44ed45a54bd9152b7477cce17be90016683c))
* **deps:** update module github.com/open-feature/flagd/core to v0.11.1 ([#1545](https://github.com/open-feature/flagd/issues/1545)) ([ca663b5](https://github.com/open-feature/flagd/commit/ca663b5832c94834f73cd5449a2f28af631d9556))
* **deps:** update module github.com/prometheus/client_golang to v1.21.0 ([#1568](https://github.com/open-feature/flagd/issues/1568)) ([a3d4162](https://github.com/open-feature/flagd/commit/a3d41625a2b79452c0732af29d0b4f320e74fe8b))
* **deps:** update module github.com/spf13/cobra to v1.9.0 ([#1564](https://github.com/open-feature/flagd/issues/1564)) ([345d2a9](https://github.com/open-feature/flagd/commit/345d2a920759e3e7046d7679a9c8a7cdb6cd3b40))
* **deps:** update module github.com/spf13/cobra to v1.9.1 ([#1566](https://github.com/open-feature/flagd/issues/1566)) ([a48cc80](https://github.com/open-feature/flagd/commit/a48cc8023963ac0ae41e70d4fd6fb0a9f453dba9))
* **deps:** update module golang.org/x/net to v0.35.0 ([#1557](https://github.com/open-feature/flagd/issues/1557)) ([13146e5](https://github.com/open-feature/flagd/commit/13146e5ac3de44e482f496b47dd3e0777d08c718))
* **deps:** update module google.golang.org/protobuf to v1.36.5 ([#1548](https://github.com/open-feature/flagd/issues/1548)) ([7b2b7cc](https://github.com/open-feature/flagd/commit/7b2b7ccbdf7cbebeb4b9c18bc7ab3321efe5a490))


### 🧹 Chore

* **deps:** update golang docker tag to v1.24 ([#1561](https://github.com/open-feature/flagd/issues/1561)) ([130904c](https://github.com/open-feature/flagd/commit/130904c212b1f8d484b96c05dd5996286c2922cd))

## [0.12.1](https://github.com/open-feature/flagd/compare/flagd/v0.12.0...flagd/v0.12.1) (2025-02-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.11.0 ([#1541](https://github.com/open-feature/flagd/issues/1541)) ([986a436](https://github.com/open-feature/flagd/commit/986a436e10e9766b897319085cf7dbbe2f10cb24))
* **deps:** update module golang.org/x/sync to v0.11.0 ([#1543](https://github.com/open-feature/flagd/issues/1543)) ([7d6c0dc](https://github.com/open-feature/flagd/commit/7d6c0dc6e6e6955af1e5225807deeb2b6797900b))

## [0.12.0](https://github.com/open-feature/flagd/compare/flagd/v0.11.8...flagd/v0.12.0) (2025-01-31)


### ⚠ BREAKING CHANGES

* flagSetMetadata in OFREP/ResolveAll, core refactors ([#1540](https://github.com/open-feature/flagd/issues/1540))

### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.18.1-20250127221518-be6d1143b690.1 ([#1535](https://github.com/open-feature/flagd/issues/1535)) ([d5ec921](https://github.com/open-feature/flagd/commit/d5ec921f02263615cd87625e211ba82b0811d041))
* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20250127221518-be6d1143b690.2 ([#1536](https://github.com/open-feature/flagd/issues/1536)) ([e23060f](https://github.com/open-feature/flagd/commit/e23060f24b2a714ae748e6b37d0d06b7caa1c95c))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.36.4-20241220192239-696330adaff0.1 ([#1529](https://github.com/open-feature/flagd/issues/1529)) ([8881a80](https://github.com/open-feature/flagd/commit/8881a804b4055da0127a16b8fc57022d24906e1b))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.36.4-20250127221518-be6d1143b690.1 ([#1537](https://github.com/open-feature/flagd/issues/1537)) ([f74207b](https://github.com/open-feature/flagd/commit/f74207bc13b75bae4275bc486df51e2da569dd41))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.8 ([#1526](https://github.com/open-feature/flagd/issues/1526)) ([fbf2ed5](https://github.com/open-feature/flagd/commit/fbf2ed527fcf3b300808c7b835a8d891df7b88a8))
* **deps:** update module google.golang.org/grpc to v1.70.0 ([#1528](https://github.com/open-feature/flagd/issues/1528)) ([79b2b0a](https://github.com/open-feature/flagd/commit/79b2b0a6bbd48676dcbdd2393feb8247529bf29c))


### ✨ New Features

* flagSetMetadata in OFREP/ResolveAll, core refactors ([#1540](https://github.com/open-feature/flagd/issues/1540)) ([b49abf9](https://github.com/open-feature/flagd/commit/b49abf95069da93bdf8369c8aa0ae40e698df760))

## [0.11.8](https://github.com/open-feature/flagd/compare/flagd/v0.11.7...flagd/v0.11.8) (2025-01-19)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.10.7 ([#1521](https://github.com/open-feature/flagd/issues/1521)) ([bf8e7e0](https://github.com/open-feature/flagd/commit/bf8e7e06d9b34e0548abb8af9cce2becb357f9c4))
* **deps:** update opentelemetry-go monorepo ([#1524](https://github.com/open-feature/flagd/issues/1524)) ([eeae9a6](https://github.com/open-feature/flagd/commit/eeae9a64caf93356fd663cc735cc422edcf9e132))
* Skip flagd banner when non-console logger in use ([#1516](https://github.com/open-feature/flagd/issues/1516)) ([bae9b6f](https://github.com/open-feature/flagd/commit/bae9b6fb3b53a9d73f4c36e7b676beb6dac03476))

## [0.11.7](https://github.com/open-feature/flagd/compare/flagd/v0.11.6...flagd/v0.11.7) (2025-01-16)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.36.3-20241220192239-696330adaff0.1 ([#1513](https://github.com/open-feature/flagd/issues/1513)) ([64c5787](https://github.com/open-feature/flagd/commit/64c57875b032edcef2e2d230e7735990e01b72b8))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.6 ([#1515](https://github.com/open-feature/flagd/issues/1515)) ([586cb62](https://github.com/open-feature/flagd/commit/586cb62e63d66c8f8371236844d506c7bcc8f123))
* **sync:** fixing missing handover of ssl configuration ([#1517](https://github.com/open-feature/flagd/issues/1517)) ([998a216](https://github.com/open-feature/flagd/commit/998a21637e512aec8d1944f842c51529e7bf6d44))

## [0.11.6](https://github.com/open-feature/flagd/compare/flagd/v0.11.5...flagd/v0.11.6) (2025-01-15)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.17.0-20241220192239-696330adaff0.1 ([#1488](https://github.com/open-feature/flagd/issues/1488)) ([8e09457](https://github.com/open-feature/flagd/commit/8e09457d739fc7d9e679e27a3449e4a13cd6c718))
* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.18.1-20241220192239-696330adaff0.1 ([#1506](https://github.com/open-feature/flagd/issues/1506)) ([b868194](https://github.com/open-feature/flagd/commit/b8681947bf5dcc24d5b0e48f5e59a61c1278c280))
* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20241220192239-696330adaff0.1 ([#1489](https://github.com/open-feature/flagd/issues/1489)) ([53add83](https://github.com/open-feature/flagd/commit/53add83a491c6e00e0d9b1b64a9461e5973edca7))
* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20241220192239-696330adaff0.2 ([#1492](https://github.com/open-feature/flagd/issues/1492)) ([9f1d94a](https://github.com/open-feature/flagd/commit/9f1d94a42ac00ecf5fc58c07a76c350e2e4ec2f6))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.36.0-20241220192239-696330adaff0.1 ([#1490](https://github.com/open-feature/flagd/issues/1490)) ([6edce72](https://github.com/open-feature/flagd/commit/6edce72e8cff01ea13cbd15d604b35ccc8337f50))
* **deps:** update module github.com/mattn/go-colorable to v0.1.14 ([#1508](https://github.com/open-feature/flagd/issues/1508)) ([87727f7](https://github.com/open-feature/flagd/commit/87727f7f8f18e4f532d152190ada5dbe3fc915b0))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.5 ([#1482](https://github.com/open-feature/flagd/issues/1482)) ([ce48cb7](https://github.com/open-feature/flagd/commit/ce48cb757659eef8807531f8522ca1b7bc80422c))
* **deps:** update module golang.org/x/net to v0.33.0 [security] ([#1486](https://github.com/open-feature/flagd/issues/1486)) ([4764077](https://github.com/open-feature/flagd/commit/476407769f47675f649c328e27e0f87860f0f79d))
* **deps:** update module golang.org/x/net to v0.34.0 ([#1498](https://github.com/open-feature/flagd/issues/1498)) ([7584f95](https://github.com/open-feature/flagd/commit/7584f95e4da50ae870014589a971b83b10c23873))
* **deps:** update module google.golang.org/grpc to v1.69.2 ([#1484](https://github.com/open-feature/flagd/issues/1484)) ([6b40ad3](https://github.com/open-feature/flagd/commit/6b40ad34c83da4a3116e7cad4139a63a6c918097))
* **deps:** update module google.golang.org/grpc to v1.69.4 ([#1510](https://github.com/open-feature/flagd/issues/1510)) ([76d6353](https://github.com/open-feature/flagd/commit/76d6353840ab8e7c93bdb0802eb1c49fc6fe1dc0))
* **deps:** update module google.golang.org/protobuf to v1.36.1 ([#1491](https://github.com/open-feature/flagd/issues/1491)) ([2c729a7](https://github.com/open-feature/flagd/commit/2c729a7c212479e689a4cb4f6a289bee8accb194))
* **deps:** update opentelemetry-go monorepo ([#1470](https://github.com/open-feature/flagd/issues/1470)) ([26b0b1a](https://github.com/open-feature/flagd/commit/26b0b1af8bc4b3a393c3453784b50f167f13f743))


### ✨ New Features

* add ssl support to sync service ([#1479](https://github.com/open-feature/flagd/issues/1479)) ([#1501](https://github.com/open-feature/flagd/issues/1501)) ([d50fcc8](https://github.com/open-feature/flagd/commit/d50fcc821c1ae043cb8cf77e464f7b738e2ff755))

## [0.11.5](https://github.com/open-feature/flagd/compare/flagd/v0.11.4...flagd/v0.11.5) (2024-12-17)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.35.2-20240906125204-0a6a901b42e8.1 ([#1451](https://github.com/open-feature/flagd/issues/1451)) ([8c6d91d](https://github.com/open-feature/flagd/commit/8c6d91d538d226b10cb954c23409902e9d245cda))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.36.0-20240906125204-0a6a901b42e8.1 ([#1475](https://github.com/open-feature/flagd/issues/1475)) ([0b11c6c](https://github.com/open-feature/flagd/commit/0b11c6cf612b244bda6bab119814647f3ce8de2e))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.4 ([#1433](https://github.com/open-feature/flagd/issues/1433)) ([d33c7a5](https://github.com/open-feature/flagd/commit/d33c7a5522d0909448c6d9d80b0a33d8511f0738))
* **deps:** update module github.com/stretchr/testify to v1.10.0 ([#1455](https://github.com/open-feature/flagd/issues/1455)) ([8c843df](https://github.com/open-feature/flagd/commit/8c843df7714b1f2d120c5cac8e40c7220cc0c05b))
* **deps:** update module golang.org/x/net to v0.31.0 ([#1446](https://github.com/open-feature/flagd/issues/1446)) ([9e35111](https://github.com/open-feature/flagd/commit/9e351117b4b2ebbb4a016d6b189077ae65a83124))
* **deps:** update module golang.org/x/net to v0.32.0 ([#1458](https://github.com/open-feature/flagd/issues/1458)) ([ac0b123](https://github.com/open-feature/flagd/commit/ac0b123ce84a0772f144ae0ae8f3153992635ea4))
* **deps:** update module golang.org/x/sync to v0.9.0 ([#1445](https://github.com/open-feature/flagd/issues/1445)) ([8893e94](https://github.com/open-feature/flagd/commit/8893e94b94ae79f80a0aa0f25cca5caf874e9d2e))
* **deps:** update module google.golang.org/grpc to v1.68.0 ([#1442](https://github.com/open-feature/flagd/issues/1442)) ([cd27d09](https://github.com/open-feature/flagd/commit/cd27d098e6d8d8b0f681ef42d26dba1ebac67d12))
* **deps:** update module google.golang.org/grpc to v1.68.1 ([#1456](https://github.com/open-feature/flagd/issues/1456)) ([0b6e2a1](https://github.com/open-feature/flagd/commit/0b6e2a1cd64910226d348c921b08a6de8013ac90))
* **deps:** update module google.golang.org/grpc to v1.69.0 ([#1469](https://github.com/open-feature/flagd/issues/1469)) ([dd4869f](https://github.com/open-feature/flagd/commit/dd4869f5e095066f80c9d82d1be83155e7504d88))
* **deps:** update module google.golang.org/protobuf to v1.35.2 ([#1450](https://github.com/open-feature/flagd/issues/1450)) ([6b9834d](https://github.com/open-feature/flagd/commit/6b9834d4b3aa9b71cf22cd0e463b1807de164ccf))
* **deps:** update module google.golang.org/protobuf to v1.36.0 ([#1474](https://github.com/open-feature/flagd/issues/1474)) ([6a8a9a9](https://github.com/open-feature/flagd/commit/6a8a9a981c29f6a4f2bc617d4219f63d765f0a5b))
* **deps:** update opentelemetry-go monorepo ([#1447](https://github.com/open-feature/flagd/issues/1447)) ([68b5794](https://github.com/open-feature/flagd/commit/68b5794180da84af9adc1f2cd80f929489969c1c))


### ✨ New Features

* add context-value flag ([#1448](https://github.com/open-feature/flagd/issues/1448)) ([7ca092e](https://github.com/open-feature/flagd/commit/7ca092e478c937eca0c91357394499763545dc1c))

## [0.11.4](https://github.com/open-feature/flagd/compare/flagd/v0.11.3...flagd/v0.11.4) (2024-10-28)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.17.0-20240906125204-0a6a901b42e8.1 ([#1409](https://github.com/open-feature/flagd/issues/1409)) ([f07d348](https://github.com/open-feature/flagd/commit/f07d348b42bfa43c3b263ed7ba69c1630a14ec10))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.35.1-20240906125204-0a6a901b42e8.1 ([#1420](https://github.com/open-feature/flagd/issues/1420)) ([1f06d5a](https://github.com/open-feature/flagd/commit/1f06d5a1837ea2b753974e96c2a1154d6cb3e582))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.3 ([#1411](https://github.com/open-feature/flagd/issues/1411)) ([a312196](https://github.com/open-feature/flagd/commit/a312196c118705d7a8eb0056fdb98480b887f7c5))
* **deps:** update module github.com/prometheus/client_golang to v1.20.5 ([#1425](https://github.com/open-feature/flagd/issues/1425)) ([583ba89](https://github.com/open-feature/flagd/commit/583ba894f2de794b36b6a1cc3bfceb9c46dc9d96))
* **deps:** update module go.uber.org/mock to v0.5.0 ([#1427](https://github.com/open-feature/flagd/issues/1427)) ([0c6fd7f](https://github.com/open-feature/flagd/commit/0c6fd7fa688db992d4e58a202889cbfea07eebf6))
* **deps:** update module golang.org/x/net to v0.30.0 ([#1417](https://github.com/open-feature/flagd/issues/1417)) ([4d5b75e](https://github.com/open-feature/flagd/commit/4d5b75eed9097c09760fcc71bfdf473cd19232ec))
* **deps:** update module google.golang.org/grpc to v1.67.1 ([#1415](https://github.com/open-feature/flagd/issues/1415)) ([85a3a6b](https://github.com/open-feature/flagd/commit/85a3a6b46233fcc7cf71a0292b46c82ac8e66d7b))


### ✨ New Features

* added custom grpc resolver ([#1424](https://github.com/open-feature/flagd/issues/1424)) ([e5007e2](https://github.com/open-feature/flagd/commit/e5007e2bcb6f049a3c54e09331065bb9abe215be))
* support azure blob sync ([#1428](https://github.com/open-feature/flagd/issues/1428)) ([5c39cfe](https://github.com/open-feature/flagd/commit/5c39cfe30a3dead4f6db2c6f9ee4c12193cd479b))

## [0.11.3](https://github.com/open-feature/flagd/compare/flagd/v0.11.2...flagd/v0.11.3) (2024-09-23)


### 🐛 Bug Fixes

* **deps:** update kubernetes package and controller runtime, fix proto lint ([#1290](https://github.com/open-feature/flagd/issues/1290)) ([94860d6](https://github.com/open-feature/flagd/commit/94860d6ceabe9eb7c1e5dd8ea139a796710d6d8b))
* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.16.2-20240906125204-0a6a901b42e8.1 ([#1399](https://github.com/open-feature/flagd/issues/1399)) ([18dd4e2](https://github.com/open-feature/flagd/commit/18dd4e279f1278938749bd21d106ecfbaf2f5fff))
* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20240906125204-0a6a901b42e8.1 ([#1400](https://github.com/open-feature/flagd/issues/1400)) ([954d972](https://github.com/open-feature/flagd/commit/954d97238210f90b650493ae76277d4a8d80788a))
* **deps:** update module connectrpc.com/connect to v1.17.0 ([#1408](https://github.com/open-feature/flagd/issues/1408)) ([e7eb691](https://github.com/open-feature/flagd/commit/e7eb691094dfbf02e37d79c41f60f556415e7640))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.2 ([#1385](https://github.com/open-feature/flagd/issues/1385)) ([3b5a818](https://github.com/open-feature/flagd/commit/3b5a818b69ffca61347a3feaa85dd1a8f8001e24))
* **deps:** update module github.com/prometheus/client_golang to v1.20.3 ([#1384](https://github.com/open-feature/flagd/issues/1384)) ([8fd16b2](https://github.com/open-feature/flagd/commit/8fd16b23b1fa8517128af36b3068ca18ebbad6c3))
* **deps:** update module github.com/prometheus/client_golang to v1.20.4 ([#1406](https://github.com/open-feature/flagd/issues/1406)) ([a0a6426](https://github.com/open-feature/flagd/commit/a0a64269b08251317676075fdea7bc65bea8a8dc))
* **deps:** update module github.com/rs/cors to v1.11.1 ([#1392](https://github.com/open-feature/flagd/issues/1392)) ([8bd549e](https://github.com/open-feature/flagd/commit/8bd549e8603b4a61cc26d0c09ef5dd13860cb3da))
* **deps:** update module github.com/rs/xid to v1.6.0 ([#1386](https://github.com/open-feature/flagd/issues/1386)) ([2317013](https://github.com/open-feature/flagd/commit/231701346abab2018cc7c495ebc7f285bb2a46d2))
* **deps:** update module golang.org/x/net to v0.29.0 ([#1398](https://github.com/open-feature/flagd/issues/1398)) ([0721e02](https://github.com/open-feature/flagd/commit/0721e02daae4c92438490169113d3d76ca4a028a))
* **deps:** update module google.golang.org/grpc to v1.66.0 ([#1393](https://github.com/open-feature/flagd/issues/1393)) ([c96e9d7](https://github.com/open-feature/flagd/commit/c96e9d764aa51caf00fbde07cdc7d2de55b98b9e))
* **deps:** update module google.golang.org/grpc to v1.66.1 ([#1402](https://github.com/open-feature/flagd/issues/1402)) ([50c9cd3](https://github.com/open-feature/flagd/commit/50c9cd3ada2f470a22374392a5a152a487636645))
* **deps:** update module google.golang.org/grpc to v1.66.2 ([#1405](https://github.com/open-feature/flagd/issues/1405)) ([69ec28f](https://github.com/open-feature/flagd/commit/69ec28fceb597bdaad63b184943b66ccdb4af0b7))
* **deps:** update module google.golang.org/grpc to v1.67.0 ([#1407](https://github.com/open-feature/flagd/issues/1407)) ([1ad6480](https://github.com/open-feature/flagd/commit/1ad6480a0f37c4677e53065ef455f615b26b1f17))
* **deps:** update opentelemetry-go monorepo ([#1387](https://github.com/open-feature/flagd/issues/1387)) ([22aef5b](https://github.com/open-feature/flagd/commit/22aef5bbf030c619e48fbe22a16d83e071b11902))
* **deps:** update opentelemetry-go monorepo ([#1403](https://github.com/open-feature/flagd/issues/1403)) ([fc4cd3e](https://github.com/open-feature/flagd/commit/fc4cd3e547f4826ea0bb8cc1bb2304807932b4e6))
* remove dep cycle with certreloader ([#1410](https://github.com/open-feature/flagd/issues/1410)) ([5244f6f](https://github.com/open-feature/flagd/commit/5244f6f6c94f310fd80c7ab84942103cc8c18a39))


### ✨ New Features

* add mTLS support to otel exporter ([#1389](https://github.com/open-feature/flagd/issues/1389)) ([8737f53](https://github.com/open-feature/flagd/commit/8737f53444016b114ee4ae52eead0b835af0e200))

## [0.11.2](https://github.com/open-feature/flagd/compare/flagd/v0.11.1...flagd/v0.11.2) (2024-08-22)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.5.1-20240215170432-1e611e2999cc.1 ([#1372](https://github.com/open-feature/flagd/issues/1372)) ([ae24595](https://github.com/open-feature/flagd/commit/ae2459504f7eccafebccec83fa1f72b08f41a978))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.1 ([#1355](https://github.com/open-feature/flagd/issues/1355)) ([8fcfb14](https://github.com/open-feature/flagd/commit/8fcfb146b0c55712c1758201ee4bc59e83b0898c))
* **deps:** update module golang.org/x/net to v0.28.0 ([#1380](https://github.com/open-feature/flagd/issues/1380)) ([239a432](https://github.com/open-feature/flagd/commit/239a432c18bf6780117b5d563443124887b38120))
* **deps:** update module golang.org/x/sync to v0.8.0 ([#1378](https://github.com/open-feature/flagd/issues/1378)) ([4804c17](https://github.com/open-feature/flagd/commit/4804c17a67ea9761079ecade34ccb3446643050b))


### 🧹 Chore

* **deps:** update dependency go to v1.22.6 ([#1297](https://github.com/open-feature/flagd/issues/1297)) ([50b92c1](https://github.com/open-feature/flagd/commit/50b92c17cfd872d3e6b95fef3b3d96444e563715))
* **deps:** update golang docker tag to v1.23 ([#1382](https://github.com/open-feature/flagd/issues/1382)) ([abb5ca3](https://github.com/open-feature/flagd/commit/abb5ca3e31308535c66a94300d6f6409fd370b95))
* improve gRPC sync service shutdown behavior ([#1375](https://github.com/open-feature/flagd/issues/1375)) ([79d9085](https://github.com/open-feature/flagd/commit/79d9085a50c49a97b70febb5f444fa3ea965220b))

## [0.11.1](https://github.com/open-feature/flagd/compare/flagd/v0.11.0...flagd/v0.11.1) (2024-07-08)


### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.4.0-20240215170432-1e611e2999cc.2 ([#1342](https://github.com/open-feature/flagd/issues/1342)) ([efdd921](https://github.com/open-feature/flagd/commit/efdd92139903b89ac986a62ff2cf4f5cfef91cde))
* **deps:** update module github.com/open-feature/flagd/core to v0.10.0 ([#1340](https://github.com/open-feature/flagd/issues/1340)) ([1e487b4](https://github.com/open-feature/flagd/commit/1e487b4bafad9814f190d0bf3a1d833def9ef5af))
* **deps:** update module golang.org/x/net to v0.27.0 ([#1353](https://github.com/open-feature/flagd/issues/1353)) ([df9834b](https://github.com/open-feature/flagd/commit/df9834bea2a7ae20c5926c98dc423ab6363ef332))
* **deps:** update module google.golang.org/grpc to v1.65.0 ([#1346](https://github.com/open-feature/flagd/issues/1346)) ([72a6b87](https://github.com/open-feature/flagd/commit/72a6b876e880ff0b43440d9b63710c7a87536988))
* **deps:** update opentelemetry-go monorepo ([#1347](https://github.com/open-feature/flagd/issues/1347)) ([37fb3cd](https://github.com/open-feature/flagd/commit/37fb3cd81d5436e9d8cd3ea490a3951ae5794130))
* invalid scoped-sync responses for empty flags ([#1352](https://github.com/open-feature/flagd/issues/1352)) ([51371d2](https://github.com/open-feature/flagd/commit/51371d25e25e1199336a5a831530506313628ff3))

## [0.11.0](https://github.com/open-feature/flagd/compare/flagd/v0.10.3...flagd/v0.11.0) (2024-06-27)


### ⚠ BREAKING CHANGES

* support emitting errors from the bulk evaluator ([#1338](https://github.com/open-feature/flagd/issues/1338))

### 🐛 Bug Fixes

* **deps:** update module buf.build/gen/go/open-feature/flagd/connectrpc/go to v1.16.2-20240215170432-1e611e2999cc.1 ([#1293](https://github.com/open-feature/flagd/issues/1293)) ([2694e7f](https://github.com/open-feature/flagd/commit/2694e7f99fc356c23b6b34265418b0b0160deb62))
* **deps:** update module buf.build/gen/go/open-feature/flagd/grpc/go to v1.4.0-20240215170432-1e611e2999cc.1 ([#1333](https://github.com/open-feature/flagd/issues/1333)) ([494062f](https://github.com/open-feature/flagd/commit/494062fed891fab0fb659352142dbbc97c8f1492))
* **deps:** update module buf.build/gen/go/open-feature/flagd/protocolbuffers/go to v1.34.2-20240215170432-1e611e2999cc.2 ([#1330](https://github.com/open-feature/flagd/issues/1330)) ([32291ad](https://github.com/open-feature/flagd/commit/32291ad93d25d79299a7a02381df70e2719c4fbc))
* **deps:** update module github.com/open-feature/flagd/core to v0.9.3 ([#1296](https://github.com/open-feature/flagd/issues/1296)) ([1f7b8bd](https://github.com/open-feature/flagd/commit/1f7b8bd938f799da98462e45e52c0e1ac6cb83e6))
* **deps:** update module github.com/rs/cors to v1.11.0 ([#1299](https://github.com/open-feature/flagd/issues/1299)) ([5f77541](https://github.com/open-feature/flagd/commit/5f775413fb33b4afed1ef98484463f07a67ed1cb))
* **deps:** update module github.com/spf13/cobra to v1.8.1 ([#1332](https://github.com/open-feature/flagd/issues/1332)) ([c62bcb0](https://github.com/open-feature/flagd/commit/c62bcb0ec68fbcac40d16df001379f117c4df37e))
* **deps:** update module github.com/spf13/viper to v1.19.0 ([#1334](https://github.com/open-feature/flagd/issues/1334)) ([1097b99](https://github.com/open-feature/flagd/commit/1097b9961b672d44a81e5b9e7a56f163e08e4909))
* **deps:** update module golang.org/x/net to v0.26.0 ([#1337](https://github.com/open-feature/flagd/issues/1337)) ([83bdbb5](https://github.com/open-feature/flagd/commit/83bdbb5e7ea1be9da51d06e6b22c997f0354ef98))
* **deps:** update opentelemetry-go monorepo ([#1314](https://github.com/open-feature/flagd/issues/1314)) ([e9f1a7a](https://github.com/open-feature/flagd/commit/e9f1a7a04828f36691e694375b3c665140bc7dee))
* readable error messages  ([#1325](https://github.com/open-feature/flagd/issues/1325)) ([7ff33ef](https://github.com/open-feature/flagd/commit/7ff33effcc47e31c5b7fdc33385d8128db2163fc))


### ✨ New Features

* support `FLAGD_DEBUG` / `--debug` / `-x` ([#1326](https://github.com/open-feature/flagd/issues/1326)) ([298bd36](https://github.com/open-feature/flagd/commit/298bd36698224a0dca8b289f4cb0b80ae2fa6e0a))
* support emitting errors from the bulk evaluator ([#1338](https://github.com/open-feature/flagd/issues/1338)) ([b9c099c](https://github.com/open-feature/flagd/commit/b9c099cb7fa002a509a82c81b467f5e784c27e82))

## [0.10.3](https://github.com/open-feature/flagd/compare/flagd/v0.10.2...flagd/v0.10.3) (2024-06-06)


### 🧹 Chore

* adapt telemetry setup error handling ([#1315](https://github.com/open-feature/flagd/issues/1315)) ([20bcb78](https://github.com/open-feature/flagd/commit/20bcb78d11dbb16aab2b14d5869bb990a0f7bca5))
* fix unit tests and ensure their execution ([#1316](https://github.com/open-feature/flagd/issues/1316)) ([25041c0](https://github.com/open-feature/flagd/commit/25041c016ae84afb01b8eb1e2b693aae3199a6ac))

## [0.10.2](https://github.com/open-feature/flagd/compare/flagd/v0.10.1...flagd/v0.10.2) (2024-05-10)


### ✨ New Features

* Create interface for eval events.  ([#1288](https://github.com/open-feature/flagd/issues/1288)) ([9714215](https://github.com/open-feature/flagd/commit/9714215cedb0fd28daddf086ce1255ec29b877d4))


### 🧹 Chore

* bump go deps to latest ([#1307](https://github.com/open-feature/flagd/issues/1307)) ([004ad08](https://github.com/open-feature/flagd/commit/004ad083dc01538791148d6233e453d2a3009fcd))

## [0.10.1](https://github.com/open-feature/flagd/compare/flagd/v0.10.0...flagd/v0.10.1) (2024-04-19)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.9.0 ([#1281](https://github.com/open-feature/flagd/issues/1281)) ([3cfb052](https://github.com/open-feature/flagd/commit/3cfb0523cc857dd2019d712c621afe81c2b41398))


### ✨ New Features

* move json logic operator registration to resolver ([#1291](https://github.com/open-feature/flagd/issues/1291)) ([b473457](https://github.com/open-feature/flagd/commit/b473457ddff28789fee1eeb6704491b6aa3525e3))

## [0.10.0](https://github.com/open-feature/flagd/compare/flagd/v0.9.2...flagd/v0.10.0) (2024-04-10)


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

* refactor evaluation core ([#1259](https://github.com/open-feature/flagd/issues/1259)) ([0e6604c](https://github.com/open-feature/flagd/commit/0e6604cd038dc13d7d40e622523320bf03efbcd0))
* update go deps ([#1279](https://github.com/open-feature/flagd/issues/1279)) ([219789f](https://github.com/open-feature/flagd/commit/219789fca8a929d552e4e8d1f6b6d5cd44505f43))

## [0.9.2](https://github.com/open-feature/flagd/compare/flagd/v0.9.1...flagd/v0.9.2) (2024-03-27)


### ✨ New Features

* OFREP support for flagd  ([#1247](https://github.com/open-feature/flagd/issues/1247)) ([9d12fc2](https://github.com/open-feature/flagd/commit/9d12fc20702a86e8385564659be88f07ad36d9e5))

## [0.9.1](https://github.com/open-feature/flagd/compare/flagd/v0.9.0...flagd/v0.9.1) (2024-03-15)


### 🐛 Bug Fixes

* **deps:** update module google.golang.org/protobuf to v1.33.0 [security] ([#1248](https://github.com/open-feature/flagd/issues/1248)) ([b2b0fa1](https://github.com/open-feature/flagd/commit/b2b0fa19a6254c02c81ef44828b643a5a25ea5b5))
* update protobuff CVE-2024-24786 ([#1249](https://github.com/open-feature/flagd/issues/1249)) ([fd81c23](https://github.com/open-feature/flagd/commit/fd81c235fb4a09dfc42289ac316ac3a1d7eff58c))


### ✨ New Features

* serve sync.proto on port 8015 ([#1237](https://github.com/open-feature/flagd/issues/1237)) ([7afdc0c](https://github.com/open-feature/flagd/commit/7afdc0cda47d080575cb87a94b35cfe051f88422))


### 🧹 Chore

* move packaging & isolate service implementations  ([#1234](https://github.com/open-feature/flagd/issues/1234)) ([b58fab3](https://github.com/open-feature/flagd/commit/b58fab3df030ef7e9e10eafa7a0141c05aa05bbd))

## [0.9.0](https://github.com/open-feature/flagd/compare/flagd/v0.8.2...flagd/v0.9.0) (2024-02-20)


### ⚠ BREAKING CHANGES

* new proto (flagd.sync.v1) for sync sources ([#1214](https://github.com/open-feature/flagd/issues/1214))

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.5 ([#1198](https://github.com/open-feature/flagd/issues/1198)) ([ce38845](https://github.com/open-feature/flagd/commit/ce388458b9c8a686a7b6ff38b532c941d43d842c))


### ✨ New Features

* new proto (flagd.sync.v1) for sync sources ([#1214](https://github.com/open-feature/flagd/issues/1214)) ([544234e](https://github.com/open-feature/flagd/commit/544234ebd9f9be5f54c2865a866575a7869a56c0))


### 🧹 Chore

* **deps:** update golang docker tag to v1.22 ([#1201](https://github.com/open-feature/flagd/issues/1201)) ([d14c69e](https://github.com/open-feature/flagd/commit/d14c69e93e56d32a37b2428f1db2d4ac79563597))

## [0.8.2](https://github.com/open-feature/flagd/compare/flagd/v0.8.1...flagd/v0.8.2) (2024-02-05)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.4 ([#1119](https://github.com/open-feature/flagd/issues/1119)) ([e998e41](https://github.com/open-feature/flagd/commit/e998e41f7c6fc8007458dff08e66aa19c7b7b0e7))
* use correct link in sources flag helper text in start cmd ([#1126](https://github.com/open-feature/flagd/issues/1126)) ([b9d30e0](https://github.com/open-feature/flagd/commit/b9d30e0a52eaf50553e1ce4c65f60bc67d931ea6))

## [0.8.1](https://github.com/open-feature/flagd/compare/flagd/v0.8.0...flagd/v0.8.1) (2024-01-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.3 ([#1104](https://github.com/open-feature/flagd/issues/1104)) ([b6c00c7](https://github.com/open-feature/flagd/commit/b6c00c7615040399b60f9085a8238d417445546d))
* **deps:** update module github.com/spf13/viper to v1.18.2 ([#1069](https://github.com/open-feature/flagd/issues/1069)) ([f0d6206](https://github.com/open-feature/flagd/commit/f0d620698abbde6ef455c2dd64b02a52eac96a89))

## [0.8.0](https://github.com/open-feature/flagd/compare/flagd/v0.7.2...flagd/v0.8.0) (2023-12-22)


### ⚠ BREAKING CHANGES

* remove deprecated flags ([#1075](https://github.com/open-feature/flagd/issues/1075))

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.2 ([#1056](https://github.com/open-feature/flagd/issues/1056)) ([81e83ea](https://github.com/open-feature/flagd/commit/81e83ea0a4aa78d853ea7700cb06bb2a0f329619))
* **deps:** update module github.com/spf13/viper to v1.18.0 ([#1060](https://github.com/open-feature/flagd/issues/1060)) ([9dfa689](https://github.com/open-feature/flagd/commit/9dfa6899ed3a25a5c34f8b0ebd152b01b1097dec))


### 🧹 Chore

* refactoring component structure ([#1044](https://github.com/open-feature/flagd/issues/1044)) ([0c7f78a](https://github.com/open-feature/flagd/commit/0c7f78a95fa4ad2a8b2afe2f6023b9c6d4fd48ed))
* remove deprecated flags ([#1075](https://github.com/open-feature/flagd/issues/1075)) ([49f6fe5](https://github.com/open-feature/flagd/commit/49f6fe5679425b31b1e1cf39a2a2e4767b2e1db9))

## [0.7.2](https://github.com/open-feature/flagd/compare/flagd/v0.7.1...flagd/v0.7.2) (2023-12-05)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.1 ([#1037](https://github.com/open-feature/flagd/issues/1037)) ([0ed9b68](https://github.com/open-feature/flagd/commit/0ed9b68341d026681c684a726b215ff910fe2a00))

## [0.7.1](https://github.com/open-feature/flagd/compare/flagd/v0.7.0...flagd/v0.7.1) (2023-11-28)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.7.0 ([#1014](https://github.com/open-feature/flagd/issues/1014)) ([deec49e](https://github.com/open-feature/flagd/commit/deec49e99ef52f62adbf278a8f58936acbb86b9d))


### 🔄 Refactoring

* Rename metrics-port to management-port ([#1012](https://github.com/open-feature/flagd/issues/1012)) ([5635e38](https://github.com/open-feature/flagd/commit/5635e38703cae835a53e9cce83d5bc42d00091e2))

## [0.7.0](https://github.com/open-feature/flagd/compare/flagd/v0.6.8...flagd/v0.7.0) (2023-11-15)


### ⚠ BREAKING CHANGES

* OFO APIs were updated to version v1beta1, since they are more stable now. Resources of the alpha versions are no longer supported in flagd or flagd-proxy.

### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/go-sdk-contrib/providers/flagd to v0.1.18 ([#1011](https://github.com/open-feature/flagd/issues/1011)) ([90d4e4e](https://github.com/open-feature/flagd/commit/90d4e4e7d9db9e21fa38d96fdecb81ab78868732))


### ✨ New Features

* support OFO v1beta1 API ([#997](https://github.com/open-feature/flagd/issues/997)) ([bb6f5bf](https://github.com/open-feature/flagd/commit/bb6f5bf0fc382ade75d80a34d209beaa2edc459d))


### 🧹 Chore

* move e2e tests to test ([#1005](https://github.com/open-feature/flagd/issues/1005)) ([a94b639](https://github.com/open-feature/flagd/commit/a94b6399e529ca03c6034eb86ec4028d7e8c2a82))

## [0.6.8](https://github.com/open-feature/flagd/compare/flagd/v0.6.7...flagd/v0.6.8) (2023-11-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.7 ([#966](https://github.com/open-feature/flagd/issues/966)) ([c038a3a](https://github.com/open-feature/flagd/commit/c038a3a3700eee82afa3e2cb2484614ec6ed566c))
* **deps:** update module github.com/open-feature/go-sdk to v1.8.0 ([#994](https://github.com/open-feature/flagd/issues/994)) ([266cf9f](https://github.com/open-feature/flagd/commit/266cf9f82ee8b4a4ba8ad1c0594388d2987a8c4b))
* **deps:** update module github.com/open-feature/go-sdk-contrib/tests/flagd to v1.3.1 ([#760](https://github.com/open-feature/flagd/issues/760)) ([30dda72](https://github.com/open-feature/flagd/commit/30dda72145c05de298140f880238ed37be73631a))
* **deps:** update module github.com/spf13/cobra to v1.8.0 ([#993](https://github.com/open-feature/flagd/issues/993)) ([05c7870](https://github.com/open-feature/flagd/commit/05c7870cc7662117f85e9c6528508327ae320b83))


### 🧹 Chore

* fix lint errors ([#987](https://github.com/open-feature/flagd/issues/987)) ([0c3af2d](https://github.com/open-feature/flagd/commit/0c3af2da01f91f6fc6d5ac78a33dd79032537ea9))


### 🔄 Refactoring

* migrate to connectrpc/connect-go ([#990](https://github.com/open-feature/flagd/issues/990)) ([7dd5b2b](https://github.com/open-feature/flagd/commit/7dd5b2b4c284481bcba5a9c45bd6c85ad1dc6d33))

## [0.6.7](https://github.com/open-feature/flagd/compare/flagd/v0.6.6...flagd/v0.6.7) (2023-10-12)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.6 ([#916](https://github.com/open-feature/flagd/issues/916)) ([1f80e4d](https://github.com/open-feature/flagd/commit/1f80e4db9f8d1ba24884a71f2f8d552499ab5fe2))
* **deps:** update module github.com/open-feature/go-sdk-contrib/providers/flagd to v0.1.17 ([#759](https://github.com/open-feature/flagd/issues/759)) ([a2a2c3c](https://github.com/open-feature/flagd/commit/a2a2c3c7effd1708136eaac5df00ae02276d5005))
* **deps:** update module github.com/spf13/viper to v1.17.0 ([#956](https://github.com/open-feature/flagd/issues/956)) ([31d015d](https://github.com/open-feature/flagd/commit/31d015d329ae9c1da3ec13878078371bcbf43fbf))
* **deps:** update module go.uber.org/zap to v1.26.0 ([#917](https://github.com/open-feature/flagd/issues/917)) ([e57e206](https://github.com/open-feature/flagd/commit/e57e206c937d5b11b81d46ee57b3e92cc454dd88))


### 🧹 Chore

* docs rework ([#927](https://github.com/open-feature/flagd/issues/927)) ([27b3193](https://github.com/open-feature/flagd/commit/27b31938210c8930d9cbb31c1c76220d185b3949))


### 📚 Documentation

* fixed typos and linting issues ([#957](https://github.com/open-feature/flagd/issues/957)) ([0bade57](https://github.com/open-feature/flagd/commit/0bade574005f8faf977de30b14ac89acbb276472))

## [0.6.6](https://github.com/open-feature/flagd/compare/flagd/v0.6.5...flagd/v0.6.6) (2023-09-14)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.5 ([#900](https://github.com/open-feature/flagd/issues/900)) ([c2ddcbf](https://github.com/open-feature/flagd/commit/c2ddcbfe49b8507fe463c11eb2b031bbc331792a))


### 🧹 Chore

* add new flagd-evaluator e2e suite ([#898](https://github.com/open-feature/flagd/issues/898)) ([37ab55d](https://github.com/open-feature/flagd/commit/37ab55d26a9902935e4f1ddfd1a6af28d3b1cfa4))

## [0.6.5](https://github.com/open-feature/flagd/compare/flagd/v0.6.4...flagd/v0.6.5) (2023-09-08)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.4 ([#880](https://github.com/open-feature/flagd/issues/880)) ([ebb543d](https://github.com/open-feature/flagd/commit/ebb543d6eec18134e44ee7fe623fd2a336a1cf8d))
* **deps:** update opentelemetry-go monorepo ([#868](https://github.com/open-feature/flagd/issues/868)) ([d48317f](https://github.com/open-feature/flagd/commit/d48317f61d7db7ba0398dc9ab7cdd174a0b87555))


### 🧹 Chore

* disable caching on integration tests ([#899](https://github.com/open-feature/flagd/issues/899)) ([16dd21e](https://github.com/open-feature/flagd/commit/16dd21e5834519af3a22ffeb989ab398f8c1ddd9))
* upgrade to go 1.20 ([#891](https://github.com/open-feature/flagd/issues/891)) ([977167f](https://github.com/open-feature/flagd/commit/977167fb8db330b62726097616dcd691267199ad))

## [0.6.4](https://github.com/open-feature/flagd/compare/flagd/v0.6.3...flagd/v0.6.4) (2023-08-30)


### 🐛 Bug Fixes

* **deps:** update module github.com/cucumber/godog to v0.13.0 ([#855](https://github.com/open-feature/flagd/issues/855)) ([5b42486](https://github.com/open-feature/flagd/commit/5b4248654f7199afc50663e73609eeb20a3d11ec))
* **deps:** update module github.com/open-feature/flagd/core to v0.6.3 ([#794](https://github.com/open-feature/flagd/issues/794)) ([9671964](https://github.com/open-feature/flagd/commit/96719649affeb1f8412e8b25f52d7292281d8230))


### 🧹 Chore

* **deps:** update golang docker tag to v1.21 ([#822](https://github.com/open-feature/flagd/issues/822)) ([effe29d](https://github.com/open-feature/flagd/commit/effe29d50e33e6c06ef40d7f83f1b3f0df6bd1a2))

## [0.6.3](https://github.com/open-feature/flagd/compare/flagd/v0.6.2...flagd/v0.6.3) (2023-08-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.2 ([#779](https://github.com/open-feature/flagd/issues/779)) ([f34de59](https://github.com/open-feature/flagd/commit/f34de59fc8e636be043ce89758950d6ea3fe7376))
* **deps:** update module go.uber.org/zap to v1.25.0 ([#786](https://github.com/open-feature/flagd/issues/786)) ([40d0aa6](https://github.com/open-feature/flagd/commit/40d0aa66cf422db6811206d777b55396a96f330f))

## [0.6.2](https://github.com/open-feature/flagd/compare/flagd/v0.6.1...flagd/v0.6.2) (2023-07-28)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.6.1 ([#745](https://github.com/open-feature/flagd/issues/745)) ([d290d8f](https://github.com/open-feature/flagd/commit/d290d8fda8aa84ed2db6454fdd26e60b028e3f7f))

## [0.6.1](https://github.com/open-feature/flagd/compare/flagd/v0.6.0...flagd/v0.6.1) (2023-07-27)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/go-sdk-contrib/tests/flagd to v1.2.3 ([#749](https://github.com/open-feature/flagd/issues/749)) ([cd63e48](https://github.com/open-feature/flagd/commit/cd63e489d681c0998a9c38072410653473ce40fc))

## [0.6.0](https://github.com/open-feature/flagd/compare/flagd/v0.5.4...flagd/v0.6.0) (2023-07-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.4 ([#693](https://github.com/open-feature/flagd/issues/693)) ([33705a6](https://github.com/open-feature/flagd/commit/33705a67300ec70760ba0baeb610f5a2e931205f))
* **deps:** update module github.com/open-feature/go-sdk-contrib/providers/flagd to v0.1.13 ([#697](https://github.com/open-feature/flagd/issues/697)) ([435448f](https://github.com/open-feature/flagd/commit/435448f449044eb5fff88c94e81883cc801c02c4))
* **deps:** update module github.com/spf13/viper to v1.16.0 ([#679](https://github.com/open-feature/flagd/issues/679)) ([798a975](https://github.com/open-feature/flagd/commit/798a975bb1a47420e814b6dd439f1cece1a263e5))


### 🔄 Refactoring

* **flagd:** update build.Dockerfile with buildkit caching ([#724](https://github.com/open-feature/flagd/issues/724)) ([3e9cc1a](https://github.com/open-feature/flagd/commit/3e9cc1a7d697b64690a8772fe0ec8e84e34ebf6c))
* **flagd:** update profile.Dockerfile with buildkit caching ([#723](https://github.com/open-feature/flagd/issues/723)) ([3f263c6](https://github.com/open-feature/flagd/commit/3f263c65a6fe8f9e1f42d105dfbc89b9497cd080))
* remove protobuf dependency from eval package ([#701](https://github.com/open-feature/flagd/issues/701)) ([34ffafd](https://github.com/open-feature/flagd/commit/34ffafd9a777da3f11bd3bfa81565e774cc63214))

## [0.5.4](https://github.com/open-feature/flagd/compare/flagd/v0.5.3...flagd/v0.5.4) (2023-06-07)


### 🧹 Chore

* update otel dependencies ([#649](https://github.com/open-feature/flagd/issues/649)) ([2114e41](https://github.com/open-feature/flagd/commit/2114e41c38951247866c0b408e5f933282902e70))


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.3 ([#634](https://github.com/open-feature/flagd/issues/634)) ([1bc7e99](https://github.com/open-feature/flagd/commit/1bc7e99473bc0c7bcacfb40030562e556d3895d6))
* **deps:** update module github.com/open-feature/go-sdk-contrib/providers/flagd to v0.1.12 ([#635](https://github.com/open-feature/flagd/issues/635)) ([fe88061](https://github.com/open-feature/flagd/commit/fe88061ed6e0f1b6119af4c96a02495c4ff8072b))
* **deps:** update module github.com/open-feature/go-sdk-contrib/tests/flagd to v1.2.2 ([#651](https://github.com/open-feature/flagd/issues/651)) ([9776973](https://github.com/open-feature/flagd/commit/9776973109a1bb45ab611ede6b2c4d2c01508455))


### ✨ New Features

* telemetry improvements ([#653](https://github.com/open-feature/flagd/issues/653)) ([ea02cba](https://github.com/open-feature/flagd/commit/ea02cba24bde982d55956fe54de1e8f27226bfc6))


### 🔄 Refactoring

* introduce additional linting rules + fix discrepancies ([#616](https://github.com/open-feature/flagd/issues/616)) ([aef0b90](https://github.com/open-feature/flagd/commit/aef0b9042dcbe5b3f9a7e97960b27366fe50adfe))

## [0.5.3](https://github.com/open-feature/flagd/compare/flagd/v0.5.2...flagd/v0.5.3) (2023-05-04)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.2 ([#613](https://github.com/open-feature/flagd/issues/613)) ([218f435](https://github.com/open-feature/flagd/commit/218f435f0212fa24483b2af25e184e154e575eb1))
* **deps:** update module github.com/spf13/cobra to v1.7.0 ([#587](https://github.com/open-feature/flagd/issues/587)) ([12b3477](https://github.com/open-feature/flagd/commit/12b34773a68f6ae7e7e605aebc9f7075eb819994))


### ✨ New Features

* Introduce connect traces ([#624](https://github.com/open-feature/flagd/issues/624)) ([28bac6a](https://github.com/open-feature/flagd/commit/28bac6a54aed79cb8d84a147ffea296c36f5bd51))

## [0.5.2](https://github.com/open-feature/flagd/compare/flagd/v0.5.1...flagd/v0.5.2) (2023-04-13)


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/flagd/core to v0.5.1 ([#579](https://github.com/open-feature/flagd/issues/579)) ([58eed62](https://github.com/open-feature/flagd/commit/58eed62f5021e5c7a01a171067b725bf3ff83965))


### ✨ New Features

* otel traces for flag evaluation ([#598](https://github.com/open-feature/flagd/issues/598)) ([1757035](https://github.com/open-feature/flagd/commit/175703548f88469f25d749e320ee48030c9f9074))

## [0.5.1](https://github.com/open-feature/flagd/compare/flagd/v0.5.0...flagd/v0.5.1) (2023-04-12)


### ✨ New Features

* flagd OTEL collector ([#586](https://github.com/open-feature/flagd/issues/586)) ([494bec3](https://github.com/open-feature/flagd/commit/494bec33dcc1ddf0fa5cd0866f06265618408f5e))


### 🐛 Bug Fixes

* fall back to default port if env var cannot be parsed ([#591](https://github.com/open-feature/flagd/issues/591)) ([1fda104](https://github.com/open-feature/flagd/commit/1fda10473dba36149e13fa0cb8bb686d6861e568))

## [0.5.0](https://github.com/open-feature/flagd/compare/flagd/v0.4.5...flagd/v0.5.0) (2023-03-30)


### ⚠ BREAKING CHANGES

* unify sources configuration handling ([#560](https://github.com/open-feature/flagd/issues/560))


### 🐛 Bug Fixes

* benchmark pipeline ([#538](https://github.com/open-feature/flagd/issues/538)) ([62cc0fc](https://github.com/open-feature/flagd/commit/62cc0fcfd6a63a6059352704117dbb78160eb689))
* **deps:** update module github.com/open-feature/flagd/core to v0.4.5 ([#552](https://github.com/open-feature/flagd/issues/552)) ([41799f6](https://github.com/open-feature/flagd/commit/41799f624c261a84599cdd406cf28f4b33e49851))


### 🧹 Chore

* refactor configuration handling for startup ([#551](https://github.com/open-feature/flagd/issues/551)) ([8dfbde5](https://github.com/open-feature/flagd/commit/8dfbde5bbffd16fb66797a750d15f0226edf54a7))

## [0.4.5](https://github.com/open-feature/flagd/compare/flagd/v0.4.4...flagd/v0.4.5) (2023-03-20)


### 📚 Documentation

* improve markdown quality ([#498](https://github.com/open-feature/flagd/issues/498)) ([c77fa37](https://github.com/open-feature/flagd/commit/c77fa37979899f95ba51f69eeee21d96b6ab239c))


### ✨ New Features

* grpc connection options to flagd configuration options ([#532](https://github.com/open-feature/flagd/issues/532)) ([aa74951](https://github.com/open-feature/flagd/commit/aa74951f43b662ff2df53e68d347fc10e6d23bb8))
* Introduce flagd kube proxy ([#495](https://github.com/open-feature/flagd/issues/495)) ([440864c](https://github.com/open-feature/flagd/commit/440864ce87174618321c9d5146221490d8f07b24))

## [0.4.4](https://github.com/open-feature/flagd/compare/flagd-v0.4.3...flagd/v0.4.4) (2023-03-10)


### ✨ New Features

* Restructure for monorepo setup ([#486](https://github.com/open-feature/flagd/issues/486)) ([ed2993c](https://github.com/open-feature/flagd/commit/ed2993cd67b8a46db3beb6bb8a360e1aa20349da))

## [0.4.2](https://github.com/open-feature/flagd/compare/v0.4.1...v0.4.2) (2023-03-09)


### 🧹 Chore

* Add targeted Flag to example config ([#467](https://github.com/open-feature/flagd/issues/467)) ([6a039ce](https://github.com/open-feature/flagd/commit/6a039cef875caae61ea6c65799f3b6dc3863d131))
* **deps:** pin dependencies ([#473](https://github.com/open-feature/flagd/issues/473)) ([679e860](https://github.com/open-feature/flagd/commit/679e8600f57ab1e03c493c4a4046bd9d7368efac))
* **deps:** update google-github-actions/release-please-action digest to e0b9d18 ([#474](https://github.com/open-feature/flagd/issues/474)) ([5b85b2a](https://github.com/open-feature/flagd/commit/5b85b2a611d9199e39735f101ed7e560257ce2e4))
* refactoring and improve coverage for K8s Sync ([#466](https://github.com/open-feature/flagd/issues/466)) ([6dc441e](https://github.com/open-feature/flagd/commit/6dc441e2f2418c1fd3a5a58dbb99f848ccbd8735))


### 🐛 Bug Fixes

* add registry login ([#476](https://github.com/open-feature/flagd/issues/476)) ([99de755](https://github.com/open-feature/flagd/commit/99de755749df43d2b1028d47487b78b0ab626a9e))
* **deps:** update module golang.org/x/crypto to v0.7.0 ([#472](https://github.com/open-feature/flagd/issues/472)) ([f53f6c8](https://github.com/open-feature/flagd/commit/f53f6c885ee90813161b99be5a273b485e064de8))
* **deps:** update module google.golang.org/protobuf to v1.29.0 ([#478](https://github.com/open-feature/flagd/issues/478)) ([f9adc8e](https://github.com/open-feature/flagd/commit/f9adc8e3746256bcec045c06c78034c45722d60c))


### ✨ New Features

* grpc tls connectivity (grpcs) ([#477](https://github.com/open-feature/flagd/issues/477)) ([228f430](https://github.com/open-feature/flagd/commit/228f430e4945173755f52b8e712b23c28314517e))
* introduce per-sync configurations ([#448](https://github.com/open-feature/flagd/issues/448)) ([1d80039](https://github.com/open-feature/flagd/commit/1d80039558b29fff117478e308fd794a1244f0e5))

## [0.4.1](https://github.com/open-feature/flagd/compare/v0.4.0...v0.4.1) (2023-03-07)


### 🔄 Refactoring

* remove unused struct field ([#458](https://github.com/open-feature/flagd/issues/458)) ([a04c0b8](https://github.com/open-feature/flagd/commit/a04c0b837dbe9e28d1e01e43ea9e378a6c0f316a))


### 🧹 Chore

* **deps:** update sigstore/cosign-installer digest to bd2d118 ([#471](https://github.com/open-feature/flagd/issues/471)) ([ee90f48](https://github.com/open-feature/flagd/commit/ee90f48317ec600f09534306503dc752254a1d09))


### 🐛 Bug Fixes

* **deps:** update module github.com/open-feature/go-sdk-contrib/providers/flagd to v0.1.10 ([#459](https://github.com/open-feature/flagd/issues/459)) ([cbdf9b0](https://github.com/open-feature/flagd/commit/cbdf9b07c30239d7d04ef770cf4461fb33422fe9))
* **deps:** update module golang.org/x/net to v0.8.0 ([#468](https://github.com/open-feature/flagd/issues/468)) ([10d5f2c](https://github.com/open-feature/flagd/commit/10d5f2c55081a25daa1f0e0fa81f96f0fffbbc7b))
* fix broken image signing ([#461](https://github.com/open-feature/flagd/issues/461)) ([05bb51c](https://github.com/open-feature/flagd/commit/05bb51c7ab30f6e976b87f54ca889e978f834211))
* fixing image delimeter  ([#463](https://github.com/open-feature/flagd/issues/463)) ([b4ee495](https://github.com/open-feature/flagd/commit/b4ee495dc8e00b032518ea42d272a36b3b662e95))
* security issues ([#464](https://github.com/open-feature/flagd/issues/464)) ([7f1e759](https://github.com/open-feature/flagd/commit/7f1e759a87a9af63e9384005c959a3f500cc474c))
* set readiness once only ([#465](https://github.com/open-feature/flagd/issues/465)) ([41a888d](https://github.com/open-feature/flagd/commit/41a888d6b60c030b913280c2a1eeff8b25e8aada))

## [0.4.0](https://github.com/open-feature/flagd/compare/v0.3.7...v0.4.0) (2023-03-02)


### ⚠ BREAKING CHANGES

* Use OTel to export metrics (metric name changes) ([#419](https://github.com/open-feature/flagd/issues/419))

### 🧹 Chore

* add additional sections to the release notes ([#449](https://github.com/open-feature/flagd/issues/449)) ([798f71a](https://github.com/open-feature/flagd/commit/798f71a92d2e2f450a53cda93b44217cbb2ad7fd))
* attach image sbom to release artefacts ([#407](https://github.com/open-feature/flagd/issues/407)) ([fb4ee50](https://github.com/open-feature/flagd/commit/fb4ee502217e2262849df09258f3a0ffa7edec13))
* **deps:** update actions/configure-pages digest to fc89b04 ([#417](https://github.com/open-feature/flagd/issues/417)) ([04014e7](https://github.com/open-feature/flagd/commit/04014e7cb37e43f5ed3726dfd31da96202abc043))
* **deps:** update amannn/action-semantic-pull-request digest to b6bca70 ([#441](https://github.com/open-feature/flagd/issues/441)) ([ce0ebe1](https://github.com/open-feature/flagd/commit/ce0ebe13dd992688a3a0464ff401f2c40651da52))
* **deps:** update docker/login-action digest to ec9cdf0 ([#437](https://github.com/open-feature/flagd/issues/437)) ([2650670](https://github.com/open-feature/flagd/commit/2650670d35166e119f9a92613d3aca81523b9faa))
* **deps:** update docker/metadata-action digest to 3343011 ([#438](https://github.com/open-feature/flagd/issues/438)) ([e7ebf32](https://github.com/open-feature/flagd/commit/e7ebf32caf0eae7449e673da0c10998f97ebf781))
* **deps:** update github/codeql-action digest to 32dc499 ([#439](https://github.com/open-feature/flagd/issues/439)) ([f91d91b](https://github.com/open-feature/flagd/commit/f91d91bf020d330f96572c5ee11a210c0c7f4311))
* **deps:** update google-github-actions/release-please-action digest to d3c71f9 ([#406](https://github.com/open-feature/flagd/issues/406)) ([6e1ffb2](https://github.com/open-feature/flagd/commit/6e1ffb27fea5e91014a6991b2afca9a59f89117f))
* disable caching tests in CI ([#442](https://github.com/open-feature/flagd/issues/442)) ([28a35f6](https://github.com/open-feature/flagd/commit/28a35f62d618539362ae83a48f11af08ca2ae245))
* fix race condition on init read ([#409](https://github.com/open-feature/flagd/issues/409)) ([0c9eb23](https://github.com/open-feature/flagd/commit/0c9eb2322df99b4216d40afd1cb3b8873b0c59ff))
* integration test stability ([#432](https://github.com/open-feature/flagd/issues/432)) ([5a6a5d5](https://github.com/open-feature/flagd/commit/5a6a5d5887badd846cffe882c8c22a35b850fa06))
* integration tests ([#312](https://github.com/open-feature/flagd/issues/312)) ([6192ac8](https://github.com/open-feature/flagd/commit/6192ac8820b0f472672ba177b7c5838244b6e277))
* reorder release note sections ([df7bfce](https://github.com/open-feature/flagd/commit/df7bfce85ec7d6abaa987f87341c5af380904b51))
* use -short flag in benchmark tests ([#431](https://github.com/open-feature/flagd/issues/431)) ([e68a6aa](https://github.com/open-feature/flagd/commit/e68a6aadb3dac46676299ab94a34a0bcc39a67af))


### 🐛 Bug Fixes

* **deps:** update kubernetes packages to v0.26.2 ([#450](https://github.com/open-feature/flagd/issues/450)) ([2885227](https://github.com/open-feature/flagd/commit/28852270f34ff81c072337b29aa17f4b6634e9cc))
* **deps:** update module github.com/bufbuild/connect-go to v1.5.2 ([#416](https://github.com/open-feature/flagd/issues/416)) ([feb7f04](https://github.com/open-feature/flagd/commit/feb7f047365263758a63d8dffea936f621a4966d))
* **deps:** update module github.com/open-feature/go-sdk-contrib/providers/flagd to v0.1.9 ([#427](https://github.com/open-feature/flagd/issues/427)) ([42d2705](https://github.com/open-feature/flagd/commit/42d270558bf9badcff9a9b352fda35491c45aebe))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.29 ([#429](https://github.com/open-feature/flagd/issues/429)) ([b7fae81](https://github.com/open-feature/flagd/commit/b7fae81b89b3a1a0793a688c32569c4284633c6a))
* **deps:** update module github.com/stretchr/testify to v1.8.2 ([#440](https://github.com/open-feature/flagd/issues/440)) ([ab3e674](https://github.com/open-feature/flagd/commit/ab3e6748abc7843c022afeaf7cb11193cdcf59c5))
* **deps:** update module golang.org/x/net to v0.7.0 ([#410](https://github.com/open-feature/flagd/issues/410)) ([c6133b6](https://github.com/open-feature/flagd/commit/c6133b6af61f3d73ae73d318a1a9f44db2540540))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.14.5 ([#454](https://github.com/open-feature/flagd/issues/454)) ([f907f11](https://github.com/open-feature/flagd/commit/f907f114f23fa2efa2637e254e829e4d53a90b51))
* remove non-error error log from parseFractionalEvaluationData ([#446](https://github.com/open-feature/flagd/issues/446)) ([34aca79](https://github.com/open-feature/flagd/commit/34aca79e6ec9876a6cced0fe49e1ceea34d83696))


### ✨ New Features

* add debug logging for merge behaviour  ([#456](https://github.com/open-feature/flagd/issues/456)) ([dc71e84](https://github.com/open-feature/flagd/commit/dc71e84f0704690b528e7f1c2b56cb4898374fbf))
* add Health and Readiness probes ([#418](https://github.com/open-feature/flagd/issues/418)) ([7f2358c](https://github.com/open-feature/flagd/commit/7f2358ce207527c890f4a2f46ce4b9e8bf697095))
* Add version to startup message ([#430](https://github.com/open-feature/flagd/issues/430)) ([8daf613](https://github.com/open-feature/flagd/commit/8daf613e7e4f4492df0c06e2ef464f4337cadaca))
* introduce flag merge behaviour ([#414](https://github.com/open-feature/flagd/issues/414)) ([524f65e](https://github.com/open-feature/flagd/commit/524f65ea7215466bb4ac24a8d0d5953dd1cfe9a0))
* introduce grpc sync for flagd ([#297](https://github.com/open-feature/flagd/issues/297)) ([33413f2](https://github.com/open-feature/flagd/commit/33413f25882a3f1cf4953da0f18e746bfb69faf4))
* refactor and improve K8s sync provider ([#443](https://github.com/open-feature/flagd/issues/443)) ([4c03bfc](https://github.com/open-feature/flagd/commit/4c03bfc812e7ceabcac0979290bd74d9efc9da15))
* Use OTel to export metrics (metric name changes) ([#419](https://github.com/open-feature/flagd/issues/419)) ([eb3982a](https://github.com/open-feature/flagd/commit/eb3982a1cb72d664022b5cb126b533cf61497001))


### 📚 Documentation

* add .net flagd provider ([73d7840](https://github.com/open-feature/flagd/commit/73d7840c9fdef9c62371c677e02c0d9773c85f95))
* configuration merge docs ([#455](https://github.com/open-feature/flagd/issues/455)) ([6cb66b1](https://github.com/open-feature/flagd/commit/6cb66b14d01b6ee1c270bbdd3e30d4016757eae5))
* documentation for creating a provider ([#413](https://github.com/open-feature/flagd/issues/413)) ([d0c099d](https://github.com/open-feature/flagd/commit/d0c099d9aba3ed4d760a1858381f5e29b6d49a9c))
* updated filepaths for schema store regex ([#344](https://github.com/open-feature/flagd/issues/344)) ([2d0e9d9](https://github.com/open-feature/flagd/commit/2d0e9d956fbc99f2775821cfecdceb2b016d2b78))

## [0.3.7](https://github.com/open-feature/flagd/compare/v0.3.6...v0.3.7) (2023-02-13)


### Bug Fixes

* **deps:** update module golang.org/x/net to v0.6.0 ([#396](https://github.com/open-feature/flagd/issues/396)) ([beb7564](https://github.com/open-feature/flagd/commit/beb756470b1e1d5ef0670b8322b6ed9cb44efa24))
* **deps:** update module google.golang.org/grpc to v1.53.0 ([#388](https://github.com/open-feature/flagd/issues/388)) ([174cd7c](https://github.com/open-feature/flagd/commit/174cd7c70fa5ae2573db2c5972b75786633e2f41))
* error handling of Serve/ServeTLS funcs ([#397](https://github.com/open-feature/flagd/issues/397)) ([8923bf2](https://github.com/open-feature/flagd/commit/8923bf2d407e18b65c188aef9bf7370fc74c3be2))
* fix race in http sync test ([#401](https://github.com/open-feature/flagd/issues/401)) ([1d0c8e1](https://github.com/open-feature/flagd/commit/1d0c8e168b73f7fbd4b27ece733041bbe08261c0))
* sbom artefact name ([#380](https://github.com/open-feature/flagd/issues/380)) ([3daef26](https://github.com/open-feature/flagd/commit/3daef263c43ed63776d604d27f7ae6b993fff143)), closes [#379](https://github.com/open-feature/flagd/issues/379)

## [0.3.6](https://github.com/open-feature/flagd/compare/v0.3.5...v0.3.6) (2023-02-06)


### Bug Fixes

* set ResolveObject reason ([#375](https://github.com/open-feature/flagd/issues/375)) ([dcf199d](https://github.com/open-feature/flagd/commit/dcf199dab9d11b86454028869a54d77a474fc4a6))

## [0.3.5](https://github.com/open-feature/flagd/compare/v0.3.4...v0.3.5) (2023-02-06)


### Features

* flagd image signing ([#338](https://github.com/open-feature/flagd/issues/338)) ([eca6a60](https://github.com/open-feature/flagd/commit/eca6a60967999a303ceef5465f1acc35c83afd6d))
* update in logging to console and Unify case usage, seperators and punctuation for logging ([#322](https://github.com/open-feature/flagd/issues/322)) ([0bdcfd2](https://github.com/open-feature/flagd/commit/0bdcfd2fecc03b15be9fc4b0489431b8fa86aed8))


### Bug Fixes

* **deps:** update module github.com/bufbuild/connect-go to v1.5.1 ([#365](https://github.com/open-feature/flagd/issues/365)) ([e25f452](https://github.com/open-feature/flagd/commit/e25f452906e034e339309270cc8db6dcd58e9973))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.28 ([#342](https://github.com/open-feature/flagd/issues/342)) ([e6df80f](https://github.com/open-feature/flagd/commit/e6df80fd25d3da342e72d2ca0e923d9bf3d3f797))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.14.2 ([#336](https://github.com/open-feature/flagd/issues/336)) ([836d3cf](https://github.com/open-feature/flagd/commit/836d3cf3c06570d59929c3464e3c8e11c9b5a2fa))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.14.3 ([#372](https://github.com/open-feature/flagd/issues/372)) ([330ac91](https://github.com/open-feature/flagd/commit/330ac91e375124826b2a7a1a22d0daa18368ab99))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.14.4 ([#374](https://github.com/open-feature/flagd/issues/374)) ([d90e561](https://github.com/open-feature/flagd/commit/d90e561bfc5b798d13d4ba8f30f523b1053f3748))
* fix unbuffered channel blocking goroutine  ([#358](https://github.com/open-feature/flagd/issues/358)) ([4f1905a](https://github.com/open-feature/flagd/commit/4f1905a9ac6d62b5edb297fba904aac8680c89cf))
* introduced RWMutex to flag state to prevent concurrent r/w of map ([#370](https://github.com/open-feature/flagd/issues/370)) ([93e356b](https://github.com/open-feature/flagd/commit/93e356b4ab0b65c71659bd52d73f618edffc96f5))
* use event.Has func for file change notification handling (increased stability across OS) ([#361](https://github.com/open-feature/flagd/issues/361)) ([09f74b9](https://github.com/open-feature/flagd/commit/09f74b9c5d15622c98da08558cbcd63fe9422754))

## [0.3.4](https://github.com/open-feature/flagd/compare/v0.3.3...v0.3.4) (2023-01-28)


### Bug Fixes

* **deps:** update goreleaser/goreleaser-action action to v4 ([#340](https://github.com/open-feature/flagd/issues/340)) ([b9fcd5c](https://github.com/open-feature/flagd/commit/b9fcd5caa67a61b447c437b651471b4603b2b272))

## [0.3.3](https://github.com/open-feature/flagd/compare/v0.3.2...v0.3.3) (2023-01-28)


### Bug Fixes

* **deps:** update module github.com/bufbuild/connect-go to v1.5.0 ([#326](https://github.com/open-feature/flagd/issues/326)) ([7f332e5](https://github.com/open-feature/flagd/commit/7f332e50ecb1cea19108d1fa2fd79da3d5864bf9))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.26 ([#331](https://github.com/open-feature/flagd/issues/331)) ([be67e5f](https://github.com/open-feature/flagd/commit/be67e5f5bc1fb7351a04ffc4180447a27d57d32a))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.27 ([#335](https://github.com/open-feature/flagd/issues/335)) ([824cf1a](https://github.com/open-feature/flagd/commit/824cf1ab0f2e18826207af16d5328b817c755c8e))
* send datasync on remove fs events ([#339](https://github.com/open-feature/flagd/issues/339)) ([4c9aaac](https://github.com/open-feature/flagd/commit/4c9aaaca77b1c8b16f59434aeb37407fced47ecf))

## [0.3.2](https://github.com/open-feature/flagd/compare/v0.3.1...v0.3.2) (2023-01-26)


### Bug Fixes

* deprecation warning fix ([#317](https://github.com/open-feature/flagd/issues/317)) ([a2630db](https://github.com/open-feature/flagd/commit/a2630dbba151f35cc985d38de9cf25bfee2b76c8))
* **deps:** update kubernetes packages to v0.26.1 ([#267](https://github.com/open-feature/flagd/issues/267)) ([26825f2](https://github.com/open-feature/flagd/commit/26825f288c56df638fd160caa93f926a8c136108))
* **deps:** update module github.com/diegoholiveira/jsonlogic/v3 to v3.2.7 ([#283](https://github.com/open-feature/flagd/issues/283)) ([2ab5a00](https://github.com/open-feature/flagd/commit/2ab5a00fa6f19c7e0fe1a4e36649fae2633ac269))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.24 ([#290](https://github.com/open-feature/flagd/issues/290)) ([38d3eba](https://github.com/open-feature/flagd/commit/38d3ebaffcb1f36a38003273c62c6317f0ee75a3))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.25 ([#324](https://github.com/open-feature/flagd/issues/324)) ([ed1d3aa](https://github.com/open-feature/flagd/commit/ed1d3aaba4ca179a89757a6b1c3f328826e787fc))
* **deps:** update module github.com/spf13/viper to v1.15.0 ([#296](https://github.com/open-feature/flagd/issues/296)) ([d43220b](https://github.com/open-feature/flagd/commit/d43220b2be58e4bce05050c5d1b36788289ae7cc))
* **deps:** update module google.golang.org/grpc to v1.52.1 ([#314](https://github.com/open-feature/flagd/issues/314)) ([ad25388](https://github.com/open-feature/flagd/commit/ad25388461816100e19bda44a8e0077770ea0ee4))
* **deps:** update module google.golang.org/grpc to v1.52.3 ([#325](https://github.com/open-feature/flagd/issues/325)) ([8013ea5](https://github.com/open-feature/flagd/commit/8013ea5c6fa311b337c7ec1b1e8e756080808948))
* Update flagd systemd config to use URI ([#315](https://github.com/open-feature/flagd/issues/315)) ([93a04b4](https://github.com/open-feature/flagd/commit/93a04b46133e9220ec6f23d833c11f195e05c13e))
* update outdated doc link in deprecation warning ([#316](https://github.com/open-feature/flagd/issues/316)) ([19695d2](https://github.com/open-feature/flagd/commit/19695d2715129d6718ca0617b6ec6922ffb79c9b))

## [0.3.1](https://github.com/open-feature/flagd/compare/v0.3.0...v0.3.1) (2023-01-12)


### Features

* file extension detection ([#257](https://github.com/open-feature/flagd/issues/257)) ([ca22541](https://github.com/open-feature/flagd/commit/ca2254117adc163b94d662b3d1fbfd868f788fcb))
* ResolveAll endpoint for bulk evaluation ([#239](https://github.com/open-feature/flagd/issues/239)) ([6437c43](https://github.com/open-feature/flagd/commit/6437c43022b5c94d2fb835a406d85a4e836f2fcf))


### Bug Fixes

* **deps:** update module github.com/bufbuild/connect-go to v1.4.1 ([#268](https://github.com/open-feature/flagd/issues/268)) ([712d7dd](https://github.com/open-feature/flagd/commit/712d7dd4a34980bf9eddad99d926cbdd5d69d624))
* **deps:** update module github.com/mattn/go-colorable to v0.1.13 ([#260](https://github.com/open-feature/flagd/issues/260)) ([5b11504](https://github.com/open-feature/flagd/commit/5b11504cdce50c137540cc79d2db94e70a21338b))
* **deps:** update module github.com/open-feature/open-feature-operator to v0.2.23 ([#261](https://github.com/open-feature/flagd/issues/261)) ([a1dd3b9](https://github.com/open-feature/flagd/commit/a1dd3b9005374b5527f12b8e138250cacddc71af))
* **deps:** update module github.com/rs/cors to v1.8.3 ([#264](https://github.com/open-feature/flagd/issues/264)) ([0e6f2f3](https://github.com/open-feature/flagd/commit/0e6f2f3e5a77dae7d491eaf1094a65e692bebe5d))
* **deps:** update module github.com/stretchr/testify to v1.8.1 ([#265](https://github.com/open-feature/flagd/issues/265)) ([2ec61c6](https://github.com/open-feature/flagd/commit/2ec61c6bc61c266451b496ff18c3dd9a74173233))
* improve invalid sync URI errror msg ([#252](https://github.com/open-feature/flagd/issues/252)) ([5939870](https://github.com/open-feature/flagd/commit/5939870b8994dbca585c53dd022485090aab2406))
* replace character slice with regex replace ([#250](https://github.com/open-feature/flagd/issues/250)) ([c92d101](https://github.com/open-feature/flagd/commit/c92d1012b0de6af694c3af2fede28053e2572b04))

## [0.3.0](https://github.com/open-feature/flagd/compare/v0.2.7...v0.3.0) (2023-01-06)


### ⚠ BREAKING CHANGES

* consolidated configuration change events into one event ([#241](https://github.com/open-feature/flagd/issues/241))

### Features

* consolidated configuration change events into one event ([#241](https://github.com/open-feature/flagd/issues/241)) ([f9684b8](https://github.com/open-feature/flagd/commit/f9684b858dfef40576e0031654b421a37e8bb1d6))
* support yaml evaluator ([#206](https://github.com/open-feature/flagd/issues/206)) ([2dbace5](https://github.com/open-feature/flagd/commit/2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0))


### Bug Fixes

* changed eventing configuration mutex to rwmutex and added missing lock ([#220](https://github.com/open-feature/flagd/issues/220)) ([5bbef9e](https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f))
* omitempty targeting field in Flag structure ([#247](https://github.com/open-feature/flagd/issues/247)) ([3f406b5](https://github.com/open-feature/flagd/commit/3f406b53bda8b5beb8b0929da3802a0368c13151))

## [0.2.7](https://github.com/open-feature/flagd/compare/v0.2.5...v0.2.7) (2022-12-02)


### ⚠ BREAKING CHANGES

* start command flag refactor ([#222](https://github.com/open-feature/flagd/issues/222))

### Features

* enable request logging via the --debug flag ([#226](https://github.com/open-feature/flagd/issues/226)) ([11954b5](https://github.com/open-feature/flagd/commit/11954b521cc6197d0dc04b163e66e38d4c288047))
* Resurrected the STATIC flag reason. Documented the caching strategy. ([#224](https://github.com/open-feature/flagd/issues/224)) ([5830592](https://github.com/open-feature/flagd/commit/5830592053c55dc9e55c16854e40c3fc8345d6d1))
* snap ([#211](https://github.com/open-feature/flagd/issues/211)) ([c619844](https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738))
* start command flag refactor ([#222](https://github.com/open-feature/flagd/issues/222)) ([14474cc](https://github.com/open-feature/flagd/commit/14474ccf65b9b92213e8c792e94c458022484df4))


### Miscellaneous Chores

* release v0.2.6 ([93cfb78](https://github.com/open-feature/flagd/commit/93cfb78d024b436fa7fb17fd41f74d1508bf8b64))
* release v0.2.7 ([4a9f6df](https://github.com/open-feature/flagd/commit/4a9f6df4e472229ff805e9d5d3aa581c7c9c0667))

## [0.2.5](https://github.com/open-feature/flagd/compare/v0.2.4...v0.2.5) (2022-10-20)


### Bug Fixes

* CVE-2022-32149 ([#198](https://github.com/open-feature/flagd/issues/198)) ([11a7b34](https://github.com/open-feature/flagd/commit/11a7b3472ab2bc39bce7c40037e8f83736065163))

## [0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4) (2022-10-14)


### Bug Fixes

* ApiVersion check fix ([#193](https://github.com/open-feature/flagd/issues/193)) ([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))

## [0.2.3](https://github.com/open-feature/flagd/compare/v0.2.2...v0.2.3) (2022-10-13)


### Features

* Eventing ([#187](https://github.com/open-feature/flagd/issues/187)) ([3f7fcd2](https://github.com/open-feature/flagd/commit/3f7fcd2f57318fad4e0bf501f202af990d3c5a79))
* fixing informer issues ([#191](https://github.com/open-feature/flagd/issues/191)) ([837b0c6](https://github.com/open-feature/flagd/commit/837b0c673e7e7d4799f100291ca520d22944f22a))
* only fire modify event when FeatureFlagConfiguration Generation field has changed ([#167](https://github.com/open-feature/flagd/issues/167)) ([e2fc7ee](https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20))

## [0.2.2](https://github.com/open-feature/flagd/compare/v0.2.1...v0.2.2) (2022-10-03)


### Bug Fixes

* updated merge functionality ([#182](https://github.com/open-feature/flagd/issues/182)) ([94d7697](https://github.com/open-feature/flagd/commit/94d7697d08a07cede4a548ef998792d00f8954a0))

## [0.2.1](https://github.com/open-feature/flagd/compare/v0.2.0...v0.2.1) (2022-09-27)


### Bug Fixes

* updated tcp listener ([#174](https://github.com/open-feature/flagd/issues/174)) ([b750ed1](https://github.com/open-feature/flagd/commit/b750ed1268b5e6efe779a34e764bad3e781f8e93))

## [0.2.0](https://github.com/open-feature/flagd/compare/v0.1.1...v0.2.0) (2022-09-26)


### ⚠ BREAKING CHANGES

* Updated service to use connect (#163)

### Features

* Updated service to use connect ([#163](https://github.com/open-feature/flagd/issues/163)) ([828d5c4](https://github.com/open-feature/flagd/commit/828d5c4c11157f5b7a77f5041806ba2523186764))


### Bug Fixes

* checkout release tag before running container and binary releases ([#171](https://github.com/open-feature/flagd/issues/171)) ([50fe46f](https://github.com/open-feature/flagd/commit/50fe46fbbf120a0657c1df35b370cdc9051d0f31))

## [0.1.1](https://github.com/open-feature/flagd/compare/v0.1.0...v0.1.1) (2022-09-23)


### Bug Fixes

* bubbles up unclean error exits ([#170](https://github.com/open-feature/flagd/issues/170)) ([9f7db02](https://github.com/open-feature/flagd/commit/9f7db0259d2d24cb880eeddaebd3b8f48758248a))
* upgrade package containing vulnerability ([#162](https://github.com/open-feature/flagd/issues/162)) ([82278c7](https://github.com/open-feature/flagd/commit/82278c7cf08cc6b50f49ab500caf6f9003fc0823))
