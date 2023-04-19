package health

import (
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

/* 删除学生以及学生相关的体检体侧晨检 */
func DeleteStudentWithCheckAndExam(studentID int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	defer tx.Rollback()
	student, err := kindergartendb.DeleteKindergartenStudentTx(tx, studentID)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	} else if student == nil {
		return nil
	}
	if err := DeleteKindergartenStudentMorningChecksByStudentTx(tx, studentID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	if err := DeleteKindergartenStudentMedicalExaminationsByStudentTx(tx, studentID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	/* 删除体侧信息 */
	if _, err := tx.Exec(`UPDATE "kindergarten_student_fitness_test" SET "deleted"=? WHERE NOT "deleted" AND "student_id" = ?)`, true, studentID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

func DeleteStudentWithCheckAndExamByClassTx(tx *pg.Tx, classID int64) error {
	/* 删除晨检信息 */
	if _, err := tx.Exec(`UPDATE "kindergarten_student_morning_check" SET "deleted"=? WHERE NOT "deleted" AND "student_id" IN (SELECT "id" FROM "kindergarten_student" WHERE NOT "deleted" AND "class_id"=?)`, true, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}

	/* 删除体检信息 */
	if _, err := tx.Exec(`UPDATE "kindergarten_student_medical_examination" SET "deleted"=? WHERE NOT "deleted" AND "student_id" IN (SELECT "id" FROM "kindergarten_student" WHERE NOT "deleted" AND "class_id"=?)`, true, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}

	/* 删除体侧信息 */
	if _, err := tx.Exec(`UPDATE "kindergarten_student_fitness_test" SET "deleted"=? WHERE NOT "deleted" AND "student_id" IN (SELECT "id" FROM "kindergarten_student" WHERE NOT "deleted" AND "class_id"=?)`, true, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}

	if r, err := tx.Exec(`UPDATE "kindergarten_student" SET "deleted"=? WHERE NOT "deleted" AND "class_id"=?`, true, classID); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	} else if r.RowsAffected() > 0 {
		/* 更新班级和幼儿园的学生数量 */
		if _, err := tx.Model(&kindergartendb.KindergartenClass{ID: classID}).WherePK().Set(`"number_of_student"="number_of_student"-?`, r.RowsAffected()).Update(); err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}

		if _, err := tx.Model(&kindergartendb.Kindergarten{}).Where(`"id"=(SELECT "kindergarten_id" FROM "kindergarten_class" WHERE "id"=?)`, classID).Set(`"number_of_student"="number_of_student"-?`, r.RowsAffected()).Update(); err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}
	}

	return nil
}
