# gobs

Config:

| ENV                    | YAML                   | Meaning                                          | Example                         |
|------------------------|------------------------|--------------------------------------------------|---------------------------------|
| CLIENT_ADDR            | client.addr            | rtsp url to client stream                        | "rtsp://192.168.1.66:8554/live" |
| CLIENT_RECONNECT_PAUSE | client.reconnect-pause | Max timeout to client reconnection process in ms | 20000                           |
| SERVER_TCP_PORT        | server.tcp-port        | Server tcp port for rtsp streaming               | 8554                            |
| SERVER_UDP_PORT        | server.udp-port        | Server udp port for rtsp streaming               | 8000                            |
| SERVER_UDP_RTCP_PORT   | server.udp-rtcp-port   | Server udp port for rtcp streaming               | 8001                            |
