# Installing

Use `go get` to install the latest version

    go get -u github.com/bigflood/timeline

# Usage

    timeline command comand-arguments ...

# Example

    $ timeline docker pull ubuntu:16.04
    O> 4.7446522s  16.04: Pulling from library/ubuntu
    O> 6.2898623s  b234f539f7a1: Already exists
    O> 6.5368239s  55172d420b43: Already exists
    O> 6.7952175s  5ba5bbeb6b91: Already exists
    O> 7.0776396s  43ae2841ad7a: Already exists
    O> 7.3549998s  f6c9c6de4190: Already exists
    O> 7.9952407s  Digest: sha256:b050c1822d37a4463c01ceda24d0fc4c679b0dd3c43e742730e2884d3c582e3a
    O> 8.3032053s  Status: Downloaded newer image for ubuntu:16.04
    
    Command time: 8.3121773s
    System CPU time: 78.125ms
    User CPU time: 156.25ms
