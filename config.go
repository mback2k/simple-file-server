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
	"github.com/mback2k/simple-file-server/fileserver"
	"github.com/spf13/viper"
)

type configAlias struct {
	Source string
	Target string
}

type configLogging struct {
	Level string
}

type config struct {
	Address          string
	DocumentRoot     string
	DirectoryListing bool
	DirectoryIndex   []string
	AliasList        []*configAlias
	Logging          *configLogging
}

func loadConfig(s *fileserver.Server) (*config, error) {
	vpr := viper.GetViper()
	vpr.SetDefault("Address", s.Address)
	vpr.SetDefault("DocumentRoot", s.DocumentRoot)
	vpr.SetDefault("DirectoryListing", s.DirectoryListing)
	vpr.SetDefault("DirectoryIndex", s.DirectoryIndex)
	vpr.SetConfigName("simple-file-server")
	vpr.AddConfigPath("/etc/simple-file-server/")
	vpr.AddConfigPath("$HOME/.simple-file-server")
	vpr.AddConfigPath(".")
	err := vpr.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg config
	err = vpr.UnmarshalExact(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
