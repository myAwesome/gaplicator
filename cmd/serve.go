package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/myAwesome/gaplicator/internal/generator"
)

var (
	serveHost string
	servePort int
)

type buildRequest struct {
	YAML   string `json:"yaml"`
	Output string `json:"output"`
}

type buildResponse struct {
	OK     bool     `json:"ok"`
	Output string   `json:"output,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a local HTTP bridge for the web schema generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})
		mux.HandleFunc("/build", handleBuild)

		addr := fmt.Sprintf("%s:%d", serveHost, servePort)
		fmt.Printf("Bridge listening on http://%s\n", addr)
		fmt.Println("POST /build with JSON: {\"yaml\":\"...\",\"output\":\"dist/my-app\"}")
		return http.ListenAndServe(addr, withCORS(mux))
	},
}

func handleBuild(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, buildResponse{OK: false, Errors: []string{"failed to read request body"}})
		return
	}

	var req buildRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, buildResponse{OK: false, Errors: []string{"invalid JSON payload"}})
		return
	}

	if strings.TrimSpace(req.YAML) == "" {
		writeJSON(w, http.StatusBadRequest, buildResponse{OK: false, Errors: []string{"yaml is required"}})
		return
	}

	cfg, err := generator.ParseConfigBytes([]byte(req.YAML))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, buildResponse{OK: false, Errors: []string{fmt.Sprintf("parse config: %v", err)}})
		return
	}

	if errs := generator.ValidateConfig(cfg); len(errs) > 0 {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, e.Error())
		}
		writeJSON(w, http.StatusBadRequest, buildResponse{OK: false, Errors: msgs})
		return
	}

	out := strings.TrimSpace(req.Output)
	if out == "" {
		out = "dist"
	}
	out = filepath.Clean(out)

	if err := runBuild(cfg, out, nil); err != nil {
		writeJSON(w, http.StatusInternalServerError, buildResponse{OK: false, Errors: []string{err.Error()}})
		return
	}

	writeJSON(w, http.StatusOK, buildResponse{
		OK:     true,
		Output: out,
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func init() {
	serveCmd.Flags().StringVar(&serveHost, "host", "127.0.0.1", "Host for the bridge server")
	serveCmd.Flags().IntVar(&servePort, "port", 8787, "Port for the bridge server")
}
