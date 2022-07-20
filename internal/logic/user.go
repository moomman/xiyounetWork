package logic

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"ttms/internal/middleware"
	"ttms/internal/pkg/token"

	"github.com/go-redis/redis/v8"

	"ttms/internal/dao"
	db "ttms/internal/dao/db/sqlc"
	"ttms/internal/global"
	"ttms/internal/model/reply"
	"ttms/internal/model/request"
	"ttms/internal/pkg/app/errcode"
	"ttms/internal/pkg/password"
	"ttms/internal/pkg/utils"
	email2 "ttms/internal/worker/email"

	"github.com/jackc/pgx/v4"

	"github.com/gin-gonic/gin"
)

type user struct{}

type tokenResult struct {
	token   string
	payload *token.Payload
	err     error
}

func createToken(resultChan chan<- tokenResult, userID int64, userName string, expireTime time.Duration) func() {
	return func() {
		defer close(resultChan)
		accessToken, pal, err := global.Maker.CreateToken(userID, userName, expireTime)
		resultChan <- tokenResult{
			token:   accessToken,
			payload: pal,
			err:     err,
		}
	}
}

func (user) Login(c *gin.Context, params *request.Login) (reply.Login, errcode.Err) {
	user, err := dao.Group.DB.GetUserByName(c, params.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Info(err.Error())
			return reply.Login{}, errcode.ErrNotFound
		}
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Login{}, errcode.ErrServer
	}

	if err := password.CheckPassword(params.Password, user.Password); err != nil {
		global.Logger.Info(err.Error())
		return reply.Login{}, errcode.ErrLogin
	}
	accessChan := make(chan tokenResult, 1)
	refreshChan := make(chan tokenResult, 1)
	// 短token
	global.Worker.SendTask(createToken(accessChan, user.ID, user.Username, global.Settings.Token.AssessTokenDuration))
	// 长token
	global.Worker.SendTask(createToken(refreshChan, user.ID, user.Username, global.Settings.Token.RefreshTokenDuration))
	accessResult := <-accessChan
	if err := accessResult.err; err != nil {
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Login{}, errcode.ErrServer
	}
	refreshResult := <-refreshChan
	if err := refreshResult.err; err != nil {
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Login{}, errcode.ErrServer
	}
	return reply.Login{
		Username:     user.Username,
		Avatar:       user.Avatar,
		UserId:       user.ID,
		Privilege:    user.Privilege,
		AccessToken:  accessResult.token,
		RefreshToken: refreshResult.token,
		PalLoad:      accessResult.payload,
	}, nil
}

func (user) IsRepeat(c *gin.Context, params *request.IsRepeat) errcode.Err {
	user, err := dao.Group.DB.GetUserByName(c, params.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
			return errcode.ErrServer
		}
		return nil
	}

	if user != nil {
		return errcode.ErrNameHasExist
	}
	return nil
}

func (user) Register(c *gin.Context, params *request.Register) (reply.Register, errcode.Err) {
	_, err := dao.Group.DB.GetUserByName(c, params.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
			return reply.Register{}, errcode.ErrServer
		}
	} else {
		return reply.Register{}, errcode.ErrNameHasExist
	}

	if isVerify := email2.Check(params.Email, params.VerifyCode); !isVerify {
		return reply.Register{}, errcode.ErrOutTimeVerify
	}

	actor := "用户"
	var inviteCode string
	if params.InviteCode != "" {
		err = dao.Group.Redis.Get(c, params.InviteCode, &inviteCode)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return reply.Register{}, errcode.ErrOutTimeInvite
			}
			return reply.Register{}, errcode.ErrRedis
		}
		if inviteCode == "" || len(inviteCode) == 0 {
			return reply.Register{}, errcode.ErrOutTimeInvite
		} else {
			actor = inviteCode
		}
	}

	// 将明文加密
	hashPassword, err := password.HashPassword(params.Password)
	if err != nil {
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Register{}, errcode.ErrServer
	}
	userParams := &db.CreateUserParams{
		Username:  params.Username,
		Password:  hashPassword,
		Avatar:    global.Settings.Rule.DefaultCoverURL,
		Email:     params.Email,
		Signature: "日常摆烂",
		Privilege: db.Privilege(actor),
	}
	num, err := dao.Group.DB.CheckUserRepeat(c, &db.CheckUserRepeatParams{
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {

			global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
			return reply.Register{}, errcode.ErrServer
		}
	}
	if num != 0 {
		return reply.Register{}, errcode.ErrNameOrEmailExist
	}

	user, err := dao.Group.DB.CreateUser(c, userParams)
	if err != nil {
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Register{}, errcode.ErrServer
	}

	// 短token
	accessToken, pal, err := global.Maker.CreateToken(user.ID, user.Username, global.Settings.Token.AssessTokenDuration)
	if err != nil {
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Register{}, errcode.ErrServer
	}
	refreshToken, _, err := global.Maker.CreateToken(user.ID, user.Username, global.Settings.Token.RefreshTokenDuration)
	if err != nil {
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return reply.Register{}, nil
	}
	return reply.Register{
		UserId:       user.ID,
		PalLoad:      pal,
		Privilege:    actor,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (user) UpdateUserAvatar(c *gin.Context, params *request.UpdateUserAvatar) errcode.Err {
	avatarParams := &db.UpdateUserAvatarParams{
		Avatar: params.NewAvatar,
		ID:     params.UserId,
	}
	if err := dao.Group.DB.UpdateUserAvatar(c, avatarParams); err != nil {
		return errcode.ErrServer
	}
	return nil
}

func (user) UpdateUserInfo(c *gin.Context, params *request.UpdateUserInfo) errcode.Err {

	userParams := db.UpdateUserParams{
		Username: params.Username,

		Email:     params.Email,
		Birthday:  time.Unix(int64(params.Birthday), 0),
		Gender:    db.Gender(params.Gender),
		Signature: params.Signature,
		Hobby:     sql.NullString{},
		Lifestate: db.Lifestate(params.LifeState),
		ID:        params.UserId,
	}

	var hobby string
	// 将hobby当做string存储，拿出的时候按照[]string取出
	for i := range params.Hobbys {
		hobby += params.Hobbys[i] + " "
	}
	userParams.Hobby.String = hobby
	if params.Hobbys == nil || len(params.Hobbys) < 1 {
		userParams.Hobby.Valid = false
	} else {
		userParams.Hobby.Valid = true
	}

	err := dao.Group.DB.UpdateUser(c, &userParams)
	if err != nil {

		return errcode.ErrServer
	}
	return nil
}

func (user) ModifyPassword(c *gin.Context, params *request.ModifyPassword) errcode.Err {
	if isVerify := email2.Check(params.Email, params.VerifyCode); !isVerify {
		return errcode.NewErr(3011, "验证码错误")
	}
	hashPassword, err := password.HashPassword(params.NewPassword)
	if err != nil {
		return errcode.ErrServer
	}
	err1 := dao.Group.DB.UpdatePassword(c, &db.UpdatePasswordParams{
		Password: hashPassword,
		Email:    params.Email,
	})
	if err1 != nil {
		return errcode.ErrServer
	}

	return nil
}

func (user) GetUsers(c *gin.Context) (errcode.Err, reply.GetUsers) {
	users, err := dao.Group.DB.GetUsers(c)
	if err != nil {
		return errcode.NewErr(3022, "在查询users列表时出错"), reply.GetUsers{}
	}

	newUsers := make([]*reply.User, len(users))
	for i := range users {
		u := &reply.User{
			ID:        users[i].ID,
			Username:  users[i].Username,
			Avatar:    users[i].Avatar,
			Lifestate: users[i].Lifestate,
			Hobby:     strings.Split(users[i].Hobby.String, " "),
			Email:     users[i].Email,
			Birthday:  users[i].Birthday,
			Gender:    users[i].Gender,
			Signature: users[i].Signature,
			Privilege: users[i].Privilege,
		}
		newUsers = append(newUsers, u)
	}
	return nil, reply.GetUsers{
		Users: newUsers,
	}
}

func (user) GetUserInfo(c *gin.Context, params *request.GetUserInfo) (errcode.Err, reply.GetUserInfo) {
	users, err := dao.Group.DB.GetUserById(c, params.UserId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows){
			global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
			return errcode.ErrServer, reply.GetUserInfo{}
		}
	}
	hobby := make([]string, 0)
	if users.Hobby.String != "" {
		hobby = strings.Split(users.Hobby.String, " ")
	}

	u := &reply.User{
		ID:        users.ID,
		Username:  users.Username,
		Avatar:    users.Avatar,
		Lifestate: users.Lifestate,
		Hobby:     hobby,
		Email:     users.Email,
		Birthday:  users.Birthday,
		Gender:    users.Gender,
		Signature: users.Signature,
		Privilege: users.Privilege,
	}

	return nil, reply.GetUserInfo{
		User: u,
	}
}

func (user) Generate(c *gin.Context, params *request.Generate) errcode.Err {
	randomString := utils.RandomString(6)
	err := dao.Group.Redis.SetTimeOut(c, randomString, params.GivedRight, global.Settings.Rule.InviteCodeTime)
	if err != nil {
		return errcode.ErrRedis
	}

	return nil
}

func (user) Refresh(c *gin.Context, params *request.RefreshParams) (errcode.Err, *reply.RefreshParams) {
	accessPayLoad, err1 := global.Maker.VerifyToken(params.AccessToken)
	if err1 != nil {
		return errcode.ErrUnauthorizedAuthNotExist, nil
	}
	RefreshPayLoad, err2 := global.Maker.VerifyToken(params.RefreshToken)
	if err2 != nil {
		return errcode.ErrUnauthorizedAuthNotExist, nil
	}

	if RefreshPayLoad.ExpiredAt.Unix() < time.Now().Unix() {
		return errcode.ErrOutTimeRefreshToken, nil
	}

	token, _, err := global.Maker.CreateToken(accessPayLoad.UserID, accessPayLoad.UserName, global.Settings.Token.AssessTokenDuration)
	if err != nil {
		return errcode.ErrGenerateToken, nil
	}

	return nil, &reply.RefreshParams{
		NewAccessToken: token,
	}

}

func (user) DeleteUser(c *gin.Context, params *request.DeleteUser) errcode.Err {
	err := dao.Group.DB.DeleteUser(c, params.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errcode.ErrUserNotExist
		}
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return errcode.ErrServer
	}

	return nil
}

func (user) ListUserInfo(c *gin.Context, params *request.ListUserInfo) (errcode.Err, *reply.ListUserInfo) {
	info, err := dao.Group.DB.ListUserInfo(c, &db.ListUserInfoParams{
		Limit:  params.PageSize,
		Offset: params.Page,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errcode.ErrUserNotExist, nil
		}
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return errcode.ErrServer, nil
	}
	num, err := dao.Group.DB.ListNum(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errcode.ErrUserNotExist, nil
		}
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return errcode.ErrServer, nil
	}
	return nil, &reply.ListUserInfo{
		UserInfos: info,
		Total:     num,
	}
}

func (user) SearchUser(c *gin.Context, params *request.SearchUser) (errcode.Err, *reply.SearchUser) {
	data, err := dao.Group.DB.SearchUserByName(c, &db.SearchUserByNameParams{
		Limit:    params.PageSize,
		Offset:   params.Page,
		Username: params.Username,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errcode.ErrUserNotExist, nil
		}
		global.Logger.Error(err.Error(), middleware.ErrLogMsg(c)...)
		return errcode.ErrServer, nil
	}
	num, err := dao.Group.DB.ListNameNum(c, params.Username)
	if err != nil {
		global.Logger.Error(err.Error())
		return errcode.ErrServer, nil
	}
	return nil, &reply.SearchUser{
		UserInfos: data,
		Total:     num,
	}

}
