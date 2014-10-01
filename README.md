### zfs

zfs golang package helps you manage zfs filesystems/snaphots on local/remote hosts

this package wraps runcmd package; i create it essentially for zbackup

(https://github.com/theairkit/zbackup)

so, I'll be glad if you find this package useful,

but I think, that much better if you will use runcmd for your cases:

https://github.com/theairkit/runcmd

all types and methods are self-explained:

http://godoc.org/github.com/theairkit/zfs

Installation:
```bash
go get github.com/theairkit/zfs
```

### Description and examples

First, create Zfs runner: this is a type, that holds runcmd.Runner,
and then use him for manage zfs filsystems/snapshots:

```go
lRunner, err := zfs.NewZfs(runcmd.NewLocalRunner())
if err != nil {
    // handle error
}
list, err := lRunner.List("zroot", zfs.FS, true)
if err != nil {
    // handle error
}
```

Useful code snippet: send zfs snapshot from local to remote host:

```
lRunner, err := zfs.NewZfs(runcmd.NewLocalRunner())
if err != nil {
    // handle error
}

rRunner, err := zfs.NewZfs(runcmd.NewRemoteKeyAuthRunner(
    user,
    host,
    key,
    ))
if err != nil {
    // handle error
}

cmdRecv, err := this.rRunner.RecvSnap(dst, snapPostfix)
if err != nil {
    // handle error
}

cmdSend, err := this.lRunner.SendSnap(src, snapCurr, snapNew, cmdRecv)
if err != nil {
    // handle error
}

if err := cmdSend.Wait(); err != nil {
    // handle error
}
// In this case EOF is not error: http://golang.org/pkg/io/
// EOF is the error returned by Read when no more input is available.
// Functions should return EOF only to signal a graceful end of input.
if err := cmdRecv.StdinPipe().Close(); err != nil && err != io.EOF {
    // handle error
}
if err := cmdRecv.Wait(); err != nil {
    // handle error
}
```

zfs_test.go - WIP
