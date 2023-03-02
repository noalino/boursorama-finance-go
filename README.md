# Boursorama-finance-go

A basic tool aiming at scraping financial assets quotes from the [boursorama](https://www.boursorama.com/bourse/) website.

Both an API and a CLI are available to use.

## Getting started

Clone the repository:

```shell
git clone https://github.com/benoitgelineau/boursorama-finance-go
cd boursorama-finance-go
```

Then build the Docker image:

```shell
docker build -t boursorama-finance-go .
```

#### To run the API:

```shell
docker run --rm --name boursorama-finance-go-api -p 8080:8080 boursorama-finance-go quotes-api
```

It starts the API on _localhost:8080_

#### To run the CLI

[Download](https://github.com/benoitgelineau/boursorama-finance-go/releases) and run the binary file.

Or you can run it inside your terminal with Docker:

```shell
docker run --rm --name boursorama-finance-go-cli boursorama-finance-go quotes
```

## How it works

### API

[OpenAPI documentation](internal/api/openapi.yml)

### CLI

Available commands:

- `search NAME | ISIN`

```text
Usage: quotes search NAME | ISIN

Flags:

  -help
    	Get help on the 'quotes search' command.
```

- `get [OPTIONS] SYMBOL`

```text
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
    	["1","7","30","365"] (default "1")
```

You first need to `search` for the asset you want to get quotes from, and if there is a result, it will return a **SYMBOL** value.

Choose the asset you were looking for and use the **SYMBOL** value in the `get` command to fetch the quotes.

Example:

```shell
$ docker run --rm boursorama-finance-go quotes search berkshire
Searching for 'berkshire'...
Results found:

|  SYMBOL  |        NAME         | MARKET |   LAST PRICE   |
|----------|---------------------|--------|----------------|
| BRK.A    | BERKSHIRE HATH RG-A | NYSE   | 462 127.74 USD |
|----------|---------------------|--------|----------------|
| BRK.B    | BERKSHIRE HATH RG-B | NYSE   | 304.71 USD     |
|----------|---------------------|--------|----------------|
| BHLB     | BERKSHIRE HILLS     | NYSE   | 29.14 USD      |
|----------|---------------------|--------|----------------|
| BGRY     | BERKSHIRE GREY RG-A | NASDAQ | 1.30 USD       |
|----------|---------------------|--------|----------------|
| 1u0HN0.L | BERKSHIRE HATH RG-A | LSE    | 0.00 USD       |
|----------|---------------------|--------|----------------|
| 1zBRYN   | BERKSHIRE HATH RG-B | XETRA  | 285.10 EUR     |
|----------|---------------------|--------|----------------|
| 1u0R37.L | BERKSHIRE HATH RG-B | LSE    | 303.58 USD     |
|----------|---------------------|--------|----------------|


$ docker run --rm boursorama-finance-go quotes get 3kberk
date,brk.a
01/03/2023,462127.74
28/02/2023,463349.99
27/02/2023,461852.25
24/02/2023,461252.50
23/02/2023,460287.00
22/02/2023,460256.19
21/02/2023,459079.94
17/02/2023,467117.62
16/02/2023,467651.13
15/02/2023,469395.00
14/02/2023,472462.50
13/02/2023,476385.06
10/02/2023,472245.00
09/02/2023,466885.00
08/02/2023,467747.50
07/02/2023,474210.00
06/02/2023,466567.50
03/02/2023,467059.38
02/02/2023,472286.12
01/02/2023,470606.00
```

## Licensing

The code in this project is licensed under GPL-3.0 License.
