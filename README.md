# sdm panic reproduction
This repo contains an example program that reproduces a panic experienced with the sdm client version 1.4.33.

```shell
$ sdm --version
sdm version 1.4.33 (75a716d03cc387b6686b643d42f168ac7ace6ca4 #1140)
$
```

This does not reproduce the panic in the client itself, but rather produces the same panic in a new program (contained in this repo). This is meant to illustrate an understanding of the conditions that lead to the panic.

## Manifestation
The panic manifests itself in the `sdm listen` systemd unit. This can be first (and most obviously) experienced when installing the sdm agent:
```shell
$ sudo ./sdm install
[sudo] password for thisguy: 
Installing strongDM listener

Checking environment
- user detected: root
- install user detected: thisguy

Creating SDM configuration directory (.sdm)
- /home/thisguy/.sdm created
- /home/thisguy/.sdm/sdm.log created
- /home/thisguy/.sdm/gui.log created
- /home/thisguy/.sdm/state.db initialized

Creating directory to install the binary at /opt/strongdm
- moving binary to /opt/strongdm/bin/sdm
- symlinking from /opt/strongdm/bin/sdm to /usr/local/bin/sdm

Installing strongDM listen in process manager
Service installed and started

Waiting for strongDM listener to start

```
This hangs indefinately (mostly, I explain later).

You can then see evidence of the panic in the logs for the sdm unit:
```shell
$ journalctl -u sdm
May 06 22:08:13 pop-os systemd[1]: Started strongDM database proxy service..
May 06 22:08:13 pop-os sdm[232156]: panic: runtime error: invalid memory address or nil pointer dereference
May 06 22:08:13 pop-os sdm[232156]: [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x146e4c6]
May 06 22:08:13 pop-os sdm[232156]: goroutine 1 [running]:
May 06 22:08:13 pop-os sdm[232156]: github.com/awnumar/memguard/core.Purge()
May 06 22:08:13 pop-os sdm[232156]:         github.com/awnumar/memguard@v0.21.0/core/exit.go:19 +0x36
May 06 22:08:13 pop-os sdm[232156]: github.com/awnumar/memguard/core.Panic(0x17a1ea0, 0xc0014d6350)
May 06 22:08:13 pop-os sdm[232156]:         github.com/awnumar/memguard@v0.21.0/core/exit.go:77 +0x22
May 06 22:08:13 pop-os sdm[232156]: github.com/awnumar/memguard/core.NewBuffer(0x20, 0xc0014d41b0, 0xc00138de28, 0x146eaed)
May 06 22:08:13 pop-os sdm[232156]:         github.com/awnumar/memguard@v0.21.0/core/buffer.go:75 +0x33a
May 06 22:08:13 pop-os sdm[232156]: github.com/awnumar/memguard/core.NewCoffer(0xc0014d6340)
May 06 22:08:13 pop-os sdm[232156]:         github.com/awnumar/memguard@v0.21.0/core/coffer.go:33 +0x48
May 06 22:08:13 pop-os sdm[232156]: github.com/awnumar/memguard/core.init.0()
May 06 22:08:13 pop-os sdm[232156]:         github.com/awnumar/memguard@v0.21.0/core/enclave.go:13 +0x36
May 06 22:08:13 pop-os sdm[232156]: runtime: note: your Linux kernel may be buggy
May 06 22:08:13 pop-os sdm[232156]: runtime: note: see https://golang.org/wiki/LinuxKernelSignalVectorBug
May 06 22:08:13 pop-os sdm[232156]: runtime: note: mlock workaround for kernel bug failed with errno 12
May 06 22:08:13 pop-os systemd[1]: sdm.service: Main process exited, code=exited, status=2/INVALIDARGUMENT
May 06 22:08:13 pop-os systemd[1]: sdm.service: Failed with result 'exit-code'.
May 06 22:08:16 pop-os systemd[1]: sdm.service: Scheduled restart job, restart counter is at 1.
May 06 22:08:16 pop-os systemd[1]: Stopped strongDM database proxy service..
May 06 22:08:16 pop-os systemd[1]: Started strongDM database proxy service..
May 06 22:08:16 pop-os sdm[232197]: panic: runtime error: invalid memory address or nil pointer dereference
May 06 22:08:16 pop-os sdm[232197]: [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x146e4c6]
May 06 22:08:16 pop-os sdm[232197]: goroutine 1 [running]:
May 06 22:08:16 pop-os sdm[232197]: github.com/awnumar/memguard/core.Purge()
May 06 22:08:16 pop-os sdm[232197]:         github.com/awnumar/memguard@v0.21.0/core/exit.go:19 +0x36
May 06 22:08:16 pop-os sdm[232197]: github.com/awnumar/memguard/core.Panic(0x17a1ea0, 0xc000e86380)
May 06 22:08:16 pop-os sdm[232197]:         github.com/awnumar/memguard@v0.21.0/core/exit.go:77 +0x22
May 06 22:08:16 pop-os sdm[232197]: github.com/awnumar/memguard/core.NewBuffer(0x20, 0xc000e9a0b0, 0x0, 0x0)
May 06 22:08:16 pop-os sdm[232197]:         github.com/awnumar/memguard@v0.21.0/core/buffer.go:75 +0x33a
May 06 22:08:16 pop-os sdm[232197]: github.com/awnumar/memguard/core.NewCoffer(0xc000e86360)
May 06 22:08:16 pop-os sdm[232197]:         github.com/awnumar/memguard@v0.21.0/core/coffer.go:35 +0x94
May 06 22:08:16 pop-os sdm[232197]: github.com/awnumar/memguard/core.init.0()
May 06 22:08:16 pop-os sdm[232197]:         github.com/awnumar/memguard@v0.21.0/core/enclave.go:13 +0x36
May 06 22:08:16 pop-os systemd[1]: sdm.service: Main process exited, code=exited, status=2/INVALIDARGUMENT
May 06 22:08:16 pop-os systemd[1]: sdm.service: Failed with result 'exit-code'.
```
This continues indefinitely, with the "restart counter" increasing over time.

Worth noting is the "runtime: note: your Linux kernel may be buggy". I've run the test program at the linked wiki entry and believe my kernel does not have this issue. For completeness my `uname -a` output is here:
```shell
$ uname -a
Linux pop-os 5.4.0-7626-generic #30~1588169883~20.04~bbe668a-Ubuntu SMP Wed Apr 29 21:00:02 UTC  x86_64 x86_64 x86_64 GNU/Linux
$
```
That message from the go runtime package is given whenever the `mlock` go issues as a mitigattion against this kernel bug errors. The given error (12) is `ENOMEM` which indicates the process has hit the ulimit for locked memory (this defaults to 65536kbs on my system), and whether the runtime is able to display that error before the program exits is a race condition against the panic in the memguard package (notice it only showed that for one of the two crashes shown above).

## This Program
The program in this repository reproduces the exact same panic in memguard. I use the same version of `github.com/awnumar/memguard` as the sdm client listed here, as well as the same version of go (1.14.1, confirmed with sdm support).

### Building
While this issue is not version specific, you should (at least at first) build/run with the same version of go as the sdm client listed here. Go makes this super easy:

Assuming `$GOPATH/bin` is in your `$PATH`:
```shell
$ go get golang.org/dl/go1.14.1
$ go1.14.1 download
$ cd sdm-panic-reproduction
$ go1.14.1 run .
```

### The panic
The `github.com/awnumar/memguard` package has an `init()` function that attempts to allocate and `mlock` some memory (not much, but some). If the `mlock` fails it attempts to back-out cleanly, and call it's `Purge()` function. However, because initialization hadn't yet completed, `Purge()` calls `.Lock()` on a mutex that has not yet been allocated; this is the `invalid memory address or nil pointer dereference` (it's a nil pointer).

What this means though, is that *before memguards `init()` is run* up to the `ulimit` in memory has already been locked. This is no doubt at least *partially* due to the go runtime's kernel bug mitigation (which it does based on the presented kernel version, since it can't test if where you're running is patched / compiled with GCC 8 and thus not affected), but this repo demonstrates another way it can be encountered: another `init()` function `mlock`ing too much memory first.

### But sometimes...
I found that very rarely the `sdm listen` agent will start successfully. But I do mean rarely (think 200+ restarts most of the time before it works). This seems very strange and indicates a race condition possibly?

### Tomfoolery
`init()`s are run in order based on dependancies, and alphabetically. My github name (`thisguycodes`) comes after `awnumar`, so I had to change the package name in `go.mod` (I prepended an "a") to get my `init()` to run first. But with go modules your repo location doesn't have to match so hey w/e it works!

If you change the package name back in `go.mod` (and update the import in `main.go`) you'll see a panic from my package instead.

# Workaround
For those experiencing this issue, I also have found an effective workaround.

Edit the systemd unit file `sdm` installs (`/etc/systemd/system/sdm.service`):
```
[Unit]
Description=strongDM database proxy service.
ConditionFileIsExecutable=/opt/strongdm/bin/sdm
Requires=multi-user.target
After=multi-user.target

[Service]
ExecStart=/opt/strongdm/bin/sdm "listen"


WorkingDirectory=/home/USER/.sdm

User=USER


Restart=always
RestartSec=3
EnvironmentFile=-/etc/sysconfig/sdm

[Install]
WantedBy=multi-user.target
```

And add the line:
```
LimitMEMLOCK=infinity
```
In the `[Service]` section.