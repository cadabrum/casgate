// Copyright 2021 The Casdoor Authors. All Rights Reserved.
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

package controllers

import (
	"encoding/json"

	"github.com/beego/beego/utils/pagination"
	"github.com/casdoor/casdoor/object"
)

// GetPermissions
// @Title GetPermissions
// @Tag Permission API
// @Description get permissions
// @Param   owner     query    string  true        "The owner of permissions"
// @Success 200 {array} object.Permission The Response object
// @Failure 500 Internal server error
// @Failure 401 Unauthorized
// @router /get-permissions [get]
func (c *ApiController) GetPermissions() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)
	
	count, err := object.GetPermissionCount(request.Owner, request.Field, request.Value)
	if err != nil {
		c.ResponseInternalServerError(err.Error())
		return
	}

	paginator := pagination.SetPaginator(c.Ctx, request.Limit, count)
	permissions, err := object.GetPaginationPermissions(request.Owner, paginator.Offset(), request.Limit, request.Field, request.Value, request.SortField, request.SortOrder)
	if err != nil {
		c.ResponseInternalServerError(err.Error())
		return
	}

	c.ResponseOk(permissions, paginator.Nums())
}

// GetPermissionsBySubmitter
// @Title GetPermissionsBySubmitter
// @Tag Permission API
// @Description get permissions by submitter
// @Success 200 {array} object.Permission The Response object
// @Failure 500 Internal server error
// @Failure 401 Unauthorized
// @router /get-permissions-by-submitter [get]
func (c *ApiController) GetPermissionsBySubmitter() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)

	user, ok := c.RequireSignedInUser()
	if !ok {
		c.ResponseUnauthorized(c.T("auth:Unauthorized operation"))
		return
	}

	permissions, err := object.GetPermissionsBySubmitter(user.Owner, user.Name)
	if err != nil {
		c.ResponseInternalServerError(err.Error())
		return
	}

	c.ResponseOk(permissions, len(permissions))
}

// GetPermissionsByRole
// @Title GetPermissionsByRole
// @Tag Permission API
// @Description get permissions by role
// @Param   id     query    string  true        "The id ( owner/name ) of the role"
// @Success 200 {array} object.Permission The Response object
// @Failure 500 Internal server error
// @router /get-permissions-by-role [get]
func (c *ApiController) GetPermissionsByRole() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)

	roleFromDb, _ := object.GetRole(request.Id)
	if roleFromDb == nil {
		c.ResponseOk()
		return
	}
	c.ValidateOrganization(roleFromDb.Owner)

	permissions, err := object.GetPermissionsByRole(request.Id)
	if err != nil {
		c.ResponseInternalServerError(err.Error())
		return
	}

	c.ResponseOk(permissions, len(permissions))
}

// GetPermission
// @Title GetPermission
// @Tag Permission API
// @Description get permission
// @Param   id     query    string  true        "The id ( owner/name ) of the permission"
// @Success 200 {object} object.Permission The Response object
// @Failure 500 Internal server error
// @router /get-permission [get]
func (c *ApiController) GetPermission() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)

	permission, err := object.GetPermission(request.Id)
	if err != nil {
		c.ResponseInternalServerError(err.Error())
		return
	}
	if permission == nil {
		c.ResponseOk()
		return
	}
	c.ValidateOrganization(permission.Owner)

	c.ResponseOk(permission)
}

// UpdatePermission
// @Title UpdatePermission
// @Tag Permission API
// @Description update permission
// @Param   id     query    string  true        "The id ( owner/name ) of the permission"
// @Param   body    body   object.Permission  true        "The details of the permission"
// @Success 200 {object} controllers.Response The Response object
// @Failure 400 Bad request
// @router /update-permission [post]
func (c *ApiController) UpdatePermission() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)

	var permission object.Permission
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &permission)
	if err != nil {
		c.ResponseBadRequest(err.Error())
		return
	}
	permFromDb, _ := object.GetPermission(request.Id)
	if permFromDb == nil {
		c.Data["json"] = wrapActionResponse(false)
		c.ServeJSON()
		return
	}
	c.ValidateOrganization(permFromDb.Owner)

	c.Data["json"] = wrapActionResponse(object.UpdatePermission(request.Id, &permission))
	c.ServeJSON()
}

// AddPermission
// @Title AddPermission
// @Tag Permission API
// @Description add permission
// @Param   body    body   object.Permission  true        "The details of the permission"
// @Success 200 {object} controllers.Response The Response object
// @Failure 400 Bad request
// @router /add-permission [post]
func (c *ApiController) AddPermission() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)

	var permission object.Permission
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &permission)
	if err != nil {
		c.ResponseBadRequest(err.Error())
		return
	}
	c.ValidateOrganization(permission.Owner)

	c.Data["json"] = wrapActionResponse(object.AddPermission(&permission))
	c.ServeJSON()
}

// DeletePermission
// @Title DeletePermission
// @Tag Permission API
// @Description delete permission
// @Param   body    body   object.Permission  true        "The details of the permission"
// @Success 200 {object} controllers.Response The Response object
// @Failure 400 Bad request
// @router /delete-permission [post]
func (c *ApiController) DeletePermission() {
	request := c.ReadRequestFromQueryParams()
	c.ContinueIfHasRightsOrDenyRequest(request)
	
	var permission object.Permission
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &permission)
	if err != nil {
		c.ResponseBadRequest(err.Error())
		return
	}

	permFromDb, _ := object.GetPermission(permission.GetId())
	if permFromDb == nil {
		c.Data["json"] = wrapActionResponse(false)
		c.ServeJSON()
		return
	}
	c.ValidateOrganization(permFromDb.Owner)

	c.Data["json"] = wrapActionResponse(object.DeletePermission(&permission))
	c.ServeJSON()
}
