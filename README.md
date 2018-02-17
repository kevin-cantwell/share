### What is this?

This is just a simple command-line tool that will spin up a local server and expose it to the internet using ngrok. It makes sharing files or text convenient.

### Dependencies

* [ngrok](https://ngrok.com/download)

### Install

```
$ go get -u github.com/kevin-cantwell/share
```

### Usage

Share a file:

```
$ share foo.txt
```

Share all files in a directory:

```
$ share .
```

Share from stdin:

```
$ echo "SECRET PASSWORD" | share
```

Any of the above commands will expose the inputs to the internet via ngrok:

```sh
ngrok by @inconshreveable                                                                                                                                         (Ctrl+C to quit)

Session Status                online
Version                       2.2.8
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    http://826c706c.ngrok.io -> localhost:57829
Forwarding                    https://826c706c.ngrok.io -> localhost:57829

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```