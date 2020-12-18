package models

import (
	"mygin_websrv/public/common"
	"time"
)

type SystemUser struct {
	Id            int       `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	Name          string    `json:"name" xorm:"not null comment('姓名') VARCHAR(50)"`
	Nickname      string    `json:"nickname" xorm:"not null default '' comment('用户登录名') unique VARCHAR(50)"`
	Password      string    `json:"password" xorm:"not null comment('密码') index VARCHAR(50)"`
	Salt          string    `json:"salt" xorm:"not null comment('盐') VARCHAR(4)"`
	Phone         string    `json:"phone" xorm:"not null default '' comment('手机号') VARCHAR(11)"`
	Avatar        string    `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(300)"`
	Introduction  string    `json:"introduction" xorm:"not null default '' comment('简介') VARCHAR(300)"`
	Status        int       `json:"status" xorm:"not null default 1 comment('状态（0 停止1启动）') TINYINT(4)"`
	Utime         time.Time `json:"utime" xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	LastLoginTime time.Time `json:"last_login_time" xorm:"not null default '0000-00-00 00:00:00' comment('上次登录时间') DATETIME"`
	LastLoginIp   string    `json:"last_login_ip" xorm:"not null default '' comment('最近登录IP') VARCHAR(50)"`
	Ctime         time.Time `json:"ctime" xorm:"not null comment('注册时间') DATETIME"`
}

type SearchUser struct {
	Name string `json:"name" xorm:"not null comment('姓名') VARCHAR(50)"`
}

var t_system_user = "system_role"

func (u *SystemUser) UserExist() bool {
	exist, err := mysqlClt.Get(u)
	if err == nil && exist {
		return true
	}
	return false
}

func (u *SystemUser) SelectAll() ([]SystemUser, error) {
	var sysUsers []SystemUser
	err := mysqlClt.Find(&sysUsers)
	return sysUsers, err
}

func (u *SystemUser) SelectByName(name string) ([]SearchUser, error) {
	var sysUsers []SearchUser
	err := mysqlClt.Table(sysUsers).Where("name like ?", name+"%").Find(&sysUsers)
	return sysUsers, err
}

func (u *SystemUser) SelectByPage(paging *common.Paging) ([]SystemUser, error) {
	var sysUsers []SystemUser
	var err error
	paging.Total, err = mysqlClt.Where("status = ?", 1).Count(u)
	paging.GetPages()
	if paging.Total < 1 {
		return sysUsers, err
	}
	err = mysqlClt.Where("status = ?", 1).Limit(int(paging.PageSize), int(paging.StartNums)).Find(&sysUsers)
	return sysUsers, err
}

func (u *SystemUser) Insert(roles []interface{}) (int, error) {
	session := mysqlClt.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}
	if _, err := session.Insert(u); err != nil {
		return 0, err
	}
	//如果没有设置权限,直接返回并提交数据库事务
	if len(roles) < 1 {
		return u.Id, session.Commit()
	}
	for _, k := range roles {
		roleModel := SystemUser{Name: k.(string)}
		if exist := roleModel.UserExist(); exist {
			continue
		}
		if roleModel.Status == 0 {
			continue
		}
		userRoleModel := SystemUserRole{SystemRoleId: roleModel.Id, SystemUserId: u.Id}
		exist, err := session.Get(&userRoleModel)
		if err != nil {
			return 0, err
		}
		if exist {
			continue
		}
		_, err = session.Insert(&userRoleModel)
		if err != nil {
			return 0, err
		}
	}
	return u.Id, session.Commit()
}

func (u *SystemUser) Update(roles []interface{}) error {
	session := mysqlClt.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := mysqlClt.Where("id = ?", u.Id).Update(u); err != nil {
		return err
	}
	roleModel := SystemUserRole{}
	if _, err := mysqlClt.Where("system_user_id=?", u.Id).Delete(&roleModel); err != nil {
		return err
	}
	//如果没有设置权限
	if len(roles) < 1 {
		return session.Commit()
	}
	for _, k := range roles {
		roleModel := SystemRole{Name: k.(string)}
		has := roleModel.GetRow()
		if !has {
			continue
		}
		if roleModel.Status == 0 {
			continue
		}
		userRoleModel := SystemUserRole{SystemRoleId: roleModel.Id, SystemUserId: u.Id}
		has, err := session.Get(&userRoleModel)
		if err != nil {
			return err
		}
		if has {
			continue
		}
		_, err = session.Insert(&userRoleModel)
		if err != nil {
			return err
		}
	}
	return session.Commit()
}

func (u *SystemUser) UpdatePasswd() error {
	if _, err := mysqlClt.Where("id = ?", u.Id).Update(u); err != nil {
		return err
	}
	return nil
}

func (u *SystemUser) Delete() error {
	if _, err := mysqlClt.Exec("update "+t_system_role+" set status=? where id=?", 0, u.Id); err != nil {
		return err
	}
	return nil
}
