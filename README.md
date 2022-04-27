# Gobinsec

This tool parses Go binary dependencies and calls [NVD database](https://nvd.nist.gov/) to produce a vulnerability report. Binaries must have been built with module support to be analyzed with Gobinsec.

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Configuration](#configuration)
4. [Cache](#cache)
    - [Memcachier](#memcachier)
    - [Memcached](#memcached)
    - [File](#file)
5. [Timeout and Expiration](#timeout-and-expiration)
6. [Versions](#versions)
7. [How to Fix Vulnerabilities](#how-to-fix-vulnerabilities)
8. [Information about vulnerabilities](#information-about-vulnerabilities)
9. [How Gobinsec works](#how-gobinsec-works)
10. [License](#license)

## Installation

To install Gobinsec, type following command:

```
$ sh -c "$(curl -L https://github.com/intercloud/gobinsec/releases/latest/download/install)"
```

Or (if you don't have *curl* installed):

```
$ sh -c "$(wget -O - https://github.com/intercloud/gobinsec/releases/latest/download/install)"
```

Alternatively, if Go is already installed on your machine, you can also build and install latest version in your *GOPATH* with:

```
$ go install github.com/intercloud/gobinsec@latest
```

Finally, you can download binary for your platform in [latest release](https://github.com/intercloud/gobinsec/releases). Rename it *gobinsec*, make it executable with `chmod +x gobinsec` and move it somewhere in your *PATH*.

## Usage

To analyze given binary:

```yaml
$ gobinsec path/to/binary
binary: VULNERABLE
dependencies:
- name:    'golang.org/x/text'
  version: 'v0.3.0'
  vulnerable: true
  vulnerabilities:
  - id: 'CVE-2020-14040'
    exposed: true
    ignored: false
    references:
    - 'https://groups.google.com/forum/#!topic/golang-announce/bXVeAmGOqz0'
    - 'https://lists.fedoraproject.org/archives/list/package-announce@lists.fedoraproject.org/message/TACQFZDPA7AUR6TRZBCX2RGRFSDYLI7O/'
    matchs:
    - 'v < 0.3.3'
    - '?'
```

You can pass more than one binary to check on command line.

Exit code is *1* if exposed vulnerabilities were found, *2* if there was an error analyzing a binary and *0* otherwise. If a binary is vulnerable, exposed vulnerabilities are printed in report.

You can pass *-verbose* option on command line to print vulnerability report, even if binary is not vulnerable and for all vulnerabilities, even if they are ignored or not exposed.

To print cache information, pass *-cache* on command line. This will print dependencies along with following symbols:

- **>>>** vulnerabilities sent in cache for given dependency
- **<<<** vulnerabilities retrieved from cache for given dependency
- **!!!** vulnerabilities missed in cache for given dependency

You can set *-strict* flag on command line so that vulnerabilities without version are considered matching dependency version. In this case, you should check vulnerability manually and disable it in configuration file if necessary.

You can pass configuration file with *-config config.yml*, see configuration section below.

## Configuration

You can pass configuration on command line with `-config` option:

```
$ gobinsec -config config.yml path/to/binary
```

Configuration file is in YAML format as follows:

```yaml
api-key: "28c6112c-a7bc-4a4e-9b14-75be6da02211"
strict: false
memcachier:
  address:    "mcx.cy.eu-central-1.ec2.memcachier.com:11211"
  username:   "username"
  password:   "password"
  timeout:    "1s"
  expiration: "24h"
memcached:
  address:    "127.0.0.1:11211"
  timeout:    "1s"
  expiration: "24h"
file:
  name:       "~/.gobinsec-cache.yml"
  expiration: "24h"
ignore:
- "CVE-2020-14040"
```

Configuration fields are the following:

- **api-key**: this is your NVD API key
- **strict**: tells if we should consider vulnerability matches without version as matching dependency
- **memcachier** is the configuration for *memcachier*, see below
- **memcached** is the configuration for *memcached*, see below
- **file** is the configuration for *file* cache, see below
- **ignore**: a list of CVE vulnerabilities to ignore

You can also set NVD API Key in your environment with variable *NVD_API_KEY*. This key may be overwritten with value in configuration file. Your API key must be set in environment to be able to run integration tests (with target *integ*).

Note that without API key, you will be limited to *10* requests in a rolling *60* second window while this limit is *100* with an API key.

## Cache

A cache is useful because if you perform more call to NVD database than allowed, your calls will significantly slow down. Gobinsec tries to build caches in this order:

### Memcachier

A cache is built with *Memcachier* if following section is found in configuration file:

```yaml
memcachier:
  address:    ...
  username:   ...
  password:   ...
  timeout:    ...
  expiration: ...
```

Else, il will look for following environment variables:

```
MEMCACHIER_ADDRESS
MEMCACHIER_USERNAME
MEMCACHIER_PASSWORD
MEMCACHIER_TIMEOUT
MEMCACHIER_EXPIRATION
```

[Memcachier](https://www.memcachier.com) is an online cache provider with free tiers.

*Timeout* and *Expiration* configuration entries are optional and their default values are *1s* and *24h*.

### Memcached

A cache is built with *Memcached* if following section is found in configuration file:

```yaml
memcached:
  address:    ...
  timeout:    ...
  expiration: ...
```

Else it will look for following environment variables:

```
MEMCACHED_ADDRESS
MEMCACHED_TIMEOUT
MEMCACHED_EXPIRATION
```

*Timeout* and *Expiration* configuration entries are optional and their default values are *1s* and *24h*.

A sample [docker-compose.yml](https://github.com/intercloud/gobinsec/blob/main/docker-compose.yml) file to start a *memcached* instance is provided in this project.

### File

A file cache is used is none of preceding options is configured. By default, database file in YAML format is stored in *~/.gobinsec-cache.yml* and cache duration is of one day (or *24h*). You can overwrite these default values with following configuration section:

```yaml
file:
  name:       "~/.gobinsec-cache.yml"
  expiration: "24h"
```

Or with these environment variables:

```
FILECACHE_FILE
FILECACHE_EXPIRATION
```

*Expiration* configuration entry is optional and its default value is *24h*.

## Timeout and Expiration

These configuration keys are durations with a time unit. Possible units are:

- **ns** for nanosecond
- **us** or **µs** for microsecond
- **ms** for millisecond
- **s** for second
- **m** for minute
- **h** for hour

Thus, you could write, for instance:

- *1s* to set timeout to *1* second
- *2h30m* to set expiration to *2* hours and *30* minutes

## Versions

Dependencies and vulnerabilities have versions. There are three types of them:

- **Tag**: which implements semantic versioning, such as `1.2.3` or `v1.2.3`
- **Commit**: such as `v0.0.0-20210410081132-afb366fc7cd1` which is made of three parts: a major version, build date and commit ID
- **Date**: in ISO format such as `2022-01-07`

A dependency may have a tag or commit version. In the first case, developer used a released version of this dependency, in the last he's using a particular commit that wasn't released.

A vulnerability might have version conditions on tag or date. For instance:

- `v < 2017-03-17` means that vulnerability will affect dependencies before *2017-03-17*.
- `1.16.0 <= v <= 1.16.4` means that vulnerability will affect dependencies from version *1.16.0* to *1.16.4*, included

Given vulnerability is exposed if dependency version is in the range of affected versions. Thus to determine if a dependency is affected we must be able to compare versions between dependency and vulnerability.

- This is possible if dependency has a commit version and vulnerability a date one
- This is possible if dependency has a tag version and vulnerability a tag one
- This is not possible if dependency has a tag version and vulnerability a date one

In this later case, the vulnerability is considered exposed. You should check manually if release date of the dependency is in the date range of the vulnerability or not. You can then ignore vulnerability adding its ID in the configuration *ignore* list.

Sometimes, vulnerabilities have no version or date range. This is the case when vulnerability affects a given software (a Linux distribution for instance). In this case, vulnerability condition appears as a question mark and we consider that dependency is not affected. You can change this behavior passing `-strict` option on command line or in configuration. In this case you will have to check manually and ignore such vulnerabilities.

## How to Fix Vulnerabilities

In some cases, you can fix a vulnerability by using latest dependency version. Let's say dependency *golang.org/x/crypto* in version *v0.0.0-20200622213623-75b288015ac9* is affected by vulnerability [CVE-2020-29652](https://nvd.nist.gov/vuln/detail/CVE-2020-29652).

We can see in vulnerability description that it was fixed after version *0.0.0-20201203163018-be400aefbc4c*. Thus latest version *v0.0.0-20220128200615-198e4374d7ed* will fix the issue. We can update this dependency to latest version with following commands:

```
$ go get -u golang.org/x/crypto
$ go mod tidy
```

Of course this is possible only if a fix was written and committed to fix the issue.

## Information about vulnerabilities

The best way to receive security announcements is to subscribe to the [golang-announce mailing list](https://groups.google.com/g/golang-announce). Any messages pertaining to a security issue will be prefixed with `[security]`. See the page about the [Go Security Policy](https://go.dev/security) for details about the process of vulnerability management.

Here is a list of sites where you can find information about vulnerabilities:

- [National Vulnerability Database](https://nvd.nist.gov/) lists vulnerabilities and provides an API to search them. See hereafter for details about querying the API.
- [https://cve.mitre.org/](https://cve.mitre.org/) provides a page to search CVE at <https://cve.mitre.org/cve/search_cve_list.html>. Note that this site is currently moving to <https://www.cve.org/>.
- [CVE Details](https://www.cvedetails.com/) also lists CVE vulnerabilities and hosts a page dedicated to Go at <https://www.cvedetails.com/vulnerability-list/vendor_id-14185/Golang.html>.
- Other sources for CVE vulnerabilities: <https://www.circl.lu/services/cve-search/> and <https://www.cve-search.org/>.

## How Gobinsec works

This tool first lists dependencies embedded in binary using [buildinfo package](https://pkg.go.dev/debug/buildinfo).

Then, it calls [National Vulnerability Database](https://nvd.nist.gov/) to lists known vulnerabilities for embedded dependencies. You can find documentation on its API at <https://nvd.nist.gov/developers/vulnerabilities> and get an API key here: <https://nvd.nist.gov/developers/request-an-api-key>.

For instance, to get vulnerabilities for library *golang.org/x/text*, we would call <https://services.nvd.nist.gov/rest/json/cves/1.0/?keyword=golang.org/x/text>, which returns following JSON payload:

```json
{
    "resultsPerPage": 1,
    "startIndex": 0,
    "totalResults": 1,
    "result": {
        "CVE_data_type": "CVE",
        "CVE_data_format": "MITRE",
        "CVE_data_version": "4.0",
        "CVE_data_timestamp": "2021-12-07T15:40Z",
        "CVE_Items": [
            {
                "cve": {
                    "data_type": "CVE",
                    "data_format": "MITRE",
                    "data_version": "4.0",
                    "CVE_data_meta": {
                        "ID": "CVE-2020-14040",
                        "ASSIGNER": "cve@mitre.org"
                    },
                    "problemtype": {
                        ...
                    },
                    "references": {
                        "reference_data": [
                            {
                                "url": "https://groups.google.com/forum/#!topic/golang-announce/bXVeAmGOqz0",
                                "name": "https://groups.google.com/forum/#!topic/golang-announce/bXVeAmGOqz0",
                                "refsource": "MISC",
                                "tags": [
                                    "Third Party Advisory"
                                ]
                            },
                            ...
                        ]
                    },
                    "description": {
                        "description_data": [
                            {
                                "lang": "en",
                                "value": "..."
                            }
                        ]
                    }
                },
                "configurations": {
                    "CVE_data_version": "4.0",
                    "nodes": [
                        {
                            "operator": "OR",
                            "children": [],
                            "cpe_match": [
                                {
                                    "vulnerable": true,
                                    "cpe23Uri": "cpe:2.3:a:golang:text:*:*:*:*:*:*:*:*",
                                    "versionEndExcluding": "0.3.3",
                                    "cpe_name": []
                                }
                            ]
                        },
                        ...
                    ]
                },
                "impact": {
                    "baseMetricV3": {
                        ...
                    },
                    "baseMetricV2": {
                        ...
                    }
                },
                "publishedDate": "2020-06-17T20:15Z",
                "lastModifiedDate": "2020-11-18T14:44Z"
            }
        ]
    }
}
```

This data is parsed to produce YAML report.

# License

This software is release under the [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html).

*Enjoy!*
