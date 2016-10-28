package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/UltimateSoftware/udocs/cli/udocs"
	"github.com/mholt/archiver"
	"golang.org/x/net/context"
)

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (s *Server) reverseProxyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	port := s.settings.Port
	if s.scheme == "https" {
		port = "443"
	}

	rootURL := &url.URL{
		Scheme: s.scheme,
		Host:   s.host + ":" + port,
		Path:   "/" + s.settings.RootRoute,
	}

	httputil.NewSingleHostReverseProxy(rootURL).ServeHTTP(w, r)
}

func (s *Server) staticHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(rice.MustFindBox("../../static").HTTPBox())).ServeHTTP(w, r)
}

func (s *Server) pageHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := s.dao.Fetch(r.URL.Path)
	if err != nil {
		logAndWriteError(w, r, http.StatusNotFound, "unable to fetch data", err)
		return
	}

	if r.URL.Query().Get("ajax") == "true" {
		logAndWriteBinaryResponse(w, r, http.StatusOK, data)
		return
	}

	if ext := filepath.Ext(r.URL.Path); ext != "" && ext != ".html" && ext != ".quip" {
		logAndWriteBinaryResponse(w, r, http.StatusOK, data)
		return
	}

	sidebar, err := udocs.LoadSidebar(s.dao)
	if err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, "failed to load sidebar", err)
		return
	}

	if err := s.tmpl.WithParameter("sidebar", sidebar).Execute(w, "document", data); err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, "failed to execute html template", err)
		return
	}

	logResponse(http.StatusOK, r)
	return

}

func (s *Server) quipBlobHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	thread := ctx.Value("thread").(string)
	id := ctx.Value("id").(string)

	blob, err := udocs.DefaultQuipClient.GetBlob(thread, id)
	if err != nil {
		logAndWriteError(w, r, http.StatusNotFound, "unable to get Quip blob", err)
		return
	}

	logAndWriteBinaryResponse(w, r, http.StatusOK, blob)
	return
}

func (s *Server) updateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	route := ctx.Value("route").(string)
	dest := filepath.Join(udocs.BuildPath(), fmt.Sprintf("%s_%d", route, time.Now().Unix()))

	docs, err := extractTarball(r.Body, dest)
	if err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, "server.updateHandler unable to extract tarball", err)
		return
	}

	if err := udocs.Validate(docs); err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, "server.updateHandler failed to validate docs directory", err)
		return
	}

	if err := udocs.Build(route, docs, s.dao); err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, "server.updateHandler unable to build docs", err)
		return
	}

	href := fmt.Sprintf("%s:%s/%s", s.settings.EntryPoint, s.settings.Port, route)
	logAndWriteJSONResponse(w, r, http.StatusCreated, http.StatusText(http.StatusCreated), href)
}

func (s *Server) destroyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	route := ctx.Value("route").(string)

	sidebar, err := udocs.LoadSidebar(s.dao)
	if err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, "server.destroyHandler failed to load sidebar", err)
		return
	}

	if err := updateSidebar(sidebar, udocs.Summary{Route: route, Header: ""}, s.dao); err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, "server.destroyHandler failed remove resource from sidebar", err)
		return
	}

	if err := s.dao.DeleteGlob(route); err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, "server.destroyHandler failed to delete resource", err)
		return
	}

	logResponse(http.StatusOK, r)
}

func (s *Server) searchHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	queryResult, err := s.dao.Query(q)
	if err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, "server.searchHandler failed to execute query", err)
		return
	}

	sidebar, err := udocs.LoadSidebar(s.dao)
	if err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, "server.searchHandler failed to load sidebar", err)
		return
	}

	tmpl := s.tmpl.WithParameter("query_result", queryResult).WithParameter("sidebar", sidebar)
	if err := tmpl.Execute(w, "search", nil); err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, "server.pageHandler failed to execute template", err)
		return
	}

	logResponse(http.StatusOK, r)
}

func extractTarball(rc io.ReadCloser, dest string) (string, error) {
	dir := filepath.Join(udocs.ArchivePath(), filepath.Base(dest))
	os.MkdirAll(dir, 0755)
	tarball := filepath.Join(udocs.ArchivePath(), filepath.Base(dest), "docs.tar.gz")
	defer os.Remove(tarball)

	tmp, err := os.Create(tarball)
	if err != nil {
		return "", fmt.Errorf("api.extractTarball unable to open tmp file: %v", err)
	}

	if _, err := io.Copy(tmp, rc); err != nil {
		return "", fmt.Errorf("api.extractTarball failed to copy tarball: %v", err)
	}
	rc.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return "", fmt.Errorf("api.extractTarball failed to make dest directory: %v", err)
	}

	if err := archiver.TarGz.Open(tarball, dest); err != nil {
		return "", fmt.Errorf("api.extractTarball failed to gunzip src tar file: %v", err)
	}

	return filepath.Join(dest, "docs"), nil
}

func generalizeStringMap(m map[string]string) map[string]interface{} {
	generalized := make(map[string]interface{})
	for k, v := range m {
		generalized[k] = v
	}
	return generalized
}

func writeBinaryResponse(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	if strings.HasSuffix(r.URL.Path, "css") {
		w.Header().Set("content-type", "text/css")
	} else if strings.HasSuffix(r.URL.Path, "js") {
		w.Header().Set("content-type", "text/javascript")
	} else if strings.HasSuffix(r.URL.Path, "html") {
		w.Header().Set("content-type", "text/html")
	} else if strings.HasSuffix(r.URL.Path, ".css.map") {
		w.Header().Set("content-type", "text/css")
	}
	w.WriteHeader(code)
	w.Write(data)
}

func logAndWriteBinaryResponse(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	writeBinaryResponse(w, r, code, data)
	logResponse(code, r)
}

func logAndWriteJSONResponse(w http.ResponseWriter, r *http.Request, code int, msg, href string) {
	resp := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Href    string `json:"href"`
	}{
		Code:    http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Href:    href,
	}
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error: logAndWriteResponse failed writing data: %v", err)
	}
	logResponse(code, r)
}

func logAndWriteError(w http.ResponseWriter, r *http.Request, code int, msg string, err error) {
	log.Printf("error: %s: %v", msg, err)
	http.Error(w, fmt.Sprintf("%d %s\n%s\n", code, http.StatusText(code), msg), code)
	logResponse(code, r)
}

func logResponse(code int, r *http.Request) {
	log.Printf("%s %d %s %s %s", r.RemoteAddr, code, r.Method, r.URL.String(), r.Proto)
}
