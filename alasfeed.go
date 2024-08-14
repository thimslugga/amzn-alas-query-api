package main

import (
	"encoding/json"
	"encoding/xml"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// alas rss feed urls
//var amazonLinuxFeed = "https://alas.aws.amazon.com/alas.rss"
var amazonLinux2Feed = "https://alas.aws.amazon.com/AL2/alas.rss"
//var amazonLinux2023Feed = "https://alas.aws.amazon.com/AL2022/alas.rss"
var amazonLinux2023Feed = "https://alas.aws.amazon.com/AL2023/alas.rss"

// regex used for alas parsing
var alasStringRegex = regexp.MustCompile("ALAS-[0-9]+-[0-9]+")
var pkgsRegex = regexp.MustCompile(": (?P<Block>.*$)")
var priorityRegex = regexp.MustCompile("\\((?P<Block>[a-z]+)\\)")
var newPkgListRegex = regexp.MustCompile("<pre>(?P<Block>.*?)</pre>")

// ALASResponse is the root of the RSS response from the ALAS feeds
type ALASResponse struct {
	Channel Channel `xml:"channel"`
}

// Channel contains the list of items in the RSS response
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Language    string `xml:"language"`
	TTL         int    `xml:"ttl"`
	Vulns       []Vuln `xml:"item"`
}

// Vuln is a single `item` in the RSS response
type Vuln struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Link        string `xml:"link"`
}

// ExpandedVuln is a Vuln with all of it's fields scraped already.
type ExpandedVuln struct {
	ALAS        string    `json:"alas"`
	CVEs        []string  `json:"cves"`
	Packages    []string  `json:"packages"`
	Priority    string    `json:"priority"`
	NewPackages []Package `json:"newPackages"`
	Link        string    `json:"link"`
	PubDate     string    `json:"pubDate"`
}

// GetALASFeed returns the unmarshaled feed for a given URL
func GetALASFeed(endpoint string) (feed ALASResponse, err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	xml.Unmarshal(data, &feed)
	return
}

// ToJSON dumps the current vuln to JSON
func (v *ExpandedVuln) ToJSON() (data []byte) {
	data, _ = json.Marshal(v)
	return
}

// Expand parses and scrapes all the available fields for a vuln. Returns
// the fully populated ExpandedVuln.
func (v *Vuln) Expand() (expanded *ExpandedVuln) {
	expanded = &ExpandedVuln{
		Link:     v.Link,
		PubDate:  v.PubDate,
		ALAS:     v.ALASString(),
		CVEs:     v.CVEList(),
		Packages: v.Packages(),
		Priority: v.Priority(),
	}
	var err error
	newPkgs, err := v.NewPackages()
	if err != nil {
		log.Printf("WARNING: Could not parse new packages for %s\n", expanded.ALAS)
	} else {
		expanded.NewPackages = make([]Package, 0)
		for _, pkg := range newPkgs {
			parsed, err := NewPackageFromString(pkg)
			if err != nil {
				log.Printf("WARNING: Could not parse NEVRA from %s: %s\n", pkg, err)
				continue
			}
			expanded.NewPackages = append(expanded.NewPackages, parsed)
		}
	}
	return
}

// ALASString returns the ALAS ID for an item in the RSS feed
func (v *Vuln) ALASString() string {
	return alasStringRegex.FindString(v.Title)
}

// CVEList returns the list of CVEs associated with an item in the RSS feed
func (v *Vuln) CVEList() (cveList []string) {
	cveList = make([]string, 0)
	trim := strings.TrimSpace(v.Description)
	if trim != "" {
		split := strings.Split(trim, ", ")
		for _, x := range split {
			cveList = append(cveList, x)
		}
	}
	return
}

// Packages returns the list of packages affected in an item from the RSS feed
func (v *Vuln) Packages() []string {
	pkgs := pkgsRegex.FindStringSubmatch(v.Title)[1]
	return strings.Split(pkgs, ", ")
}

// Priority returns the priority of the item in the RSS feed.
func (v *Vuln) Priority() string {
	return priorityRegex.FindStringSubmatch(v.Title)[1]
}

// NewPackages will scrape the description page for the ALAS and return the
// list of updated packages.
func (v *Vuln) NewPackages() (pkgs []string, err error) {
	pkgs = make([]string, 0)
	resp, err := http.Get(v.Link)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	raw := newPkgListRegex.FindStringSubmatch(string(body))[1]
	unescaped := html.UnescapeString(raw)
	for _, x := range strings.Split(unescaped, "<br />") {
		if trim := strings.TrimSpace(x); trim != "" {
			if string(trim[len(trim)-1]) != ":" {
				// its not the arch line
				pkgs = append(pkgs, trim)
			}
		}
	}
	return
}
