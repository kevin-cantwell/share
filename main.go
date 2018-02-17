package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	port := openPort()
	go func() {
		var handler http.Handler
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Serve the body of stdin
			body, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalln("Failed to read from stdin:", err)
			}
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(body)
			})
		} else {
			if len(os.Args) < 2 {
				log.Fatalln("Usage: `share <filepath>`")
			}
			arg := os.Args[1]
			fs, err := fileServer(arg)
			if err != nil {
				// Serve the argument as a value
				handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(arg))
				})
			} else {
				// Serve the filepath
				handler = fs
			}
		}
		log.Println("Listening on", "http://localhost:"+port)
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			log.Fatalln(err)
		}
	}()

	cmd := exec.Command("ngrok", "http", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func openPort() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	addr := listener.Addr().String()
	return strings.Split(addr, ":")[1]
}

func fileServer(path string) (http.Handler, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return http.FileServer(http.Dir(path)), nil
	case mode.IsRegular():
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			filename := filepath.Base(path)
			switch r.URL.Path {
			case "/":
				http.Redirect(w, r, "/"+filename, http.StatusFound)
			case "/" + filename:
				http.ServeFile(w, r, path)
			default:
				http.NotFound(w, r)
			}
		}), nil
	default:
		return nil, errors.New("fs: not a file or directory: " + path)
	}
}
