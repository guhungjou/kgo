package kindergarten

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	systemdb "gitlab.com/ykgk/kgo/db/system"
	"gitlab.com/ykgk/kgo/x"
)

func GetKindergarten(id int64) (*Kindergarten, error) {
	kindergarten := Kindergarten{ID: id}
	if err := db.PG().Model(&kindergarten).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &kindergarten, nil
}

/* 分页获取幼儿园列表 */
func FindKindergartens(query, districtID string, page int, pageSize int) ([]*Kindergarten, int, error) {
	kindergartens := make([]*Kindergarten, 0)

	q := db.PG().Model(&kindergartens).Relation(`District`).
		Relation(`District.Parent`).Relation(`District.Parent.Parent`)
	if query != "" {
		qry := fmt.Sprintf("%%%s%%", query)
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`"name" ILIKE ?`, qry)
			q = q.WhereOr(`"remark" ILIKE ?`, qry)
			return q, nil
		})
	}
	if districtID != "" {
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`"district"."id"=?`, districtID)
			q = q.WhereOr(`"district__parent"."id"=?`, districtID)
			q = q.WhereOr(`"district__parent__parent"."id"=?`, districtID)
			return q, nil
		})
	}

	total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).Order(`id DESC`).SelectAndCountEstimate(100000)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	}
	/* 查询幼儿园的园长 */
	for _, kindergarten := range kindergartens {
		manager := KindergartenTeacher{}
		if err := db.PG().Model(&manager).Where(`"kindergarten_id"=?`, kindergarten.ID).Where(`"role"=?`, KindergartenTeacherRoleManager).Limit(1).Select(); err != nil {
			if err != db.ErrNoRows {
				log.Errorf("DB Error: %v", err)
				return nil, 0, err
			}
		} else {
			kindergarten.Manager = &manager
		}
	}
	return kindergartens, total, nil
}

/* 创建幼儿园，同时创建幼儿园园长帐号 */
func CreateKindergarten(districtID, name, remark, managerUsername, managerPassword, managerName, managerGender, managerPhone string, createdBy int64) (*Kindergarten, *KindergartenTeacher, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, nil, err
	}
	defer tx.Rollback()

	district := systemdb.District{ID: districtID}
	if err := tx.Model(&district).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, nil, err
	}

	kg := Kindergarten{
		DistrictID:      district.ID,
		Name:            name,
		Remark:          remark,
		NumberOfStudent: 0,
		NumberOfTeacher: 1,
		CreatedBy:       createdBy,
	}
	if _, err := tx.Model(&kg).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, nil, err
	}

	teacher := KindergartenTeacher{
		Name:     managerName,
		Gender:   managerGender,
		Phone:    managerPhone,
		Username: managerUsername,

		Role:           KindergartenTeacherRoleManager,
		KindergartenID: kg.ID,

		Salt:  x.RandomString(12),
		PType: randomPType(),
	}

	teacher.Password = encrypt(teacher.Salt, managerPassword, teacher.PType)

	if _, err := tx.Model(&teacher).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, nil, err
	}

	return &kg, &teacher, nil
}

/* 更新幼儿园基本信息 */
func UpdateKindergarten(id int64, districtID, name, remark string, createdBy int64) (*Kindergarten, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	district := systemdb.District{ID: districtID}
	if err := tx.Model(&district).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	kg := Kindergarten{ID: id}
	if err := tx.Model(&kg).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(&kg).WherePK().Set(`"district_id"=?`, district.ID).
		Set(`"name"=?`, name).Set(`"remark"=?`, remark).Set(`"updated_at"=CURRENT_TIMESTAMP`).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &kg, nil
}
