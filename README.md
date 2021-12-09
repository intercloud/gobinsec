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
  - id: 'CVE-2020-14040'
    exposed: true
    ignored: false
    references:
    - 'https://groups.google.com/forum/#!topic/golang-announce/bXVeAmGOqz0'
    - 'https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/TACQFZDPA7AUR6TRZBCX2RGRFSDYLI7O/'
    versions:
    - <: '0.3.3'
```

Exit code is *1* if binary is vulnerable and *2* if there was an error analyzing binary and *0* otherwise. If binary is vulnerable, exposed vulnerabilities are printed in report.

You can pass *-verbose* option on command line to print vulnerability report, even if binary is not vulnerable and for all vulnerabilities, even if they are not exposed.

## Configuration

You can pass configuration on command line with `-config` option:

```
$ gobinsec -config config.yml path/to/binary
```

Configuration file is in YAML format as follows:

```yaml
api-key: "28c6112c-a7bc-4a4e-9b14-75be6da02211"
ignore:
- "CVE-2020-14040"
```

It has two entries:

- **api-key**: this is your NVD API key
- **ignore**: a list of CVE vulnerabilities to ignore

Note that without API key,you will be limited to *10* requests in a rolling *60* second window; the rate limit with an API key is *100* requests in a rolling *60* second window.

## Data source

This tool calls [National Vulnerability Database](https://nvd.nist.gov/) that lists known vulnerabilities. You can find documentation on its API at <https://nvd.nist.gov/developers/vulnerabilities> and get an API key here: <https://nvd.nist.gov/developers/request-an-api-key>.

*Enjoy!*
