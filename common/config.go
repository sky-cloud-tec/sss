// Simple syslog server.
// Copyright (C) 2019  sky-cloud.net
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package common

// AppConfig app configurations
type AppConfig struct {
	LogCfg *LogConfig
	SrvCfg *ServerConfig
}

// ServerConfig syslog server configuration
type ServerConfig struct {
	TCPAddr string // listen ip
	UDPAddr string
	RFC3164 string // port num for syslog format RFC3164
	RFC5424 string //
	Unknown string
	Filters []string
}

// AppConfigInstance ...
var AppConfigInstance *AppConfig

func init() {
	AppConfigInstance = &AppConfig{
		LogCfg: &LogConfig{},
		SrvCfg: &ServerConfig{Filters: make([]string, 0)},
	}
}
