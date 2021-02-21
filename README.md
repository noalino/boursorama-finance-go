# Quotes

A basic tool aiming at scraping financial assets quotes from the [boursorama](https://www.boursorama.com/bourse/) website.

Both an API and a CLI are available to use.

## Getting started

Clone the repository:

```shell
git clone https://github.com/benoitgelineau/go-fetch-quotes
cd go-fetch-quotes
```

Then build the Docker image:

```shell
docker build -t go-fetch-quotes .
```

#### To run the API:

```shell
docker run --rm --name go-fetch-quotes-api -p 8080:8080 go-fetch-quotes quotes-api
```

It starts the API on _localhost:8080_

#### To run the CLI

```shell
docker run --rm --name go-fetch-quotes-cli go-fetch-quotes quotes
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

You first need to `search` for the asset you want to get quotes from, and if there is a result, it will return a __SYMBOL__ value.

Choose the asset you were looking for and use the __SYMBOL__ value in the `get` command to fetch the quotes.

Example:

```shell
$ docker run --rm --name go-fetch-quotes-cli go-fetch-quotes quotes search berkshire
Searching for 'berkshire'...
Results found:

| SYMBOL |        NAME         | LAST PRICE |
|--------|---------------------|------------|
| BRK.B  | BERKSHIRE HATH RG-B | 241.690 $  |
|        | NYSE                |            |
|--------|---------------------|------------|
| BHLB   | BERKSHIRE HILLS     | 20.030 $   |
|        | NYSE                |            |
|--------|---------------------|------------|
| 3kBERK | BERKSHIRE BANCOR    | 9.450 $    |
|        | OTCBB               |            |
|--------|---------------------|------------|


$ docker run --rm --name go-fetch-quotes go-fetch-quotes quotes get 3kberk      
date,3kberk
19/02/2021,9.45
18/02/2021,9.41
16/02/2021,10.00
11/02/2021,9.31
10/02/2021,9.71
09/02/2021,10.00
08/02/2021,9.25
04/02/2021,9.21
03/02/2021,9.15
02/02/2021,9.95
28/01/2021,9.95
26/01/2021,10.00
25/01/2021,9.15
22/01/2021,10.20
19/01/2021,9.70
```

## Licensing

The code in this project is licensed under GPL-3.0 License.

