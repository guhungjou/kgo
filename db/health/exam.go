package health

import (
	// gomain "gitlab.com/ykgk/kgo"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

func GetKindergartenStudentMedicalExamination(id int64) (*KindergartenStudentMedicalExamination, error) {

	exam := KindergartenStudentMedicalExamination{ID: id}

	if err := db.PG().Model(&exam).Relation(`Student`).Relation(`Kindergarten`).Relation(`Student.Class`).WherePK().Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	MyContent := "标准身高为"

	//获取exam中的height和weight、BMI、Hemoglobin、SightLS、SightRS、SightLC、SightRC、ToothCount、ALT、EyeLeft、EyeRight
	height := ""
	if exam.Height == 0 {
		height = "未测试"
	} else {
		height = strconv.FormatFloat(exam.Height, 'f', 2, 64)

		stHeight, _ := GetStandardScaleHL("身高", exam.Student.Gender, float64(int(exam.Student.Age())))
		MyContent += fmt.Sprintf("%.2f", stHeight.Max) +
			"~" + fmt.Sprintf("%.2f", stHeight.Min) +
			"cm，他为" + height +
			"cm。"
	}

	weight := ""
	if exam.Weight == 0 {
		weight = "未测试"
	} else {
		weight = strconv.FormatFloat(exam.Weight, 'f', 2, 64)
		stWeight, _ := GetStandardScaleHL("体重", exam.Student.Gender, float64(int(exam.Student.Age())))
		MyContent += "标准体重为" + fmt.Sprintf("%.2f", stWeight.Max) +
			"~" + fmt.Sprintf("%.2f", stWeight.Min) +
			"kg，他为" + weight + "kg。"
	}

	BMI := ""
	if exam.BMI == 0 {
		BMI = "未测试"
	} else {
		BMI = strconv.FormatFloat(exam.BMI, 'f', 2, 64)
		stBMI, _ := GetStandardScaleHL("BMI", exam.Student.Gender, float64(int(exam.Student.Age())))
		MyContent += "标准BMI为" + fmt.Sprintf("%.2f", stBMI.Max) +
			"~" + fmt.Sprintf("%.2f", stBMI.Min) +
			"，他为" + BMI + "。"
	}

	Hemoglobin := ""
	if exam.Hemoglobin == 0 {
		Hemoglobin = "未测试"
	} else {
		Hemoglobin = strconv.FormatFloat(exam.Hemoglobin, 'f', 2, 64)
		stHemoglobin, _ := GetStandardScaleHL("血红蛋白", "", 0)
		MyContent += "标准血红蛋白为" + fmt.Sprintf("%.2f", stHemoglobin.Max) +
			"~" + fmt.Sprintf("%.2f", stHemoglobin.Min) +
			"，他为" + Hemoglobin + "。"
	}
	if exam.SightLC != "" || exam.SightRC != "" {
		MyContent += "当视力仪测的球镜度不为normal时，有近视的可能性，"
	}

	if exam.SightLS == "" {
	} else {
		MyContent += "他的左眼球镜度为" + exam.SightLSStatus + "，"
	}

	if exam.SightRS == "" {
	} else {
		MyContent += "右眼球镜度为" + exam.SightRSStatus + "。"
	}

	if exam.SightLC != "" || exam.SightRC != "" {
		MyContent += "当视力仪测的柱镜度不为normal时，有散光的可能性，"
	}
	if exam.SightLC == "" {
	} else {
		MyContent += "他的左眼柱镜度为" + exam.SightLCStatus + "，"

	}

	if exam.SightRC == "" {
	} else {
		MyContent += "右眼柱镜度为" + exam.SightRCStatus + "。"
	}

	ToothCount := ""
	if exam.ToothCount == 0 {
		ToothCount = "未测试"
	} else {
		ToothCount = strconv.Itoa(exam.ToothCount)
		MyContent += "牙齿数一般为20颗左右，他为" + ToothCount + "，他有" + strconv.Itoa(exam.ToothCariesCount) + "颗龋齿。"

	}
	ALT := ""
	if exam.ALT == 0 {
		ALT = "未测试"
	} else {
		ALT = strconv.FormatFloat(exam.ALT, 'f', 2, 64)
		stALT, _ := GetStandardScaleHL("谷丙转氨酶", "", 0)
		MyContent += "标准谷丙转氨酶为" + fmt.Sprintf("%.2f", stALT.Max) +
			"~" + fmt.Sprintf("%.2f", stALT.Min) +
			"，他为" + ALT + "。"
	}

	if exam.EyeLeftStatus != "" || exam.EyeRightStatus != "" {
		MyContent += "通过视力表测的视力如下："
	}

	if exam.EyeLeftStatus == "" {

	} else {
		if exam.EyeLeftStatus == "very_bad" {
			MyContent += "他的左眼视力为重度近视，"
		} else if exam.EyeLeftStatus == "little_bad" {
			MyContent += "他的左眼视力为轻度近视，"
		} else if exam.EyeLeftStatus == "良好" {
			MyContent += "他的左眼视力为良好，"
		}

	}

	if exam.EyeRightStatus == "" {

	} else {
		if exam.EyeRightStatus == "very_bad" {
			MyContent += "他的右眼视力为重度近视，"
		} else if exam.EyeRightStatus == "little_bad" {
			MyContent += "他的右眼视力为轻度近视，"
		} else if exam.EyeRightStatus == "良好" {
			MyContent += "他的右眼视力为良好，"
		}

	}
	MyContent += "使用官方文件语言风格，严格按照上述每一项检测结果提出精简的建议，不少于250字，但不多于300字。"
	// BMI := exam.BMI
	// Hemoglobin := exam.Hemoglobin
	// SightLS := exam.SightLS
	// SightRS := exam.SightRS
	// SightLC := exam.SightLC
	// SightRC := exam.SightRC
	// ToothCount := exam.ToothCount
	// ALT := exam.ALT
	// EyeLeft := exam.EyeLeft
	// EyeRight := exam.EyeRight
	//得到标准身高

	//得到标准体重

	//得到标准BMI

	//得到标准血红蛋白

	//得到标准谷丙转氨酶

	//得到标准球镜度

	//得到标准柱镜度

	fmt.Println(MyContent)
	for i := 0; i < 6; i++ {
		config := openai.DefaultConfig("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")
		proxyUrl, err := url.Parse("http://localhost:7890")
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}

		client := openai.NewClientWithConfig(config)
		// client := openai.NewClient("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")
		// fmt.Println(context.Background())

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Temperature: 0.4,
				MaxTokens:   1000,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						//exam.Student.Age()转换成string类型
						Content: MyContent,
					},
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "不要用第一人称回答。你的回答中不要对饮酒、隐形眼镜等进行建议。",
					},
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "你的回答中不要出现：饮酒、normal、根据检测结果、按照、他、您、建议等字眼。请注意，小孩本身使用电脑等电子产品的频率就少，因此对于视力有问题的情况，轻度近视请从“减少电子产品使用频率，规范用眼、培养良好阅读和学习姿势”等方面考虑，可以加以扩展；重度近视请从复查、治疗等方面考虑，可以加以扩展。",
					},
				},
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			// time.Sleep(1 * time.Second)
			fmt.Print("重发请求...")
		} else {
			fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++" + resp.Choices[0].Message.Content)
			exam.Advise = exam.Student.Name + ":" + resp.Choices[0].Message.Content
			break
		}
	}

	return &exam, nil
}

func GetKindergartenStudentMedicalExaminationToday(studentID int64) (*KindergartenStudentMedicalExamination, error) {
	y, m, d := time.Now().Local().Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	exam := KindergartenStudentMedicalExamination{}
	if err := db.PG().Model(&exam).Relation(`Student`).Where(`"date"=?`, date).Where(`"student_id"=?`, studentID).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &exam, nil
}

func getKindergartenStudentMedicalExaminationTodayTx(tx *pg.Tx, studentID int64) (*KindergartenStudentMedicalExamination, error) {
	student := kindergartendb.KindergartenStudent{ID: studentID}
	if err := tx.Model(&student).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	y, m, d := time.Now().Local().Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	exam := KindergartenStudentMedicalExamination{}
	if err := tx.Model(&exam).Where(`"date"=?`, date).Where(`"student_id"=?`, student.ID).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		exam = KindergartenStudentMedicalExamination{
			KindergartenID: student.KindergartenID,
			StudentID:      student.ID,
			Date:           date,
		}
		if _, err := tx.Model(&exam).Insert(); err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	}
	exam.Student = &student
	return &exam, nil
}

/* 上传 身高 信息 */
func CreateKindergartenStudentMedicalExaminationHeight(studentID int64, height float64) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}
	status := GetHeightStatus(exam.Student.Age(), exam.Student.Gender, height)

	if _, err := tx.Model(exam).WherePK().Set(`"height"=?`, height).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"height_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"height_status"=?`, status).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentMedicalExaminationBMITx(tx, exam); err != nil {
		log.Errorf("updateKindergartenStudentMedicalExaminationBMITx Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 更新体检的BMI指数 */
func updateKindergartenStudentMedicalExaminationBMITx(tx *pg.Tx, exam *KindergartenStudentMedicalExamination) error {
	if exam.HeightUpdatedAt.IsZero() || exam.WeightUpdatedAt.IsZero() || exam.Height <= 0 || exam.Weight <= 0 {
		return nil
	}
	bmi := exam.Weight / exam.Height / exam.Height * 10000.0
	status := GetBMIStatus(exam.Student.Age(), exam.Student.Gender, bmi)

	if _, err := tx.Model(exam).WherePK().
		Set(`"bmi"=?`, bmi).Set(`"bmi_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"bmi_status"=?`, status).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

/* 上传 体重 信息 */
func CreateKindergartenStudentMedicalExaminationWeight(studentID int64, weight float64) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	status := GetWeightStatus(exam.Student.Age(), exam.Student.Gender, weight)

	if _, err := tx.Model(exam).WherePK().Set(`"weight"=?`, weight).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"weight_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"weight_status"=?`, status).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := updateKindergartenStudentMedicalExaminationBMITx(tx, exam); err != nil {
		log.Errorf("updateKindergartenStudentMedicalExaminationBMITx Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 上传 血红蛋白 信息 */
func CreateKindergartenStudentMedicalExaminationHemoglobin(studentID int64, hemoglobin float64, remark string) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	status := GetHemoglobinStatus(hemoglobin)

	if _, err := tx.Model(exam).WherePK().Set(`"hemoglobin"=?`, hemoglobin).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"hemoglobin_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"hemoglobin_remark"=?`, remark).Set(`"hemoglobin_status"=?`, status).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 上传 视力 信息 */
func CreateKindergartenStudentMedicalExaminationSight(studentID int64, ls, lc, rs, rc string, lRemark, rRemark string) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	lsstatus := GetSightSStatus(x.ParseFloat64(ls))
	lcstatus := GetSightCStatus(x.ParseFloat64(lc))
	rsstatus := GetSightSStatus(x.ParseFloat64(rs))
	rcstatus := GetSightCStatus(x.ParseFloat64(rc))

	if _, err := tx.Model(exam).WherePK().Set(`"sight_l_s"=?`, ls).Set(`"sight_l_c"=?`, lc).
		Set(`"sight_r_s"=?`, rs).Set(`"sight_r_c"=?`, rc).
		Set(`"sight_l_remark"=?`, lRemark).Set(`"sight_r_remark"=?`, rRemark).
		Set(`"sight_l_s_status"=?`, lsstatus).Set(`"sight_l_c_status"=?`, lcstatus).
		Set(`"sight_r_s_status"=?`, rsstatus).Set(`"sight_r_c_status"=?`, rcstatus).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"sight_updated_at"=CURRENT_TIMESTAMP`).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 上传 视力2 左眼 信息 */
func CreateKindergartenStudentMedicalExaminationNewLSight(studentID int64, EyeLeft float32) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	// lsstatus := GetSightSStatus(x.ParseFloat64(ls))
	// lcstatus := GetSightCStatus(x.ParseFloat64(lc))
	// rsstatus := GetSightSStatus(x.ParseFloat64(rs))
	// rcstatus := GetSightCStatus(x.ParseFloat64(rc))
	var test string
	if EyeLeft < 4.5 {
		test = "very_bad"
	} else if EyeLeft < 4.9 {
		test = "little_bad"
	} else {
		test = "good"
	}
	if _, err := tx.Model(exam).WherePK().Set(`"eye_left"=?`, EyeLeft).Set(`"eye_left_status"=?`, test).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 上传 视力2 右眼 信息 */
func CreateKindergartenStudentMedicalExaminationNewRSight(studentID int64, EyeRight float32) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	// lsstatus := GetSightSStatus(x.ParseFloat64(ls))
	// lcstatus := GetSightCStatus(x.ParseFloat64(lc))
	// rsstatus := GetSightSStatus(x.ParseFloat64(rs))
	// rcstatus := GetSightCStatus(x.ParseFloat64(rc))
	var test string
	if EyeRight < 4.5 {
		test = "very_bad"
	} else if EyeRight < 4.9 {
		test = "little_bad"
	} else {
		test = "good"
	}
	if _, err := tx.Model(exam).WherePK().Set(`"eye_right"=?`, EyeRight).Set(`"eye_right_status"=?`, test).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 上传 牙齿 信息 */
func CreateKindergartenStudentMedicalExaminationTooth(studentID int64, tooth, caries int, remark string) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	if _, err := tx.Model(exam).WherePK().Set(`"tooth_count"=?`, tooth).
		Set(`"tooth_caries_count"=?`, caries).
		Set(`"tooth_remark"=?`, remark).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"tooth_updated_at"=CURRENT_TIMESTAMP`).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

func FindKindergartenStudentMedicalExaminations(query string, kindergartenID, classID, studentID int64,
	heightFilters, weightFilters, hemoglobinFilters, sightFilters, altFilters, bmiFilters []string,
	startTime, endTime time.Time, page, pageSize int) ([]*KindergartenStudentMedicalExamination, int, error) {

	exams := make([]*KindergartenStudentMedicalExamination, 0)
	q := db.PG().Model(&exams).Relation(`Student`).Relation(`Kindergarten`).Relation(`Student.Class`).Where(`NOT "kindergarten_student_medical_examination"."deleted"`)

	if query != "" {
		q = q.Where(`Student.name ILIKE ?`, fmt.Sprintf("%%%s%%", query))
	}
	if kindergartenID > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."kindergarten_id"=?`, kindergartenID)
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
		q = q.Where(`"kindergarten_student_medical_examination"."date">=?`, startTime)
	}
	if !endTime.IsZero() {
		y, m, d := endTime.Local().Date()
		endTime = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		q = q.Where(`"kindergarten_student_medical_examination"."date"<?`, endTime)
	}

	if len(heightFilters) > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."height_status" IN (?)`, pg.In(heightFilters))
	}

	if len(weightFilters) > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."weight_status" IN (?)`, pg.In(weightFilters))
	}

	if len(hemoglobinFilters) > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."hemoglobin_status" IN (?)`, pg.In(hemoglobinFilters))
	}

	if len(sightFilters) > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."sight_updated_at" IS NOT NULL`)
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, filter := range sightFilters {
				if filter == "normal" {
					q = q.WhereOr(`
					"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND
					"kindergarten_student_medical_examination"."sight_l_c_status" = ? AND
					"kindergarten_student_medical_examination"."sight_r_s_status" = ? AND
					"kindergarten_student_medical_examination"."sight_r_c_status" = ?`,
						StandardStatusNormal, StandardStatusNormal, StandardStatusNormal, StandardStatusNormal)
				} else if filter == "short" { // 近视
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND "kindergarten_student_medical_examination"."sight_r_s_status" = ?`, StandardStatusLow, StandardStatusLow)
				} else if filter == "lshort" { // 左眼近视
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND "kindergarten_student_medical_examination"."sight_r_s_status" = ?`, StandardStatusLow, StandardStatusNormal)
				} else if filter == "rshort" { // 右眼近视
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND "kindergarten_student_medical_examination"."sight_r_s_status" = ?`, StandardStatusNormal, StandardStatusLow)
				} else if filter == "long" { // 远视
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND "kindergarten_student_medical_examination"."sight_r_s_status" = ?`, StandardStatusHigh, StandardStatusHigh)
				} else if filter == "llong" { // 左眼远视
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND "kindergarten_student_medical_examination"."sight_r_s_status" = ?`, StandardStatusHigh, StandardStatusNormal)
				} else if filter == "rlong" { // 右眼远视
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_s_status" = ? AND "kindergarten_student_medical_examination"."sight_r_s_status" = ?`, StandardStatusNormal, StandardStatusHigh)
				} else if filter == "ast" { // 散光
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_c_status" != ? AND "kindergarten_student_medical_examination"."sight_r_c_status" != ?`, StandardStatusNormal, StandardStatusNormal)
				} else if filter == "last" { // 左眼散光
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_c_status" != ? AND "kindergarten_student_medical_examination"."sight_r_c_status" = ?`, StandardStatusNormal, StandardStatusNormal)
				} else if filter == "rast" { // 右眼散光
					q = q.WhereOr(`"kindergarten_student_medical_examination"."sight_l_c_status" = ? AND "kindergarten_student_medical_examination"."sight_r_c_status" != ?`, StandardStatusNormal, StandardStatusNormal)
				}
			}
			return q, nil
		})
	}

	if len(altFilters) > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."alt_status" IN (?)`, pg.In(altFilters))
	}

	if len(bmiFilters) > 0 {
		q = q.Where(`"kindergarten_student_medical_examination"."bmi_status" IN (?)`, pg.In(bmiFilters))
	}

	q = q.Order(`kindergarten_student_medical_examination.updated_at DESC`)

	if page > 0 && pageSize > 0 {

		total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).SelectAndCountEstimate(100000)
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, 0, err
		}
		return exams, total, nil
	} else {
		err := q.Limit(10000).Select()
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return nil, 0, err
		}
		return exams, 0, nil
	}
}

func DeleteKindergartenStudentMedicalExaminationsByStudentTx(tx *pg.Tx, studentID int64) error {
	if _, err := tx.Model(&KindergartenStudentMedicalExamination{}).Where(`"student_id"=?`, studentID).
		Where(`NOT "deleted"`).Set(`"deleted"=?`, true).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

type KindergartenStudentMedicalExaminationALTData struct {
	NO     string  `json:"no" xlsx:"学号"`
	Name   string  `json:"name" xlsx:"姓名"`
	ALT    float64 `json:"alt" xlsx:"ALT"`
	Remark string  `json:"remark" xlsx:"样本号"`
}

/* 批量上传 谷丙转氨酶 信息 */
func BatchCreateKindergartenStudentMedicalExaminationALT(classID int64, datas []*KindergartenStudentMedicalExaminationALTData) error {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	defer tx.Rollback()

	for _, data := range datas {
		student := kindergartendb.KindergartenStudent{}
		var err error
		if data.NO != "" {
			err = tx.Model(&student).Where(`"class_id"=?`, classID).Where(`"no"=?`, data.NO).Select()
		} else {
			err = tx.Model(&student).Where(`"class_id"=?`, classID).Where(`"name"=?`, data.Name).Select()
		}
		if err != nil {
			log.Errorf("DB Error: %v", err)
			return err
		}
		if _, err = CreateKindergartenStudentMedicalExaminationALTTx(tx, student.ID, data.ALT, data.Remark); err != nil {
			log.Errorf("CreateKindergartenStudentMedicalExaminationALTTx Error: %v", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

/* 上传 谷丙转氨酶 信息 */
func CreateKindergartenStudentMedicalExaminationALTTx(tx *pg.Tx, studentID int64, alt float64, remark string) (*KindergartenStudentMedicalExamination, error) {
	exam, err := getKindergartenStudentMedicalExaminationTodayTx(tx, studentID)
	if err != nil {
		log.Errorf("getKindergartenStudentMedicalExaminationByToday() Error: %v", err)
		return nil, err
	}

	status := GetALTStatus(alt)

	if _, err := tx.Model(exam).WherePK().Set(`"alt"=?`, alt).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Set(`"alt_updated_at"=CURRENT_TIMESTAMP`).
		Set(`"alt_remark"=?`, remark).Set(`"alt_status"=?`, status).
		Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

/* 上传 谷丙转氨酶 信息 */
func CreateKindergartenStudentMedicalExaminationALT(studentID int64, alt float64, remark string) (*KindergartenStudentMedicalExamination, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	exam, err := CreateKindergartenStudentMedicalExaminationALTTx(tx, studentID, alt, remark)
	if err != nil {
		log.Errorf("CreateKindergartenStudentMedicalExaminationALTTx Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return exam, nil
}

type KindergartenStudentMedicalExaminationHeightVisionData struct {
	Height int    `json:"height" pg:"h"`
	Gender string `json:"gender" pg:"gender"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentMedicalExaminationHeightVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentMedicalExaminationHeightVisionData, error) {
	datas := make([]*KindergartenStudentMedicalExaminationHeightVisionData, 0)

	sql := `
	SELECT FLOOR("height"/10) AS "h","student"."gender" AS "gender", COUNT(1) AS "count"
	FROM "kindergarten_student_medical_examination" AS "exam"
	INNER JOIN "kindergarten_student" AS "student" ON "student"."id"="exam"."student_id"
	WHERE NOT "exam"."deleted" AND "exam"."height_updated_at" IS NOT NULL
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"exam"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"exam"."kindergarten_id"=?`)
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

type KindergartenStudentMedicalExaminationWeightVisionData struct {
	Weight int    `json:"weight" pg:"w"`
	Gender string `json:"gender" pg:"gender"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentMedicalExaminationWeightVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentMedicalExaminationWeightVisionData, error) {
	datas := make([]*KindergartenStudentMedicalExaminationWeightVisionData, 0)

	sql := `
	SELECT FLOOR("weight"/5) AS "w","student"."gender" AS "gender", COUNT(1) AS "count"
	FROM "kindergarten_student_medical_examination" AS "exam"
	INNER JOIN "kindergarten_student" AS "student" ON "student"."id"="exam"."student_id"
	WHERE NOT "exam"."deleted" AND "exam"."weight_updated_at" IS NOT NULL
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"exam"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"exam"."kindergarten_id"=?`)
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

func FindKindergartenStudentMedicalExaminationDates(kindergartenID, classID int64) ([]time.Time, error) {

	type T struct {
		Date time.Time `json:"date" pg:"date"`
	}
	ts := make([]*T, 0)

	sql := `
		SELECT "date" FROM "kindergarten_student_medical_examination" AS "exam"
		INNER JOIN "kindergarten_student" AS "student" ON "student"."id"="exam"."student_id"
		WHERE NOT "exam"."deleted"
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)
	if kindergartenID > 0 {
		wheres = append(wheres, `"exam"."kindergarten_id"=?`)
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

func FindKindergartenStudentMedicalExaminationBMIVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentMedicalExaminationStatusVisionData, error) {
	return FindKindergartenStudentMedicalExaminationStatusVision("bmi", kindergartenID, classID, date)
}

func FindKindergartenStudentMedicalExaminationHemoglobinVision(kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentMedicalExaminationStatusVisionData, error) {
	return FindKindergartenStudentMedicalExaminationStatusVision("hemoglobin", kindergartenID, classID, date)
}

type KindergartenStudentMedicalExaminationStatusVisionData struct {
	Status string `json:"status" pg:"status"`
	Count  int    `json:"count" pg:"count"`
}

func FindKindergartenStudentMedicalExaminationStatusVision(fieldName string, kindergartenID, classID int64, date time.Time) ([]*KindergartenStudentMedicalExaminationStatusVisionData, error) {
	datas := make([]*KindergartenStudentMedicalExaminationStatusVisionData, 0)

	sql := `
		SELECT "` + fieldName + `_status" AS "status",COUNT(1) AS "count" FROM
			"kindergarten_student_medical_examination" AS "exam" INNER JOIN
			"kindergarten_student" AS "student" ON "student"."id"="exam"."student_id"
			WHERE NOT "exam"."deleted" AND NOT "student"."deleted" AND "exam"."` + fieldName + `_updated_at" IS NOT NULL
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"exam"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"exam"."kindergarten_id"=?`)
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

type FindKindergartenStudentMedicalExaminationSightVisionData struct {
	LSStatus string `json:"lsstatus" pg:"lsstatus"`
	LCStatus string `json:"lcstatus" pg:"lcstatus"`
	RSStatus string `json:"rsstatus" pg:"rsstatus"`
	RCStatus string `json:"rcstatus" pg:"rcstatus"`
	Status   string `json:"status"`
	Count    int    `json:"count" pg:"count"`
}

func FindKindergartenStudentMedicalExaminationSightVision(kindergartenID, classID int64, date time.Time) ([]*FindKindergartenStudentMedicalExaminationSightVisionData, error) {
	datas := make([]*FindKindergartenStudentMedicalExaminationSightVisionData, 0)

	sql := `
		SELECT "sight_l_s_status" AS "lsstatus","sight_r_s_status" AS "rsstatus",
			"sight_l_c_status" AS "lcstatus","sight_r_c_status" AS "rcstatus",
			COUNT(1) AS "count" FROM
			"kindergarten_student_medical_examination" AS "exam" INNER JOIN
			"kindergarten_student" AS "student" ON "student"."id"="exam"."student_id"
			WHERE NOT "exam"."deleted" AND NOT "student"."deleted" AND "exam"."sight_updated_at" IS NOT NULL
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"exam"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"exam"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}
	sql += ` AND ` + strings.Join(wheres, ` AND `)

	sql += ` GROUP BY "lsstatus","rsstatus","lcstatus","rcstatus" ORDER BY "lsstatus" DESC, "rsstatus" DESC`

	if _, err := db.PG().Query(&datas, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return datas, nil
}

type FindKindergartenStudentMedicalExaminationEyeVisionData struct {
	EyeLeftStatus  string `json:"lstatus" pg:"lstatus"`
	EyeRightStatus string `json:"rstatus" pg:"rstatus"`
	Status         string `json:"status"`
	Count          int    `json:"count" pg:"count"`
}

func FindKindergartenStudentMedicalExaminationEyeVision(kindergartenID, classID int64, date time.Time) ([]*FindKindergartenStudentMedicalExaminationEyeVisionData, error) {
	datas := make([]*FindKindergartenStudentMedicalExaminationEyeVisionData, 0)

	sql := `
		SELECT "eye_left_status" AS "lstatus","eye_right_status" AS "rstatus",
			COUNT(1) AS "count" FROM
			"kindergarten_student_medical_examination" AS "exam" INNER JOIN
			"kindergarten_student" AS "student" ON "student"."id"="exam"."student_id"
			WHERE NOT "exam"."deleted" AND NOT "student"."deleted" AND ( "exam"."eye_left_status" !='' OR "exam"."eye_right_status" !='' )
	`
	wheres := make([]string, 0)
	params := make([]interface{}, 0)

	wheres = append(wheres, `"exam"."date"=?`)
	params = append(params, date)

	if kindergartenID > 0 {
		wheres = append(wheres, `"exam"."kindergarten_id"=?`)
		params = append(params, kindergartenID)
	}
	if classID != 0 {
		wheres = append(wheres, `"student"."class_id"=?`)
		params = append(params, classID)
	}
	sql += ` AND ` + strings.Join(wheres, ` AND `)

	sql += ` GROUP BY "lstatus","rstatus" ORDER BY "lstatus" DESC, "rstatus" DESC`

	if _, err := db.PG().Query(&datas, sql, params...); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return datas, nil

}
