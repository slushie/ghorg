# ghorg

Github Organization Statistics Tool

## Introduction

This tool provides basic statistics about repositories within a Github organization. 

```bash
$ ghorg stars ruby
Stars  Name               URL
-----  ----               ---
14128  ruby               https://github.com/ruby/ruby
1231   rake               https://github.com/ruby/rake
664    www.ruby-lang.org  https://github.com/ruby/www.ruby-lang.org
459    rdoc               https://github.com/ruby/rdoc
387    psych              https://github.com/ruby/psych
```

## Installation

If you are a Go developer or if you have the `go` tool configured, 
you can install ghorg using `go get`:

```bash
$ go get github.com/slushie/ghorg
``` 

Alternatively, you can install a [binary release from Github](https://github.com/slushie/ghorg/releases)
into your system's `$PATH`.

## Usage

    
     e88~~\  888   |                       / 
    d888     888___|  e88~-_  888-~\ e88~88e 
    8888 __  888   | d888   i 888    888 888 
    8888   | 888   | 8888   | 888    "88_88" 
    Y888   | 888   | Y888   ' 888     /      
     "88__/  888   |  "88_-~  888    Cb      
                                      Y8""8D 
    
    This tool shows basic statistics for your Github organization.
    
    Usage:
      ghorg [command]
    
    Available Commands:
      contrib     List repos by PRs
      forks       List repos by forks
      help        Help about any command
      pulls       List repos by PRs
      stars       List repos by stargazers
    
    Flags:
      -T, --access-token string   Github OAuth2 access token used to authenticate REST calls.
          --config string         Path to ghorg config file
      -h, --help                  help for ghorg
      -N, --organization string   Organization name
    
    Use "ghorg [command] --help" for more information about a command.


Note that Github enforces heavy rate limiting for unauthenticated API access. To avoid
errors related to rate limiting, be sure to set a Github access token either via the 
`ACCESS_TOKEN` environment variable, or by using the `--access-token` option on the
command line.  

NB: This tool will warn you when no access token has been specified. You can obtain a 
[personal access token](https://github.com/settings/tokens) from the Github web UI.

## Development

To hack on this code, clone the repo (or use `go get`) and build by running:

```bash
$ make depends
$ make
```

To run the automated test suite, run:

```bash
$ make dev-depends
$ make test
``` 

This will watch for code changes, run unit tests, and report results to your browser. See 
[GoConvey](https://github.com/smartystreets/goconvey) for more details on adding unit tests.

## Roadmap

This is an MVP release of the `ghorg` tool. In future iterations, expect to see 
improvements such as:

* Persistent HTTP caching via the local filesystem
* Progress bars during long running operations
* Configurable timeout values
* Integrated usage analytics
* OAuth2 client authentication (aka, the three legged OAuth flow)
* Configurable column output
* Colors!

## Author

Josh Leder <josh@ha.cr>