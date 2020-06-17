# grpc2way
golang grpc two way communication with Multiplexing (yamux)
https://github.com/hashicorp/yamux

real time communication

-  server bind on port
  - client connect to server port
   - server accept
    - client create yamux server session
      - server create yamux client session
        - client create server on yamux session
          - server create client on yamux session
          
    
