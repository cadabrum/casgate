// Copyright 2023 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"fmt"
	"regexp"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/xorm"
)

func GetSession(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := AppOrmer.Engine.Prepare()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if owner != "" {
		session = session.And("LOWER(owner) = LOWER(?)", owner)
	}

	match, mErr := regexp.MatchString("^[a-z_]+$", util.SnakeString(field))
	if !match || mErr != nil {
		field = ""
	}

	if field != "" {
		filterValue := ""
		if value != "" {
			filterValue = fmt.Sprintf("%%%s%%", value)
		}
		if util.IsFieldValueAllowedForDB(field) {
			session = session.And(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", util.SnakeString(field)), filterValue)
		}
	}

	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}

	match, mErr = regexp.MatchString("^[a-z_]+$", util.SnakeString(sortField))
	if !match || mErr != nil {
		sortField = "created_time"
	}

	if sortOrder == "ascend" {
		session = session.Asc(util.SnakeString(sortField))
	} else {
		session = session.Desc(util.SnakeString(sortField))
	}

	return session
}

func GetSessionForUser(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := AppOrmer.Engine.Prepare()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if owner != "" {
		if offset == -1 {
			session = session.And("LOWER(owner) = LOWER(?)", owner)
		} else {
			session = session.And("LOWER(a.owner) = LOWER(?)", owner)
		}
	}

	match, mErr := regexp.MatchString("^[a-z_]+$", util.SnakeString(field))
	if !match || mErr != nil {
		field = ""
	}

	if field != "" && value != "" {
		if util.IsFieldValueAllowedForDB(field) {
			if offset != -1 {
				field = fmt.Sprintf("a.%s", field)
			}
			session = session.And(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", util.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
		}
	}

	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}

	match, mErr = regexp.MatchString("^[a-z_]+$", util.SnakeString(sortField))
	if !match || mErr != nil {
		sortField = "created_time"
	}

	tableNamePrefix := conf.GetConfigString("tableNamePrefix")
	tableName := tableNamePrefix + "user"

	if offset == -1 {
		if sortOrder == "ascend" {
			session = session.Asc(util.SnakeString(sortField))
		} else {
			session = session.Desc(util.SnakeString(sortField))
		}
	} else {
		if sortOrder == "ascend" {
			session = session.Alias("a").
				Join("INNER", []string{tableName, "b"}, "LOWER(a.owner) = LOWER(b.owner) and LOWER(a.name) = LOWER(b.name)").
				Select("b.*").
				Asc("a." + util.SnakeString(sortField))
		} else {
			session = session.Alias("a").
				Join("INNER", []string{tableName, "b"}, "LOWER(a.owner) = LOWER(b.owner) and LOWER(a.name) = LOWER(b.name)").
				Select("b.*").
				Desc("a." + util.SnakeString(sortField))
		}
	}

	notUserAccesToken := builder.Not{builder.Like{"a.tag", "<access-token>"}}
	session = session.And(notUserAccesToken)

	return session
}
