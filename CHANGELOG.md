# Changelog

All notable changes to this project will be documented in this file. See
[Conventional Commits](https://www.conventionalcommits.org/) for commit message
guidelines.

This file is maintained by [release-please](https://github.com/googleapis/release-please);
do not edit it by hand.

## [0.5.2](https://github.com/amiwrpremium/go-derive/compare/v0.5.1...v0.5.2) (2026-05-10)


### Documentation

* **examples:** cover the 30 methods + best_quotes channel from v0.5.0 ([#74](https://github.com/amiwrpremium/go-derive/issues/74)) ([70820c2](https://github.com/amiwrpremium/go-derive/commit/70820c2050dbb81e950ad8ec53321fd7b1449468))

## [0.5.1](https://github.com/amiwrpremium/go-derive/compare/v0.5.0...v0.5.1) (2026-05-10)


### Documentation

* **types:** point schema annotations at docs.derive.xyz ([#73](https://github.com/amiwrpremium/go-derive/issues/73)) ([77b7139](https://github.com/amiwrpremium/go-derive/commit/77b7139e772a1f189b94372e633443a769653e30))


### Tests

* **methods:** decode coverage for the 30 new methods from v0.5.0 ([#71](https://github.com/amiwrpremium/go-derive/issues/71)) ([7de30cd](https://github.com/amiwrpremium/go-derive/commit/7de30cdbfb0d5cef1261a80abe1c4173212232c1))

## [0.5.0](https://github.com/amiwrpremium/go-derive/compare/v0.4.1...v0.5.0) (2026-05-10)


### ⚠ BREAKING CHANGES

* bump for SubAccount breaking change in #68 ([#69](https://github.com/amiwrpremium/go-derive/issues/69))

### Features

* bump for SubAccount breaking change in [#68](https://github.com/amiwrpremium/go-derive/issues/68) ([#69](https://github.com/amiwrpremium/go-derive/issues/69)) ([2bf845f](https://github.com/amiwrpremium/go-derive/commit/2bf845f841aea91fa9f898fd247ae863cd3c977a))

## [0.4.1](https://github.com/amiwrpremium/go-derive/compare/v0.4.0...v0.4.1) (2026-05-10)


### Features

* **channels:** cover the 4 cockpit channels the SDK was missing ([#67](https://github.com/amiwrpremium/go-derive/issues/67)) ([5fbf8a4](https://github.com/amiwrpremium/go-derive/commit/5fbf8a4799643ece0331edaac38a49b41a85c23c))


### Documentation

* **examples:** cover the 21 typed endpoints from v0.4.0 ([#65](https://github.com/amiwrpremium/go-derive/issues/65)) ([8ccde10](https://github.com/amiwrpremium/go-derive/commit/8ccde1030e2b7106911b2571f5ec66b6a2891085))

## [0.4.0](https://github.com/amiwrpremium/go-derive/compare/v0.3.0...v0.4.0) (2026-05-09)


### ⚠ BREAKING CHANGES

* **methods:** rename GetOrderHistory→GetOrders; restore GetOrderHistory as time-window query ([#64](https://github.com/amiwrpremium/go-derive/issues/64))
* **types:** tighten four bare-string fields to documented enums ([#63](https://github.com/amiwrpremium/go-derive/issues/63))

### Features

* **methods:** add 21 typed read + trigger-order endpoints ([#61](https://github.com/amiwrpremium/go-derive/issues/61)) ([7c95b93](https://github.com/amiwrpremium/go-derive/commit/7c95b93c37c67185c745c15c261ae62ebba6f905))
* **methods:** rename GetOrderHistory→GetOrders; restore GetOrderHistory as time-window query ([#64](https://github.com/amiwrpremium/go-derive/issues/64)) ([ac7b2c7](https://github.com/amiwrpremium/go-derive/commit/ac7b2c70390e5a59307f54aa98bf02e897fc2be3))
* **types:** tighten four bare-string fields to documented enums ([#63](https://github.com/amiwrpremium/go-derive/issues/63)) ([9e029ca](https://github.com/amiwrpremium/go-derive/commit/9e029ca289ab9d9b64684959b88a8c9ccc6caef1))

## [0.3.0](https://github.com/amiwrpremium/go-derive/compare/v0.2.17...v0.3.0) (2026-05-09)


### ⚠ BREAKING CHANGES

* **methods:** type every json.RawMessage return value, reorganise extras.go, fix doc-audit drift ([#59](https://github.com/amiwrpremium/go-derive/issues/59))

### Features

* **methods:** type every json.RawMessage return value, reorganise extras.go, fix doc-audit drift ([#59](https://github.com/amiwrpremium/go-derive/issues/59)) ([a97e4af](https://github.com/amiwrpremium/go-derive/commit/a97e4af083d985124334f4bfc68bcafba4d8f76a))

## [0.2.17](https://github.com/amiwrpremium/go-derive/compare/v0.2.16...v0.2.17) (2026-05-08)


### Continuous Integration

* pin staticcheck via semver tag for Scorecard compliance ([#56](https://github.com/amiwrpremium/go-derive/issues/56)) ([60b06c1](https://github.com/amiwrpremium/go-derive/commit/60b06c1708fe9bbf5198b2c11fd4edc8f88848a7))

## [0.2.16](https://github.com/amiwrpremium/go-derive/compare/v0.2.15...v0.2.16) (2026-05-08)


### Dependencies

* **deps:** update dependency go to v1.26.3 ([#54](https://github.com/amiwrpremium/go-derive/issues/54)) ([159a4ef](https://github.com/amiwrpremium/go-derive/commit/159a4ef093a2ece72efed688ccc16dc8c683d5c1))

## [0.2.15](https://github.com/amiwrpremium/go-derive/compare/v0.2.14...v0.2.15) (2026-05-08)


### Documentation

* add hero banner (dark + light variants) ([#51](https://github.com/amiwrpremium/go-derive/issues/51)) ([c246b48](https://github.com/amiwrpremium/go-derive/commit/c246b4802ed2b7ec758329e98168fcf4091b4241))
* regroup README badges and add margin below banner ([#53](https://github.com/amiwrpremium/go-derive/issues/53)) ([a4a4068](https://github.com/amiwrpremium/go-derive/commit/a4a40684fdb5a3db1ed2c8024eeed38092b25c7c))

## [0.2.14](https://github.com/amiwrpremium/go-derive/compare/v0.2.13...v0.2.14) (2026-05-07)


### Bug Fixes

* **ci:** pin Go to 1.25.10 + 1.26.2 (close GO-2026-4971 + GO-2026-4918) ([#49](https://github.com/amiwrpremium/go-derive/issues/49)) ([cecc84d](https://github.com/amiwrpremium/go-derive/commit/cecc84d8e4c42714254938ac9ca198dfde07d1d3))


### Documentation

* refresh stale data, add scanner badges, plug doc gaps ([#48](https://github.com/amiwrpremium/go-derive/issues/48)) ([9d13b07](https://github.com/amiwrpremium/go-derive/commit/9d13b07e198f29bf08b98ef97577850384bd305e))

## [0.2.13](https://github.com/amiwrpremium/go-derive/compare/v0.2.12...v0.2.13) (2026-05-07)


### Bug Fixes

* **ci:** restore SLSA generator tag pin + exclude it from Renovate ([#47](https://github.com/amiwrpremium/go-derive/issues/47)) ([81e4cc2](https://github.com/amiwrpremium/go-derive/commit/81e4cc255172599385577ce6a5ea42ede36cf41e))


### Dependencies

* **deps:** update actions/checkout action to v6.0.2 ([#39](https://github.com/amiwrpremium/go-derive/issues/39)) ([02bc595](https://github.com/amiwrpremium/go-derive/commit/02bc595fab49a8a355c3a91ac0cc878b6483cdf1))
* **deps:** update github actions (major) ([#43](https://github.com/amiwrpremium/go-derive/issues/43)) ([8e9c424](https://github.com/amiwrpremium/go-derive/commit/8e9c424777da5bcfaad0c3e42f661e8c44af44d4))

## [0.2.12](https://github.com/amiwrpremium/go-derive/compare/v0.2.11...v0.2.12) (2026-05-07)


### Bug Fixes

* **ci:** allow deps scope + stop Renovate from bumping go directive ([#44](https://github.com/amiwrpremium/go-derive/issues/44)) ([69042ab](https://github.com/amiwrpremium/go-derive/commit/69042ab25df1c490b3eef49fb29b49bbd5a5f4bd))


### Dependencies

* **deps:** update github actions (minor) ([#41](https://github.com/amiwrpremium/go-derive/issues/41)) ([0ad5194](https://github.com/amiwrpremium/go-derive/commit/0ad5194be2b34c3364953eac8546de0082d8a17e))

## [0.2.11](https://github.com/amiwrpremium/go-derive/compare/v0.2.10...v0.2.11) (2026-05-07)


### Bug Fixes

* **ci:** enable GFM in remark-lint config (close 20 false-positive alerts) ([#37](https://github.com/amiwrpremium/go-derive/issues/37)) ([ee1c3eb](https://github.com/amiwrpremium/go-derive/commit/ee1c3eb4517e12c27746cccf9d3ce7f97c2ec550))

## [0.2.10](https://github.com/amiwrpremium/go-derive/compare/v0.2.9...v0.2.10) (2026-05-07)


### Documentation

* replicate full package overview into every non-canonical file ([#35](https://github.com/amiwrpremium/go-derive/issues/35)) ([040694b](https://github.com/amiwrpremium/go-derive/commit/040694b504c5a4bfd26d8086544ae7576514d59f))

## [0.2.9](https://github.com/amiwrpremium/go-derive/compare/v0.2.8...v0.2.9) (2026-05-07)


### Documentation

* group README badges into 6 sections + add 5 new badges ([#33](https://github.com/amiwrpremium/go-derive/issues/33)) ([f46f79b](https://github.com/amiwrpremium/go-derive/commit/f46f79bdcd00de712f00e7daffa173afaa4cedae))

## [0.2.8](https://github.com/amiwrpremium/go-derive/compare/v0.2.7...v0.2.8) (2026-05-07)


### Documentation

* surface OpenSSF Best Practices badge (project 12775) ([#31](https://github.com/amiwrpremium/go-derive/issues/31)) ([64372ab](https://github.com/amiwrpremium/go-derive/commit/64372ab732884f922f97010caaa6d2acc0866dc8))

## [0.2.7](https://github.com/amiwrpremium/go-derive/compare/v0.2.6...v0.2.7) (2026-05-07)


### Documentation

* rewrite per-file Package headers to satisfy Revive form check ([#29](https://github.com/amiwrpremium/go-derive/issues/29)) ([e0d5efd](https://github.com/amiwrpremium/go-derive/commit/e0d5efd7a5fe4885e104933b859b5f9e5c481196))

## [0.2.6](https://github.com/amiwrpremium/go-derive/compare/v0.2.5...v0.2.6) (2026-05-07)


### Documentation

* add per-file Package headers (Codacy Revive strict mode) ([#27](https://github.com/amiwrpremium/go-derive/issues/27)) ([6929ac0](https://github.com/amiwrpremium/go-derive/commit/6929ac0f8cc29b6378f0051eb4f4e31fb7314ec5))

## [0.2.5](https://github.com/amiwrpremium/go-derive/compare/v0.2.4...v0.2.5) (2026-05-07)


### Bug Fixes

* **ci:** pin go install tool versions (Scorecard PinnedDependencies) ([#25](https://github.com/amiwrpremium/go-derive/issues/25)) ([57a03b0](https://github.com/amiwrpremium/go-derive/commit/57a03b0a622d7924f37ef7aef610d551c8c97266))
* **ci:** replace curl|bash codacy coverage with pinned Action ([#26](https://github.com/amiwrpremium/go-derive/issues/26)) ([8a79a78](https://github.com/amiwrpremium/go-derive/commit/8a79a786eea55d1eeb9692fac1d9ddddb275ea23))


### Documentation

* move package-comments to alphabetically-first file in each package ([#23](https://github.com/amiwrpremium/go-derive/issues/23)) ([03147fe](https://github.com/amiwrpremium/go-derive/commit/03147fed5023d09b072d2d494514906c33caf87d))

## [0.2.4](https://github.com/amiwrpremium/go-derive/compare/v0.2.3...v0.2.4) (2026-05-07)


### Bug Fixes

* **ci:** bypass implicit success() filter on slsa provenance job ([#21](https://github.com/amiwrpremium/go-derive/issues/21)) ([45eedf6](https://github.com/amiwrpremium/go-derive/commit/45eedf6e840ce3572ebbff4ee062d1e7cdd13165))

## [0.2.3](https://github.com/amiwrpremium/go-derive/compare/v0.2.2...v0.2.3) (2026-05-06)


### Bug Fixes

* **ci:** unblock SLSA provenance after release-event skip propagation ([#19](https://github.com/amiwrpremium/go-derive/issues/19)) ([dba2df9](https://github.com/amiwrpremium/go-derive/commit/dba2df95365dc0225d4b4b84cf9bd07011bf75f2))

## [0.2.2](https://github.com/amiwrpremium/go-derive/compare/v0.2.1...v0.2.2) (2026-05-05)


### Bug Fixes

* **ci:** replace path suppressions with proper rule config ([#15](https://github.com/amiwrpremium/go-derive/issues/15)) ([766ef7b](https://github.com/amiwrpremium/go-derive/commit/766ef7b34f78c6e1e0b0594e0c011796b67bb69d))
* **security:** validate workflow_dispatch inputs + improve security docs ([#18](https://github.com/amiwrpremium/go-derive/issues/18)) ([c1a7eb0](https://github.com/amiwrpremium/go-derive/commit/c1a7eb03030bd06b0bd5e7d17224cec2c1bc55c7))

## [0.2.1](https://github.com/amiwrpremium/go-derive/compare/v0.2.0...v0.2.1) (2026-05-05)


### Bug Fixes

* **ci:** scope GITHUB_TOKEN permissions per-job (Scorecard TokenPermissions) ([#14](https://github.com/amiwrpremium/go-derive/issues/14)) ([a5d6ead](https://github.com/amiwrpremium/go-derive/commit/a5d6eadd2cb21ef8327ed2bb6dd6a550e8c7a0af))
* **ci:** suppress 32 Codacy markdown/CI alerts via .codacy.yml ([#12](https://github.com/amiwrpremium/go-derive/issues/12)) ([7e3a934](https://github.com/amiwrpremium/go-derive/commit/7e3a9349026a9f02841020cc80629e8433c71621))

## [0.2.0](https://github.com/amiwrpremium/go-derive/compare/v0.1.3...v0.2.0) (2026-05-05)


### ⚠ BREAKING CHANGES

* **types:** remove deprecated Trade.Realized alias ([#10](https://github.com/amiwrpremium/go-derive/issues/10))

### Features

* **types:** remove deprecated Trade.Realized alias ([#10](https://github.com/amiwrpremium/go-derive/issues/10)) ([29ccf2e](https://github.com/amiwrpremium/go-derive/commit/29ccf2ea4f5edeec1cc6040764761340c31a5109))


### Code Refactoring

* **ws:** put context.Context first in closeWatcher and readPump ([#9](https://github.com/amiwrpremium/go-derive/issues/9)) ([f77350f](https://github.com/amiwrpremium/go-derive/commit/f77350fa8988d3a3cd27d134c1e6a19e25a8049b))

## [0.1.3](https://github.com/amiwrpremium/go-derive/compare/v0.1.2...v0.1.3) (2026-05-05)


### Bug Fixes

* **ci:** align ruleset check names with what GitHub actually publishes ([1a24b2b](https://github.com/amiwrpremium/go-derive/commit/1a24b2b825d20020195976a36488959b5d6a3598))
* **ci:** branch-master ruleset adjustments for free GitHub plans ([3817045](https://github.com/amiwrpremium/go-derive/commit/381704573e267672f08d5cc0de05dca310975bb9))
* **ci:** drop _comment inside pull_request.parameters (rejected by GitHub schema) ([51aa16d](https://github.com/amiwrpremium/go-derive/commit/51aa16dfd172da09da6a75b5c0f68c52e08afa8a))
* **ci:** drop required_signatures from both rulesets (incompatible with bot PRs) ([dd8eeac](https://github.com/amiwrpremium/go-derive/commit/dd8eeacb1b7102797521740ac68c4df702f916bb))
* **ci:** drop review requirement on master ruleset (solo-maintained project) ([4e03977](https://github.com/amiwrpremium/go-derive/commit/4e039770d1b2c971abab25bdead03b941ab00fac))
* **ci:** tag ruleset — drop Integration bypass (repo-level rulesets reject it) ([e356495](https://github.com/amiwrpremium/go-derive/commit/e35649580a59fbf521427a060ef874e79a0ed990))
* **deps:** drop deprecated top-level "go" key from renovate.json ([91d6999](https://github.com/amiwrpremium/go-derive/commit/91d69991a44d3893fcaf05f188696139e265a31f))

## [0.1.2](https://github.com/amiwrpremium/go-derive/compare/v0.1.1...v0.1.2) (2026-05-05)


### Features

* **ci:** use RELEASE_PLEASE_TOKEN PAT so releases auto-fire release.yml ([b8eca3c](https://github.com/amiwrpremium/go-derive/commit/b8eca3cc763faafe11e9d0c5c20befbd7d4497e6))

## [0.1.1](https://github.com/amiwrpremium/go-derive/compare/v0.1.0...v0.1.1) (2026-05-05)


### Bug Fixes

* **ci:** manually attach SLSA provenance — generator's auto-upload is broken ([6d074d4](https://github.com/amiwrpremium/go-derive/commit/6d074d409d1bc5e04378f068ad13c08b9802769c))
* **ci:** use tag ref for slsa-github-generator (its security model requires it) ([e5a2aba](https://github.com/amiwrpremium/go-derive/commit/e5a2ababb8f4089cced11667553871d93404e913))
* **ci:** wire SLSA subjects + ungate release uploads on workflow_dispatch ([a04741a](https://github.com/amiwrpremium/go-derive/commit/a04741aa757fb6725a4edda130ebf6bf30e6a6f4))
* **docs:** drop legacy CHANGELOG section that clashed with auto-generated style ([865a5dd](https://github.com/amiwrpremium/go-derive/commit/865a5dd15884899399800f9f0628c477ba184f92))

## [0.1.0](https://github.com/amiwrpremium/go-derive/compare/v0.1.0...v0.1.0) (2026-05-05)


### Features

* **auth:** EIP-712 signing, EIP-191 auth headers, and session-key support ([5dd0d31](https://github.com/amiwrpremium/go-derive/commit/5dd0d3196ad54ed8f248fbc7589243ee19fa5064))
* **channels:** typed subscription descriptors for public + private streams ([5d8bcc0](https://github.com/amiwrpremium/go-derive/commit/5d8bcc0328e94817652a07684eadb905bff25147))
* **codec:** ABI/decimal codec helpers for the signing path ([69ef723](https://github.com/amiwrpremium/go-derive/commit/69ef723eae09a107c3842864f61029fcfd3fe753))
* **contracts:** on-chain helper interfaces for deposit, withdraw, session keys ([dc4ccac](https://github.com/amiwrpremium/go-derive/commit/dc4ccac2cce8dbedff1e1b4f9dd1fcde50cd284d))
* **derive:** top-level facade bundling REST + WebSocket clients ([917574a](https://github.com/amiwrpremium/go-derive/commit/917574aa2a76499d1c6731e4a46e9af5d51035d7))
* **enums:** wire enum bindings for every Derive string type ([811120c](https://github.com/amiwrpremium/go-derive/commit/811120cca99a6deffee77bcee4262e52b13268f5))
* **errors:** typed errors with sentinels, APIError, and the 136-code catalogue ([8153d27](https://github.com/amiwrpremium/go-derive/commit/8153d27e7b8fb20cefafc94f70cc080113be3b49))
* **jsonrpc:** JSON-RPC 2.0 envelope, ID generator, and subscription routing ([4e7e242](https://github.com/amiwrpremium/go-derive/commit/4e7e242fd77ed494b7948f28a862245910158e2c))
* **methods:** typed RPC method bindings for every Derive endpoint ([e1ca44d](https://github.com/amiwrpremium/go-derive/commit/e1ca44d45c05575f7b5a1a39202c7a106965a1b4))
* **netconf:** per-network endpoints, chain ids, and EIP-712 domain bundles ([775ec4e](https://github.com/amiwrpremium/go-derive/commit/775ec4e17d08e7cb39714c5ae3aab9aaab1b4e29))
* **rest:** REST client over the HTTP transport with functional options ([2848ee9](https://github.com/amiwrpremium/go-derive/commit/2848ee90d78314d0cd0618dfea65fd0c84e1dd28))
* **retry:** exponential backoff with jitter ([c755d37](https://github.com/amiwrpremium/go-derive/commit/c755d3760215290a661813b8a4fcb3cc5ec86dfd))
* root package with Version constant and UserAgent helper ([7bb53fb](https://github.com/amiwrpremium/go-derive/commit/7bb53fbba5133ead33a9b01095d5bf20bbe4b2a0))
* **transport:** HTTP and WebSocket JSON-RPC transports with rate limiting ([9bba95d](https://github.com/amiwrpremium/go-derive/commit/9bba95d3915b18eaab21578220a8a8d86b0f1c1c))
* **types:** core domain types for prices, accounts, orders, and market data ([7000126](https://github.com/amiwrpremium/go-derive/commit/7000126b7a1e5b7fe99769e0bb5ead036b3d7f53))
* **ws:** WebSocket client with generic typed subscriptions ([0c56820](https://github.com/amiwrpremium/go-derive/commit/0c56820d61f42371ee887e285d9991ceea9dd67e))


### Bug Fixes

* **ci:** convert Go coverage to LCOV before Codacy upload ([3c053f8](https://github.com/amiwrpremium/go-derive/commit/3c053f8d34744c4bc60ad674a5ea1cc83ea4a2cc))
* **ci:** disambiguate Codacy SARIF runs for GitHub Code Scanning ([59ff1a2](https://github.com/amiwrpremium/go-derive/commit/59ff1a2fd409d252e0308882c8fd45b0927b44a2))
* **ci:** handle BASE==HEAD on initial pushes in trufflehog ([80f2efa](https://github.com/amiwrpremium/go-derive/commit/80f2efa66d9b70be0382a50c75270d5ffa4a970d))
* **ci:** ignore go-ethereum in license-check (dual-licensed module) ([dd50933](https://github.com/amiwrpremium/go-derive/commit/dd509337bb60a00d163b635aaa682e526efbf1d6))
* **ci:** repin Semgrep to current digest on the canonical namespace ([d041cd5](https://github.com/amiwrpremium/go-derive/commit/d041cd508bf80db36601d644969df027b890feec))
* **ci:** silence final SC2001 shellcheck warning in pin-check ([b5f386d](https://github.com/amiwrpremium/go-derive/commit/b5f386db111645913ad77ef328fc29bd469ea7fa))
* **ci:** silence shellcheck noise + extend typos dictionary ([5a32f99](https://github.com/amiwrpremium/go-derive/commit/5a32f99020fe48a10ed7e6215d087bd3a89ec273))
* **ci:** switch osv-scanner to the maintainer-recommended reusable workflow ([d6887c3](https://github.com/amiwrpremium/go-derive/commit/d6887c3fd5b5acfcbc378fa561c9901450fbaf7f))
* **ci:** tag Codacy SARIF runs with toolName+index for uniqueness ([4c2e8c0](https://github.com/amiwrpremium/go-derive/commit/4c2e8c07602a02df53a5cf96a64e7cbe0360d94e))
* **ci:** trust trufflehog action's auto-detect branch (revert prev fix) ([cc2b691](https://github.com/amiwrpremium/go-derive/commit/cc2b691c43a3c62beb1e7f110d7043e2f80a54c9))
* **docs:** rewrite line that markdown read as a bare list item ([8585672](https://github.com/amiwrpremium/go-derive/commit/8585672f96216602ab173cae63e378a7f593bbb9))
* **examples:** gofmt the perp-impact TWAP example ([007f770](https://github.com/amiwrpremium/go-derive/commit/007f770474510325d8bc0c22b88bdbc721515ae5))


### Documentation

* governance, contribution, and security policy ([1f55749](https://github.com/amiwrpremium/go-derive/commit/1f55749e8a117860d5258faec19ca0a988929143))
* OpenSSF Security Insights manifest ([c7cdf81](https://github.com/amiwrpremium/go-derive/commit/c7cdf81d5f3a0ba97c45df235478cbe0834aa5cd))
* **readme:** unhide Codacy grade + coverage badges ([9509916](https://github.com/amiwrpremium/go-derive/commit/950991668abb446093fe27e01b6d03632cff0641))
* write the SDK documentation set under docs/ ([daaa685](https://github.com/amiwrpremium/go-derive/commit/daaa685d21cac28a6656627bb39227b1a308aeb8))


### Tests

* **integration:** live testnet integration tests gated by build tag ([fb7416c](https://github.com/amiwrpremium/go-derive/commit/fb7416c57b655d26578ecfd256cfde695ffb5345))
* **testutil:** in-memory mocks for HTTP, WebSocket, and Transport ([fdb8a9b](https://github.com/amiwrpremium/go-derive/commit/fdb8a9be6bb925e3d1c8d3b52910adeac28c35cb))


### Build System

* annotate Version constant for release-please bumping ([aae1fc4](https://github.com/amiwrpremium/go-derive/commit/aae1fc424e2ad110d63e6844c606568d760c85e6))
* **go:** bootstrap go module ([2c124fb](https://github.com/amiwrpremium/go-derive/commit/2c124fb77321269d86fbad7d331a4b1f413b40fb))
* Makefile, lefthook hooks, and code-quality tool configs ([7806822](https://github.com/amiwrpremium/go-derive/commit/7806822ed2bed014c540a6ba8a10663de4077e30))
* Renovate + Dependabot dependency-management configuration ([7f37f4b](https://github.com/amiwrpremium/go-derive/commit/7f37f4b503339fe717145c43df101d73e42d7d68))


### Continuous Integration

* Codecov and Codacy coverage-reporting configuration ([f3e8d8f](https://github.com/amiwrpremium/go-derive/commit/f3e8d8ffd2e7db405b82aa5b4a591c88009193a7))
* declarative repo settings, branch + tag rulesets, templates ([741419f](https://github.com/amiwrpremium/go-derive/commit/741419facc06d7925712a39491338c9b85b46957))
* GitHub Actions workflows for testing, security, and release ([2eb2486](https://github.com/amiwrpremium/go-derive/commit/2eb2486b67547ccc1cd974128cef629f33245711))


### Miscellaneous

* **examples:** 91 runnable example programs covering every public surface ([b56a378](https://github.com/amiwrpremium/go-derive/commit/b56a378ac9c0c2f048c58689b21941ae821d7c0c))
* scaffold repository ([e66b14c](https://github.com/amiwrpremium/go-derive/commit/e66b14c93ca121e637c877d939bd827435223d88))
