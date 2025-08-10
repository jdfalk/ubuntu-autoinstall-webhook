<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [System Architecture Diagram](#system-architecture-diagram)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# System Architecture Diagram

```ascii
                                +-----------------+
                                |                 |
                                |  Administrator  |
                                |                 |
                                +--------+--------+
                                         |
                                         | HTTPS
                                         v
+------------------+          +----------+---------+
|                  |          |                    |
|  New System      +--------->+    Webserver      |
|  (PXE Boot)      |  HTTP    |    Service        |
|                  |          |                    |
+------------------+          +---------+----------+
                                        ^
                                        | gRPC
                                        |
         +-----------------------------++-+----------------------------+
         |                             |  |                            |
         v                             v  v                            v
+--------+---------+         +---------+--+------+          +---------+---------+
|                  |         |                   |          |                   |
|  File Editor     |<------->+  Configuration    |<-------->+  Certificate      |
|  Service         |  gRPC   |  Service          |  gRPC    |  Issuer Service   |
|                  |         |                   |          |                   |
+--------+---------+         +---------+---------+          +---------+---------+
         ^                             ^                             ^
         |                             |                             |
         | gRPC                        | gRPC                        | gRPC
         |                             |                             |
+--------+---------+         +---------+---------+          +--------+----------+
|                  |         |                   |          |                   |
|  Database        |<------->+  DNSMasq Watcher  |          |  Admin CLI        |
|  Service         |  gRPC   |  Service          |          |  Command          |
|                  |         |                   |          |                   |
+------------------+         +-------------------+          +-------------------+
```
