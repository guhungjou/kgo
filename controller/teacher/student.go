package teacher

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tealeg/xlsx/v3"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/errors"
	"gitlab.com/ykgk/kgo/db"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type CreateKindergartenStudentRequest struct {
	NO       string    `json:"no" form:"no" query:"no"`
	Name     string    `json:"name" form:"name" query:"name" validate:"required"`
	Gender   string    `json:"gender" form:"gender" query:"gender" validate:"required"`
	Remark   string    `json:"remark" form:"remark" query:"remark"`
	Birthday time.Time `json:"birthday" form:"birthday" query:"birthday" validate:"required"`
	Device   string    `json:"device" form:"device" query:"device" validate:"required"`
	ClassID  int64     `json:"class_id" form:"class_id" query:"class_id"`
}

/* 创建学生 */
func CreateKindergartenStudent(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenStudentRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Device = x.FormatMacAddress(req.Device)

	student, err := kindergartendb.GetKindergartenStudentByDevice(req.Device)
	if err != nil {
		return ctx.InternalServerError()
	} else if student != nil {
		/* 设备重复 */
		return ctx.Fail(errors.ErrKindergartenStudentDeviceDuplicated, nil)
	}

	if req.NO != "" {
		s, err := kindergartendb.GetKindergartenStudentByNO(req.ClassID, req.NO)
		if err != nil {
			return ctx.InternalServerError()
		} else if s != nil {
			/* 学号重复 */
			return ctx.Fail(errors.ErrKindergartenStudentNODuplicated, nil)
		}
	}

	// if !x.ValidateMacAddress(req.Device) {
	// 	return ctx.Fail(errors.ErrKindergartenStudentDeviceInvalid, nil)
	// }

	student, err = kindergartendb.CreateKindergartenStudent(req.Name, req.NO, req.Remark, req.Gender, req.Birthday, req.Device, ctx.Teacher.KindergartenID, req.ClassID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(student)
}

type UpdateKindergartenStudentRequest struct {
	NO       string    `json:"no" form:"no" query:"no"`
	Name     string    `json:"name" form:"name" query:"name" validate:"required"`
	Gender   string    `json:"gender" form:"gender" query:"gender" validate:"required"`
	Birthday time.Time `json:"birthday" form:"birthday" query:"birthday" validate:"required"`
	Remark   string    `json:"remark" form:"remark" query:"remark"`
	Device   string    `json:"device" form:"device" query:"device" validate:"required"`
	ClassID  int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func UpdateKindergartenStudent(c echo.Context) error {
	ctx := c.(*context.Context)
	req := UpdateKindergartenStudentRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	id := ctx.IntParam(`id`)

	student, err := kindergartendb.GetKindergartenStudent(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil || student.KindergartenID != ctx.Teacher.KindergartenID {
		return ctx.NotFound()
	}

	if req.NO != "" {
		s, err := kindergartendb.GetKindergartenStudentByNO(req.ClassID, req.NO)
		if err != nil {
			return ctx.InternalServerError()
		} else if s != nil && s.ID != student.ID {
			/* 学号重复 */
			return ctx.Fail(errors.ErrKindergartenStudentNODuplicated, nil)
		}
	}

	req.Device = x.FormatMacAddress(req.Device)

	if student.Device != req.Device {
		/* 要更改设备，则检查新设备是否已存在 */
		t, err := kindergartendb.GetKindergartenStudentByDevice(req.Device)
		if err != nil {
			return ctx.InternalServerError()
		} else if t != nil {
			/* 设备重复 */
			return ctx.Fail(errors.ErrKindergartenStudentDeviceDuplicated, nil)
		}
		/* 检查设备格式 */
		// if !x.ValidateMacAddress(req.Device) {
		// 	return ctx.Fail(errors.ErrKindergartenStudentDeviceInvalid, nil)
		// }
	}

	student, err = kindergartendb.UpdateKindergartenStudent(id, req.Name, req.NO, req.Remark, req.Gender, req.Birthday, req.Device, req.ClassID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(student)
}

type FindKindergartenStudentsRequest struct {
	Query    string `json:"query" form:"query" query:"query"`
	Gender   string `json:"gender" form:"gender" query:"gender"`
	ClassID  int64  `json:"class_id" form:"class_id" query:"class_id"`
	Page     int    `json:"page" form:"page" query:"page"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindKindergartenStudents(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	students, total, err := kindergartendb.FindKindergartenStudents(req.Query, req.Gender, ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(students, req.Page, req.PageSize, total)
}

/* 导出学生 */
func ExportKindergartenStudents(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	students, _, err := kindergartendb.FindKindergartenStudents(req.Query, req.Gender, ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), 0, 0)
	if err != nil {
		return ctx.InternalServerError()
	}
	headers := []string{"ID", "学号", "姓名", "性别", "出生年月", "年龄", "班级", "设备", "备注", "创建时间"}
	rows := make([][]interface{}, 0)
	for _, s := range students {
		row := make([]interface{}, 0)
		row = append(row, s.ID)
		row = append(row, s.NO)
		row = append(row, s.Name)
		row = append(row, s.GenderName())
		row = append(row, s.Birthday)
		row = append(row, s.Age())
		if s.Class != nil {
			row = append(row, s.Class.Name)
		} else {
			row = append(row, "(无)")
		}
		row = append(row, s.Device)
		row = append(row, s.Remark)
		row = append(row, s.CreatedAt)
		rows = append(rows, row)
	}
	return ctx.XLSX("学生列表", headers, rows)
}

/* 下载上传学生的模板 */
func DownloadKindergartenStudentTemplate(c echo.Context) error {
	ctx := c.(*context.Context)

	headers := []string{"学号", "姓名", "性别", "出生年月", "备注", "班级", "MAC地址"}
	rows := make([][]interface{}, 0)

	return ctx.XLSX("学生模板", headers, rows)
}

type LoadKindergartenStudentResult struct {
	NO        string    `json:"no" xlsx:"学号"`
	Name      string    `json:"name" xlsx:"姓名" validate:"required"`
	Gender    string    `json:"gender" xlsx:"性别" validate:"required"`
	ClassName string    `json:"class_name" xlsx:"班级"`
	Remark    string    `json:"remark" xlsx:"备注"`
	Device    string    `json:"device" xlsx:"MAC地址"`
	Birthday  time.Time `json:"birthday" xlsx:"出生年月,2006/1/2" validate:"required"`

	ClassID int64                             `json:"class_id" validate:"gt=0"`
	Class   *kindergartendb.KindergartenClass `json:"class,omitempty"`

	Status []string `json:"status"`
}

/* 解析XLSX文件，返回老师信息 */
func LoadKindergartenStudent(c echo.Context) error {
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

	results := make([]*LoadKindergartenStudentResult, 0)

	if err := x.ParseXLSXSheet(sheet, &results); err != nil {
		return ctx.BadRequest()
	}

	devices := make(map[string]bool)
	for _, r := range results {
		r.Status = []string{}

		if r.Name == "" {
			r.Status = append(r.Status, "NoName")
		}
		if r.Gender == "男" {
			r.Gender = "male"
		} else if r.Gender == "女" {
			r.Gender = "female"
		} else {
			r.Status = append(r.Status, "NoGender")
		}

		if r.ClassName != "" {
			class, err := kindergartendb.GetKindergartenClassByName(ctx.Teacher.KindergartenID, r.ClassName)
			if err != nil {
				return ctx.InternalServerError()
			} else if class == nil {
				r.Status = append(r.Status, "NoClass")
			} else {
				r.Class = class
				r.ClassID = class.ID
				if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager && ctx.Teacher.ClassID != class.ID {
					r.Status = append(r.Status, "WrongClass")
				}
			}
		} else {
			r.Status = append(r.Status, "NoClass")
		}

		r.Device = x.FormatMacAddress(r.Device)
		if r.Device == "" {
			r.Status = append(r.Status, "NoDevice")
		} else if devices[r.Device] {
			r.Status = append(r.Status, "DuplicatedDevice")
		} else {
			student, err := kindergartendb.GetKindergartenStudentByDevice(r.Device)
			if err != nil {
				return ctx.InternalServerError()
			} else if student != nil {
				r.Status = append(r.Status, "DuplicatedDevice")
			}
			// else if !x.ValidateMacAddress(r.Device) {
			// 	r.Status = append(r.Status, "InvalidDevice")
			// }
		}

		if r.Birthday.IsZero() {
			r.Status = append(r.Status, "NoBirthday")
		}

		if len(r.Status) == 0 {
			devices[r.Device] = true
		}
	}
	return ctx.Success(results)
}

type CreateKindergartenStudentLoadRequest struct {
	Students []*LoadKindergartenStudentResult `json:"students" validate:"gt=0"`
}

func CreateKindergartenStudentLoad(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenStudentLoadRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	tx, err := db.Begin()
	if err != nil {
		return ctx.InternalServerError()
	}
	defer tx.Rollback()

	for _, t := range req.Students {
		// if !x.ValidateMacAddress(t.Device) {
		// 	/* 检查设备格式 */
		// 	return ctx.Fail(errors.ErrKindergartenStudentDeviceInvalid, nil)
		// }
		t.Device = x.FormatMacAddress(t.Device)
		student, err := kindergartendb.GetKindergartenStudentByDevice(t.Device)
		if err != nil {
			return ctx.InternalServerError()
		} else if student != nil {
			/* 设备重复 */
			return ctx.Fail(errors.ErrKindergartenStudentDeviceDuplicated, nil)
		}

		if t.NO != "" {
			student, err := kindergartendb.GetKindergartenStudentByNOTx(tx, t.ClassID, t.NO)
			if err != nil {
				return ctx.InternalServerError()
			} else if student != nil {
				/* 学号重复 */
				return ctx.Fail(errors.ErrKindergartenStudentNODuplicated, nil)
			}
		}

		if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager && ctx.Teacher.ClassID != t.ClassID {
			return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
		}

		student, err = kindergartendb.CreateKindergartenStudentTx(tx, t.Name, t.NO, t.Remark, t.Gender, t.Birthday, t.Device, ctx.Teacher.KindergartenID, t.ClassID)
		if err != nil {
			return ctx.InternalServerError()
		}
	}

	if err := tx.Commit(); err != nil {
		return ctx.InternalServerError()
	}

	return ctx.Success(nil)
}

func DeleteKindergartenStudent(c echo.Context) error {
	ctx := c.(*context.Context)
	id := ctx.IntParam(`id`)

	student, err := kindergartendb.GetKindergartenStudent(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if student == nil {
		return ctx.NotFound()
	} else if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager && ctx.Teacher.ClassID != student.ClassID {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	if err := healthdb.DeleteStudentWithCheckAndExam(student.ID); err != nil {
		return ctx.InternalServerError()
	}

	return ctx.Success(student)
}
