package health

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	systemdb "gitlab.com/ykgk/kgo/db/system"
)

func GetKindergartenStudentFitnessTestToday(studentID int64) (*KindergartenStudentFitnessTest, error) {
	y, m, d := time.Now().Local().Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	test := KindergartenStudentFitnessTest{}
	if err := db.PG().Model(&test).Relation(`Student`).Where(`"date"=?`, date).Where(`"student_id"=?`, studentID).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &test, nil
}

func getKindergartenStudentFitnessTestTodayTx(tx *pg.Tx, studentID int64) (*KindergartenStudentFitnessTest, error) {
	student := kindergartendb.KindergartenStudent{ID: studentID}
	if err := tx.Model(&student).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	y, m, d := time.Now().Local().Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	test := KindergartenStudentFitnessTest{}
	if err := tx.Model(&test).Where(`"date"=?`, date).Where(`"student_id"=?`, student.ID).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		test = KindergartenStudentFitnessTest{
			KindergartenID: student.KindergartenID,
			StudentID:      student.ID,
			Date:           date,
		}
		if _, err := tx.Model(&test).Insert(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	}
	test.Student = &student
	return &test, nil
}

func updateKindergartenStudentFitnessTestTotalScoreTx(tx *pg.Tx, test *KindergartenStudentFitnessTest) error {
	if _, err := tx.Model(test).WherePK().Set(`"total_score"="shuttle_run_10_score"+"standing_long_jump_score"+"standing_long_jump_score"+"baseball_throw_score"+"bunny_hopping_score"+"sit_and_reach_score"+"balance_beam_score"+"height_and_weight_score"`).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	if !test.HeightUpdatedAt.IsZero() && !test.WeightUpdatedAt.IsZero() && !test.ShuttleRun10UpdatedAt.IsZero() && !test.BaseballThrowUpdatedAt.IsZero() &&
		!test.BunnyHoppingUpdatedAt.IsZero() && !test.BalanceBeamUpdatedAt.IsZero() && !test.SitAndReachUpdatedAt.IsZero() && !test.StandingLongJumpUpdatedAt.IsZero() {
		/* 全部评测 */
		if test.TotalScore > 31 {
			test.TotalStatus = "excellent"
		} else if test.TotalScore <= 31 && test.TotalScore >= 28 {
			test.TotalStatus = "good"
		} else if test.TotalScore >= 20 && test.TotalScore <= 27 {
			test.TotalStatus = "okay"
		} else {
			test.TotalStatus = "fail"
		}
	}
	return nil
}
func CreateKindergartenStudentFitnessTestHeight(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}

	q := tx.Model(test).WherePK().Set(`"height"=?`, v).Set(`"height_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"updated_at"=CURRENT_TIMESTAMP`)
	if !test.WeightUpdatedAt.IsZero() && test.Weight > 0 {
		score, err := systemdb.GetStandardSclaeHWScoreTx(tx, test.Student.Gender, v, test.Weight)
		if err != nil {
			log.Errorf("GetStandardSclaeHWScoreTx() Error: %v", err)
			return nil, err
		}
		q = q.Set(`"height_and_weight_score"=?`, score)
	}

	if _, err := q.Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

func CreateKindergartenStudentFitnessTestWeight(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}

	q := tx.Model(test).WherePK().Set(`"weight"=?`, v).Set(`"weight_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"updated_at"=CURRENT_TIMESTAMP`)
	if !test.HeightUpdatedAt.IsZero() && test.Height > 0 {
		score, err := systemdb.GetStandardSclaeHWScoreTx(tx, test.Student.Gender, test.Height, v)
		if err != nil {
			log.Errorf("GetStandardSclaeHWScoreTx() Error: %v", err)
			return nil, err
		}
		q = q.Set(`"height_and_weight_score"=?`, score)
	}

	if _, err := q.Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

func CreateKindergartenStudentFitnessTestShuttleRun10(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}
	score, err := systemdb.GetStandardScaleScoreTx(tx, "10米折返跑(秒)", test.Student.Gender, test.Student.Age(), v)
	if err != nil {
		log.Errorf("GetStandardScaleScoreTx() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(test).WherePK().Set(`"shuttle_run_10"=?`, v).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"shuttle_run_10_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"shuttle_run_10_score"=?`, score).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

func CreateKindergartenStudentFitnessTestStandingLongJump(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}
	score, err := systemdb.GetStandardScaleScoreTx(tx, "立定跳远(厘米)", test.Student.Gender, test.Student.Age(), v)
	if err != nil {
		log.Errorf("GetStandardScaleScoreTx() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(test).WherePK().Set(`"standing_long_jump"=?`, v).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"standing_long_jump_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"standing_long_jump_score"=?`, score).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

/* 网球掷远 */
func CreateKindergartenStudentFitnessTestBaseballThrow(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}
	score, err := systemdb.GetStandardScaleScoreTx(tx, "网球掷远(米)", test.Student.Gender, test.Student.Age(), math.Round(v/0.5)*0.5)
	if err != nil {
		log.Errorf("GetStandardScaleScoreTx() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(test).WherePK().Set(`"baseball_throw"=?`, v).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"baseball_throw_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"baseball_throw_score"=?`, score).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

/* 双脚连续跳(秒) */
func CreateKindergartenStudentFitnessTestBunnyHopping(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}
	score, err := systemdb.GetStandardScaleScoreTx(tx, "双脚连续跳(秒)", test.Student.Gender, test.Student.Age(), math.Round(v/0.5)*0.5)
	if err != nil {
		log.Errorf("GetStandardScaleScoreTx() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(test).WherePK().Set(`"bunny_hopping"=?`, v).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"bunny_hopping_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"bunny_hopping_score"=?`, score).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

/* 坐位体前屈(厘米) */
func CreateKindergartenStudentFitnessTestSitAndReach(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}
	score, err := systemdb.GetStandardScaleScoreTx(tx, "坐位体前屈(厘米)", test.Student.Gender, test.Student.Age(), math.Round(v/0.5)*0.5)
	if err != nil {
		log.Errorf("GetStandardScaleScoreTx() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(test).WherePK().Set(`"sit_and_reach"=?`, v).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"sit_and_reach_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"sit_and_reach_score"=?`, score).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

/* 坐位体前屈(厘米) */
func CreateKindergartenStudentFitnessTestBalanceBeam(studentID int64, v float64) (*KindergartenStudentFitnessTest, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	v = math.Floor(v*10) / 10 // 精确到一位小数

	test, err := getKindergartenStudentFitnessTestTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentFitnessTestTodayTx() Error: %v", err)
		return nil, err
	}
	score, err := systemdb.GetStandardScaleScoreTx(tx, "走平衡木(秒)", test.Student.Gender, test.Student.Age(), math.Round(v/0.5)*0.5)
	if err != nil {
		log.Errorf("GetStandardScaleScoreTx() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(test).WherePK().Set(`"balance_beam"=?`, v).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"balance_beam_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"balance_beam_score"=?`, score).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentFitnessTestTotalScoreTx(tx, test); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return test, nil
}

func DeleteKindergartenStudentFitnessTestsByStudentTx(tx *pg.Tx, studentID int64) error {
	if _, err := tx.Model(&KindergartenStudentFitnessTest{}).Where(`"student_id"=?`, studentID).
		Where(`NOT "deleted"`).Set(`"deleted"=?`, true).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

func FindKindergartenStudentFitnessTests(query string, kindergartenID, classID, studentID int64,
	startTime, endTime time.Time, heightWeightFilters, shuttleRun10Filters, standingLongJumpFilters, baseballThrowFilters,
	bunnyHoppingFilters, sitAndReachFilters, balanceBeamFilters []int, totalStatusFilters []string,
	page, pageSize int) ([]*KindergartenStudentFitnessTest, int, error) {
	tests := make([]*KindergartenStudentFitnessTest, 0)
	q := db.PG().Model(&tests).Relation(`Student`).Relation(`Kindergarten`).Relation(`Student.Class`).
		Where(`NOT "kindergarten_student_fitness_test"."deleted"`).Where(`NOT Student.deleted`)
	if query != "" {
		q = q.Where(`Student.name ILIKE ?`, fmt.Sprintf("%%%s%%", query))
	}
	if kindergartenID > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."kindergarten_id"=?`, kindergartenID)
	}
	if classID != 0 {
		q = q.Where(`Student.class_id=?`, classID)
	}
	if studentID > 0 {
		q = q.Where(`Student.id=?`, studentID)
	}

	if !startTime.IsZero() {
		/* PSQL中的 date 不包含时区，必须使用UTC时间来查询 */
		y, m, d := startTime.Local().Date()
		startTime = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		q = q.Where(`"kindergarten_student_fitness_test"."date">=?`, startTime)
	}
	if !endTime.IsZero() {
		y, m, d := endTime.Local().Date()
		endTime = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		q = q.Where(`"kindergarten_student_fitness_test"."date"<?`, endTime)
	}

	if len(heightWeightFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."height_and_weight_score" IN (?)`, pg.In(heightWeightFilters)).
			Where(`"kindergarten_student_fitness_test"."height_updated_at" IS NOT NULL`).
			Where(`"kindergarten_student_fitness_test"."weight_updated_at" IS NOT NULL`)
	}

	if len(shuttleRun10Filters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."shuttle_run_10_score" IN (?)`, pg.In(shuttleRun10Filters)).
			Where(`"kindergarten_student_fitness_test"."shuttle_run_10_updated_at" IS NOT NULL`)
	}

	if len(standingLongJumpFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."standing_long_jump_score" IN (?)`, pg.In(standingLongJumpFilters)).
			Where(`"kindergarten_student_fitness_test"."standing_long_jump_updated_at" IS NOT NULL`)
	}

	if len(baseballThrowFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."baseball_throw_score" IN (?)`, pg.In(baseballThrowFilters)).
			Where(`"kindergarten_student_fitness_test"."baseball_throw_updated_at" IS NOT NULL`)
	}

	if len(bunnyHoppingFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."bunny_hopping_score" IN (?)`, pg.In(bunnyHoppingFilters)).
			Where(`"kindergarten_student_fitness_test"."bunny_hopping_updated_at" IS NOT NULL`)
	}

	if len(sitAndReachFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."sit_and_reach_score" IN (?)`, pg.In(sitAndReachFilters)).
			Where(`"kindergarten_student_fitness_test"."sit_and_reach_updated_at" IS NOT NULL`)
	}

	if len(balanceBeamFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."balance_beam_score" IN (?)`, pg.In(balanceBeamFilters)).
			Where(`"kindergarten_student_fitness_test"."balance_beam_updated_at" IS NOT NULL`)
	}

	if len(totalStatusFilters) > 0 {
		q = q.Where(`"kindergarten_student_fitness_test"."total_status" IN (?)`, pg.In(totalStatusFilters))
	}

	q = q.Order(`kindergarten_student_fitness_test.updated_at DESC`)

	if page > 0 && pageSize > 0 {
		total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).SelectAndCountEstimate(100000)
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, 0, err
		}
		return tests, total, nil
	} else {
		err := q.Limit(10000).Select()
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, 0, err
		}
		return tests, 0, nil
	}
}

func FindKindergartenStudentFitnessTestDates(kindergartenID, classID int64) ([]time.Time, error) {

	type T struct {
		Date time.Time `json:"date" pg:"date"`
	}
	ts := make([]*T, 0)

	sql := `
		SELECT "date" FROM "kindergarten_student_fitness_test" AS "test"
		INNER JOIN "kindergarten_student" AS "student" ON "student"."id"="test"."student_id"
		WHERE NOT "test"."deleted"
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)
	if kindergartenID > 0 {
		wheres = append(wheres, `"test"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}

	if len(wheres) > 0 {
		sql += ` AND ` + strings.Join(wheres, ` AND `)
	}

	sql += ` GROUP BY "date" ORDER BY "date" DESC LIMIT 100`

	if _, err := db.PG().Query(&ts, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	dates := make([]time.Time, 0)
	for _, t := range ts {
		dates = append(dates, t.Date)
	}
	return dates, nil
}

type KindergartenStudentFitnessTestScoreVisionData struct {
	Score float64 `json:"score" pg:"score"`
	Count int     `json:"count" pg:"count"`
}

func FindKindergartenStudentFitnessTestScoreVision(fieldName string, kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentFitnessTestScoreVisionData, error) {
	datas := make([]*KindergartenStudentFitnessTestScoreVisionData, 0)

	sql := `
	SELECT "` + fieldName + `_score" AS "score",COUNT(1) AS "count" FROM
		"kindergarten_student_fitness_test" AS "test" INNER JOIN
		"kindergarten_student" AS "student" ON "student"."id"="test"."student_id"
		WHERE NOT "test"."deleted" AND NOT "student"."deleted" AND "test"."` + fieldName + `_updated_at" IS NOT NULL
`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"test"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"test"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}
	sql += ` AND ` + strings.Join(wheres, ` AND `)

	sql += ` GROUP BY "score" ORDER BY "score" DESC`

	if _, err := db.PG().Query(&datas, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return datas, nil
}

type KindergartenStudentFitnessTestHeightVisionData struct {
	Height int    `json:"height" pg:"h"`
	Gender string `json:"gender" pg:"gender"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentFitnessTestHeightVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentFitnessTestHeightVisionData, error) {
	datas := make([]*KindergartenStudentFitnessTestHeightVisionData, 0)

	sql := `
	SELECT FLOOR("height"/10) AS "h","student"."gender" AS "gender", COUNT(1) AS "count"
	FROM "kindergarten_student_fitness_test" AS "test"
	INNER JOIN "kindergarten_student" AS "student" ON "student"."id"="test"."student_id"
	WHERE NOT "test"."deleted" AND "test"."height_updated_at" IS NOT NULL
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"test"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"test"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}
	sql += ` AND ` + strings.Join(wheres, ` AND `)

	sql += ` GROUP BY "gender","h" ORDER BY "h"`

	if _, err := db.PG().Query(&datas, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return datas, nil
}

type KindergartenStudentFitnessTestWeightVisionData struct {
	Weight int    `json:"weight" pg:"w"`
	Gender string `json:"gender" pg:"gender"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentFitnessTestWeightVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentFitnessTestWeightVisionData, error) {
	datas := make([]*KindergartenStudentFitnessTestWeightVisionData, 0)

	sql := `
	SELECT FLOOR("weight"/5) AS "w","student"."gender" AS "gender", COUNT(1) AS "count"
	FROM "kindergarten_student_fitness_test" AS "test"
	INNER JOIN "kindergarten_student" AS "student" ON "student"."id"="test"."student_id"
	WHERE NOT "test"."deleted" AND "test"."weight_updated_at" IS NOT NULL
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"test"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"test"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}
	sql += ` AND ` + strings.Join(wheres, ` AND `)

	sql += ` GROUP BY "gender","w" ORDER BY "w"`

	if _, err := db.PG().Query(&datas, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return datas, nil
}

type KindergartenStudentFitnessTestStatusVisionData struct {
	Status string `json:"status" pg:"status"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentFitnessTestStatusVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentFitnessTestStatusVisionData, error) {
	datas := make([]*KindergartenStudentFitnessTestStatusVisionData, 0)

	sql := `
	SELECT "test"."total_status" AS "status",COUNT(1) AS "count" FROM
		"kindergarten_student_fitness_test" AS "test" INNER JOIN
		"kindergarten_student" AS "student" ON "student"."id"="test"."student_id"
		WHERE NOT "test"."deleted" AND NOT "student"."deleted" AND "test"."total_status" != ''
`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"test"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"test"."kindergarten_id"=?`)
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
