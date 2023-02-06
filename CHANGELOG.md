# Changelog

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
