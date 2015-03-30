gmond2fluentd
============

This program extracts metrics from Ganglia Monitoring Daemon (gmond) to Fluentd.



## Usage


```
Usage:
  gmond2fluentd file <json_file>  [options]  
  gmond2fluentd tcp  [options]  
  gmond2fluentd --help
  gmond2fluentd --version

Options:

  -s <gmond>, --src <gmond>         gmond source host:port
                                      [default: 127.0.0.1:8649].

  -d <fluentd>, --dest <fluentd>    fluentd in_forward TCP host:port
                                      [default: 127.0.0.1:24224].

  -t <tag>, --tag <tag>             tag sending to Fluentd's in_forward plugin
                                      [default: ganglia].

  -p <seconds>, --period <seconds>  interval of metric query [default: 60]

  --stdout                          also dump to stdout.
```



## Build

Build the executable for your platform (before compiling, please make sure that you have [Go](https://golang.org/) compiler installed):

```
$ go install github.com/docopt/docopt-go
$ go build
```

Or, build the *linux-amd64* executables for Docker:

```
$ ./build.sh
```

It will place the `gmond2fluentd_linux-amd64` and `gmond2fluentd_static_linux-amd64` executables into the `docker` directory.

## Demo

Please go to the the `docker` directory for more details.


### History

- 0.3 - Static Go binary + DASH (without GPLv2 parts).
- 0.2 - Static Go binary + pure busybox.
- 0.1 - Initial release. 


## License

Apache License V2.0.  See [LICENSE](LICENSE) file for details.