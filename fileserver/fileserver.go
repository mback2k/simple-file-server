/*
	simple-file-server - A simple file server to host static files.
	Copyright (C) 2019  Marc Hoersken <info@marc-hoersken.de>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package fileserver

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Server data structure
type Server struct {
	http.Handler

	Address          string
	DocumentRoot     string
	DirectoryListing bool
	DirectoryIndex   []string
	AliasMap         map[string]string
}

// FileServer creates a new Server struct
func FileServer() *Server {
	s := &Server{
		Address:          ":http-alt",
		DocumentRoot:     ".",
		DirectoryListing: false,
		DirectoryIndex:   []string{"index.html"},
		AliasMap:         map[string]string{},
	}
	return s
}

func (s *Server) getURLPath(r *http.Request) string {
	urlpath := r.URL.Path
	log.Debugln("getURLPath", urlpath)
	if !strings.HasPrefix(urlpath, "/") {
		urlpath = "/" + urlpath
		r.URL.Path = urlpath
	}
	urlpath = path.Clean(r.URL.Path)
	for alias := range s.AliasMap {
		if strings.HasPrefix(urlpath, alias) {
			urlpath = path.Join(s.AliasMap[alias], urlpath[len(alias):])
		}
	}
	log.Debugln("getURLPath", urlpath)
	return urlpath
}

func (s *Server) getURLFile(urlpath string) string {
	urlfile := filepath.Clean(filepath.FromSlash(urlpath))
	log.Debugln("getURLFile", urlfile)
	if length := len(urlfile); length > 0 && os.IsPathSeparator(urlfile[length-1]) {
		urlfile = urlfile[:length-1]
	}
	log.Debugln("getURLFile", urlfile)
	return urlfile
}

func (s *Server) findDirectoryFile(file string) (string, error) {
	log.Debugln("findDirectoryFile", file)
	stat, err := os.Stat(file)
	if stat != nil {
		if stat.IsDir() {
			// Found directory, checking for index files
			for i := range s.DirectoryIndex {
				index := filepath.Join(file, s.DirectoryIndex[i])
				index, err = s.findDirectoryFile(index)
				if err == nil {
					file = index
					break
				}
			}
			// Ignore missing directory index if listing is enabled
			if s.DirectoryListing {
				err = nil
			}
		}
	}
	log.Debugln("findDirectoryFile", file)
	return file, err
}

func (s *Server) getFileFromURLFile(urlfile string) (string, error) {
	file := filepath.Join(s.DocumentRoot, urlfile)
	return s.findDirectoryFile(file)
}

func (s *Server) sendStatus(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlpath := s.getURLPath(r)
	urlfile := s.getURLFile(urlpath)

	file, err := s.getFileFromURLFile(urlfile)
	if err != nil {
		log.Error(err)
		s.sendStatus(w, http.StatusNotFound)
		return
	}

	log.Debugln("ServeFile", file)
	http.ServeFile(w, r, file)
	return
}

func (s *Server) checkDocumentRoot() error {
	if !filepath.IsAbs(s.DocumentRoot) {
		root, err := filepath.Abs(s.DocumentRoot)
		if err != nil {
			return err
		}
		s.DocumentRoot = root
	}
	if !strings.HasSuffix(s.DocumentRoot, "/") {
		s.DocumentRoot = s.DocumentRoot + "/"
	}

	stat, err := os.Stat(s.DocumentRoot)
	if err != nil {
		return err
	}
	if stat == nil {
		return &os.PathError{Path: s.DocumentRoot}
	}
	if !stat.IsDir() {
		return &os.PathError{Path: s.DocumentRoot}
	}
	return nil
}

// ListenAndServe listens on the network address and runs the HTTP server
func (s *Server) ListenAndServe() error {
	log.Println("Checking document root")
	err := s.checkDocumentRoot()
	if err != nil {
		return err
	}

	log.Println("Starting HTTP server at", s.Address)
	return http.ListenAndServe(s.Address, s)
}
