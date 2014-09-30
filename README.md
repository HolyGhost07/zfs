### zfs

zfs golang package helps you manage zfs filesystems/snaphots on local/remote hosts

this package wraps runcmd package; i create it essentially for my zbackup program 
(https://github.com/theairkit/zbackup)

so, I'll be glad if you find this package useful,

but I think, that much better if you will use runcmd for your cases:

https://github.com/theairkit/runcmd


all types and methods are self-explained:

http://godoc.org/github.com/theairkit/runcmd

Installation:
```bash
go get github.com/theairkit/zfs
```

### Description and examples

First, create Zfs runner: this is a type, that holds runcmd.Runner,
and use him for manage zfs filsystems/snapshots

```go
//TODO
```

Useful code snippet: send zfs snapshot from local to remote host:

```
//TODO
```

zfs_test.go - WIP
