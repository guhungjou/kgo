package teacher

import (
	con "context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	openai "github.com/sashabaranov/go-openai"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenStudentMedicalExaminationsRequest struct {
	Query     string    `json:"query" form:"query" query:"query"`
	StudentID int64     `json:"student_id" form:"student_id" query:"student_id"`
	ClassID   int64     `json:"class_id" form:"class_id" query:"class_id"`
	StartTime time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime   time.Time `json:"end_time" form:"end_time" query:"end_time"`
	Page      int       `json:"page" form:"page" query:"page"`
	PageSize  int       `json:"page_size" form:"page_size" query:"page_size"`

	HeightFilters     []string `json:"height_filters" form:"height_filters" query:"height_filters"`
	WeightFilters     []string `json:"weight_filters" form:"weight_filters" query:"weight_filters"`
	HemoglobinFilters []string `json:"hemoglobin_filters" form:"hemoglobin_filters" query:"hemoglobin_filters"`
	SightFilters      []string `json:"sight_filters" form:"sight_filters" query:"sight_filters"`
	ALTFilters        []string `json:"alt_filters" form:"alt_filters" query:"alt_filters"`
	BMIFilters        []string `json:"bmi_filters" form:"bmi_filters" query:"bmi_filters"`
}

func FindKindergartenStudentMedicalExaminations(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentMedicalExaminationsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)
	// fmt.Print("=====================", (req.HemoglobinFilters))
	exams, total, err := healthdb.FindKindergartenStudentMedicalExaminations(
		req.Query, ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), req.StudentID,
		req.HeightFilters, req.WeightFilters, req.HemoglobinFilters, req.SightFilters, req.ALTFilters, req.BMIFilters,
		req.StartTime, req.EndTime, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(exams, req.Page, req.PageSize, total)
}

func ExportKindergartenStudentMedicalExaminations(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentMedicalExaminationsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	exams, _, err := healthdb.FindKindergartenStudentMedicalExaminations(
		req.Query, ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), req.StudentID,
		req.HeightFilters, req.WeightFilters, req.HemoglobinFilters, req.SightFilters, req.ALTFilters, req.BMIFilters,
		req.StartTime, req.EndTime, 0, 0)
	if err != nil {
		return ctx.InternalServerError()
	}

	headers := []string{"日期", "班级", "学生", "性别", "身高(cm)", "体重(kg)", "血红蛋白(Hb)", "谷丙转氨酶(U/L)", "视力(左)", "视力(右)", "牙齿", "龋齿", "创建时间"}
	rows := make([][]interface{}, 0)
	for _, exam := range exams {
		row := make([]interface{}, 0)
		row = append(row, fmt.Sprintf("%04d-%02d-%02d", exam.Date.Year(), exam.Date.Month(), exam.Date.Day()))
		row = append(row, exam.Student.Class.Name)
		row = append(row, exam.Student.Name)
		row = append(row, exam.Student.GenderName())
		if !exam.HeightUpdatedAt.IsZero() {
			row = append(row, exam.Height)
			row = append(row, exam.Weight)
		} else {
			row = append(row, "###")
			row = append(row, "###")
		}
		if !exam.HemoglobinUpdatedAt.IsZero() {
			row = append(row, exam.Hemoglobin)
		} else {
			row = append(row, "###")
		}
		if !exam.ALTUpdatedAt.IsZero() {
			row = append(row, exam.ALT)
		} else {
			row = append(row, "###")
		}
		if !exam.SightUpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%s/%s", exam.SightLS, exam.SightLC))
			row = append(row, fmt.Sprintf("%s/%s", exam.SightRS, exam.SightRC))
		} else {
			row = append(row, "###")
			row = append(row, "###")
		}
		if !exam.ToothUpdatedAt.IsZero() {
			row = append(row, exam.ToothCount)
			row = append(row, exam.ToothCariesCount)
		} else {
			row = append(row, "###")
			row = append(row, "###")
		}

		row = append(row, exam.CreatedAt)
		rows = append(rows, row)
	}
	return ctx.XLSX("体检记录", headers, rows)
}

func GetKindergartenAnswer(c echo.Context) error {
	ctx := c.(*context.Context)
	// ctx获得参数question
	question := ctx.FormValue(`question`)
	// question := ctx.Param(`question`)
	// fmt.Print("111111111111111" + question)
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
		con.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 0.4,
			MaxTokens:   1000,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					//exam.Student.Age()转换成string类型
					Content: question,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "不要建议定期检测！！！不要用第一人称回答。请注意，小孩本身使用电脑等电子产品的频率就少，因此对于视力有问题的情况，轻度近视请从“减少电子产品使用频率，规范用眼、培养良好阅读和学习姿势”等方面考虑，可以加以扩展；重度近视请从复查、治疗等方面考虑，可以加以扩展。",
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "对于身高、体重，描述数字或区间即可，不需要对高矮胖瘦进行评价，但是要对其他项目（如果有）进行评价，！！！",
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你的回答中不要对饮酒、隐形眼镜等进行建议，你的回答中不要出现：饮酒、normal、根据检测结果、按照、他、您、建议等字眼。",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		// time.Sleep(1 * time.Second)
		fmt.Print("重发请求...")
	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++" + resp.Choices[0].Message.Content)
	return ctx.Success(resp.Choices[0].Message.Content)

}

func GetKindergartenStudentMedicalExamination(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	exam, err := healthdb.GetKindergartenStudentMedicalExamination(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if exam == nil {
		return ctx.NotFound()
	}

	return ctx.Success(exam)
}

type FindKindergartenStudentMedicalExaminationHeightVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentMedicalExaminationHeightVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMedicalExaminationHeightVisionRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			req.ClassID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}

	datas, err := healthdb.FindKindergartenStudentMedicalExaminationHeightVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

/* 获取最近一百个有体检信息的日期 */
func FindKindergartenStudentMedicalExaminationDates(c echo.Context) error {
	ctx := c.(*context.Context)

	var classID int64
	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			classID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}
	dates, err := healthdb.FindKindergartenStudentMedicalExaminationDates(ctx.Teacher.KindergartenID, classID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(dates)
}

type FindKindergartenStudentMedicalExaminationWeightVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentMedicalExaminationWeightVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMedicalExaminationWeightVisionRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			req.ClassID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}

	datas, err := healthdb.FindKindergartenStudentMedicalExaminationWeightVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

type FindKindergartenStudentMedicalExaminationStatusVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentMedicalExaminationStatusVision(c echo.Context, field string) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMedicalExaminationStatusVisionRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			req.ClassID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}

	datas, err := healthdb.FindKindergartenStudentMedicalExaminationStatusVision(field, ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

func FindKindergartenStudentMedicalExaminationBMIVision(c echo.Context) error {
	return FindKindergartenStudentMedicalExaminationStatusVision(c, "bmi")
}

func FindKindergartenStudentMedicalExaminationHemoglobinVision(c echo.Context) error {
	return FindKindergartenStudentMedicalExaminationStatusVision(c, "hemoglobin")
}

func FindKindergartenStudentMedicalExaminationALTVision(c echo.Context) error {
	return FindKindergartenStudentMedicalExaminationStatusVision(c, "alt")
}

func FindKindergartenStudentMedicalExaminationSightVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMedicalExaminationStatusVisionRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			req.ClassID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}

	datas, err := healthdb.FindKindergartenStudentMedicalExaminationSightVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

func FindKindergartenStudentMedicalExaminationEyeVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMedicalExaminationStatusVisionRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			req.ClassID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}

	datas, err := healthdb.FindKindergartenStudentMedicalExaminationEyeVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}
