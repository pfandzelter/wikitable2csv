package wiki

import (
	"net/url"
	"regexp"
)

// CreateApiUrl transforms a MediaWiki page URL into an API query URL.
// URL schemes: https://www.mediawiki.org/wiki/Manual:Short_URL
// Example input:
// https://en.wikipedia.org/wiki/Lists_of_earthquakes
func CreateApiUrl(rawURL string) (title string, queryUrl string, err error) {
	var apiInRoot bool

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	pathname := parsed.Path
	search := parsed.RawQuery

	var match []string

	switch {
	// 1. http://example.org/w/index.php/Page_title (recent versions of MediaWiki, without CGI support)
	// -> parser.pathname: /w/index.php/Page_title
	case regexp.MustCompile(`^/w/index\.php/.+$`).MatchString(pathname):
		re := regexp.MustCompile(`^/w/index\.php/([^&#]+).*$`)
		match = re.FindStringSubmatch(pathname)

		// 2. http://example.org/w/index.php?title=Page_title (recent versions of MediaWiki, with CGI support)
		// -> parser.search: ?title=Wikip%C3%A9dia:Rapports/Nombre_de_pages_par_namespace&action=view
		// -> parser.pathname: /w/index.php
	case pathname == "/w/index.php":
		re := regexp.MustCompile(`^title=([^&#]+).*$`)
		match = re.FindStringSubmatch(search)

		// 3. http://example.org/wiki/Page_title This is the most common configuration, same as in Wikipedia, though not the default because it requires server side modifications
		// 4. http://example.org/view/Page_title
		// -> parser.pathname: /wiki/Lists_of_earthquakes
		// -> short url must begin with lowercase letter after first slash
	case regexp.MustCompile(`^/[a-z_-]+/[^&#]+.*$`).MatchString(pathname):
		re := regexp.MustCompile(`^/[a-z_-]+/([^&#]+).*$`)
		match = re.FindStringSubmatch(pathname)

		// 5. http://example.org/Page_title (not recommended)
		// --> parser.pathname: /Page_title
	case regexp.MustCompile(`^/.+$`).MatchString(pathname):
		re := regexp.MustCompile(`^/(.+)$`)
		match = re.FindStringSubmatch(pathname)
		apiInRoot = true
	}

	if len(match) > 1 {
		title = match[1]

		apiSlug := "w/"
		if apiInRoot {
			apiSlug = ""
		}

		queryUrl = parsed.Scheme + "://" + parsed.Host + "/" +
			apiSlug +
			"api.php?action=parse&redirects=true&format=json&errorformat=plaintext&origin=*&page=" +
			title +
			"&prop=text"
	}

	return title, queryUrl, nil
}
