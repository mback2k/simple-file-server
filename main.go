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

package main

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/mback2k/simple-file-server/fileserver"
)

func setupServer() *fileserver.Server {
	log.Println("Loading server defaults")
	s := fileserver.FileServer()

	log.Println("Loading configuration")
	c, err := loadConfig(s)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Applying configuration")
	s.Address = c.Address
	s.DocumentRoot = c.DocumentRoot
	s.DirectoryListing = c.DirectoryListing
	s.DirectoryIndex = c.DirectoryIndex
	for i := range c.AliasList {
		alias := c.AliasList[i]
		s.AliasMap[alias.Source] = alias.Target
	}

	l, err := log.ParseLevel(c.Logging.Level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)

	return s
}

func main() {
	s := setupServer()

	runtime.GC()

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
