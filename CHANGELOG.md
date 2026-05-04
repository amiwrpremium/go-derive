# Changelog

All notable changes to this project will be documented in this file. See
[Conventional Commits](https://www.conventionalcommits.org/) for commit message
guidelines.

This file is maintained by [release-please](https://github.com/googleapis/release-please);
do not edit it by hand.

## [0.1.0] - 2026-05-02

### Features

- Initial public release.
- REST client (`pkg/rest`) covering markets, orders, positions, collateral, transactions, subaccounts, RFQ, MMP.
- WebSocket client (`pkg/ws`) with JSON-RPC calls and generics-typed subscriptions.
- Typed channel descriptors in `pkg/channels` (public + private).
- EIP-712 action signing and EIP-191 auth headers (`pkg/auth`) with `LocalSigner` and `SessionKeySigner`.
- Domain types (`pkg/types`), enums (`pkg/enums`), errors (`pkg/errors`).
- Top-level facade in `pkg/derive`.
- Stubbed `pkg/contracts` package.
- Examples for public REST, private REST, public WS subscribe, private WS subscribe, and place-order-via-WS.
