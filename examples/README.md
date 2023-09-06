# Examples

## Get quotes for multiple assets

You will find a [script](./quotes.sh) to get quotes for multiple assets. This is a simple example you can get inspiration from.

### Requirements

You should have these CLI tools in your **PATH** to make it work:

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
echo "brk.a\nbrk.b" | ./quotes.sh assets | xsv table
```

will output:

```shell
date        brk.a      brk.a
28/07/2023  533627.50  533627.50
04/08/2023  533196.47  533196.47
11/08/2023  543579.95  543579.95
18/08/2023  534994.01  534994.01
25/08/2023  540124.97  540124.97
01/09/2023  549632.48  549632.48
```

Please note that the arguments of **quotes** are provided inside the script.
