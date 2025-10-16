# Sui Go Client

A Go client SDK for the public Sui gRPC (currently) endpoints. It wraps the generated protobuf stubs with helpers, pagination, coin selection, and transaction utils so you can interact with Sui fullnodes.

## Installation

```
go get github.com/0xdraco/go-sui
```

## Features

- Dial helpers for mainnet/devnet/testnet and custom endpoints.
- Strongly typed service accessors (`LedgerClient`, `StateClient`, etc.).
- Convenience helpers for common read APIs:
  - `GetObject`, `BatchGetObjects`, `GetTransaction`, checkpoint & epoch helpers.
  - Automatic pagination for `ListOwnedObjects`, `ListBalances`, `ListDynamicFields`, and package versions.
- Coin selection utilities (`SelectCoins`, `SelectUpToNLargestCoins`) for gas/payment flows.
- Transaction helpers:
  - `SimulateTransaction` with optional gas selection.
  - `ExecuteTransactionAndWait` / `ExecuteSignedTransactionAndWait` that block until the transaction appears in a checkpoint.

## Getting Started

```go
ctx := context.Background()
client, err := grpc.NewMainnetClient(ctx)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

obj, err := client.GetObject(ctx, "0x123...", nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println("object version", obj.GetVersion())
```

For coin selection + transaction execution see `grpc/coin_selection.go` and `grpc/transaction.go` for examples.

## Testing

Coming soon.

## Regenerating Protos

```
make proto
```

This expects `protoc` and `protoc-gen-go` / `protoc-gen-go-grpc` to be installed.

## Contributing

1. Fork and clone this repository.
2. Install Go 1.24+ and the protobuf toolchain (`protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`).
3. Run `go test ./...` before submitting a PR. If you add helpers around live RPCs, include unit tests that use fakes alongside any integration tests.
4. Regenerate gRPC code with `make proto` when proto files change, and run `gofmt` on modified Go files.
5. Create a pull request describing the change.

Issues and feature requests are welcome via GitHub Issues.

## License

Apache-2.0 (matching upstream Sui SDK licensing).
