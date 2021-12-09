package main

type Response struct {
	ResultsPerPage int    `json:"resultsPerPage"`
	StartIndex     int    `json:"startIndex"`
	TotalResults   int    `json:"totalResults"`
	Result         Result `json:"result"`
}

type Result struct {
	Items []Item `json:"CVE_Items"`
}

type Item struct {
	CVE            CVE            `json:"cve"`
	Configurations Configurations `json:"configurations"`
}

type CVE struct {
	Metadata   Metadata  `json:"CVE_data_meta"`
	References Reference `json:"references"`
}

type Reference struct {
	ReferenceData []ReferenceData `json:"reference_data"`
}

type ReferenceData struct {
	URL string `json:"url"`
}

type Metadata struct {
	ID string `json:"ID"`
}

type Configurations struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Operator string  `json:"operator"`
	Match    []Match `json:"cpe_match"`
}

type Match struct {
	Vulnerable          bool   `json:"vulnerable"`
	VersionEndExcluding string `json:"versionEndExcluding"`
	VersionEndIncluding string `json:"versionEndIncluding"`
}
