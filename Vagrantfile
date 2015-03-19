Vagrant.configure(2) do |config|

  config.vm.box = "williamyeh/ubuntu-trusty64-docker"

  # for web frontend
  config.vm.network "forwarded_port", guest: 80, host: 8000

  # for Fluentd in_forward plugin
  config.vm.network "forwarded_port", guest: 24224, host: 24224

  # for Ganglia's Monitoring Daemon (gmond)
  config.vm.network "forwarded_port", guest: 8649, host: 8649

end
