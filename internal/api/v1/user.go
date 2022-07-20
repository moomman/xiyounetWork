package v1

import (
	"ttms/internal/global"
	"ttms/internal/logic"
	"ttms/internal/model/common"
	"ttms/internal/model/request"
	"ttms/internal/pkg/app"
	"ttms/internal/pkg/app/errcode"

	"github.com/gin-gonic/gin"
)

type user struct {
}

// Register
// @Tags      user
// @Summary   用户注册
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.Register  true  "用户名，密码，邮箱,邀请码"
// @Success   200   {object}  common.State{data=reply.Register}
// @Router    /user/register [post]
func (user) Register(c *gin.Context) {

	rly := app.NewResponse(c)
	var registerPar request.Register
	if err := c.ShouldBindJSON(&registerPar); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	reply, err := logic.Group.User.Register(c, &registerPar)
	rly.Reply(err, reply)
}

//

// IsRepeat
// @Tags      user
// @Summary   验证用户输入的用户名是否存在
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.IsRepeat  true  "用户名"
// @Success   200   {object}  common.State{}
// @Router    /user/isRepeat [get]
func (user) IsRepeat(c *gin.Context) {
	rly := app.NewResponse(c)
	var params request.IsRepeat
	if err := c.ShouldBindQuery(&params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	err := logic.Group.User.IsRepeat(c, &params)
	if err != nil {
		rly.Reply(err, nil)
		return
	}

	rly.Reply(nil, nil)
}

// Login
// @Tags      user
// @Summary   用户登录
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.Login  true  "用户名，密码"
// @Success   200   {object}  common.State{data=reply.Login}
// @Router    /user/login [post]
func (user) Login(c *gin.Context) {
	rly := app.NewResponse(c)
	var loginPar request.Login
	if err := c.ShouldBindJSON(&loginPar); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	reply, err := logic.Group.User.Login(c, &loginPar)
	rly.Reply(err, reply)
}

// ModifyPassword
// @Tags      user
// @Summary   更改用户密码
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.ModifyPassword  true  "密码，邮箱"
// @Success   200   {object}  common.State{}
// @Router    /user/modifyPwd [put]
func (user) ModifyPassword(c *gin.Context) {
	rly := app.NewResponse(c)
	var params request.ModifyPassword
	if err := c.ShouldBindJSON(&params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	if err := logic.Group.User.ModifyPassword(c, &params); err != nil {
		rly.Reply(errcode.ErrServer, nil)
		return
	}

	rly.Reply(nil, nil)
}

// UpdateUserInfo
// @Tags      user
// @Summary   更新用户信息
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UpdateUserInfo  true  "用户id,用户名，邮箱，生日，性别，爱好，生活状态，个性签名"
// @Success   200   {object}  common.State{}
// @Router    /user/info/modify [put]
func (user) UpdateUserInfo(c *gin.Context) {
	rly := app.NewResponse(c)
	var params request.UpdateUserInfo
	if err := c.ShouldBindJSON(&params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}
	if err := params.Judge(); err != nil {
		rly.Reply(err, nil)
		return
	}

	if err := logic.Group.User.UpdateUserInfo(c, &params); err != nil {
		rly.Reply(errcode.ErrServer, nil)
		return
	}

	rly.Reply(nil, nil)
}

// UpdateUserAvatar
// @Tags      user
// @Summary   更新用户头像
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.UpdateUserAvatar  true  "用户id，新头像url"
// @Success   200   {object}  common.State{}
// @Router    /user/updateAvatar [put]
func (user) UpdateUserAvatar(c *gin.Context) {
	rly := app.NewResponse(c)
	var params request.UpdateUserAvatar
	if err := c.ShouldBindJSON(&params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	if err := logic.Group.User.UpdateUserAvatar(c, &params); err != nil {
		rly.Reply(errcode.ErrServer, nil)
		return
	}
	rly.Reply(nil, nil)
}

// GetUsers
// @Tags      user
// @Summary   获取所有用户
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  common.State{data=reply.GetUsers}
// @Router    /user/list [get]
func (user) GetUsers(c *gin.Context) {
	response := app.NewResponse(c)

	err, data := logic.Group.User.GetUsers(c)
	if err != nil {
		response.ReplyList(err, nil)
		return
	}
	response.ReplyList(err, data)

}

// GetUserInfo
// @Tags      user
// @Summary   获取用户信息
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetUserInfo  true  "用户id"
// @Success   200   {object}  common.State{data=reply.GetUserInfo}
// @Router    /user/get [get]
func (user) GetUserInfo(c *gin.Context) {
	response := app.NewResponse(c)
	var params request.GetUserInfo
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	err, data := logic.Group.User.GetUserInfo(c, &params)
	if err != nil {
		response.Reply(errcode.ErrServer, nil)
		return
	}

	response.Reply(nil, data)

}

// Generate
// @Tags      user
// @Summary   生成邀请码
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.Generate  true  "用户id,要给予的权限(这个写成管理员就行)"
// @Success   200   {object}  common.State{}
// @Router    /user/generate [post]
func (user) Generate(c *gin.Context) {
	response := app.NewResponse(c)
	var params request.Generate
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	if err := logic.Group.User.Generate(c, &params); err != nil {
		response.Reply(err, nil)
	}

	response.Reply(nil, nil)
}

// Refresh
// @Tags      user
// @Summary   刷新token
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.RefreshParams  true  "accessToken,refreshToken"
// @Success   200   {object}  common.State{data=reply.RefreshParams}
// @Router    /user/refresh [post]
func (user) Refresh(c *gin.Context) {
	response := app.NewResponse(c)
	var params request.RefreshParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	err, data := logic.Group.User.Refresh(c, &params)
	if err != nil {
		response.Reply(err, nil)
		return
	}

	response.Reply(nil, data)
}

// Delete
// @Tags      user
// @Summary   根据id删除用户
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.DeleteUser  true  "用户id"
// @Success   200   {object}  common.State{}
// @Router    /user/delete [post]
func (user) Delete(c *gin.Context) {
	response := app.NewResponse(c)
	var params request.DeleteUser
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}

	err := logic.Group.User.DeleteUser(c, &params)
	if err != nil {
		response.Reply(err, nil)
		return
	}

	response.Reply(nil, nil)
}

// ListUserInfo
// @Tags      user
// @Summary   展示所有用户
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.ListUserInfo  true  "分页"
// @Success   200   {object}  common.State{data=reply.ListUserInfo}
// @Router    /user/listInfo [get]
func (user) ListUserInfo(c *gin.Context) {
	response := app.NewResponse(c)
	var params request.ListUserInfo
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(c)

	err, data := logic.Group.User.ListUserInfo(c, &request.ListUserInfo{
		Pager: common.Pager{
			Page:     offset,
			PageSize: limit,
		},
	})
	if err != nil {
		response.Reply(err, nil)
		return
	}

	response.Reply(nil, data)
}

// SearchUser
// @Tags      user
// @Summary   搜索用户
// @Security  BasicAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.SearchUser  true  "分页"
// @Success   200   {object}  common.State{data=reply.ListUserInfo}
// @Router    /user/search [get]
func (user) SearchUser(c *gin.Context) {
	response := app.NewResponse(c)
	var params request.SearchUser
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()), nil)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(c)

	err, data := logic.Group.User.SearchUser(c, &request.SearchUser{
		Username: "%" + params.Username + "%",
		Pager: common.Pager{
			Page:     offset,
			PageSize: limit,
		},
	})
	if err != nil {
		response.Reply(err, nil)
		return
	}

	response.Reply(nil, data)
}
