package device

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

type CreateKindergartenStudentMedicalExaminationHeightRequest struct {
	StudentID int64   `json:"student_id" form:"student_id" query:"student_id" `
	Height    float64 `json:"height" form:"height" query:"height" `
}

func CreateKindergartenStudentMedicalExaminationHeight(c echo.Context) error {
	fmt.Print("=========================================================")
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationHeightRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationHeight(req.StudentID, req.Height)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationWeightRequest struct {
	StudentID int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	Weight    float64 `json:"weight" form:"weight" query:"weight" validate:"gt=0"`
}

func CreateKindergartenStudentMedicalExaminationWeight(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationWeightRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationWeight(req.StudentID, req.Weight)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationHemoglobinRequest struct {
	StudentID  int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	Hemoglobin float64 `json:"hemoglobin" form:"hemoglobin" query:"hemoglobin" validate:"gt=0"`
	Remark     string  `json:"remark" form:"remark" query:"remark"`
}

func CreateKindergartenStudentMedicalExaminationHemoglobin(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationHemoglobinRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationHemoglobin(req.StudentID, req.Hemoglobin, req.Remark)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationSightRequest struct {
	StudentID int64 `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	// SightL       float64 `json:"sight_l" form:"sight_l" query:"sight_l" validate:"gt=0"`
	SightLS      string `json:"sight_l_s" form:"sight_l_s" query:"sight_l_s" validate:"gt=0"`
	SightLC      string `json:"sight_l_c" form:"sight_l_c" query:"sight_l_c" validate:"gt=0"`
	SightLRemark string `json:"sight_l_remark" form:"sight_l_remark" query:"sight_l_remark"`
	// SightR       float64 `json:"sight_r" form:"sight_r" query:"sight_r" validate:"gt=0"`
	SightRS      string `json:"sight_r_s" form:"sight_r_s" query:"sight_r_s" validate:"gt=0"`
	SightRC      string `json:"sight_r_c" form:"sight_r_c" query:"sight_r_c" validate:"gt=0"`
	SightRRemark string `json:"sight_r_remark" form:"sight_r_remark" query:"sight_r_remark"`
}

func CreateKindergartenStudentMedicalExaminationSight(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationSightRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationSight(req.StudentID, req.SightLS, req.SightLC, req.SightRS, req.SightRC,
		req.SightLRemark, req.SightRRemark)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationNewLSightRequest struct {
	StudentID int64 `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	// SightL       float64 `json:"sight_l" form:"sight_l" query:"sight_l" validate:"gt=0"`
	EyeLeft float32 `json:"eye_left" form:"eye_left" query:"eye_left" validate:"gt=0"`
}

func CreateKindergartenStudentMedicalExaminationNewLSight(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationNewLSightRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationNewLSight(req.StudentID, req.EyeLeft)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationNewRSightRequest struct {
	StudentID int64 `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	// SightL       float64 `json:"sight_l" form:"sight_l" query:"sight_l" validate:"gt=0"`
	EyeRight float32 `json:"eye_right" form:"eye_right" query:"eye_right" validate:"gt=0"`
}

func CreateKindergartenStudentMedicalExaminationNewRSight(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationNewRSightRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationNewRSight(req.StudentID, req.EyeRight)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationToothRequest struct {
	StudentID int64  `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	Tooth     int    `json:"tooth" form:"tooth" query:"tooth" validate:"gt=0"`
	Caries    int    `json:"caries" form:"caries" query:"caries"`
	Remark    string `json:"remark" form:"remark" query:"remark"`
}

func CreateKindergartenStudentMedicalExaminationTooth(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationToothRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationTooth(req.StudentID, req.Tooth, req.Caries, req.Remark)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

type CreateKindergartenStudentMedicalExaminationALTRequest struct {
	StudentID int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	ALT       float64 `json:"alt" form:"alt" query:"alt"`
	Remark    string  `json:"remark" form:"remark" query:"remark"`
}

func CreateKindergartenStudentMedicalExaminationALT(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentMedicalExaminationALTRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	student, err := kindergartendb.GetKindergartenStudent(req.StudentID)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	}

	exam, err := healthdb.CreateKindergartenStudentMedicalExaminationALT(req.StudentID, req.ALT, req.Remark)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(exam)
}

func GetKindergartenStudentMedicalExaminationTodayStatus(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	exam, err := healthdb.GetKindergartenStudentMedicalExaminationToday(id)
	if err != nil {
		return ctx.InternalServerError()
	}

	if exam == nil {
		exam = &healthdb.KindergartenStudentMedicalExamination{
			StudentID: id,
		}
	}
	return ctx.Success(exam)
}
