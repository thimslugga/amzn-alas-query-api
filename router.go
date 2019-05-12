package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	version "github.com/hashicorp/go-version"
	"github.com/julienschmidt/httprouter"
	rpmutils "github.com/sassoftware/go-rpmutils"
)

// NewRouter creates and returns a new httprouter
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/vulns", getVulns)
	return router
}

// ErrorResponse represents an error response to send back to a client
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetVulnsResponse represents the root of a response for a GET /vulns
type GetVulnsResponse struct {
	Results map[string]VulnResponseWithErrors `json:"results"`
}

// VulnResponseWithErrors is the collection of vulnerabilities and/or
// any errors raised for a package sent to GET /vulns
type VulnResponseWithErrors struct {
	Vulns  []ExpandedVuln `json:"vulns"`
	Errors []string       `json:"errors"`
}

// notReady is called for a request where the db has not been initialized yet
func notReady(w http.ResponseWriter) {
	res, _ := json.MarshalIndent(
		ErrorResponse{
			Error: "The database has not finished initializing",
		},
		"",
		"  ",
	)
	http.Error(w, string(res), http.StatusServiceUnavailable)
	return
}

// badRequest is called for any error raised during the processing of a request
func badRequest(w http.ResponseWriter, err error) {
	res, _ := json.MarshalIndent(
		ErrorResponse{
			Error: err.Error(),
		},
		"",
		"  ",
	)
	http.Error(w, string(res), http.StatusBadRequest)
	return
}

// getVulns is the main callable for a GET /vulns request
func getVulns(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if !db.Ready {
		notReady(w)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var input []string
	err = json.Unmarshal(body, &input)
	if err != nil {
		badRequest(w, err)
		return
	}
	makeGetVulnsResponse(w, input)
	return
}

// makeGetVulnsResponse will create and write the response for a given input
// to GET /vulns
func makeGetVulnsResponse(w http.ResponseWriter, input []string) {
	res := GetVulnsResponse{
		Results: make(map[string]VulnResponseWithErrors, 0),
	}
	for _, x := range input {
		results := VulnResponseWithErrors{}
		expandedVulns, errors := getExpandedVulnsForPackage(x)
		if len(expandedVulns) == 0 {
			results.Vulns = []ExpandedVuln{}
		} else {
			results.Vulns = expandedVulns
		}
		if len(errors) == 0 {
			results.Errors = []string{}
		} else {
			results.Errors = errors
		}
		res.Results[x] = results
	}
	body, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Header().Add("Content", "application/json")
	w.Write(body)
	return
}

// getExpandedVulnsForPackage will get the requested information for a given
// rpm package string. First the list of ALAS strings is queried for the base
// package name.
//
// Then for each of those ALAS's, check if the package
// referenced therein for the same architecture, name, and release is newer then
// the one provided on the input. If it is newer, we consider the input one vulnerable.
// We then provide any vulnerabilties and/or errors raised to be added to the response.
func getExpandedVulnsForPackage(pkgStr string) (expandedVulns []ExpandedVuln, errors []string) {
	errors = make([]string, 0)
	pkg, err := NewPackageFromString(pkgStr)
	if err != nil {
		errors = append(errors, err.Error())
		return
	}
	log.Println("Looking up vulnerabilities for:", pkg.Raw)
	vulns, err := db.GetVulnsByPackage(pkg.Name)
	if err != nil {
		errors = append(errors, err.Error())
		return
	}
	pkgVersion, err := version.NewVersion(pkg.Version)
	if err != nil {
		errors = append(errors, err.Error())
		return
	}
	expandedVulns = make([]ExpandedVuln, 0)
	for _, alas := range vulns {
		expanded, err := db.GetALAS(alas)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		for _, newPkg := range expanded.NewPackages {

			// If we are comparing the exact same package
			if newPkg.Arch == pkg.Arch && newPkg.Name == pkg.Name && sameRelease(pkg, newPkg) {

				// Epoch always wins
				if newPkg.Epoch > pkg.Epoch {
					log.Println(pkg.Raw, "is older than", newPkg.Raw, "by epoch constraint")
					expandedVulns = append(expandedVulns, stripNonRelated(pkg, expanded))
					break
				}

				newPkgVersion, err := version.NewVersion(newPkg.Version)
				if err != nil {
					errors = append(errors, err.Error())
					continue
				}
				// Version wins if epoch is the same
				if pkgVersion.LessThan(newPkgVersion) {
					log.Println(pkg.Raw, "is older than", newPkg.Raw, "by version constraint")
					expandedVulns = append(expandedVulns, stripNonRelated(pkg, expanded))
					break
				}

				// Release wins after everything else
				// I'll probably switch to use this implementation for the others too
				if rpmutils.Vercmp(pkg.Release, newPkg.Release) < 0 {
					log.Println(pkg.Raw, "is older than", newPkg.Raw, "by release constraint")
					expandedVulns = append(expandedVulns, stripNonRelated(pkg, expanded))
					break
				}
			}
		}
	}
	return
}

// stripNonRelated takes a package being queried, and an expanded potential
// vulnerability, then strips non-related packages associated with the vuln
// from the response.
func stripNonRelated(pkg Package, expanded ExpandedVuln) (stripped ExpandedVuln) {
	strippedPackages := make([]Package, 0)
	for _, x := range expanded.NewPackages {
		if pkg.Name == x.Name && pkg.Arch == x.Arch && sameRelease(pkg, x) {
			strippedPackages = append(strippedPackages, x)
		}
	}
	stripped = ExpandedVuln{
		ALAS:        expanded.ALAS,
		CVEs:        expanded.CVEs,
		Packages:    expanded.Packages,
		NewPackages: strippedPackages,
		Priority:    expanded.Priority,
		Link:        expanded.Link,
		PubDate:     expanded.PubDate,
	}
	return
}
