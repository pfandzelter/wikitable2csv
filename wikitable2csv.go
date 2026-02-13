package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pfandzelter/wikitable2csv-cli/pkg/wiki"
)

// update this version when making changes by tagging the commit
// compile with make to get all this information automatically
// OR go build -ldflags "-X main.version=$(git describe --tags --always) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.builtBy=goreleaser".
// OR goreleaser will do this automatically
var version = "unknown"
var commit = "unknown"
var date = "unknown"

func main() {

	var printVersion *bool
	var url *string
	var tableSelector *string
	var cssClassNamesToExcludeUnparsed *string
	var excludeTablesUnparsed *string
	var includeTablesUnparsed *string
	var userAgent *string
	var noTrimCells *bool
	var noIncludeLineBreaks *bool
	var silent *bool

	printVersion = flag.Bool("version", false, "print wikitable2csv version and exit")

	url = flag.String("url", "", "(required) the Wiki page URL, e.g., https://en.wikipedia.org/wiki/Lists_of_earthquakes")
	tableSelector = flag.String("table-selector", ".wikitable", "table selector")
	cssClassNamesToExcludeUnparsed = flag.String("exclude-css-class-name", "reference", "CSS class name to exclude, can be multiple separated by comma")
	excludeTablesUnparsed = flag.String("exclude-tables", "", "table numbers to exclude from parsing to csv, seperated by comma (specifying this and --include-tables is possible but useless)")
	includeTablesUnparsed = flag.String("include-tables", "", "table numbers to include for parsing to csv, seperated by comma, leave empty for all tables (default)")
	userAgent = flag.String("user-agent", "wikitable2csv-cli", "user agent string to send")
	noTrimCells = flag.Bool("no-trim-cells", false, "disable trimming of cells")
	noIncludeLineBreaks = flag.Bool("no-include-line-breaks", false, "disable inclusion of line breaks")
	silent = flag.Bool("silent", false, "no output to stdout, errors may still be printed to stderr")

	flag.Parse()

	if *printVersion {
		fmt.Printf("wikitables2csv\nversion %s\nbuilt %s\ncommit %s\n", version, date, commit)
		os.Exit(0)
	}

	if *url == "" {
		fmt.Fprintln(os.Stderr, "no Wiki page URL given")
		os.Exit(1)
	}

	cssClassNamesToExclude := strings.Split(*cssClassNamesToExcludeUnparsed, ",")

	includeTables := make(map[int]struct{})

	for _, s := range strings.Split(*includeTablesUnparsed, ",") {
		if s == "" {
			continue
		}

		n, err := strconv.Atoi(s)

		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse %s as table index to include\n", s)
		}

		includeTables[n] = struct{}{}
	}

	excludeTables := make(map[int]struct{})

	for _, s := range strings.Split(*excludeTablesUnparsed, ",") {
		if s == "" {
			continue
		}

		n, err := strconv.Atoi(s)

		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse %s as table index to exclude\n", s)
		}

		excludeTables[n] = struct{}{}
	}

	// parse url and create API url
	title, queryUrl, err := wiki.CreateApiUrl(*url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse url: %s", err.Error())
		os.Exit(1)
	}
	if !*silent {

		fmt.Printf("getting tables for page %s...\n", title)
	}
	// fetch data
	wikiData, err := wiki.Fetch(queryUrl, *userAgent)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not fetch page from wiki: %s", err.Error())
		os.Exit(1)
	}

	// parse data
	if !*silent {

		fmt.Printf("parsing data for page %s...\n", title)
	}
	parsed, err := wiki.Parse(*wikiData, *tableSelector, cssClassNamesToExclude, *noTrimCells, *noIncludeLineBreaks)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse tables: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("parsed %d tables\n", len(parsed))

	for i, t := range parsed {
		humanIndex := i + 1

		if len(includeTables) > 0 {
			if _, ok := includeTables[humanIndex]; !ok {
				if !*silent {
					fmt.Printf("skipping table %d (not included)\n", humanIndex)
				}
				continue
			}
		}

		if _, ok := excludeTables[humanIndex]; ok {
			if !*silent {
				fmt.Printf("skipping table %d (excluded)\n", humanIndex)
			}
			continue
		}

		fName := fmt.Sprintf("%s-%d.csv", title, humanIndex)

		if !*silent {
			fmt.Printf("writing table %d to %s\n", humanIndex, fName)
		}

		f, err := os.Create(fName)

		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create file %s: %s", fName, err.Error())
			os.Exit(1)
		}

		err = wiki.WriteCSV(t, f)

		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write to file %s: %s", fName, err.Error())
			os.Exit(1)
		}
	}

	if !*silent {
		fmt.Printf("wrote %d tables to csv\n", len(parsed))
	}
}
