gmond2fluentd
============

Repository name in Docker Hub: **[williamyeh/gmond2fluentd](https://registry.hub.docker.com/u/williamyeh/gmond2fluentd/)**

This program extracts metrics from Ganglia Monitoring Daemon (gmond) to Fluentd.



## Usage

Same as native executable, excluding `-s` and `-d` options since they are replaced by Docker's container linking mechanism.

```
Usage:

  docker run williamyeh/gmond2fluentd  \
      --link ganglia:ganglia  \
      --link fluentd:fluentd  \
      tcp  [options]

  docker run williamyeh/gmond2fluentd  \
      --link ganglia:ganglia  \
      --volumes-from fluentd  \
      file <json_file>  [options]

  docker run williamyeh/gmond2fluentd  --help

  docker run williamyeh/gmond2fluentd  --version


Options:

  -t <tag>, --tag <tag>             tag sending to Fluentd's in_forward plugin
                                      [default: ganglia].

  -p <seconds>, --period <seconds>  interval of metric query [default: 60]

  --stdout                          also dump to stdout.
```



## Demo


1. Send to Fluentd's [in_forward](http://docs.fluentd.org/articles/in_forward) (TCP port 24224):

   ```bash
   $ docker-compose  up  -d

   $ docker logs -f docker_demotcp_1
   ```


2. Send to Fluentd's [in_tail](http://docs.fluentd.org/articles/in_forward) (plaintext file)

   ```bash
   $ docker-compose  -f docker-compose-file.yml  up  -d

   $ docker logs -f docker_demofile_1
   ```


### Dependencies

- [progrium/busybox](https://registry.hub.docker.com/u/progrium/busybox/).


### History

- 0.1 - Initial release. 



## License

Apache License V2.0.  See [LICENSE](../LICENSE) file for details.