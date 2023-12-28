# Boursorama-finance-go

A basic tool aiming at scraping financial assets historical data from the [boursorama](https://www.boursorama.com/bourse/) website.

Both an API and a CLI are available to use.

## Getting started

Clone the repository:

```shell
git clone https://github.com/noalino/boursorama-finance-go
cd boursorama-finance-go
```

Then build the Docker image:

```shell
docker build -t bfinance-go .
```

#### To run the API

```shell
docker run --rm -p 8080:8080 bfinance-go bfinance-api
```

It starts the API on _localhost:8080_.

#### To run the CLI

[Download](https://github.com/noalino/boursorama-finance-go/releases) and run the binary file.

Or you can run it inside your terminal with Docker:

```shell
docker run --rm bfinance-go bfinance
```

## How it works

### API

[OpenAPI documentation](api/openapi.yml)

### CLI

Available commands:

- `search [options] ASSET`

```text
NAME:
   bfinance search - Search for a financial asset

USAGE:
   bfinance search [options] ASSET

OPTIONS:
   --page value, -P value  load specific page (default: 1)
   --pretty, -p            prettify the output (default: false)
   --verbose, -v           show more info (default: false)
   --help, -h              show help
```

- `get [options] SYMBOL`

```text
NAME:
   bfinance get - Return historical data

USAGE:
   bfinance get [options] SYMBOL

OPTIONS:
   --duration value, -d value  Specify the duration, it should be one of the following values:
                               [1M, 2M, 3M, 4M, 5M, 6M, 7M, 8M, 9M, 10M, 11M, 12M, 1Y, 2Y, 3Y] (default: "3M")
   --from value, -f value      Specify the start date, it must be in the following format:
                               DD/MM/YYYY (default: "21/11/2023")
   --period value, -p value    Specify the period, it should be one of the following values:
                               [daily, weekly, monthly, yearly] (default: "daily")
   --help, -h                  show help
```

You first need to `search` for the asset you want to get historical data from, and if there is a result, it will return a **SYMBOL** value.

Choose the asset you were looking for and use the **SYMBOL** value in the `get` command to fetch the data.

Example:

```shell
$ bfinance search --pretty --verbose apple
Searching for 'apple'...
Results found (page 1/702):
|-----------------|--------------------|------------------|-------------|
|     SYMBOL      |        NAME        |      MARKET      | CLOSE PRICE |
|-----------------|--------------------|------------------|-------------|
| AAPL            | APPLE              | NASDAQ           | 194.83 USD  |
|-----------------|--------------------|------------------|-------------|
| 1u0R2V.L        | APPLE              | LSE              | 199.00 USD  |
|-----------------|--------------------|------------------|-------------|
| 2aAAPL          | APPLE              | Swiss EBS Stocks | 193.53 EUR  |
|-----------------|--------------------|------------------|-------------|
| 1zAPC           | APPLE              | XETRA            | 178.58 EUR  |
|-----------------|--------------------|------------------|-------------|
| 1rPW94CB        | APPLE135.9SPLOPENB | Euronext Paris   | 4.16 EUR    |
|-----------------|--------------------|------------------|-------------|
| 1rPRJ5CB        | APPLE139.4SPLOPENB | Euronext Paris   | 0.79 EUR    |
|-----------------|--------------------|------------------|-------------|
| 1rPPS6CB        | APPLE152.7TPIOPENB | Euronext Paris   | 0.36 EUR    |
|-----------------|--------------------|------------------|-------------|
| 1rPX2QDB        | APPLE153.7SPSOPENB | Euronext Paris   | 1.18 EUR    |
|-----------------|--------------------|------------------|-------------|
| 1rPPP5DB        | APPLE141.2TCIOPENB | Euronext Paris   | 1.95 EUR    |
|-----------------|--------------------|------------------|-------------|
| 3rPFRSGE001LGZ2 | Tracker : Apple    | Dir Emet         | 1.09 EUR    |
|-----------------|--------------------|------------------|-------------|


$ bfinance get aapl
date,close,performance,high,low,open
20/11/2023,191.45,0.00%,191.90,189.88,189.88
21/11/2023,190.64,-0.42%,191.50,189.74,191.47
22/11/2023,191.31,+0.35%,192.93,190.83,191.47
24/11/2023,189.97,-0.70%,190.90,189.25,190.90
27/11/2023,189.79,-0.09%,190.67,188.90,189.90
28/11/2023,190.40,+0.32%,191.08,189.40,189.71
29/11/2023,189.37,-0.54%,192.09,188.97,190.98
30/11/2023,189.95,+0.31%,190.32,188.19,189.85
01/12/2023,191.24,+0.68%,191.56,189.23,190.32
04/12/2023,189.43,-0.95%,190.01,187.46,190.00
05/12/2023,193.42,+2.11%,194.40,190.21,190.22
06/12/2023,192.32,-0.57%,194.76,192.12,194.47
07/12/2023,194.27,+1.01%,195.00,193.59,193.68
08/12/2023,195.71,+0.74%,195.99,193.67,194.10
11/12/2023,193.18,-1.29%,193.49,191.43,193.02
12/12/2023,194.71,+0.79%,194.72,191.72,193.00
13/12/2023,197.96,+1.67%,198.00,194.88,195.00
14/12/2023,198.11,+0.08%,199.62,196.16,198.07
15/12/2023,197.57,-0.27%,198.40,197.02,197.38
18/12/2023,195.89,-0.85%,196.63,194.40,196.09
19/12/2023,196.94,+0.54%,196.95,195.89,196.08
20/12/2023,194.83,-1.07%,197.68,194.83,196.97

$ bfinance get --from 01/01/2023 --period monthly --duration 6M aapl
date,close,performance,high,low,open
30/12/2022,129.93,0.00%,129.95,127.43,128.32
31/01/2023,144.29,+11.05%,144.34,142.28,142.63
28/02/2023,147.41,+2.16%,149.08,146.87,146.87
31/03/2023,164.90,+11.86%,165.00,162.15,162.36
28/04/2023,169.68,+2.90%,169.85,167.88,168.58
31/05/2023,177.25,+4.46%,179.35,176.77,177.30
30/06/2023,193.97,+9.43%,194.48,191.27,191.65
31/07/2023,196.45,+1.28%,196.49,195.26,196.00
```

See [examples](./examples/README.md) if you want to know how to get historical data for multiple assets.

## Licensing

The code in this project is licensed under GPL-3.0 License.
