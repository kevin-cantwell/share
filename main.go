package main

import (
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
	log.Println("Listening on", "http://localhost:"+port)
	go func() {
		if err := http.ListenAndServe(":"+port, handler()); err != nil {
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

func handler() http.Handler {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return stdinServer()
	} else {
		return fileServer()
	}
}

func openPort() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalln("Could not obtain an open port:", err)
	}
	defer listener.Close()
	addr := listener.Addr().String()
	return strings.Split(addr, ":")[1]
}

func stdinServer() http.Handler {
	// Serve the body of stdin
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln("Failed to read from stdin:", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})
}

func fileServer() http.Handler {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: `share <filepath>`")
	}
	path := os.Args[1]
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalln("Could not get the file mode:", err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// Serve a directory
		return http.FileServer(http.Dir(path))
	case mode.IsRegular():
		// Serve a single file
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
		})
	default:
		log.Fatalln("Not a file or directory:", path)
	}

	return nil
}
