package health

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

type KindergartenStudentMorningCheckStat struct {
	Date               time.Time `pg:"date" json:"date"`
	Count              int       `pg:"count" json:"count"`                             // 体检总数
	TemperatureOK      int       `pg:"temperature_ok" json:"temperature_ok"`           // 体温正常人数
	TemperatureUnusual int       `pg:"temperature_unusual" json:"temperature_unusual"` // 体温异常人数

	HandUnusual  int `pg:"hand_unusual" json:"hand_unusual"`
	MouthUnusual int `pg:"mouth_unusual" json:"mouth_unusual"`
	EyeUnusual   int `pg:"eye_unusual" json:"eye_unusual"`

	Unusual int `pg:"unusual" json:"unusual"`
}

func GetKindergartenStudentMorningCheckStat(kindergartenID, classID int64, date time.Time) (*KindergartenStudentMorningCheckStat, error) {
	stat := KindergartenStudentMorningCheckStat{}
	stat.Date = date

	q := db.PG().Model(&KindergartenStudentMorningCheck{}).Relation(`Student`).
		Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID).
		Where(`NOT "kindergarten_student_morning_check"."deleted"`).
		Where(`"kindergarten_student_morning_check"."date"=?`, date)
	if classID > 0 {
		q = q.Where(`"student"."class_id"=?`, classID)
	}
	c, err := q.Count()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	stat.Count = c

	/* 任意异常 */
	q = db.PG().Model(&KindergartenStudentMorningCheck{}).Relation(`Student`).
		Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID).
		Where(`NOT "kindergarten_student_morning_check"."deleted"`).
		Where(`"kindergarten_student_morning_check"."date"=?`, date)
	q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		q = q.WhereOr(`"kindergarten_student_morning_check"."temperature_status"!=?`, "normal")
		q = q.WhereOr(`"kindergarten_student_morning_check"."hand"=?`, "异常")
		q = q.WhereOr(`"kindergarten_student_morning_check"."mouth"=?`, "异常")
		q = q.WhereOr(`"kindergarten_student_morning_check"."eye"=?`, "异常")
		return q, nil

	})
	if classID > 0 {
		q = q.Where(`"student"."class_id"=?`, classID)
	}
	c, err = q.Count()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	stat.Unusual = c
	/* 体温 */
	q = db.PG().Model(&KindergartenStudentMorningCheck{}).Relation(`Student`).
		Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID).
		Where(`NOT "kindergarten_student_morning_check"."deleted"`).
		Where(`"kindergarten_student_morning_check"."temperature_status"=?`, "normal").
		Where(`"kindergarten_student_morning_check"."date"=?`, date)
	if classID > 0 {
		q = q.Where(`"student"."class_id"=?`, classID)
	}
	c, err = q.Count()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	stat.TemperatureOK = c
	stat.TemperatureUnusual = stat.Count - stat.TemperatureOK

	/* 手眼口 TODO */
	q = db.PG().Model(&KindergartenStudentMorningCheck{}).Relation(`Student`).
		Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID).
		Where(`NOT "kindergarten_student_morning_check"."deleted"`).
		Where(`"kindergarten_student_morning_check"."hand"!='异常'`).
		Where(`"kindergarten_student_morning_check"."date"=?`, date)
	if classID > 0 {
		q = q.Where(`"student"."class_id"=?`, classID)
	}
	c, err = q.Count()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	// stat.TemperatureOK = c
	stat.HandUnusual = c
	// log.Errorf("+++++++++++++++++++++++++", c)

	q = db.PG().Model(&KindergartenStudentMorningCheck{}).Relation(`Student`).
		Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID).
		Where(`NOT "kindergarten_student_morning_check"."deleted"`).
		Where(`"kindergarten_student_morning_check"."mouth"=?`, "异常").
		Where(`"kindergarten_student_morning_check"."date"=?`, date)
	if classID > 0 {
		q = q.Where(`"student"."class_id"=?`, classID)
	}
	c, err = q.Count()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	// stat.TemperatureOK = c
	stat.MouthUnusual = c

	q = db.PG().Model(&KindergartenStudentMorningCheck{}).Relation(`Student`).
		Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID).
		Where(`NOT "kindergarten_student_morning_check"."deleted"`).
		Where(` "kindergarten_student_morning_check"."eye"=?`, "异常").
		Where(`"kindergarten_student_morning_check"."date"=?`, date)
	if classID > 0 {
		q = q.Where(`"student"."class_id"=?`, classID)
	}
	c, err = q.Count()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	// stat.TemperatureOK = c
	stat.EyeUnusual = c
	// fmt.Print("+++++++++++++++++++++++++++++",c)
	return &stat, nil
}

// func FindKindergartenStudentMorningCheckStats(kindergartenID, classID int64, startDate, endDate time.Time) ([]*KindergartenStudentMorningCheckStat, error) {
// 	stats := make([]*KindergartenStudentMorningCheckStat, 0)
// 	for !startDate.After(endDate) {
// 		stat, err := GetKindergartenStudentMorningCheckStat(kindergartenID, classID, startDate)
// 		if err != nil {
// 			log.Errorf("DB Error: %v", err)
// 			return nil, err
// 		}
// 		stats = append(stats, stat)
// 		startDate = startDate.AddDate(0, 0, 1)
// 	}
// 	return stats, nil
// }

func FindKindergartenStudentMorningChecks(query string,
	kindergartenID, classID, studentID int64, temperatureFilters []string,
	startTime, endTime time.Time, page, pageSize int) ([]*KindergartenStudentMorningCheck, int, error) {
	checks := make([]*KindergartenStudentMorningCheck, 0)
	q := db.PG().Model(&checks).Relation(`Student`).Relation(`Kindergarten`).Relation(`Student.Class`).Where(`NOT "kindergarten_student_morning_check"."deleted"`)

	if query != "" {
		q = q.Where(`Student.name ILIKE ?`, fmt.Sprintf("%%%s%%", query))
	}
	if kindergartenID > 0 {
		q = q.Where(`"kindergarten_student_morning_check"."kindergarten_id"=?`, kindergartenID)
	}
	if classID != 0 {
		q = q.Where(`Student.class_id=?`, classID)
	}
	if studentID > 0 {
		q = q.Where(`Student.id=?`, studentID)
	}
	if len(temperatureFilters) > 0 {
		q = q.Where(`"kindergarten_student_morning_check"."temperature_status" IN (?)`, pg.In(temperatureFilters))
	}
	if !startTime.IsZero() {
		/* PSQL中的 date 不包含时区，必须使用UTC时间来查询 */
		y, m, d := startTime.Local().Date()
		startTime = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		q = q.Where(`"kindergarten_student_morning_check"."date">=?`, startTime)
	}
	if !endTime.IsZero() {
		y, m, d := endTime.Local().Date()
		endTime = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		q = q.Where(`"kindergarten_student_morning_check"."date"<?`, endTime)
	}

	q = q.Order(`kindergarten_student_morning_check.id DESC`)
	if page > 0 && pageSize > 0 {
		total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).SelectAndCountEstimate(100000)
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, 0, err
		}
		return checks, total, nil
	} else {
		err := q.Limit(10000).Select()
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, 0, err
		}
		return checks, 0, nil
	}
}

func CreateKindergartenStudentMorningCheck(studentID int64, temperature float64, hand, mouth, eye string) (*KindergartenStudentMorningCheck, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	y, m, d := time.Now().Local().Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	student := kindergartendb.KindergartenStudent{ID: studentID}
	if err := tx.Model(&student).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	status := GetTemperatureStatus(temperature)

	check := KindergartenStudentMorningCheck{}
	if err := tx.Model(&check).Where(`"student_id"=?`, student.ID).Where(`"date"=?`, date).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		/* 晨检不存在，创建 */
		check = KindergartenStudentMorningCheck{
			KindergartenID:    student.KindergartenID,
			StudentID:         student.ID,
			Date:              date,
			Temperature:       temperature,
			TemperatureStatus: status,
			Hand:              hand,
			Mouth:             mouth,
			Eye:               eye,
		}

		if _, err := tx.Model(&check).Insert(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	} else {
		/* 晨检已存在，覆盖 */
		if _, err := tx.Model(&check).WherePK().Set(`"temperature"=?`, temperature).Set(`"temperature_status"=?`, status).
			Set(`"hand"=?`, hand).Set(`"mouth"=?`, mouth).Set(`"eye"=?`, eye).
			Set(`"updated_at"=CURRENT_TIMESTAMP`).Returning(`*`).Update(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return &check, nil
}

func DeleteKindergartenStudentMorningChecksByStudentTx(tx *pg.Tx, studentID int64) error {
	if _, err := tx.Model(&KindergartenStudentMorningCheck{}).Where(`"student_id"=?`, studentID).
		Where(`NOT "deleted"`).Set(`"deleted"=?`, true).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

type KindergartenStudentMorningCheckTemperatureVisionData struct {
	Status string `json:"status" pg:"status"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentMorningCheckTemperatureVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentMorningCheckTemperatureVisionData, error) {
	datas := make([]*KindergartenStudentMorningCheckTemperatureVisionData, 0)

	sql := `
		SELECT "temperature_status" AS "status",COUNT(1) AS "count" FROM
			"kindergarten_student_morning_check" AS "check" INNER JOIN
			"kindergarten_student" AS "student" ON "student"."id"="check"."student_id"
			WHERE NOT "check"."deleted" AND NOT "student"."deleted"
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"check"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"check"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}
	sql += ` AND ` + strings.Join(wheres, ` AND `)

	sql += ` GROUP BY "status" ORDER BY "status" DESC`

	if _, err := db.PG().Query(&datas, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return datas, nil
}
