<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [ubuntu-autoinstall-webhook](#ubuntu-autoinstall-webhook)
  - [Technical Requirements](#technical-requirements)
  - [Microservices / Commands that are part of the binary](#microservices--commands-that-are-part-of-the-binary)
  - [Client](#client)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# ubuntu-autoinstall-webhook

We are going to design a new microservices application. This document is a rough outline that discuss the overall technical requirements and provides examples of files.

## Technical Requirements

1. This will be a single golang binary but have multiple commands.
2. The server application will be several microservices.
3. All components will communicate using ProtoBufs and gRPC.
4. It will use the latest version and revision of ProtoBufs.
5. It will use the latest version and revision of gRPC.
6. The binary will use viper and cobra libraries.
7. All functions should use interfaces wherever logical and possible to facilitate accurate testing in a restricted environment.
8. All tests should not read or write to the actual filesystem or database, they should go to mocks and/or fake file systems to ensure that nothing is broken during testing.
9. This project should be organized using modern golang convention.
10. It will have an appropriate .gitignore that ignores angular / javascript, golang, vscode files, and mac os specific files.
11. All paths, options, and other areas should be configurable via the config.yaml.
12. The microservices should check upon starting if they're running in a kubernetes environment and make whatever changes are necessary to run there.
13. The web front end should use RBAC for allowing access to it's components.
    1. Roles can contain other roles.
    2. e.g. A user must have the logging role to view the logs. A user must have the ide role to use the web based ide, etc.
14. Clients using grpc can use mutual tls, pre-shared secret, ip matching, mac address matching to authenticate to the grpc server.
    1. Pre-shared secrets are written to the configuration file in the format: base64(user:password:role). The role will be client.
    2. Mutual tls is supported, and if setup in the client configuration it will automatically generate a CSR
15. All microservices with the exception of the install client, and possibly the database microservice should support having multiple replicas. If needed they can utilize some method to acquire a lock for the service that all replicas can see.
    1. The database microservice will support multiple instances when configured to use a cockroachdb backend. If using sqlite3, it must be a single instance.
    2. If running in kubernetes it should utilize a configmap or a lease <https://kubernetes.io/docs/concepts/architecture/leases/> if supported.
16. All microservices should be fully instrumented using opentelemetry, providing tracing, metrics, logs over http and grpc, and if desired to the console.
17. There should be an admin command that has several subcommands to manage everything from the commandline. An admin should be able to edit configs, create new machines, anything you would be able to do in the web interface should be available under the admin command.

## Microservices / Commands that are part of the binary

1. `file-editor` microservice that will validate, open, read and write changes to the filesystem
   1. There must be only one leader for this server so we don't get file conflicts. Use a database or kubernetes based lock/lease to ensure only one leader is processing data at a time. Multiple instances are not supported when using sqlite3.
   2. It will write to the `/var/www/html/ipxe/boot/` folder, it will write files that have the following naming template `/var/www/html/ipxe/boot/mac-{MAC_ADDRESS}.ipxe`.
   3. It will also manage the folders and symlinks under `/var/www/html/cloud-init/`.
      1. It will create two folders on discovery of a new system pxe booting, following the following templates  `/var/www/html/cloud-init/{MAC_ADDRESS}/`, and `/var/www/html/cloud-init/{MAC_ADDRESS}_install/`.
      2. It will create a symlink from `/var/www/html/cloud-init/{HOSTNAME}/` to `/var/www/html/cloud-init/{MAC_ADDRESS}/` using a provided hostname.
      3. It will create a symlink from `/var/www/html/cloud-init/{HOSTNAME}_install/` to `/var/www/html/cloud-init/{MAC_ADDRESS}_install/` using a provided hostname.
      4. It will verify that we aren't trying to assign the same hostname to a different {MAC_ADDRESS} folder. If this happens it should report an error with duplicate hostnames.
   4. It should provide full validation for the ipxe files, before writing them.
   5. It should provide full validation for all of the cloud-init auto_install files, use the cloud-init libraries if necessary.
   6. In the `/var/www/html/cloud-init/{HOSTNAME}_install` and `/var/www/html/cloud-init/{HOSTNAME}` folders they should have the following files at a minimum:
      1. `meta-data`
      2. `network-config`
      3. `user-data`
2. `database` microservice. This will perform all communication with the database. It will read/write/etc to the DB. All database functions should be contained in this microservice.
   1. The database will support 2 backends, by default it will use sqlite3. If specified in the database type option in the config, it will use cockroachdb.
   2. It will process any schema updates, and ensure indexes are correctly created regardless of the database backend.
   3. The user will be able to switch between the two database backends, and in that case any data stored in one should be replicated to the other. I.E. if you start using sqlite3 and move to cockroachdb, it should migrate the data from the sqlite3 database to cockroachdb.
   4. It will provide a way to purge old data.
   5. It will use a UUID field as the primary key for all tables. If sqlite3 doesn't support UUIDs then it will generate an RFC compliant UUID and write it either as a string or a blob whichever would be faster.
3. `configuration` microservice. This will generate, validate, and send the configuration changes to the `database` microservice, and the `file-editor` microservice.
   1. It will cache the configs for faster retrieval on first use and write the changes back to the database using the database microservice.
   2. It should support multiple instances as it should be stateless, only caching data and sending it to the file writer and database services respectively.
   3. It will create configs for newly discovered systems using several templates and send them to the file writer/database microservices:
      1. `meta-data` which will contain a randomly generated `instance-id` prefixed with `jf-`, and `local-hostname` which will contain the hostname. Any other important or useful fields should be included as well.
      2. `network-config` which should be empty as it conflicts with the `user-data` file's `networking-config` section.
      3. `user-data` which contains the user-data that the target system loads via cloud-init.
      4. `mac-{MAC_ADDRESS}.ipxe` which will contain the current defaults for that mac address. This will be changed from the autoinstall boot to setting it to use the boot from hard drive menu option. An example booting autoinstall:

      ```ipxe
        #!ipxe
        echo
        echo Booting len-serv-003 ubuntu autoinstall
        set menu-default auto-install-ubuntu
        set jfhostname len-serv-003
        chain --replace --autofree ${menu-url}
      ```

      5. `variables.sh` which will need to be adjusted per system, particularly HOSTNAME, NET_ET_ADDRESS, and WIFI_ADDRESS. Here is an example:

        ```bash
        #!/bin/bash
        # Variables file for jinstall.sh
        # Adjust these values per system

        # General system settings
        DISK="/dev/nvme0n1"
        TIMEZONE="America/New_York"
        DEBOOTSTRAP_RELEASE="oracular"
        HOSTNAME="len-serv-003"
        REPORT_STATUS="/usr/local/bin/report-status.sh"

        # Security & Authentication
        LUKS_KEY="defaultLUKSkey123"
        ROOT_PASSWORD="defaultPassword123"

        # Network configuration (Ethernet)
        NET_ET_INTERFACE="enp1s0f0"
        NET_ET_ADDRESS="172.16.3.96/23"
        NET_ET_GATEWAY="172.16.2.1"
        NET_ET_SEARCH="jf.local"
        NET_ET_NAMESERVERS=("172.16.2.1" "1.1.1.1" "8.8.8.8")

        # Network configuration (WiFi)
        WIFI_INTERFACE="wlp2s0"
        WIFI_ADDRESS="172.16.3.97/23"
        WIFI_SSID="VentureIndustries"
        WIFI_PASSWORD="8harley4IVY2\$60SILVER"

        # CockroachDB configuration
        COCKROACH_ADVERTISE="172.16.2.55:26257"
        COCKROACH_JOIN="172.16.2.45:36257,172.16.2.45:36357,172.16.2.47:36257,172.16.2.47:36357,172.16.2.30:36257,172.16.2.30:36357"
        COCKROACH_CACHE=".25"
        COCKROACH_MAX_SQL=".25"
        COCKROACH_LOCALITY="region=us,cluster-unit=lenovo"

        # Tang server URLs for Clevis binding
        TANG_URL1="http://172.16.2.45"
        TANG_URL2="http://172.16.2.46"
        TANG_URL3="http://172.16.2.47"

        # Sleep durations (in seconds)
        SLEEP_DURATION=3
        ZED_SLEEP_DURATION=30
        ```

   4. It'll retrieve and serve configs to the webserver microservice. Especially when using the web based IDE so that when changes are saved, they are validated and correctly written to the filesystem and database.
   5. It will also read the application config.yaml and watch it for updates.
      1. On startup it will read and store the configuration
      2. Other microservices will contact it for their initial configuration. They will keep retrying until they get a config or a 10 minute timeout is reached.
      3. On a configuration change it will push out updates to the relevant services so they can reload their configuration.
      4. It will validate the config.yaml.
      5. It will write a new version of `config.yaml.example` whenever a new option added, or an option changes / is updated. It will provide examples, and possible values in addition to clearly documenting every part of the configuration.
      6. The config file location should be configurable via viper/cobra.
4. `dnsmasq-watcher` this will watch dns-masq or query via some sort of API to catch MAC addresses and IP addresses so that they can be used to generate a minimal config.
   1. When a new mac address is detected that isn't already in the database, it will send a message to the `configuration` microservice noting the new system found, the mac address, ip address, and hostname if it's in the logs/messages. Otherwise it will create a randomly generated hostname based on the mac address.
   2. There must be only one leader for this server so we don't get duplicates. Use a database or kubernetes based lock/lease to ensure only one leader is processing data at a time.
   3. It should write a timestamp to the database for the logs processed, that way if it dies and the secondary takes over, it can read that timestamp and know exactly where to resume monitoring the logs.
   4. It should be able to fetch logs from the file system, journald, or kubernetes.
   5. The pxeclient boots twice, once to load the ipxe firmware, and once it's loaded that, it restarts the boot process as noted by the `vender class` and `user class`
   6. Here's an example of the logs in journald:

   ```log
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 vendor class: PXEClient:Arch:00007:UNDI:003016
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 PXE(enp8s0f0) 6c:4b:90:bc:f7:f4 proxy
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 tags: efi64, enp8s0f0
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 bootfile name: ipxe.efi
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 next server: 172.16.2.30
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 broadcast response
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 sent size:  1 option: 53 message-type  2
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 sent size:  4 option: 54 server-identifier  172.16.2.30
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 sent size:  9 option: 60 vendor-class  50:58:45:43:6c:69:65:6e:74
    Mar 06 11:57:52 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 sent size: 17 option: 97 client-machine-id  00:00:03:b3:c5:a2:fa:e9:11:84:54:60:69:2c...
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 3838357051 vendor class: PXEClient:Arch:00007:UNDI:003016
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 vendor class: PXEClient:Arch:00007:UNDI:003016
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 PXE(enp8s0f0) 6c:4b:90:bc:f7:f4 proxy
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 tags: efi64, enp8s0f0
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 bootfile name: ipxe.efi
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 next server: 172.16.2.30
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 sent size:  1 option: 53 message-type  5
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 sent size:  4 option: 54 server-identifier  172.16.2.30
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 sent size:  9 option: 60 vendor-class  50:58:45:43:6c:69:65:6e:74
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 sent size: 17 option: 97 client-machine-id  00:00:03:b3:c5:a2:fa:e9:11:84:54:60:69:2c...
    Mar 06 11:57:55 unimatrixzero dnsmasq-dhcp[9980]: 2578703609 sent size: 25 option: 43 vendor-encap  06:01:08:0a:13:02:42:6f:6f:74:69:6e:67:20...
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 vendor class: PXEClient
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 user class: iPXE
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 PXE(enp8s0f0) 6c:4b:90:bc:f7:f4 proxy
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 tags: ipxe, efi64, enp8s0f0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 bootfile name: http://172.16.2.30/ipxe/boot.ipxe
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 next server: 172.16.2.30
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 broadcast response
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 sent size:  1 option: 53 message-type  2
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 sent size:  4 option: 54 server-identifier  172.16.2.30
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 sent size:  9 option: 60 vendor-class  50:58:45:43:6c:69:65:6e:74
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 sent size: 17 option: 97 client-machine-id  00:00:03:b3:c5:a2:fa:e9:11:84:54:60:69:2c...
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 vendor class: PXEClient
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 3422164510 user class: iPXE
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 vendor class: PXEClient
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 user class: iPXE
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 PXE(enp8s0f0) 6c:4b:90:bc:f7:f4 proxy
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 tags: ipxe, efi64, enp8s0f0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 bootfile name: http://172.16.2.30/ipxe/boot.ipxe
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 next server: 172.16.2.30
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 broadcast response
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 sent size:  1 option: 53 message-type  2
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 sent size:  4 option: 54 server-identifier  172.16.2.30
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 sent size:  9 option: 60 vendor-class  50:58:45:43:6c:69:65:6e:74
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 sent size: 17 option: 97 client-machine-id  00:00:03:b3:c5:a2:fa:e9:11:84:54:60:69:2c...
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 vendor class: PXEClient
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 4218405448 user class: iPXE
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 vendor class: PXEClient
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 user class: iPXE
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 PXE(enp8s0f0) 6c:4b:90:bc:f7:f4 proxy
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 tags: ipxe, efi64, enp8s0f0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 bootfile name: http://172.16.2.30/ipxe/boot.ipxe
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 next server: 172.16.2.30
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 broadcast response
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 sent size:  1 option: 53 message-type  2
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 sent size:  4 option: 54 server-identifier  172.16.2.30
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 sent size:  9 option: 60 vendor-class  50:58:45:43:6c:69:65:6e:74
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 sent size: 17 option: 97 client-machine-id  00:00:03:b3:c5:a2:fa:e9:11:84:54:60:69:2c...
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 vendor class: PXEClient
    Mar 06 11:58:00 unimatrixzero dnsmasq-dhcp[9980]: 2274531367 user class: iPXE
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 vendor class: PXEClient
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 user class: iPXE
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 PXE(enp8s0f0) 6c:4b:90:bc:f7:f4 proxy
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 tags: ipxe, efi64, enp8s0f0
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 bootfile name: http://172.16.2.30/ipxe/boot.ipxe
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 next server: 172.16.2.30
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 broadcast response
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 sent size:  1 option: 53 message-type  2
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 sent size:  4 option: 54 server-identifier  172.16.2.30
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 sent size:  9 option: 60 vendor-class  50:58:45:43:6c:69:65:6e:74
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 sent size: 17 option: 97 client-machine-id  00:00:03:b3:c5:a2:fa:e9:11:84:54:60:69:2c...
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 available DHCP subnet: 172.16.2.0/255.255.254.0
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 vendor class: PXEClient
    Mar 06 11:58:05 unimatrixzero dnsmasq-dhcp[9980]: 83351049 user class: iPXE
    ```

5. `cert-issuer` this either communicates with cert-manager, or by using cfssl or boringssl to parse requests, generate certificates that include the full chain, and return a stored certificate.
   1. On initial startup unless the root/intermediate certificates are configured and available in the database, or available via cert-manager, etc, the server should generate a self-signed root-ca with a 20 year validity, and a intermediate ca with a 2 year validity.
      1. If running in kubernetes and using cert-manager, it should generate the appropriate kubernetes crds/objects for cert-manager to generate the self-signed certs, and the intermediate certs. It should also generate the certificate object when a system is approved. The database should also have a field that references these certificates are available via kubernetes and to not write the certificate data to the database.
      2. Initially all server components should be allowed to connect via gRPC using a pre-shared key. They will then request a mutual tls certificate and automatically be issued one. If running on a server and not in kubernetes, it should write that certificate to the config folder for the microservice. If running in kubernetes then cert-manager should've already generated a secret object to be mounted in the container.
   2. Certificates for clients should be short lived 24 hours at most, and when a node restarts and requests a certificate the previously used one should be offered unless it has less than 12 hours of validity left. Then it should be revoked, and a new one issued.
   3. Certificates for sever components should have a 3 month validity, and halfway through they should request a renewal.
   4. If running directly on a server, the config folder for each microservice should only contain a basic certs.yaml and the certificate files. The certs.yaml says use certs, and gives the name of the certs for the microservice to load.
6. `webserver` this will host the web front end and be the grpc endpoint for the install clients.
   1. The web server will return an angular web application which will be the web frontend for admins to interact with the applications.
   2. It will have api endpoints for the angular app. It will use  grpc-gateway <https://github.com/grpc-ecosystem/grpc-gateway> code to convert restful api requests to protobufs and send them to the appropriate microservices.
   3. It will also have a direct gRPC front end for the install clients to connect to.
   4. It will support operation in both HTTPS and HTTP, defaulting to insecure HTTP. This was if the communication is secured externally like by Istio, we won't double encrypt our data.
   5. It will support multiple concurrent instances, allowing for load balancing.
   6. It will process the incoming logs, status updates, hardware info, etc from the install clients, and send the requests to the appropriate microservice.
   7. It will have full support for editing the config file with error checking right in the web based IDE which will then be written to the disk. As always ENV values will override the config file, which is the default with viper.
   8. It will support authentication, through several methods:
      1. a static user:password:role mapping in the config file,
      2. a dex based oauth2/oidc authentication and role mapping,
      3. a users table in the database that can be edited by users with the admin role.
      4. The admin should be able to independently enable each different method without them conflicting with the others.
      5. The admin should be able to create additional roles and map users to roles all through the web interface.
   9. When a client initially connects to the grpc server for the first time, it sends it's identifying information regardless of what method of authentication.
      1. It should send Hostname, IP Addresses, Mac Addresses, smbios_uuid, etc. Any other fields that would uniquely identify this system, like a processor guid, or whatever else is available.
      2. This should be allowed without any authentication.
      3. The webserver will have interactive alerts that will pop up notifying the user that a new node has connected, and to approve or deny the node.
      4. The webserver will have a page listing all the nodes currently approved, denied, pending. And allow a user with the appropriate RBAC access to approve or deny a system. It should also let a user remove a system or move it to approve or deny.
      5. If approved, the identifying information will be stored in the database, on next boot if it has the same MAC address, and other unique characteristics it's automatically approved and a new certificate is issued.
      6. The webserver will use the certificate microservice to generate CSR's, approve them, and generate certificates that are pushed to the client.
7. Angular Components
    1. The webpage design should be clean and modern.
    2. It should use Angular v19+. Check that all functions, features, etc are validated against the latest syntax in v19.
    3. It should use Material Design 3, version 19+. Validate the design against Material Design best practices.
    4. It should have a functional layout allowing the user to quickly get to the information, and to drill down to find more information.
    5. It will have several pages:
       1. An initial welcome page which describes the other pages.
       2. A page for approving/denying connection requests/certificate requests.
       3. A page listing all systems that we have data for. Allowing us to easily click to see their logs, cloud-init, and configs. It should also allow you to quickly edit the configs/cloud-init by opening it in the web based IDE.
       4. Use ngx-monaco-editor-v2, the monaco editor for the web based IDE. We should be able to write new configs or edit the configs stored in the database or the config.yaml.
       5. A install status page showing the current status of any in progress install.
       6. Any other pages that seem useful.


## Client

The client will connect back to the server, request a certificate, switch to using mTLS, get a config for the node, then perform the steps that are currently being handled by these scripts:
jinstall.sh
```
#!/bin/bash
set -ex

# Trap handler: on exit, if the script exited with a nonzero code, upload diagnostic logs.
handle_exit() {
  exit_code=$?
  if [ $exit_code -ne 0 ]; then
    log "ERROR" "Script terminated abnormally with exit code $exit_code. Uploading diagnostic logs."
    # Attempt to upload the target OS logs if available.
    if [ -f /mnt/targetos/var/log/cloud-init.log ]; then
      upload_logs "/mnt/targetos/var/log/cloud-init.log"
    else
      log "WARNING" "Target OS cloud-init.log not found."
    fi
    if [ -f /mnt/targetos/var/log/cloud-init-output.log ]; then
      upload_logs "/mnt/targetos/var/log/cloud-init-output.log"
    else
      log "WARNING" "Target OS cloud-init-output.log not found."
    fi
  fi
}
trap handle_exit EXIT

# Source the variables file – adjust the path if needed
source ./variables.sh

# Source the new reporting functions.
source ./reporting.sh

#######################################
# Logging function: Logs messages with timestamp and level.
# Logs to stdout and via logger.
#######################################
log() {
  local level="$1"
  local msg="$2"
  local timestamp
  timestamp=$(date +'%Y-%m-%d %H:%M:%S')
  echo "$timestamp [$level] $msg"
  logger -t jinstall.sh "[$level] $msg"
}

#######################################
# Retry wrapper: Retries a function up to 3 times.
#######################################
retry() {
  local func="$1"
  local max_attempts=3
  local attempt=1
  local exit_code=0
  while [ $attempt -le $max_attempts ]; do
    log "INFO" "Attempt $attempt for $func"
    if $func; then
      return 0
    else
      exit_code=$?
      log "ERROR" "$func failed with exit code $exit_code"
      attempt=$((attempt+1))
      sleep 2
    fi
  done
  return $exit_code
}

#######################################
# (No local send_status_update function needed; use the one from reporting.sh)
#######################################

#######################################
# pre_install: Update packages and install initial tools.
#######################################
pre_install() {
  log "INFO" "Starting pre-install tasks"
  apt-get update || { log "ERROR" "apt-get update failed"; return 1; }
  apt-get -y install jq || { log "ERROR" "Installing jq failed"; return 1; }
  log "INFO" "Installing squid-deb-proxy-client (host)"
  apt install --yes squid-deb-proxy-client || { log "ERROR" "Installing squid-deb-proxy-client failed"; return 1; }
  send_status_update "installing" 0 "Starting cloud-init process"
  return 0
}

#######################################
# update_package_lists: Refresh APT lists.
#######################################
update_package_lists() {
  log "INFO" "Updating package lists"
  send_status_update "installing" 5 "Updating package lists"
  apt-get update || { log "ERROR" "apt-get update failed"; return 1; }
  return 0
}

#######################################
# install_required_packages: Install base utilities.
#######################################
install_required_packages() {
  log "INFO" "Installing required packages"
  send_status_update "installing" 10 "Installing required packages"
  apt-get -y install parted tftp-hpa debootstrap gdisk zfsutils-linux jq lshw || { log "ERROR" "Installing required packages failed"; return 1; }
  return 0
}

#######################################
# stop_services: Stop unnecessary services.
#######################################
stop_services() {
  log "INFO" "Stopping unnecessary services"
  send_status_update "installing" 15 "Stopping unnecessary services"
  systemctl stop zed || { log "ERROR" "Failed to stop zed service"; return 1; }
  return 0
}

#######################################
# configure_timezone_ntp: Set the timezone and enable NTP.
#######################################
configure_timezone_ntp() {
  log "INFO" "Configuring timezone and NTP"
  send_status_update "installing" 35 "Configuring timezone and NTP"
  timedatectl set-timezone "$TIMEZONE" || { log "ERROR" "Failed to set timezone"; return 1; }
  timedatectl set-ntp on || { log "ERROR" "Failed to enable NTP"; return 1; }
  systemctl restart systemd-timesyncd || { log "ERROR" "Failed to restart systemd-timesyncd"; return 1; }
  sleep "$SLEEP_DURATION"
  return 0
}

#######################################
# validate_time: Verify that the time settings are correct.
#######################################
validate_time() {
  log "INFO" "Validating time configuration"
  send_status_update "installing" 40 "Validating time configuration"
  echo "Current system time $(date)"
  echo "Timezone set to $(timedatectl show --property=Timezone --value)"
  echo "NTP synchronization status $(timedatectl show --property=NTP --value)"
  if [ "$(timedatectl show --property=Timezone --value)" != "$TIMEZONE" ]; then
    log "ERROR" "Timezone not set correctly!"
    return 1
  fi
  if [ "$(timedatectl show --property=NTP --value)" != "yes" ]; then
    log "ERROR" "NTP synchronization not enabled!"
    return 1
  fi
  return 0
}

#######################################
# partition_disk: Wipe the disk, destroy any existing zpools, and create partitions.
#######################################
partition_disk() {
  log "INFO" "Starting disk partitioning"
  send_status_update "installing" 45 "Starting disk partitioning"
  echo "Our disk is ${DISK}"
  wipefs -a "$DISK" || { log "ERROR" "wipefs failed"; return 1; }
  blkdiscard -f "$DISK" || { log "ERROR" "blkdiscard failed"; return 1; }
  sgdisk --zap-all "$DISK" || { log "ERROR" "sgdisk zap failed"; return 1; }
  # Destroy any existing zpools.
  if zpool list -H -o name 2>/dev/null | grep .; then
    for pool in $(zpool list -H -o name); do
      log "INFO" "Destroying existing zpool: $pool"
      zpool destroy "$pool" || log "WARNING" "Failed to destroy zpool: $pool"
    done
  fi
  echo "Creating GPT partition table..."
  parted -s "$DISK" mklabel gpt || { log "ERROR" "Creating GPT partition table failed"; return 1; }
  echo "Creating partitions..."
  send_status_update "installing" 50 "Creating partitions"
  parted -s "$DISK" mkpart ESP fat32 1MiB 513MiB || return 1
  parted -s "$DISK" set 1 boot on || return 1
  parted -s "$DISK" set 1 esp on || return 1
  parted -s "$DISK" mkpart RESET fat32 513MiB 4609MiB || return 1
  parted -s "$DISK" mkpart BPOOL 4609MiB 6657MiB || return 1
  parted -s "$DISK" mkpart LUKS 6657MiB 7681MiB || return 1
  parted -s "$DISK" mkpart RPOOL 7681MiB 100% || return 1
  return 0
}

#######################################
# format_partitions: Format ESP and RESET partitions.
#######################################
format_partitions() {
  log "INFO" "Formatting partitions"
  send_status_update "installing" 55 "Formatting partitions"
  mkfs.fat -F32 "${DISK}p1" || { log "ERROR" "Formatting partition p1 failed"; return 1; }
  mkfs.fat -F32 "${DISK}p2" || { log "ERROR" "Formatting partition p2 failed"; return 1; }
  return 0
}

#######################################
# configure_luks: Set up LUKS encryption and create XFS filesystem.
#######################################
configure_luks() {
  log "INFO" "Configuring LUKS encryption"
  send_status_update "installing" 60 "Configuring LUKS encryption"
  echo "$LUKS_KEY" | cryptsetup luksFormat --batch-mode "${DISK}p4" || { log "ERROR" "LUKS format failed"; return 1; }
  echo "$LUKS_KEY" | cryptsetup open "${DISK}p4" luks || { log "ERROR" "Opening LUKS failed"; return 1; }
  log "INFO" "Creating XFS filesystem on LUKS partition"
  mkfs.xfs -f -b size=4096 /dev/mapper/luks || { log "ERROR" "XFS creation failed"; return 1; }
  echo "Partitioning and formatting complete with optimal alignment!"
  return 0
}

#######################################
# mount_and_create_zfs: Mount LUKS and prepare for ZFS pools.
#######################################
mount_and_create_zfs() {
  log "INFO" "Mounting LUKS partition and preparing for ZFS pools"
  send_status_update "installing" 65 "Creating ZFS pools"
  mkdir -p /mnt/luks || { log "ERROR" "Failed to create /mnt/luks"; return 1; }
  mount /dev/mapper/luks /mnt/luks || { log "ERROR" "Mounting LUKS failed"; return 1; }
  dd if=/dev/random of=/mnt/luks/zfs.key bs=32 count=1 || { log "ERROR" "Generating ZFS key failed"; return 1; }
  chmod 600 /mnt/luks/zfs.key || { log "ERROR" "Setting ZFS key permissions failed"; return 1; }
  mkdir /mnt/targetos || { log "ERROR" "Creating /mnt/targetos failed"; return 1; }
  return 0
}

#######################################
# generate_uuid: Generate a unique identifier for dataset naming.
#######################################
generate_uuid() {
  export UUID=$(dd if=/dev/urandom bs=1 count=100 2>/dev/null | tr -dc 'a-z0-9' | cut -c-6)
  echo "UUID=${UUID}" > /mnt/targetos/uuid || { log "ERROR" "Writing UUID failed"; return 1; }
  echo "UUID is ${UUID}"
  echo "DISK=${DISK}" >> /mnt/targetos/uuid
  cat /mnt/targetos/uuid
}

#######################################
# create_bpool: Create the bpool zpool.
#######################################
create_bpool() {
  log "INFO" "Creating bpool"
  zpool create \
    -o ashift=12 \
    -o autotrim=on \
    -o cachefile=/etc/zfs/zpool.cache \
    -o compatibility=grub2 \
    -o feature@livelist=enabled \
    -o feature@zpool_checkpoint=enabled \
    -O devices=off \
    -O acltype=posixacl -O xattr=sa \
    -O compression=lz4 \
    -O normalization=formD \
    -O relatime=on \
    -O canmount=off -O mountpoint=/boot -R /mnt/targetos \
      bpool "${DISK}p3" || { log "ERROR" "Creating bpool failed"; return 1; }
}

#######################################
# create_rpool: Create the rpool zpool.
#######################################
create_rpool() {
  log "INFO" "Creating rpool"
  zpool create \
    -o ashift=12 \
    -o autotrim=on \
    -O encryption=on -O keylocation=file:///mnt/luks/zfs.key -O keyformat=raw \
    -O acltype=posixacl -O xattr=sa -O dnodesize=auto \
    -O compression=lz4 \
    -O normalization=formD \
    -O relatime=on \
    -O canmount=off -O mountpoint=/ -R /mnt/targetos \
    rpool "${DISK}p5" || { log "ERROR" "Creating rpool failed"; return 1; }
}

#######################################
# create_bpool_datasets: Create datasets for bpool.
#######################################
create_bpool_datasets() {
  log "INFO" "Creating bpool datasets"
  zfs create -o canmount=off -o mountpoint=none bpool/BOOT || { log "ERROR" "Creating bpool/BOOT failed"; return 1; }
  zfs create -o mountpoint=/boot bpool/BOOT/ubuntu_$UUID || { log "ERROR" "Creating bpool/BOOT/ubuntu_$UUID failed"; return 1; }
}

#######################################
# create_rpool_datasets: Create datasets for rpool.
#######################################
create_rpool_datasets() {
  log "INFO" "Creating rpool datasets"
  zfs create -o canmount=off -o mountpoint=none rpool/ROOT || { log "ERROR" "Creating rpool/ROOT failed"; return 1; }
  zfs create -o mountpoint=/ -o com.ubuntu.zsys:bootfs=yes -o com.ubuntu.zsys:last-used=$(date +%s) rpool/ROOT/ubuntu_$UUID || { log "ERROR" "Creating rpool/ROOT/ubuntu_$UUID failed"; return 1; }
  zfs create -o com.ubuntu.zsys:bootfs=no -o canmount=off rpool/ROOT/ubuntu_$UUID/usr || { log "ERROR" "Creating rpool/ROOT/ubuntu_$UUID/usr failed"; return 1; }
  zfs create -o com.ubuntu.zsys:bootfs=no -o canmount=off rpool/ROOT/ubuntu_$UUID/var || { log "ERROR" "Creating rpool/ROOT/ubuntu_$UUID/var failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib || { log "ERROR" "Creating var/lib failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/log || { log "ERROR" "Creating var/log failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/spool || { log "ERROR" "Creating var/spool failed"; return 1; }
  zfs create -o canmount=off -o mountpoint=/ rpool/USERDATA || { log "ERROR" "Creating rpool/USERDATA failed"; return 1; }
  zfs create -o com.ubuntu.zsys:bootfs-datasets=rpool/ROOT/ubuntu_$UUID -o canmount=on -o mountpoint=/root rpool/USERDATA/root_$UUID || { log "ERROR" "Creating root folder failed"; return 1; }
  chmod 700 /mnt/targetos/root || { log "ERROR" "Setting permissions on /mnt/targetos/root failed"; return 1; }
  log "INFO" "Creating additional rpool datasets in var"
  zfs create rpool/ROOT/ubuntu_$UUID/var/cache || { log "ERROR" "Creating var/cache failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib/nfs || { log "ERROR" "Creating var/lib/nfs failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/tmp || { log "ERROR" "Creating var/tmp failed"; return 1; }
  chmod 1777 /mnt/targetos/var/tmp || { log "ERROR" "Setting permissions on /mnt/targetos/var/tmp failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib/apt || { log "ERROR" "Creating var/lib/apt failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib/dpkg || { log "ERROR" "Creating var/lib/dpkg failed"; return 1; }
  zfs create -o com.ubuntu.zsys:bootfs=no rpool/ROOT/ubuntu_$UUID/srv || { log "ERROR" "Creating srv failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/usr/local || { log "ERROR" "Creating usr/local failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/games || { log "ERROR" "Creating var/games failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib/AccountsService || { log "ERROR" "Creating var/lib/AccountsService failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib/NetworkManager || { log "ERROR" "Creating var/lib/NetworkManager failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/lib/docker || { log "ERROR" "Creating var/lib/docker failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/mail || { log "ERROR" "Creating var/mail failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/snap || { log "ERROR" "Creating var/snap failed"; return 1; }
  zfs create rpool/ROOT/ubuntu_$UUID/var/www || { log "ERROR" "Creating var/www failed"; return 1; }
}

#######################################
# install_base_system: Bootstrap the base OS.
# Uses debootstrap with --include to preinstall required packages.
#######################################
install_base_system() {
  log "INFO" "Installing base system"
  send_status_update "installing" 80 "Installing base system"
  debootstrap --include=cryptsetup,dosfstools,openssh-server,vim,htop,curl,tmux,zsh "$DEBOOTSTRAP_RELEASE" /mnt/targetos || { log "ERROR" "debootstrap failed"; return 1; }
  log "INFO" "Copying zpool.cache"
  mkdir -p /mnt/targetos/etc/zfs || { log "ERROR" "Creating /mnt/targetos/etc/zfs failed"; return 1; }
  cp /etc/zfs/zpool.cache /mnt/targetos/etc/zfs/ || { log "ERROR" "Copying zpool.cache failed"; return 1; }
  return 0
}

#######################################
# configure_hostname_network: Set hostname and netplan configuration.
#######################################
configure_hostname_network() {
  log "INFO" "Configuring system hostname and network"
  send_status_update "installing" 85 "Configuring system hostname and network"
  hostnamectl set-hostname "$HOSTNAME" || { log "ERROR" "Setting hostname failed"; return 1; }
  echo "$HOSTNAME" > /mnt/targetos/etc/hostname || { log "ERROR" "Writing hostname failed"; return 1; }
  echo "127.0.1.1       $HOSTNAME" >> /mnt/targetos/etc/hosts || { log "ERROR" "Updating /etc/hosts failed"; return 1; }
  cat <<EOF > /mnt/targetos/etc/netplan/00-installer-config.yaml
network:
  version: 2
  renderer: systemd-networkd
  ethernets:
    ${NET_ET_INTERFACE}:
      addresses:
        - ${NET_ET_ADDRESS}
      routes:
        - to: default
          via: ${NET_ET_GATEWAY}
      nameservers:
        search:
          - ${NET_ET_SEARCH}
        addresses:
          - ${NET_ET_NAMESERVERS[0]}
          - ${NET_ET_NAMESERVERS[1]}
          - ${NET_ET_NAMESERVERS[2]}
  wifis:
    ${WIFI_INTERFACE}:
      activation-mode: off
      optional: true
      addresses:
        - ${WIFI_ADDRESS}
      access-points:
        "${WIFI_SSID}":
          password: ${WIFI_PASSWORD}
EOF
  return 0
}

#######################################
# configure_apt_sources: Write the APT sources list.
#######################################
configure_apt_sources() {
  log "INFO" "Configuring APT sources"
  cat <<EOF > /mnt/targetos/etc/apt/sources.list.d/ubuntu.sources
Types: deb
URIs: http://us.archive.ubuntu.com/ubuntu/
Suites: ${DEBOOTSTRAP_RELEASE} ${DEBOOTSTRAP_RELEASE}-updates ${DEBOOTSTRAP_RELEASE}-backports
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

Types: deb
URIs: http://security.ubuntu.com/ubuntu
Suites: ${DEBOOTSTRAP_RELEASE}-security
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg
EOF
# Set APT proxy to ensure apt works in the chroot
echo 'Acquire::http::Proxy "http://172.16.2.30:3142";' >> /mnt/targetos/etc/apt/apt.conf
  return 0
}

#######################################
# create_cloud_init_scripts: Create cloud-init one-time scripts.
#######################################
create_cloud_init_scripts() {
  log "INFO" "Creating cloud-init scripts"
  mkdir -p /mnt/targetos/var/lib/cloud/scripts/per-once || { log "ERROR" "Creating cloud-init scripts directory failed"; return 1; }

  cat <<EOF > /mnt/targetos/var/lib/cloud/scripts/per-once/01-setup-rsyslog.sh
#!/bin/bash
set -e
apt install --yes rsyslog rsyslog-relp
sed -i '/#  Default rules for rsyslog\./a \
module(load="omrelp")
*.*  action(type="omrelp" target="172.16.2.30" port="2514")' /etc/rsyslog.d/50-default.conf
systemctl enable rsyslog
systemctl restart rsyslog
EOF
  chmod +x /mnt/targetos/var/lib/cloud/scripts/per-once/01-setup-rsyslog.sh

  cat <<EOF > /mnt/targetos/var/lib/cloud/scripts/per-once/04-setup-ssh.sh
#!/bin/bash
set -e
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
systemctl restart sshd
EOF
  chmod +x /mnt/targetos/var/lib/cloud/scripts/per-once/04-setup-ssh.sh

  cat <<EOF > /mnt/targetos/var/lib/cloud/scripts/per-once/05-setup-ntp.sh
#!/bin/bash
set -e
timedatectl set-timezone "$TIMEZONE"
timedatectl set-ntp on
systemctl restart systemd-timesyncd
EOF
  chmod +x /mnt/targetos/var/lib/cloud/scripts/per-once/05-setup-ntp.sh

  cat <<EOF > /mnt/targetos/var/lib/cloud/scripts/per-once/07-hosts.sh
#!/bin/bash
set -e
echo -e "172.16.2.30 unimatrixzero.local unimatrixzero minio-01.local minio-01\n172.16.2.45 rpi-serv-001.local rpi-serv-001 minio-02.local minio-02\n172.16.2.46 rpi-serv-002.local rpi-serv-002 minio-03.local minio-03\n172.16.2.47 rpi-serv-003.local rpi-serv-003 minio-04.local minio-04" >> /etc/hosts
EOF
  chmod +x /mnt/targetos/var/lib/cloud/scripts/per-once/07-hosts.sh

  cat <<EOF > /mnt/targetos/var/lib/cloud/scripts/per-once/08-apt-update.sh
#!/bin/bash
set -e
apt-get update
apt-get -y full-upgrade
EOF
  chmod +x /mnt/targetos/var/lib/cloud/scripts/per-once/08-apt-update.sh
  return 0
}

#######################################
# configure_chroot_installscript: Create the chroot install script.
# Installs packages that were removed from debootstrap.
#######################################
configure_chroot_installscript() {
  log "INFO" "Configuring chroot install script"
  tee /mnt/targetos/installscript.sh /dev/null <<'THISISTHEEND'
#!/bin/bash
echo "sourcing uuid"
source /uuid
echo "$UUID"
echo "Updating package lists and upgrading packages"
apt update && apt upgrade -y
echo "Installing clevis packages"
apt install --yes clevis-initramfs clevis-luks clevis-systemd clevis-tpm2 clevis
echo "Making /boot/efi"
mkdosfs -F 32 -s 1 -n EFI ${DISK}p1
mkdir /boot/efi
echo /dev/disk/by-uuid/$(blkid -s UUID -o value ${DISK}p1)  /boot/efi vfat defaults 0 0 >> /etc/fstab
mount /boot/efi
echo "making /boot/efi/grub"
mkdir /boot/efi/grub /boot/grub
echo "/boot/efi/grub /boot/grub none defaults,bind 0 0" >> /etc/fstab
mount /boot/grub
echo "Installing the linux signed image for secure boot"
apt install --yes grub-efi-amd64 grub-efi-amd64-signed linux-image-generic shim-signed zfs-initramfs zsys
echo "setting root password"
echo "root:$ROOT_PASSWORD" | chpasswd
echo "Installing ssh-server"
apt install --yes openssh-server
echo "Updating boot/initramfs/grub"
grub-probe /boot
update-initramfs -c -k all
update-grub
echo "Running grub-install"
grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=ubuntu --no-floppy
mkdir /etc/zfs/zfs-list.cache
touch /etc/zfs/zfs-list.cache/bpool
touch /etc/zfs/zfs-list.cache/rpool
zed -F &
sleep ${ZED_SLEEP_DURATION}
cat /etc/zfs/zfs-list.cache/bpool
cat /etc/zfs/zfs-list.cache/rpool
for PSUID in $(ps aux | grep "zed -F" | tr -s ' ' | cut -d ' ' -f2); do kill -9 $PSUID; done
sed -Ei "s|/mnt/targetos/?|/|" /etc/zfs/zfs-list.cache/*
apt-get update
apt install -y 7zip-rar tpm2-tools prometheus-node-exporter yubikey-manager yubikey-personalization
echo "$LUKS_KEY" | clevis luks bind -y -k - -d ${DISK}p4 sss '{"t":2,"pins":{"tang":[{"url":"'"${TANG_URL1}"'"},{"url":"'"${TANG_URL2}"'"},{"url":"'"${TANG_URL3}"'"}]}}'
curl -sSfL https://mirror.cockroachdb.com/cockroach-latest.linux-amd64.tgz | tar xz && mv cockroach-latest.linux-amd64/cockroach /usr/local/bin/
rm -rf cockroach-latest.linux-amd64
useradd -r -m -d /var/lib/cockroach cockroach
mkdir -p /var/lib/cockroach/certs
chown -R cockroach:cockroach /var/lib/cockroach/certs
chmod 700 /var/lib/cockroach/certs
systemctl daemon-reload
curl 172.16.2.30/cloud-init/reporting.sh -o /root/reporting.sh
source /root/reporting.sh
send_status_update "installing" 88 "finished chroot install script"
exit
THISISTHEEND
  chmod +x /mnt/targetos/installscript.sh || { log "ERROR" "Failed to chmod chroot installscript"; return 1; }
  return 0
}

#######################################
# mount_targetos_filesystems: Mount /mnt/targetos run, proc, sys, and dev.
# Rename any existing /etc/resolv.conf to .bak then set new DNS.
#######################################
mount_targetos_filesystems() {
  log "INFO" "Mounting target OS filesystems for chroot"
  mkdir -p /mnt/targetos/{run,proc,sys,dev} || { log "ERROR" "Failed to create mount directories"; return 1; }
  mount -t tmpfs tmpfs /mnt/targetos/run || { log "ERROR" "Failed to mount /mnt/targetos/run"; return 1; }
  mount -t proc /proc /mnt/targetos/proc || { log "ERROR" "Failed to mount /mnt/targetos/proc"; return 1; }
  mount -t sysfs /sys /mnt/targetos/sys || { log "ERROR" "Failed to mount /mnt/targetos/sys"; return 1; }
  mount --bind /dev /mnt/targetos/dev || { log "ERROR" "Failed to bind mount /mnt/targetos/dev"; return 1; }
}

#######################################
# verify_chroot_network: Run network diagnostics inside the chroot.
#######################################
verify_chroot_network() {
  log "INFO" "Verifying network connectivity inside chroot"
  chroot /mnt/targetos /usr/bin/env bash -c '
    echo "Installing dnsutils (if not present)"
    apt-get update && apt-get install -y dnsutils;
    echo "----- NICs and IP Addresses (ip addr show) -----"
    ip addr show;
    echo "----- Routing Table (ip route show) -----"
    ip route show;
    echo "----- Default Gateway -----"
    ip route show default;
    echo "----- DNS Configuration (/etc/resolv.conf) -----"
    cat /etc/resolv.conf;
    echo "----- Active Listening Sockets (ss -tulwn) -----"
    ss -tulwn;
    echo "----- Ping Default Gateway (172.16.2.1) -----"
    ping -c 4 172.16.2.1;
    echo "----- Dig Query using @8.8.8.8 for google.com -----"
    dig @8.8.8.8 google.com;
  ' || { log "ERROR" "Chroot network verification failed"; return 1; }
  return 0
}

#######################################
# run_chroot_installscript: Chroot into /mnt/targetos and execute installscript.sh.
#######################################
run_chroot_installscript() {
  log "INFO" "Chrooting into target OS and running install script"
  chroot /mnt/targetos /usr/bin/env DISK="$DISK" UUID="$UUID" bash --login /installscript.sh || { log "ERROR" "Chroot installscript failed"; return 1; }
  return 0
}

#######################################
# log_cloud_init_logs: Log cloud-init and cloud-init-output logs from the target OS.
#######################################
log_cloud_init_logs() {
  log "INFO" "Logging cloud-init logs from target OS"
  if [ -f /mnt/targetos/var/log/cloud-init.log ]; then
    log "INFO" "----- /mnt/targetos/var/log/cloud-init.log -----"
    cat /mnt/targetos/var/log/cloud-init.log | while read -r line; do log "INFO" "$line"; done
  else
    log "WARNING" "/mnt/targetos/var/log/cloud-init.log not found"
  fi
  if [ -f /mnt/targetos/var/log/cloud-init-output.log ]; then
    log "INFO" "----- /mnt/targetos/var/log/cloud-init-output.log -----"
    while read -r line; do
      logger -t cloud-init-output "$line"
      echo "$(date +'%Y-%m-%d %H:%M:%S') [CLOUD-INIT-OUTPUT] $line"
    done < /mnt/targetos/var/log/cloud-init-output.log
  else
    log "WARNING" "/mnt/targetos/var/log/cloud-init-output.log not found"
  fi
  return 0
}

#######################################
# finalize_install: Unmount filesystems and export ZFS pools.
# Always returns success to ensure a clean state for subsequent runs.
#######################################
finalize_install() {
  log "INFO" "Finalizing installation: unmounting filesystems and exporting ZFS"
  echo "Unmounting filesystems"
  mount | grep -v zfs | tac | awk '/\/mnt/ {print $3}' | xargs -i{} umount -lf {} 2>/dev/null
  umount -l -f /mnt/targetos/boot 2>/dev/null
  umount -l -f /mnt/targetos 2>/dev/null
  echo "Exporting ZFS filesystem"
  zpool export -a 2>/dev/null || log "WARNING" "Exporting ZFS filesystem encountered an error; continuing anyway."
  return 0
}

#######################################
# main: Execute all steps in order with retry support.
#######################################
main() {
  log "INFO" "Starting jinstall.sh process"

  retry pre_install
  retry update_package_lists
  retry install_required_packages
  retry stop_services
  retry configure_timezone_ntp
  retry validate_time
  retry partition_disk
  retry format_partitions
  retry configure_luks
  retry mount_and_create_zfs

  # Split ZFS pool creation into separate functions.
  retry create_bpool
  retry create_rpool
  retry generate_uuid
  retry create_bpool_datasets
  retry create_rpool_datasets

  retry install_base_system
  retry configure_hostname_network
  retry configure_apt_sources
  retry create_cloud_init_scripts
  retry configure_chroot_installscript
  retry mount_targetos_filesystems
  retry verify_chroot_network
  retry run_chroot_installscript
  retry log_cloud_init_logs
  if [ "$SKIP_FINALIZE" != "1" ]; then
    retry finalize_install
  else
    log "WARNING" "SKIP_FINALIZE flag set, skipping finalize_install"
  fi

  # Send the final detailed report.
  send_final_report

  log "INFO" "jinstall.sh process completed successfully"
}

# Execute main function.
main
```

reporting.sh
```
#!/bin/bash
# reporting.sh - Reporting functions for jinstall.sh
# This script is intended to be sourced by jinstall.sh.
# It provides functions to:
#   - Send status updates.
#   - Upload logs.
#   - Send hardware information.
#   - Generate a comprehensive hardware report as JSON.
#   - Send a final report with collected logs and system information.

#######################################
# send_status_update
# Sends a status update with full event fields.
#######################################
send_status_update() {
  local STATUS="${1:-pending}"
  local PROGRESS="${2:-0}"
  local MESSAGE="${3:-}"
  local WEBHOOK_URL="http://172.16.2.30:25000/api/webhook"
  local SOURCE_IP
  SOURCE_IP=$(hostname -I | awk '{print $1}')
  local HOSTNAME
  HOSTNAME=$(hostname)
  local TIMESTAMP
  TIMESTAMP=$(date +%s)

  # Construct the JSON payload with all required fields.
  local json_payload
  json_payload=$(jq -n \
    --arg source_ip "$SOURCE_IP" \
    --argjson timestamp "$TIMESTAMP" \
    --arg origin "cloud-init" \
    --arg description "$MESSAGE" \
    --arg name "$HOSTNAME" \
    --arg result "" \
    --arg event_type "status_update" \
    --arg status "$STATUS" \
    --argjson progress "$PROGRESS" \
    --arg message "$MESSAGE" \
    '{
         "source_ip": $source_ip,
         "timestamp": $timestamp,
         "origin": $origin,
         "description": $description,
         "name": $name,
         "result": $result,
         "event_type": $event_type,
         "status": $status,
         "progress": $progress,
         "message": $message,
         "files": []
    }')

  # Validate the JSON payload.
  echo "$json_payload" | jq empty >/dev/null 2>&1
  if [ $? -ne 0 ]; then
    echo "Invalid JSON detected, aborting status update." >&2
    return 1
  fi

  # Send the JSON payload.
  curl -s -X POST "$WEBHOOK_URL" -H "Content-Type: application/json" -d "$json_payload"
}

#######################################
# upload_logs
# Uploads a given file (or system logs if no argument provided) to the reporting webhook.
#######################################
upload_logs() {
  local WEBHOOK_URL="http://172.16.2.30:25000/api/webhook"
  local HOSTNAME
  HOSTNAME=$(hostname)
  local TIMESTAMP
  TIMESTAMP=$(date +%s)
  local SOURCE_IP
  SOURCE_IP=$(hostname -I | awk '{print $1}')

  validate_json() {
    local json_data="$1"
    echo "$json_data" | jq empty >/dev/null 2>&1 || {
      echo "Invalid JSON detected, skipping upload." >&2
      return 1
    }
  }

  upload_file() {
    local file_path="$1"
    if [ -f "$file_path" ]; then
      local file_content encoded_content log_json
      file_content=$(cat "$file_path")
      encoded_content=$(echo -n "$file_content" | base64 -w 0)
      log_json=$(jq -n \
        --arg source_ip "$SOURCE_IP" \
        --argjson timestamp "$TIMESTAMP" \
        --arg origin "cloud-init" \
        --arg description "Log file upload" \
        --arg name "$HOSTNAME" \
        --arg result "" \
        --arg event_type "log_update" \
        --arg path "$file_path" \
        --arg encoding "base64" \
        --arg content "$encoded_content" \
        '{ "source_ip": $source_ip, "timestamp": $timestamp, "origin": $origin, "description": $description, "name": $name, "result": $result, "event_type": $event_type, "files": [ { "path": $path, "encoding": $encoding, "content": $content } ] }')
      validate_json "$log_json" || return
      curl -s -X POST "$WEBHOOK_URL" -H "Content-Type: application/json" -d "$log_json"
      echo "Uploaded log: $file_path"
    else
      echo "File $file_path not found!"
    fi
  }

  if [ -n "$1" ]; then
    upload_file "$1"
  else
    echo "Uploading system logs..."
    local dmesg_output journal_output
    dmesg_output=$(dmesg)
    if [ -n "$dmesg_output" ]; then
      echo "$dmesg_output" > /tmp/dmesg.log
      upload_file "/tmp/dmesg.log"
      rm -f /tmp/dmesg.log
    fi
    journal_output=$(journalctl --no-pager)
    if [ -n "$journal_output" ]; then
      echo "$journal_output" > /tmp/journal.log
      upload_file "/tmp/journal.log"
      rm -f /tmp/journal.log
    fi
    for logfile in /var/log/*.log; do
      upload_file "$logfile"
    done
    echo "All logs uploaded successfully!"
  fi
}

#######################################
# send_hardware_info
# Sends hardware info (network interface details) as JSON.
#######################################
send_hardware_info() {
  local WEBHOOK_URL="http://172.16.2.30:25000/api/hardware-info"
  local HOSTNAME
  HOSTNAME=$(hostname)
  local iface
  iface=$(ls /sys/class/net/ | head -n 1)
  local MAC
  MAC=$(cat /sys/class/net/$iface/address)
  local DRIVER
  DRIVER=$(ethtool -i $iface 2>/dev/null | awk '/driver/{print $2}')
  local CHIPSET="unknown"

  local json_payload
  json_payload=$(jq -n \
    --arg client_id "$HOSTNAME" \
    --arg mac_address "$MAC" \
    --arg interface_name "$iface" \
    --arg chipset "$CHIPSET" \
    --arg driver "$DRIVER" \
    '{ "client_id": $client_id, "mac_address": $mac_address, "interface_name": $interface_name, "chipset": $chipset, "driver": $driver }')

  echo "$json_payload" | jq empty >/dev/null 2>&1 || { echo "Invalid JSON detected, aborting hardware info upload." >&2; return 1; }
  curl -s -X POST "$WEBHOOK_URL" -H "Content-Type: application/json" -d "$json_payload"
}

#######################################
# send_cloud_init_info
# Sends cloud-init user data as JSON.
#######################################
send_cloud_init_info() {
  local WEBHOOK_URL="http://172.16.2.30:25000/api/cloud-init"
  local HOSTNAME
  HOSTNAME=$(hostname)
  local iface
  iface=$(ls /sys/class/net/ | head -n 1)
  local MAC
  MAC=$(cat /sys/class/net/$iface/address)
  local USER_DATA
  USER_DATA=$(cat /etc/cloud/cloud.cfg 2>/dev/null || echo "default cloud-init user data")

  local json_payload
  json_payload=$(jq -n \
    --arg client_id "$HOSTNAME" \
    --arg mac_address "$MAC" \
    --arg user_data "$USER_DATA" \
    '{ "client_id": $client_id, "mac_address": $mac_address, "user_data": $user_data }')

  echo "$json_payload" | jq empty >/dev/null 2>&1 || { echo "Invalid JSON detected, aborting cloud-init info upload." >&2; return 1; }
  curl -s -X POST "$WEBHOOK_URL" -H "Content-Type: application/json" -d "$json_payload"
}

#######################################
# report_collect_hardware_info
# Prints network interface details.
#######################################
report_collect_hardware_info() {
  echo "Collecting network interfaces info:"
  for iface in $(ls /sys/class/net/); do
    echo "Interface: $iface"
    echo "MAC Address: $(cat /sys/class/net/$iface/address)"
    echo "Chipset: $(ethtool -i $iface 2>/dev/null | grep driver)"
  done
}

#######################################
# report_retrieve_smbios_info
# Outputs the SMBIOS UUID and motherboard serial number.
#######################################
report_retrieve_smbios_info() {
  echo "SMBIOS UUID: $(cat /sys/class/dmi/id/product_uuid 2>/dev/null || echo 'null')"
  echo "Motherboard Serial: $(cat /sys/class/dmi/id/board_serial 2>/dev/null || echo 'null')"
}

#######################################
# report_generate_hardware_report
# Gathers network interfaces, SMBIOS data, and lshw JSON output,
# then outputs a comprehensive JSON report.
#######################################
report_generate_hardware_report() {
  network_interfaces=$(for iface in $(ls /sys/class/net/); do
    mac=$(cat /sys/class/net/$iface/address)
    driver=$(ethtool -i $iface 2>/dev/null | grep driver | awk '{print $2}')
    echo "{\"interface\": \"$iface\", \"mac_address\": \"$mac\", \"chipset\": \"$driver\"}"
  done | jq -s .)
  smbios_uuid=$(cat /sys/class/dmi/id/product_uuid 2>/dev/null || echo 'null')
  motherboard_serial=$(cat /sys/class/dmi/id/board_serial 2>/dev/null || echo 'null')
  hardware_info=$(lshw -json 2>/dev/null || echo "{}")
  report_json=$(jq -n \
    --argjson network_interfaces "$network_interfaces" \
    --arg smbios_uuid "$smbios_uuid" \
    --arg motherboard_serial "$motherboard_serial" \
    --argjson hardware_info "$hardware_info" \
    '{ "network_interfaces": $network_interfaces, "smbios_uuid": $smbios_uuid, "motherboard_serial": $motherboard_serial, "hardware_info": $hardware_info }')
  echo "$report_json"
}

#######################################
# send_final_report
# Gathers the final hardware report and the last 50 lines of cloud-init-output.log,
# then sends them as JSON to the final reporting endpoint.
#######################################
send_final_report() {
  local WEBHOOK_URL="http://172.16.2.30:25000/api/finalreport"
  local HOSTNAME
  HOSTNAME=$(hostname)
  local TIMESTAMP
  TIMESTAMP=$(date +%s)
  local SYS_INFO
  SYS_INFO=$(report_generate_hardware_report)
  local LOG_CONTENT
  if [ -f /var/log/cloud-init-output.log ]; then
    LOG_CONTENT=$(tail -n 50 /var/log/cloud-init-output.log | base64 -w 0)
  else
    LOG_CONTENT=""
  fi

  local json_payload
  json_payload=$(jq -n \
    --arg client_id "$HOSTNAME" \
    --argjson timestamp "$TIMESTAMP" \
    --arg logs "$LOG_CONTENT" \
    --argjson sys_info "$SYS_INFO" \
    '{ client_id: $client_id, timestamp: $timestamp, logs: $logs, system_info: $sys_info }')

  echo "$json_payload" | jq empty >/dev/null 2>&1 || {
    echo "Final report JSON invalid, aborting final report." >&2
    return 1
  }
  curl -s -X POST "$WEBHOOK_URL" -H "Content-Type: application/json" -d "$json_payload"
}

# End of reporting.sh
```

variables.sh
```
#!/bin/bash
# Variables file for jinstall.sh
# Adjust these values per system

# General system settings
DISK="/dev/nvme0n1"
TIMEZONE="America/New_York"
DEBOOTSTRAP_RELEASE="oracular"
HOSTNAME="len-serv-003"
REPORT_STATUS="/usr/local/bin/report-status.sh"

# Security & Authentication
LUKS_KEY="defaultLUKSkey123"
ROOT_PASSWORD="defaultPassword123"

# Network configuration (Ethernet)
NET_ET_INTERFACE="enp1s0f0"
NET_ET_ADDRESS="172.16.3.96/23"
NET_ET_GATEWAY="172.16.2.1"
NET_ET_SEARCH="jf.local"
NET_ET_NAMESERVERS=("172.16.2.1" "1.1.1.1" "8.8.8.8")

# Network configuration (WiFi)
WIFI_INTERFACE="wlp2s0"
WIFI_ADDRESS="172.16.3.97/23"
WIFI_SSID="VentureIndustries"
WIFI_PASSWORD="WIFIcool"

# CockroachDB configuration
COCKROACH_ADVERTISE="172.16.2.55:26257"
COCKROACH_JOIN="172.16.2.45:36257,172.16.2.45:36357,172.16.2.47:36257,172.16.2.47:36357,172.16.2.30:36257,172.16.2.30:36357"
COCKROACH_CACHE=".25"
COCKROACH_MAX_SQL=".25"
COCKROACH_LOCALITY="region=us,cluster-unit=lenovo"

# Tang server URLs for Clevis binding
TANG_URL1="http://172.16.2.45"
TANG_URL2="http://172.16.2.46"
TANG_URL3="http://172.16.2.47"

# Sleep durations (in seconds)
SLEEP_DURATION=3
ZED_SLEEP_DURATION=30
```
