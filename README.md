# goChatTest
Implementation of the test task, chat in the goland language

## Introductory assignment

implement a tcp server for fragmented packets
server functions:
1. send messages to all clients connected.
2. send messages to clients specified (you should tag the client)
(This part is the logical function that should to implement)

implement a client connect to the server

message format is :
| 2 bytes | x  bytes |
| content length | content|

* server and client use same message format for communication

![image](https://github.com/DFofanov/goChatTest/blob/main/doc/image.png?raw=true)

test case:
(There are 10 clients, all of which can get broadcast messages
Client 1 sends a directed broadcast message to client 2)

## License
Licensed under the GPL-3.0 License.