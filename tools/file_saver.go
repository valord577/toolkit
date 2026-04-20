package tools

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/valord577/clix"
)

var (
	base string
	port string
)

func init() {
	FileSaver.FlagStringVar(&base, "base", "", "default path to save the file")
	FileSaver.FlagStringVar(&port, "port", "58081", "http server's port")
}

var FileSaver = &clix.Command{
	Name: "filesaver",

	Summary: "Save File to FS",

	ShowDefValue: true,
	Run: func(*clix.Command, []string) (err error) {
		if base == "" {
			if base, err = os.Getwd(); err != nil {
				slog.Error("failed @os.Getwd(), errmsg: " + err.Error())
				return
			}
		}
		slog.Info("base dir: " + base)

		// Use (net/http).DefaultServeMux
		http.HandleFunc("/", filesaver)

		serv := &http.Server{
			Addr:    ":" + port,
			Handler: http.DefaultServeMux,
		}
		go func() {
			e := serv.ListenAndServe()
			if e != nil && e != http.ErrServerClosed {
				panic(err)
			}
		}()

		// Block and listen for signals
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-sig
		slog.Info("recv signal: " + s.String())

		err = serv.Shutdown(context.Background())
		return
	},
}

func filesaver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only `POST` allowed", http.StatusMethodNotAllowed)
		return
	}
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "missing filename", http.StatusBadRequest)
		return
	}

	filepath := path.Join(base, filename)
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = io.Copy(file, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("save file: " + filepath)
	http.Error(w, filepath, http.StatusOK)
}
