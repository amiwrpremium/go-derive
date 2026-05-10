# go-derive examples

Every public method and helper has its own runnable program. Each lives in
its own directory so the tree mirrors the SDK's surface area; running any
one is a single `go run`.

## Layout

```text
examples/
  example/                  shared helper package (env loading, client builders)

  derive/                   top-level facade
    new_client/             construct a Client
    rest_only/              use c.REST only
    ws_only/                use c.WS only (Connect + GetTicker)
    full/                   Connect + Login + REST + WS together

  auth/                     signing primitives
    local_signer/           NewLocalSigner from env
    session_key_signer/     NewSessionKeySigner with separate owner
    sign_action/            EIP-712 ActionData signing
    sign_auth_header/       EIP-191 timestamp signing
    http_headers/           HTTPHeaders builder
    nonce/                  NonceGen.Next loop
    action_data/            ActionData.Hash
    trade_module_data/      TradeModuleData.Hash
    transfer_module_data/   TransferModuleData.Hash
    signature/              Signature.Hex

  rest/
    public/                 (no auth)
      get_time/                       c.GetTime
      get_currencies/                 c.GetCurrencies
      get_currency/                   c.GetCurrency (DERIVE_CURRENCY)
      get_instruments/                c.GetInstruments(currency, type)
      get_instrument/                 c.GetInstrument(name)
      get_all_instruments/            c.GetAllInstruments
      get_ticker/                     c.GetTicker
      get_tickers/                    c.GetTickers
      get_orderbook/                  c.GetOrderBook
      get_trade_history/              c.GetPublicTradeHistory
      get_option_settlement_prices/   c.GetOptionSettlementPrices (DERIVE_CURRENCY)
      get_live_incidents/             c.GetLiveIncidents
      get_index_chart_data/           c.GetIndexChartData
      get_tradingview_chart_data/     c.GetTradingViewChartData
      get_spot_feed_history_candles/  c.GetSpotFeedHistoryCandles
      get_interest_rate_history/      c.GetInterestRateHistory
      get_liquidation_history/        c.GetPublicLiquidationHistory
      get_maker_programs/             c.GetMakerPrograms
      get_maker_program_scores/       c.GetMakerProgramScores (DERIVE_PROGRAM_NAME, DERIVE_EPOCH_START)
      get_referral_performance/       c.GetReferralPerformance (DERIVE_REFERRAL_CODE)
      get_vault_balances/             c.GetVaultBalances (DERIVE_WALLET)
      get_vault_share/                c.GetVaultShare (DERIVE_VAULT_NAME)
      get_vault_statistics/           c.GetVaultStatistics
      get_vault_assets/               c.GetVaultAssets
      get_vault_pools/                c.GetVaultPools
      get_vault_rates/                c.GetVaultRates (DERIVE_VAULT_TYPE)
      order_quote/                    c.OrderQuotePublic
      margin_watch/                   c.MarginWatch (DERIVE_SUBACCOUNT)
      all_statistics/                 c.GetAllStatistics
      all_user_statistics/            c.GetAllUserStatistics
      user_statistics/                c.GetUserStatistics (DERIVE_WALLET)
      get_detailed_maker_snapshot_history/  c.GetDetailedMakerSnapshotHistory (DERIVE_PROGRAM_NAME, DERIVE_EPOCH_START)
      get_all_points/                 c.GetAllPoints (DERIVE_PROGRAM_NAME)
      get_points/                     c.GetPoints (DERIVE_PROGRAM_NAME, DERIVE_WALLET)
      get_points_leaderboard/         c.GetPointsLeaderboard (DERIVE_PROGRAM_NAME)
      get_all_referral_codes/         c.GetAllReferralCodes
      get_referral_code/              c.GetReferralCode (DERIVE_WALLET)
      get_invite_code/                c.GetInviteCode (DERIVE_WALLET)
      validate_invite_code/           c.ValidateInviteCode (DERIVE_INVITE_CODE)
      get_asset/                      c.GetAsset (DERIVE_ASSET_NAME)
      get_assets/                     c.GetAssets (DERIVE_CURRENCY)
      get_descendant_tree/            c.GetDescendantTree (DERIVE_WALLET_OR_INVITE_CODE)
      get_tree_roots/                 c.GetTreeRoots
      get_bridge_balances/            c.GetBridgeBalances
      get_stdrv_snapshots/            c.GetStDRVSnapshots (DERIVE_WALLET, DERIVE_FROM_SEC, DERIVE_TO_SEC)
    private/                (DERIVE_SESSION_KEY + DERIVE_SUBACCOUNT)
      get_subaccount/         c.GetSubaccount
      get_subaccounts/        c.GetSubaccounts
      get_all_portfolios/     c.GetAllPortfolios
      get_collateral/         c.GetCollateral
      get_positions/          c.GetPositions
      get_open_orders/        c.GetOpenOrders
      get_order/              c.GetOrder (DERIVE_ORDER_ID)
      get_orders/             c.GetOrders
      get_order_history/      c.GetOrderHistory
      get_trade_history/      c.GetTradeHistory
      get_deposit_history/    c.GetDepositHistory
      get_withdrawal_history/ c.GetWithdrawalHistory
      get_liquidator_history/ c.GetLiquidatorHistory
      get_algo_orders/        c.GetAlgoOrders
      get_trigger_orders/     c.GetTriggerOrders
      orders/
        place/                       c.PlaceOrder (DERIVE_RUN_LIVE_ORDERS=1)
        place_algo/                  c.PlaceAlgoOrder (DERIVE_RUN_LIVE_ORDERS=1)
        place_trigger/               c.PlaceTriggerOrder (DERIVE_RUN_LIVE_ORDERS=1)
        cancel/                      c.CancelOrder (DERIVE_ORDER_ID)
        cancel_all/                  c.CancelAll
        cancel_by_label/             c.CancelByLabel (DERIVE_LABEL)
        cancel_by_instrument/        c.CancelByInstrument
        cancel_trigger_order/        c.CancelTriggerOrder (DERIVE_ORDER_ID)
        cancel_all_trigger_orders/   c.CancelAllTriggerOrders (DERIVE_RUN_LIVE_ORDERS=1)
        cancel_algo_order/           c.CancelAlgoOrder (DERIVE_ORDER_ID)
        cancel_all_algo_orders/      c.CancelAllAlgoOrders (DERIVE_RUN_LIVE_ORDERS=1)
      contact_info/
        create/              c.CreateContactInfo (DERIVE_CONTACT_TYPE, DERIVE_CONTACT_VALUE)
        get/                 c.GetContactInfo (optional DERIVE_CONTACT_TYPE)
        update/              c.UpdateContactInfo (DERIVE_CONTACT_ID, DERIVE_CONTACT_VALUE)
        delete/              c.DeleteContactInfo (DERIVE_CONTACT_ID)
      rfq/
        send/                c.SendRFQ
        poll/                c.PollRFQs
        cancel/              c.CancelRFQ (DERIVE_RFQ_ID)
        replace_quote/       c.ReplaceQuote (DERIVE_RUN_LIVE_ORDERS=1)
      mmp/
        set/                 c.SetMMPConfig
        reset/               c.ResetMMP

  ws/
    public/
      connect/              ws.Connect lifecycle
      get_time/ ... get_trade_history/   same RPC set as rest/public/, over WS
      subscribe/
        orderbook/                c.SubscribeOrderBook
        trades/                   c.SubscribeTrades (per-instrument)
        trades_by_type/           c.SubscribeTradesByType (per type+currency)
        trades_by_type_settled/   c.SubscribeTradesByTypeWithStatus
        ticker/                   c.SubscribeTickerSlim (compact)
        ticker_full/              c.SubscribeTicker (full)
        spot_feed/                c.SubscribeSpotFeed
        margin_watch/             c.SubscribeMarginWatch
        auctions_watch/           c.SubscribeAuctionsWatch
    private/
      login/                public/login RPC
      get_subaccount/ ... get_withdrawal_history/   private RPC set
      orders/               place / cancel / cancel_all / cancel_by_label / cancel_by_instrument
      rfq/                  send / poll / cancel
      mmp/                  set / reset
      subscribe/            subaccount channels
        orders/             c.SubscribeOrders
        balances/           c.SubscribeBalances
        trades/             c.SubscribeSubaccountTrades
        trades_settled/     c.SubscribeSubaccountTradesByStatus
        rfqs/               c.SubscribeRFQs
        quotes/             c.SubscribeQuotes
        best_quotes/        c.SubscribeBestQuotes
```

## Running one

```bash
# Public (no creds):
go run ./examples/rest/public/get_ticker
go run ./examples/ws/public/subscribe/orderbook

# Private (creds required):
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
  go run ./examples/rest/private/get_open_orders

# Live order placement (testnet only):
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
DERIVE_BASE_ASSET=0x... DERIVE_RUN_LIVE_ORDERS=1 \
  go run ./examples/rest/private/orders/place
```

## Building everything

```bash
go build ./examples/...
```

Should compile every program in the tree. CI runs this on every PR.

## Environment variables

Common to most examples, loaded by `examples/example/example.go`:

| Var | Default | Used by |
|---|---|---|
| `DERIVE_NETWORK` | `testnet` | network selection (`mainnet` or `testnet`) |
| `DERIVE_INSTRUMENT` | `BTC-PERP` | the instrument in market-data examples |
| `DERIVE_SESSION_KEY` | (unset) | hex private key, required for any `private/` example |
| `DERIVE_OWNER` | (unset) | smart-account owner — required when using a session-key signer |
| `DERIVE_SUBACCOUNT` | (unset) | numeric subaccount id, required for `private/` |
| `DERIVE_BASE_ASSET` | (unset) | on-chain asset address required by the `orders/place/` example |
| `DERIVE_RUN_LIVE_ORDERS` | (unset) | set to `1` to actually place an order |
| `DERIVE_ORDER_ID` | (unset) | for `get_order/`, `cancel/`, `cancel_trigger_order/` |
| `DERIVE_RFQ_ID` | (unset) | for the `rfq/cancel/` example |
| `DERIVE_LABEL` | (unset) | for `cancel_by_label/` |
| `DERIVE_CURRENCY` | (varies) | for `get_currency/` (ETH), `get_option_settlement_prices/` (BTC), index/spot chart examples (BTC) |
| `DERIVE_PROGRAM_NAME` | (unset) | required by `get_maker_program_scores/` |
| `DERIVE_EPOCH_START` | (unset) | required by `get_maker_program_scores/` (Unix seconds) |
| `DERIVE_REFERRAL_CODE` | (unset) | optional filter for `get_referral_performance/` |
| `DERIVE_WALLET` | (unset) | optional smart-contract wallet for `get_vault_balances/` |
| `DERIVE_VAULT_NAME` | (unset) | required by `get_vault_share/` |
| `DERIVE_VAULT_TYPE` | (unset) | optional for `get_vault_rates/` (e.g. `lbtc`, `weeth`) |
| `DERIVE_ASSET_NAME` | (unset) | required by `get_asset/` |
| `DERIVE_INVITE_CODE` | (unset) | required by `validate_invite_code/` |
| `DERIVE_WALLET_OR_INVITE_CODE` | (unset) | required by `get_descendant_tree/` |
| `DERIVE_FROM_SEC` / `DERIVE_TO_SEC` | (unset) | required by `get_stdrv_snapshots/` (Unix seconds) |
| `DERIVE_CONTACT_TYPE` | (unset) | required by `contact_info/create/`, optional filter for `contact_info/get/` |
| `DERIVE_CONTACT_VALUE` | (unset) | required by `contact_info/{create,update}/` |
| `DERIVE_CONTACT_ID` | (unset) | required by `contact_info/{update,delete}/` |

Never paste a real mainnet key — testnet is the default for a reason.
