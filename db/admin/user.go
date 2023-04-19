package admin

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	"gitlab.com/ykgk/kgo/x"
)

func GetAdminUserByUsername(username string) (*AdminUser, error) {
	user := AdminUser{}
	if err := db.PG().Model(&user).Where(`"username"=?`, username).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

func GetAdminUser(id int64) (*AdminUser, error) {
	pg := db.PG()

	user := AdminUser{ID: id}
	if err := pg.Model(&user).WherePK().Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

func GetSuperAdminUser() (*AdminUser, error) {
	user := AdminUser{}
	if err := db.PG().Model(&user).Where(`"is_superuser"=?`, true).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

/* 创建超级管理员，只有在没有管理员帐号时才能创建 */
func CreateSuperAdminUser(username, name, phone, password string) (*AdminUser, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	user := AdminUser{}
	if err := tx.Model(&user).Where(`"is_superuser"=?`, true).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	} else {
		/* 超级用户已经存在 */
		return &user, nil
	}

	user = AdminUser{
		Username:    username,
		Name:        name,
		Phone:       phone,
		IsSuperuser: true,
		Status:      AdminUserStatusOK,
		CreatedBy:   0,

		PType: randomPType(),
		Salt:  x.RandomString(12),
	}

	user.Password = encrypt(user.Salt, password, user.PType)

	if _, err := tx.Model(&user).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &user, nil
}

func CreateAdminUser(username, name, phone, password string, createdBy int64) (*AdminUser, error) {
	user := AdminUser{
		Username:    username,
		Name:        name,
		Phone:       phone,
		IsSuperuser: false,
		Status:      AdminUserStatusOK,
		CreatedBy:   createdBy,

		PType: randomPType(),
		Salt:  x.RandomString(12),
	}

	user.Password = encrypt(user.Salt, password, user.PType)

	if _, err := db.PG().Model(&user).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return &user, nil
}

/* 更新管理账号的密码 */
func UpdateAdminUserPassword(id int64, newpassword string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	defer tx.Rollback()

	user := AdminUser{ID: id}
	if err := tx.Model(&user).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}

	user.PType = randomPType()
	user.Salt = x.RandomString(12)
	user.Password = encrypt(user.Salt, newpassword, user.PType)

	if _, err := tx.Model(&user).WherePK().Set(`"ptype"=?ptype`).Set(`"salt"=?salt`).Set(`"password"=?password`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

func FindAdminUsers(query string, status string, page, pageSize int) ([]*AdminUser, int, error) {
	users := make([]*AdminUser, 0)

	q := db.PG().Model(&users).Order(`id DESC`)
	if query != "" {
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			qry := fmt.Sprintf("%%%s%%", query)
			q = q.WhereOr(`"name" ILIKE ?`, qry)
			q = q.WhereOr(`"username" ILIKE ?`, qry)
			return q, nil
		})
	}
	if status != "" {
		q = q.Where(`"status"=?`, status)
	}

	total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).SelectAndCountEstimate(10000)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	}
	return users, total, nil
}

func UpdateAdminUser(id int64, name, status, phone string) (*AdminUser, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	user := AdminUser{ID: id}
	if err := tx.Model(&user).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	q := tx.Model(&user).WherePK().Set(`"name"=?`, name).Set(`"status"=?`, status).Set(`"phone"=?`, phone)

	if _, err := q.Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return &user, nil
}
