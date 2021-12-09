# Gobinsec

This tool parses Go binary dependencies, calls NVD database to produce a vulnerability report for this binary.

## Usage

To analyze given binary:

```yaml
$ gobinsec path/to/binary
binary: 'path/to/binary'
vulnerable: true
dependencies:
- name:    'golang.org/x/text'
  version: '0.3.0'
  vulnerable: true
  vulnerabilities:
  - ID: 'CVE-2020-14040'
    exposed: true
    references:
    - 'https://groups.google.com/forum/#!topic/golang-announce/bXVeAmGOqz0'
    - 'https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/TACQFZDPA7AUR6TRZBCX2RGRFSDYLI7O/'
    versions:
    - <: '0.3.3'
```

Exit code is *1* if binary is vulnerable and *2* if there was an error analyzing binary and *0* otherwise. If binary is vulnerable, exposed vulnerabilities are printed in report.

You can pass *-verbose* option on command line to print vulnerability report, even if binary is not vulnerable and for all vulnerabilities, even if they are not exposed.

## Data source

This tool calls [National Vulnerability Database](https://nvd.nist.gov/) that lists known vulnerabilities. You can find documentation on its API at <https://nvd.nist.gov/developers/vulnerabilities> and get an API key here: <https://nvd.nist.gov/developers/request-an-api-key>.

*Enjoy!*
