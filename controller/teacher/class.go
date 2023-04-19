package teacher

import (
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/errors"
	"gitlab.com/ykgk/kgo/db"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/tealeg/xlsx/v3"
)

type FindKindergartenClassesRequest struct {
	Query    string `json:"query" form:"query" query:"query"`
	Page     int    `json:"page" form:"page" query:"page"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindKindergartenClasses(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenClassesRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	classes, total, err := kindergartendb.FindKindergartenClasses(req.Query, ctx.Teacher.KindergartenID, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(classes, req.Page, req.PageSize, total)
}

type CreateKindergartenClassRequest struct {
	Name   string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Remark string `json:"remark" form:"remark" query:"remark"`
}

func CreateKindergartenClass(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenClassRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	class, err := kindergartendb.GetKindergartenClassByName(ctx.Teacher.KindergartenID, req.Name)
	if err != nil {
		return ctx.InternalServerError()
	} else if class != nil {
		return ctx.Fail(errors.ErrKindergartenClassNameDuplicated, nil)
	}

	class, err = kindergartendb.CreateKindergartenClass(req.Name, req.Remark, ctx.Teacher.KindergartenID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(class)
}

type UpdateKindergartenClassRequest struct {
	Name   string `json:"name" form:"name" query:"name" validate:"required"`
	Remark string `json:"remark" form:"remark" query:"remark"`
}

func UpdateKindergartenClass(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	req := UpdateKindergartenClassRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	class, err := kindergartendb.GetKindergartenClass(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if class == nil || class.KindergartenID != ctx.Teacher.KindergartenID {
		return ctx.NotFound()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	if class.Name != req.Name { /* 修改了班级名，检查班级名是否重复 */
		old, err := kindergartendb.GetKindergartenClassByName(ctx.Teacher.KindergartenID, req.Name)
		if err != nil {
			return ctx.InternalServerError()
		} else if old != nil {
			return ctx.Fail(errors.ErrKindergartenClassNameDuplicated, nil)
		}
	}

	class, err = kindergartendb.UpdateKindergartenClass(class.ID, req.Name, req.Remark)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(class)
}

func GetKindergartenClass(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	class, err := kindergartendb.GetKindergartenClass(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if class == nil {
		return ctx.NotFound()
	}
	return ctx.Success(class)
}

/* 下载上传班级的模板 */
func DownloadKindergartenClassTemplate(c echo.Context) error {
	ctx := c.(*context.Context)

	headers := []string{"班级名", "备注"}
	rows := make([][]interface{}, 0)

	return ctx.XLSX("班级模板", headers, rows)
}

type LoadKindergartenClassResult struct {
	Name   string   `json:"name" xlsx:"班级名" form:"name" validate:"gt=0"`
	Remark string   `json:"remark" xlsx:"备注" form:"name"`
	Status []string `json:"status"`
}

/* 解析XLSX文件，返回班级信息 */
func LoadKindergartenClass(c echo.Context) error {
	ctx := c.(*context.Context)

	fileheader, err := ctx.FormFile("file")
	if err != nil || fileheader == nil {
		return ctx.BadRequest()
	}
	file, err := fileheader.Open()
	if err != nil || file == nil {
		return ctx.BadRequest()
	}
	defer file.Close()

	wb, err := xlsx.OpenReaderAt(file, fileheader.Size)
	if err != nil || wb == nil || len(wb.Sheets) == 0 {
		return ctx.BadRequest()
	}
	sheet := wb.Sheets[0]

	results := make([]*LoadKindergartenClassResult, 0)

	if err := x.ParseXLSXSheet(sheet, &results); err != nil {
		return ctx.BadRequest()
	}

	names := make(map[string]bool)

	for _, r := range results {
		r.Status = []string{}
		if r.Name == "" {
			r.Status = append(r.Status, "NoName")
		} else {
			class, err := kindergartendb.GetKindergartenClassByName(ctx.Teacher.KindergartenID, r.Name)
			if err != nil {
				return ctx.InternalServerError()
			} else if class != nil {
				r.Status = append(r.Status, "Duplicated")
			} else if names[r.Name] {
				r.Status = append(r.Status, "Duplicated")
			}
		}
		if len(r.Status) == 0 {
			names[r.Name] = true
		}
	}
	return ctx.Success(results)
}

type CreateKindergartenClassLoadRequest struct {
	Classes []*LoadKindergartenClassResult `json:"classes" form:"classes" query:"classes" validate:"gt=0"`
}

/* 批量创建班级 */
func CreateKindergartenClassLoad(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenClassLoadRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	for _, r := range req.Classes {
		class, err := kindergartendb.GetKindergartenClassByName(ctx.Teacher.KindergartenID, r.Name)
		if err != nil {
			return ctx.InternalServerError()
		} else if class != nil {
			return ctx.Fail(errors.ErrKindergartenClassNameDuplicated, nil)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return ctx.InternalServerError()
	}
	defer tx.Rollback()

	for _, r := range req.Classes {
		_, err = kindergartendb.CreateKindergartenClassTx(tx, r.Name, r.Remark, ctx.Teacher.KindergartenID)
		if err != nil {
			return ctx.InternalServerError()
		}
	}

	if err := tx.Commit(); err != nil {
		return ctx.InternalServerError()
	}

	return ctx.Success(nil)
}

func FindSelfClasses(c echo.Context) error {
	ctx := c.(*context.Context)

	if ctx.Teacher.Role == kindergartendb.KindergartenTeacherRoleManager {
		return FindKindergartenClasses(c)
	}
	if ctx.Teacher.Class != nil {
		return ctx.List([]*kindergartendb.KindergartenClass{ctx.Teacher.Class}, 1, 10, 1)
	}
	return ctx.List([]*kindergartendb.KindergartenClass{}, 1, 10, 0)
}

type DeleteKindergartenClassRequest struct {
	WithStudent bool `json:"with_student" form:"with_student" query:"with_student"`
	WithTeacher bool `json:"with_teacher" form:"with_teacher" query:"with_teacher"`
}

func DeleteKindergartenClass(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	req := DeleteKindergartenClassRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	tx, err := db.Begin()
	if err != nil {
		return ctx.InternalServerError()
	}
	defer tx.Rollback()

	class, err := kindergartendb.DeleteKindergartenClassTx(tx, id)
	if err != nil {
		return ctx.InternalServerError()
	} else if class == nil {
		return ctx.Success(nil)
	}

	/* 同时删除学生 */
	if req.WithStudent {
		if err := healthdb.DeleteStudentWithCheckAndExamByClassTx(tx, class.ID); err != nil {
			return ctx.InternalServerError()
		}
	} else {
		/* 不删除学生，把学生的班级设置为空 */
		if err := kindergartendb.ClearKindergartenStudentClassTx(tx, class.ID); err != nil {
			return ctx.InternalServerError()
		}
	}
	if req.WithTeacher {
		if err := kindergartendb.DeleteKindergartenTeacherByClassTx(tx, class.ID); err != nil {
			return ctx.InternalServerError()
		}
	} else {
		/* 不删除老师，将老师的班级设置成空 */
		if err := kindergartendb.ClearKindergartenTeacherClassTx(tx, class.ID); err != nil {
			return ctx.InternalServerError()
		}
	}

	if err := tx.Commit(); err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(nil)
}

type BatchCreateKindergartenStudentMedicalExaminationALTRequest struct {
	ClassID int64 `json:"class_id" form:"class_id" query:"class_id" validate:"gt=0"`
}

/* 批量创建ALT体检数据 */
func BatchCreateKindergartenStudentMedicalExaminationALT(c echo.Context) error {
	ctx := c.(*context.Context)

	req := BatchCreateKindergartenStudentMedicalExaminationALTRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	fileheader, err := ctx.FormFile("file")
	if err != nil || fileheader == nil {
		return ctx.BadRequest()
	}
	file, err := fileheader.Open()
	if err != nil || file == nil {
		return ctx.BadRequest()
	}
	defer file.Close()

	wb, err := xlsx.OpenReaderAt(file, fileheader.Size)
	if err != nil || wb == nil || len(wb.Sheets) == 0 {
		return ctx.BadRequest()
	}
	sheet := wb.Sheets[0]

	datas := make([]*healthdb.KindergartenStudentMedicalExaminationALTData, 0)

	if err := x.ParseXLSXSheet(sheet, &datas); err != nil {
		return ctx.BadRequest()
	}

	for _, data := range datas {
		if data.NO != "" {
			if student, err := kindergartendb.GetKindergartenStudentByNO(req.ClassID, data.NO); err != nil {
				return ctx.InternalServerError()
			} else if student == nil {
				return ctx.Fail(errors.ErrKindergartenStudentNONotFound, data.NO)
			}
		} else if data.Name == "" {
			return ctx.Fail(errors.ErrKindergartenStudentNameNotFound, nil)
		} else {
			if student, err := kindergartendb.GetKindergartenStudentByName(req.ClassID, data.Name); err != nil {
				if err != pg.ErrMultiRows {
					return ctx.InternalServerError()
				}
				return ctx.Fail(errors.ErrKindergartenStudentNameDuplicated, data.Name)
			} else if student == nil {
				return ctx.Fail(errors.ErrKindergartenStudentNameNotFound, data.Name)
			}
		}
	}

	err = healthdb.BatchCreateKindergartenStudentMedicalExaminationALT(req.ClassID, datas)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(nil)
}
