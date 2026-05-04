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
      get_time/             c.GetTime
      get_currencies/       c.GetCurrencies
      get_instruments/      c.GetInstruments(currency, type)
      get_instrument/       c.GetInstrument(name)
      get_ticker/           c.GetTicker
      get_orderbook/        c.GetOrderBook
      get_trade_history/    c.GetPublicTradeHistory
    private/                (DERIVE_SESSION_KEY + DERIVE_SUBACCOUNT)
      get_subaccount/       c.GetSubaccount
      get_subaccounts/      c.GetSubaccounts
      get_collateral/       c.GetCollateral
      get_positions/        c.GetPositions
      get_open_orders/      c.GetOpenOrders
      get_order/            c.GetOrder (DERIVE_ORDER_ID)
      get_order_history/    c.GetOrderHistory
      get_trade_history/    c.GetTradeHistory
      get_deposit_history/  c.GetDepositHistory
      get_withdrawal_history/  c.GetWithdrawalHistory
      orders/
        place/              c.PlaceOrder (DERIVE_RUN_LIVE_ORDERS=1)
        cancel/             c.CancelOrder (DERIVE_ORDER_ID)
        cancel_all/         c.CancelAll
        cancel_by_label/    c.CancelByLabel (DERIVE_LABEL)
        cancel_by_instrument/  c.CancelByInstrument
      rfq/
        send/               c.SendRFQ
        poll/               c.PollRFQs
        cancel/             c.CancelRFQ (DERIVE_RFQ_ID)
      mmp/
        set/                c.SetMMPConfig
        reset/              c.ResetMMP

  ws/
    public/
      connect/              ws.Connect lifecycle
      get_time/ ... get_trade_history/   same RPC set as rest/public/, over WS
      subscribe/
        orderbook/          ws.Subscribe[types.OrderBook]
        trades/             ws.Subscribe[[]types.Trade]
        ticker/             ws.Subscribe[types.Ticker]
        instruments/        ws.Subscribe[[]types.Instrument]
    private/
      login/                public/login RPC
      get_subaccount/ ... get_withdrawal_history/   private RPC set
      orders/               place / cancel / cancel_all / cancel_by_label / cancel_by_instrument
      rfq/                  send / poll / cancel
      mmp/                  set / reset
      subscribe/            subaccount channels
        orders/
        positions/
        balances/
        trades/
        rfqs/
        quotes/
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
| `DERIVE_ORDER_ID` | (unset) | for `get_order/` and `cancel/` examples |
| `DERIVE_RFQ_ID` | (unset) | for the `rfq/cancel/` example |
| `DERIVE_LABEL` | (unset) | for `cancel_by_label/` |

Never paste a real mainnet key — testnet is the default for a reason.
