package models

type SystemUserRole struct {
	Id           int `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	SystemUserId int `json:"system_user_id" xorm:"not null comment('用户主键') index(system_user_id) INT(11)"`
	SystemRoleId int `json:"system_role_id" xorm:"not null comment('角色主键') index(system_user_id) INT(11)"`
}

var t_system_user_role = "system_user_role"

func (u *SystemUserRole) GetRow() bool {
	has, err := mysqlClt.Get(u)
	if err == nil && has {
		return true
	}
	return false
}
func (u *SystemUserRole) GetRowByUid() ([]string, error) {
	var role []string
	err := mysqlClt.Table(t_system_user_role).Select(t_system_role+".name").
		Join("INNER", t_system_role, t_system_user_role+".system_role_id = "+t_system_role+".id").
		Where(t_system_role+".status = ?", 1).
		Where(t_system_user_role+".system_user_id = ?", u.SystemUserId).
		Find(&role)
	return role, err
}
func (u *SystemUserRole) Add() (int64, error) {
	session := mysqlClt.NewSession()
	// add Begin() before any action
	if err := session.Begin(); err != nil {
		return 500, err
	}
	//var uid int64
	uid, err := session.Insert(u)
	if err != nil {
		return 500, err
	}
	session.Commit()
	return uid, err
}
