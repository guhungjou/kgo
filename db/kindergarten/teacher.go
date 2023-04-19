package kindergarten

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	"gitlab.com/ykgk/kgo/x"
)

func GetKindergartenTeacher(id int64) (*KindergartenTeacher, error) {
	teacher := KindergartenTeacher{ID: id}
	if err := db.PG().Model(&teacher).Relation(`Kindergarten`).Relation(`Class`).WherePK().Where(`NOT "kindergarten_teacher"."deleted"`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &teacher, nil
}

func GetKindergartenTeacherByUsername(username string) (*KindergartenTeacher, error) {
	teacher := KindergartenTeacher{}
	if err := db.PG().Model(&teacher).Relation(`Kindergarten`).Relation(`Class`).Where(`"kindergarten_teacher"."username"=?`, username).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &teacher, nil
}

func FindKindergartenTeachers(query string, kindergartenID, classID int64, role string, page, pageSize int) ([]*KindergartenTeacher, int, error) {
	teachers := make([]*KindergartenTeacher, 0)

	q := db.PG().Model(&teachers).Relation(`Class`).Relation(`Kindergarten`).Where(`NOT "kindergarten_teacher"."deleted"`)

	if query != "" {
		qry := fmt.Sprintf("%%%s%%", query)
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`"kindergarten_teacher"."name" ILIKE ?`, qry)
			q = q.WhereOr(`"kindergarten_teacher"."phone" ILIKE ?`, qry)
			return q, nil
		})
	}

	if kindergartenID > 0 {
		q = q.Where(`"kindergarten_teacher"."kindergarten_id"=?`, kindergartenID)
	}

	if classID != 0 {
		q = q.Where(`"kindergarten_teacher"."class_id"=?`, classID)
	}

	if role != "" {
		q = q.Where(`"kindergarten_teacher"."role"=?`, role)
	}

	total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).Order(`kindergarten_teacher.id DESC`).SelectAndCount()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	}
	return teachers, total, nil
}

func UpdateKindergartenTeacherPassword(id int64, password string) (*KindergartenTeacher, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	teacher := KindergartenTeacher{ID: id}
	if err := tx.Model(&teacher).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	teacher.Salt = x.RandomString(12)
	teacher.PType = randomPType()
	teacher.Password = encrypt(teacher.Salt, password, teacher.PType)

	if _, err := tx.Model(&teacher).WherePK().Set(`"salt"=?salt`).Set(`"ptype"=?ptype`).
		Set(`"password"=?password`).Set(`"updated_at"=CURRENT_TIMESTAMP`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &teacher, nil
}

func UpdateKindergartenTeacher(id int64, name, gender, phone string, classID int64) (*KindergartenTeacher, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	teacher := KindergartenTeacher{ID: id}
	if err := tx.Model(&teacher).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if classID > 0 {
		class := KindergartenClass{ID: classID}
		if err := tx.Model(&class).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}

		if class.KindergartenID != teacher.KindergartenID {
			log.Errorf("班级和老师不在同一个幼儿园: %d %d", teacher.ID, class.ID)
			return nil, fmt.Errorf("班级和老师不在同一个幼儿园: %d %d", teacher.ID, class.ID)
		}
	}

	if classID != teacher.ClassID {
		/* 班级发生变化，新班级老师数量+1，旧班级老师数量-1 */
		if classID > 0 {
			if _, err := tx.Model(&KindergartenClass{ID: classID}).WherePK().Set(`"number_of_teacher"="number_of_teacher"+1`).Update(); err != nil {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		}
		if teacher.ClassID > 0 {
			if _, err := tx.Model(&KindergartenClass{ID: teacher.ClassID}).WherePK().Set(`"number_of_teacher"="number_of_teacher"-1`).Update(); err != nil {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		}
	}

	if _, err := tx.Model(&teacher).WherePK().Set(`"name"=?`, name).Set(`"gender"=?`, gender).
		Set(`"phone"=?`, phone).Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"class_id"=?`, classID).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &teacher, nil
}

func CreateKindergartenTeacherTx(tx *pg.Tx, username, password, name, gender, phone, role string, kindergartenID, classID int64) (*KindergartenTeacher, error) {
	kindergarten := Kindergarten{ID: kindergartenID}
	if err := tx.Model(&kindergarten).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	/* 幼儿园的老师数量 +1 */
	if _, err := tx.Model(&kindergarten).WherePK().Set(`"number_of_teacher"="number_of_teacher"+1`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if classID > 0 {
		class := KindergartenClass{ID: classID}
		if err := tx.Model(&class).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}

		/* 班级的老师数量+1 */
		if _, err := tx.Model(&class).WherePK().Set(`"number_of_teacher"="number_of_teacher"+1`).Update(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	}

	teacher := KindergartenTeacher{
		Name:     name,
		Gender:   gender,
		Phone:    phone,
		Username: username,

		Role:           role,
		KindergartenID: kindergartenID,
		ClassID:        classID,

		Salt:  x.RandomString(12),
		PType: randomPType(),
	}

	teacher.Password = encrypt(teacher.Salt, password, teacher.PType)

	if _, err := tx.Model(&teacher).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &teacher, nil
}

/* 创建老师帐号 */
func CreateKindergartenTeacher(username, password, name, gender, phone, role string, kindergartenID, classID int64) (*KindergartenTeacher, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	teacher, err := CreateKindergartenTeacherTx(tx, username, password, name, gender, phone, role, kindergartenID, classID)

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return teacher, nil
}

func DeleteKindergartenTeacherTx(tx *pg.Tx, id int64) (*KindergartenTeacher, error) {
	teacher := KindergartenTeacher{ID: id}

	if err := tx.Model(&teacher).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}

	/* 不能删除园长帐号 */
	if teacher.Role == KindergartenTeacherRoleManager {
		return &teacher, nil
	}

	if _, err := tx.Model(&teacher).WherePK().Set(`"deleted"=?`, true).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	if _, err := tx.Model(&KindergartenClass{ID: teacher.ClassID}).WherePK().Set(`"number_of_teacher"="number_of_teacher"-1`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	if _, err := tx.Model(&Kindergarten{ID: teacher.KindergartenID}).WherePK().Set(`"number_of_teacher"="number_of_teacher"-1`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return &teacher, nil
}

func DeleteKindergartenTeacher(id int64) (*KindergartenTeacher, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	teacher, err := DeleteKindergartenTeacherTx(tx, id)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return teacher, nil
}

/* 删除班级下的老师 */
func DeleteKindergartenTeacherByClassTx(tx *pg.Tx, classID int64) error {
	if r, err := tx.Exec(`UPDATE "kindergarten_teacher" SET "deleted"=? WHERE NOT "deleted" AND "class_id"=?`, true, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	} else if r.RowsAffected() > 0 {
		if _, err := tx.Exec(`UPDATE "kindergarten_class" SET "number_of_teacher"="number_of_teacher"-? WHERE "id"=?`, r.RowsAffected(), classID); err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}

		if _, err := tx.Exec(`UPDATE "kindergarten" SET "number_of_teacher"="number_of_teacher"-? WHERE "id"=(SELECT "kindergarten_id" FROM "kindergarten_class" WHERE "id"=?)`, r.RowsAffected(), classID); err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}

	}
	return nil
}

/* 将老师的班级设置成空 */
func ClearKindergartenTeacherClassTx(tx *pg.Tx, classID int64) error {
	if r, err := tx.Exec(`UPDATE "kindergarten_teacher" SET "class_id"=0 WHERE NOT "deleted" AND "class_id"=?`, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	} else if r.RowsAffected() > 0 {
		if _, err := tx.Exec(`UPDATE "kindergarten_class" SET "number_of_teacher"="number_of_teacher"-? WHERE "id"=?`, r.RowsAffected(), classID); err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}
	}
	return nil
}
