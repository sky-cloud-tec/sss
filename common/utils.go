// NetD makes network device operations easy.
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

import (
	"github.com/songtianyi/rrframework/logs"
)

var (
	// MapStringToLevel trans string to logs level enum
	MapStringToLevel = map[string]int{
		"EMERGENCY": logs.LevelEmergency,
		"ALERT":     logs.LevelAlert,
		"CRITICAL":  logs.LevelCritical,
		"ERROR":     logs.LevelError,
		"WARNING":   logs.LevelWarning,
		"NOTICE":    logs.LevelNotice,
		"INFO":      logs.LevelInformational,
		"DEBUG":     logs.LevelDebug,
	}
)
