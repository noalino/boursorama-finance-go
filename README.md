# Boursorama-finance-go

A basic tool aiming at scraping financial assets quotes from the [boursorama](https://www.boursorama.com/bourse/) website.

Both an API and a CLI are available to use.

## Getting started

Clone the repository:

```shell
git clone https://github.com/noalino/boursorama-finance-go
cd boursorama-finance-go
```

Then build the Docker image:

```shell
docker build -t boursorama-finance-go .
```

#### To run the API

```shell
docker run --rm --name boursorama-finance-go-api -p 8080:8080 boursorama-finance-go quotes-api
```

It starts the API on _localhost:8080_.

#### To run the CLI

[Download](https://github.com/noalino/boursorama-finance-go/releases) and run the binary file.

Or you can run it inside your terminal with Docker:

```shell
docker run --rm --name boursorama-finance-go-cli boursorama-finance-go quotes
```

## How it works

### API

[OpenAPI documentation](api/openapi.yml)

### CLI

Available commands:

- `search [NAME | ISIN]`

```text
Quotes search - Search a financial asset

Search a financial asset by name or ISIN and return the following information:
Symbol, Name, Category, Last price

Usage: quotes search [NAME | ISIN]

Flags:

  -help
     Get help on the 'quotes search' command.
  -page uint
     Select page. (default 1)
  -pretty
     Display output in a table.
  -verbose
     Log more info.
```

- `get [OPTIONS] SYMBOL`

```text
Quotes get - Return quotes

Usage: quotes get [OPTIONS] SYMBOL

Flags:

  -duration string
     Specify the duration, it should be one of the following values:
     ["1M","2M","3M","4M","5M","6M","7M","8M","9M","10M","11M","1Y","2Y","3Y"] (default "3M")
  -from string
     Specify the start date, it must be in the following format:
      DD/MM/YYYY (default "a month from now")
  -help
     Get help on the 'quotes get' command.
  -period string
     Specify the period, it should be one the following values:
     ["daily","weekly","monthly","yearly"] (default "daily")
```

You first need to `search` for the asset you want to get quotes from, and if there is a result, it will return a **SYMBOL** value.

Choose the asset you were looking for and use the **SYMBOL** value in the `get` command to fetch the quotes.

Example:

```shell
$ quotes search --pretty --verbose berkshire
Searching for 'berkshire'...
Results found:
|----------|---------------------|--------------------|----------------|
|  SYMBOL  |        NAME         |       MARKET       |   LAST PRICE   |
|----------|---------------------|--------------------|----------------|
| BHLB     | BERKSHIRE HILLS     | NYSE               | 21.78 USD      |
|----------|---------------------|--------------------|----------------|
| BRK.B    | BERKSHIRE HATH RG-B | NYSE               | 362.39 USD     |
|----------|---------------------|--------------------|----------------|
| BRK.A    | BERKSHIRE HATH RG-A | NYSE               | 549 632.48 USD |
|----------|---------------------|--------------------|----------------|
| 1u0HN0.L | BERKSHIRE HATH RG-A | LSE                | 0.00 USD       |
|----------|---------------------|--------------------|----------------|
| 1u0R37.L | BERKSHIRE HATH RG-B | LSE                | 362.50 USD     |
|----------|---------------------|--------------------|----------------|
| 1zBRYN   | BERKSHIRE HATH RG-B | XETRA              | 337.60 EUR     |
|----------|---------------------|--------------------|----------------|
| 1rAJW63B | BERKSHIRE HA/BNP WT | Euronext Amsterdam | 0.00 EUR       |
|----------|---------------------|--------------------|----------------|
| 1rANG58B | BERKSHIRE /BNP P-WT | Euronext Amsterdam | 0.20 EUR       |
|----------|---------------------|--------------------|----------------|
| 1rAJW02B | BERKSHIRE HA/BNP WT | Euronext Amsterdam | 0.00 EUR       |
|----------|---------------------|--------------------|----------------|
| 1rAP649N | BERKSHIRE /AAB P-WT | Euronext Amsterdam | 0.55 EUR       |
|----------|---------------------|--------------------|----------------|


$ quotes get brk.b
date,close,open
04/08/2023,350.24,354.60
07/08/2023,362.49,355.98
08/08/2023,363.81,359.87
09/08/2023,357.97,364.01
10/08/2023,357.09,359.65
11/08/2023,358.31,356.02
14/08/2023,358.53,358.87
15/08/2023,354.37,357.10
16/08/2023,354.09,354.69
17/08/2023,352.99,354.23
18/08/2023,352.67,351.12
21/08/2023,352.36,354.24
22/08/2023,350.53,353.33
23/08/2023,354.27,351.11
24/08/2023,354.21,354.21
25/08/2023,355.95,354.35
28/08/2023,355.55,357.06
29/08/2023,358.44,355.24
30/08/2023,361.11,358.50
31/08/2023,359.97,362.32
01/09/2023,362.39,362.39
05/09/2023,360.33,363.80

$ quotes get --from 01/01/2023 --period monthly --duration 6M brk.b
date,close,open
30/12/2022,308.93,306.34
31/01/2023,311.55,307.61
28/02/2023,305.02,304.89
31/03/2023,308.25,305.74
28/04/2023,328.44,325.70
31/05/2023,321.20,321.87
30/06/2023,341.18,338.69
31/07/2023,352.01,350.73
```

See [examples](./examples/README.md) if you want to know how to get quotes for multiple assets.

## Licensing

The code in this project is licensed under GPL-3.0 License.
