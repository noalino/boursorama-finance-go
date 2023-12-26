# Examples

## Get historical data for multiple assets

You will find a [script](./quotes.sh) to get historical data for multiple assets. This is a simple example you can get inspiration from.

### Requirements

You should have these CLI tools in your **PATH** to make it work:

- [GNU sed](https://www.gnu.org/software/sed/): replace CSV header
- [GNU parallel](https://www.gnu.org/software/parallel/): run commands in parallel
- [xsv](https://github.com/BurntSushi/xsv): handle CSV from CLI
- [quotes](https://github.com/noalino/boursorama-finance-go/releases) binary file

### Usage

The script:

- reads symbols from the stdin
- create CSV files for each symbol in the provided output folder (only argument to the script)
- merge these files from the _date_ column
- output the result in stdout.

For instance, running this:

```shell
echo "aapl\nmsft\ntsla" | ./quotes.sh assets | xsv table
```

will output:

```shell
date        aapl    msft    tsla
17/11/2023  189.69  369.85  234.30
24/11/2023  189.97  377.43  235.45
01/12/2023  191.24  374.51  238.83
08/12/2023  195.71  374.23  243.84
15/12/2023  197.57  370.73  253.50
```

Please note that the arguments of **quotes** are provided inside the script.
