package kindergarten

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	"gitlab.com/ykgk/kgo/x"
)

func GetKindergartenStudentByNOTx(tx *pg.Tx, classID int64, no string) (*KindergartenStudent, error) {
	student := KindergartenStudent{}
	if err := tx.Model(&student).Where(`NOT "deleted"`).Where(`"class_id"=?`, classID).Where(`"no"=?`, no).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &student, nil
}

func GetKindergartenStudentByName(classID int64, name string) (*KindergartenStudent, error) {
	student := KindergartenStudent{}
	if err := db.PG().Model(&student).Where(`NOT "deleted"`).Where(`"class_id"=?`, classID).Where(`"name"=?`, name).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &student, nil
}

func GetKindergartenStudentByNO(classID int64, no string) (*KindergartenStudent, error) {
	student := KindergartenStudent{}
	if err := db.PG().Model(&student).Where(`NOT "deleted"`).Where(`"class_id"=?`, classID).Where(`"no"=?`, no).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &student, nil
}

func GetKindergartenStudent(id int64) (*KindergartenStudent, error) {
	student := KindergartenStudent{ID: id}
	if err := db.PG().Model(&student).WherePK().Relation(`Kindergarten`).Relation(`Class`).
		Where(`NOT "kindergarten_student"."deleted"`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &student, nil
}

func GetKindergartenStudentByDevice(device string) (*KindergartenStudent, error) {
	device = x.FormatMacAddress(device)
	student := KindergartenStudent{}
	if err := db.PG().Model(&student).Relation(`Class`).Relation(`Kindergarten`).
		Where(`"kindergarten_student"."device"=?`, device).
		Where(`NOT "kindergarten_student"."deleted"`).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &student, nil
}

func CreateKindergartenStudentTx(tx *pg.Tx, name, no, remark, gender string, birthday time.Time, device string, kindergartenID, classID int64) (*KindergartenStudent, error) {
	kindergarten := Kindergarten{ID: kindergartenID}
	if err := tx.Model(&kindergarten).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	/* 检查学号是否重复 */
	if no != "" {
		student := KindergartenStudent{}
		if err := tx.Model(&student).Where(`NOT "deleted"`).Where(`"no"=?`, no).Where(`"class_id"=?`, classID).Limit(1).Select(); err != nil {
			if err != db.ErrNoRows {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		} else {
			/* 学号已经存在 */
			log.Errorf(`学号已经存在 %d %s`, kindergarten.ID, no)
			return nil, fmt.Errorf(`学号已经存在 %d %s`, kindergarten.ID, no)
		}
	}

	if classID > 0 {
		class := KindergartenClass{ID: classID}
		if err := tx.Model(&class).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		if class.KindergartenID != kindergarten.ID {
			log.Errorf(`班级不在该幼儿园内 %d %d`, class.ID, kindergarten.ID)
			return nil, fmt.Errorf(`班级不在该幼儿园内 %d %d`, class.ID, kindergarten.ID)
		}
		/* 班级学生数 +1 */
		if _, err := tx.Model(&class).WherePK().Set(`"number_of_student"="number_of_student"+1`).Update(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	}

	/* 检查设备是否已存在 */
	student := KindergartenStudent{}
	if err := tx.Model(&student).Where(`"device"=?`, device).Where(`NOT "deleted"`).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		/* 不存在 */
	} else {
		/* 设备已存在，创建失败 */
		log.Errorf(`设备已存在 %s`, device)
		return nil, fmt.Errorf(`设备已存在 %s`, device)
	}

	if _, err := tx.Model(&kindergarten).WherePK().Set(`"number_of_student"="number_of_student"+1`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	student = KindergartenStudent{
		Name:     name,
		NO:       no,
		Remark:   remark,
		Gender:   gender,
		Device:   device,
		Birthday: birthday,

		KindergartenID: kindergarten.ID,
		ClassID:        classID,
	}

	if _, err := tx.Model(&student).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &student, nil
}

func CreateKindergartenStudent(name, no, remark, gender string, birthday time.Time, device string, kindergartenID, classID int64) (*KindergartenStudent, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	device = x.FormatMacAddress(device)

	student, err := CreateKindergartenStudentTx(tx, name, no, remark, gender, birthday, device, kindergartenID, classID)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return student, nil
}

func UpdateKindergartenStudent(id int64, name, no, remark, gender string, birthday time.Time, device string, classID int64) (*KindergartenStudent, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	device = x.FormatMacAddress(device)

	student := KindergartenStudent{ID: id}
	if err := tx.Model(&student).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	/* 检查学号是否重复 */
	if no != "" {
		s := KindergartenStudent{}
		if err := tx.Model(&s).Where(`NOT "deleted"`).Where(`"no"=?`, no).Where(`"class_id"=?`, classID).Limit(1).Select(); err != nil {
			if err != db.ErrNoRows {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		} else if s.ID != student.ID {
			/* 学号已经存在 */
			log.Errorf(`学号已经存在 %d %s`, student.KindergartenID, no)
			return nil, fmt.Errorf(`学号已经存在 %d %s`, student.KindergartenID, no)
		}
	}

	if classID != student.ClassID {
		if classID > 0 {
			class := KindergartenClass{ID: classID}
			if err := tx.Model(&class).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
			if class.KindergartenID != student.KindergartenID {
				log.Errorf(`班级不在该幼儿园内 %d %d`, class.ID, student.KindergartenID)
				return nil, fmt.Errorf(`班级不在该幼儿园内 %d %d`, class.ID, student.KindergartenID)
			}
			/* 班级变更，新班级的学生数+1，旧班级的学生数-1 */
			if _, err := tx.Model(&class).WherePK().Set(`"number_of_student"="number_of_student"+1`).Update(); err != nil {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		}
		if student.ClassID > 0 {
			if _, err := tx.Model(&KindergartenClass{ID: student.ClassID}).WherePK().Set(`"number_of_student"="number_of_student"-1`).Update(); err != nil {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		}
	}

	q := tx.Model(&student).WherePK().Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"name"=?`, name).Set(`"no"=?`, no).
		Set(`"remark"=?`, remark).Set(`"gender"=?`, gender).Set(`"class_id"=?`, classID).Set(`"birthday"=?`, birthday)

	if student.Device != device {
		/* 更改设备， 检查设备是否已存在 */
		t := KindergartenStudent{}
		if err := tx.Model(&t).Where(`"device"=?`, device).Where(`NOT "deleted"`).Limit(1).Select(); err != nil {
			if err != db.ErrNoRows {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
			q = q.Set(`"device"=?`, device)
		} else {
			/* 设备已存在，创建失败 */
			log.Errorf(`设备已存在 %s`, device)
			return nil, fmt.Errorf(`设备已存在 %s`, device)
		}
	}

	if _, err := q.Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &student, nil
}

func FindKindergartenStudents(query, gender string, kindergartenID, classID int64, page, pageSize int) ([]*KindergartenStudent, int, error) {
	students := make([]*KindergartenStudent, 0)

	q := db.PG().Model(&students).Relation(`Class`).Relation(`Kindergarten`).Where(`NOT "kindergarten_student"."deleted"`)
	if query != "" {
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`"kindergarten_student"."name" ILIKE ?`, fmt.Sprintf("%%%s%%", query))
			q = q.WhereOr(`"kindergarten_student"."device"=?`, query)
			q = q.WhereOr(`"kindergarten_student"."no"=?`, query)
			return q, nil
		})
	}
	if gender != "" {
		q = q.Where(`"kindergarten_student"."gender"=?`, gender)
	}
	if kindergartenID > 0 {
		q = q.Where(`"kindergarten_student"."kindergarten_id"=?`, kindergartenID)
	}
	if classID != 0 {
		q = q.Where(`"kindergarten_student"."class_id"=?`, classID)
	}

	if page > 0 && pageSize > 0 {
		q = q.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	total, err := q.Order(`kindergarten_student.id DESC`).SelectAndCountEstimate(100000)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	}
	return students, total, nil
}

func DeleteKindergartenStudentTx(tx *pg.Tx, id int64) (*KindergartenStudent, error) {
	student := KindergartenStudent{ID: id}
	if err := tx.Model(&student).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	if _, err := tx.Model(&student).WherePK().Set(`"deleted"=?`, true).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(&Kindergarten{ID: student.KindergartenID}).WherePK().Set(`"number_of_student"="number_of_student"-1`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(&KindergartenClass{ID: student.ClassID}).WherePK().Set(`"number_of_student"="number_of_student"-1`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &student, nil
}

func DeleteKindergartenStudent(id int64) (*KindergartenStudent, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	student, err := DeleteKindergartenStudentTx(tx, id)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return student, nil
}

/* 将学生的班级设置成空 */
func ClearKindergartenStudentClassTx(tx *pg.Tx, classID int64) error {
	if r, err := tx.Exec(`UPDATE "kindergarten_student" SET "class_id"=0 WHERE NOT "deleted" AND "class_id"=?`, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	} else if r.RowsAffected() > 0 {
		if _, err := tx.Exec(`UPDATE "kindergarten_class" SET "number_of_student"="number_of_student"-? WHERE "id"=?`, r.RowsAffected(), classID); err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}
	}
	return nil
}
