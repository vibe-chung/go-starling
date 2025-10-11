# go-starling
Starling CLI in go

A command line interface, using cobra, to interface with the starling bank api using the following schema:

https://api.starlingbank.com/api/openapi.yaml

## Installation

```bash
go build
```

## Usage

### Login

Save your Starling Bank API access token:

```bash
./go-starling login
```

### List Accounts

List all your Starling Bank accounts:

```bash
./go-starling list-accounts
```

This will display your account details including:
- Account UID
- Default Category UID
- Account name and type
- Currency

### List Transactions

List transactions for a specific account:

```bash
./go-starling list-transactions \
  --account-uid <ACCOUNT_UID> \
  --category-uid <CATEGORY_UID>
```

Optional flags:
- `--changes-since`: Timestamp to get transactions since (RFC3339 format). Defaults to 30 days ago.

Example:

```bash
./go-starling list-transactions \
  --account-uid aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee \
  --category-uid ffffffff-0000-1111-2222-333333333333
```

With custom date range:

```bash
./go-starling list-transactions \
  --account-uid aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee \
  --category-uid ffffffff-0000-1111-2222-333333333333 \
  --changes-since 2025-01-01T00:00:00Z
```

