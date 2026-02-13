# wikitable2csv-cli

A CLI tool to extract tables from Wiki pages and convert them to CSV.
This tool is a fork of [`gambolputty/wikitable2csv`](https://github.com/gambolputty/wikitable2csv), rewritten in Go. Most of the source logic has been converted verbatim, frontend code has been removed, CLI code has been added.

A web tool to extract tables from Wiki pages and convert them to CSV. Use it online [here](https://wikitable2csv.ggor.de/).

## Install

Install with go (>=1.22) installed: `go install github.com/pfandzelter/wikitable2csv-cli@latest`

Alternatively, you may download the binary for your platform from the release page.

You may also clone this repository and run `go install`.

## Usage

```text
Usage of ./wikitable2csv:
  -exclude-css-class-name string
     CSS class name to exclude, can be multiple separated by comma (default "reference")
  -exclude-tables string
     table numbers to exclude from parsing to csv, seperated by comma (specifying this and --include-tables is possible but useless)
  -include-tables string
     table numbers to include for parsing to csv, seperated by comma, leave empty for all tables (default)
  -no-include-line-breaks
     disable inclusion of line breaks
  -no-trim-cells
     disable trimming of cells
  -silent
     no output to stdout, errors may still be printed to stderr
  -table-selector string
     table selector (default ".wikitable")
  -url string
     (required) the Wiki page URL, e.g., https://en.wikipedia.org/wiki/Lists_of_earthquakes
  -user-agent string
     user agent string to send (default "wikitable2csv-cli")
  -version
     print wikitable2csv version and exit
```

- Download all tables from a Wikipedia page:

    ```sh
    wikitable2csv -url https://en.wikipedia.org/wiki/Lists_of_earthquakes
    ```

    For each table in the page, this will create a CSV file in your current directory.

- Download a specific table from a Wikipedia page:

    ```sh
    wikitable2csv -include-tables 2 -url https://en.wikipedia.org/wiki/Lists_of_earthquakes
    ```

    This will only download the second table on the page to your current directory.

- Download a specific table from a Wikipedia page:

    ```sh
    wikitable2csv -include-tables 2,4,5 -url https://en.wikipedia.org/wiki/Lists_of_earthquakes
    ```

    This will only download the second, fourth, and fifth tables on the page to your current directory.
    If there are fewer tables on the page than you specify, no error is given.

- Download all tables but specific tables:

    ```sh
    wikitable2csv -exclude-tables 1,3 -url https://en.wikipedia.org/wiki/Lists_of_earthquakes
    ```

    This will only download all tables but the first and third table.

- Specify a user agent:

    ```sh
    wikitable2csv -user-agent "wikitable2csv-custom" -url https://en.wikipedia.org/wiki/Lists_of_earthquakes
    ```

    This changes the user agent of your request, useful to let Wikipedia know who is calling.

## License

[MIT](https://github.com/pfandzelter/wikitable2csv-cli/blob/main/LICENSE) © Gregor Weichbrodt, Tobias Pfandzelter
