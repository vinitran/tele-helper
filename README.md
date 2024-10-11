
# Telegram Tool with Chromedp

This project is a tool that automates logging into Telegram and retrieving `query_id` values from the Telegram app using `Chromedp`. It supports managing multiple Telegram accounts with customizable concurrency levels.

## Features

- **Login to Telegram**: Automates the login process for Telegram accounts.
- **Retrieve `query_id`**: Extracts `query_id` values from the Telegram app with support for multiple threads.
- **Import/Export User Data**: Allows exporting and importing user data to/from a zip file.

## Prerequisites

- Go 1.23+ installed on your system.
- Telegram account credentials.
- `chromedp` package installed (`go get -u github.com/chromedp/chromedp`).

## Installation

Clone the repository:

```bash
git https://github.com/vinitran/tele-helper.git
cd tele-helper
```

Install dependencies:

```bash
go mod tidy
```

## Usage

### 1. Login to Telegram

To log in to a Telegram account, run:

```bash
go run cmd/main login --name <account_name>
```

Replace `<account_name>` with the name of the account you want to log in to. This will initiate the login process using Chromedp.

### 2. Get `query_id`

To retrieve the `query_id` from the Telegram app, run:

```bash
go run cmd/main queryid --name <account_name> --threads <number_of_threads>
```

- Replace `<account_name>` with the name of the specific account.
- Use the `--threads` flag to specify the number of concurrent threads (default is 2).

Example:

```bash
go run cmd/main queryid --name user123 --threads 3
```

### 3. Export User Data

To export user data to a zip file:

```bash
go run cmd/main export --name <account_name>
```

This will export the data associated with the account to a zip file.

### 4. Import User Data

To import user data from a zip file:

```bash
go run cmd/main import
```

## Flags

- `--name`: Specifies the name of the Telegram account.
- `--threads`: Specifies the number of concurrent threads to use for retrieving `query_id` (default: 2).
- `--app`: App name (required for the `queryid` command).

## Example

```bash
go run cmd/main login --name user123
go run cmd/main queryid --name user123 --threads 4 --app blum
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.