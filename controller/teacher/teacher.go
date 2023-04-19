package teacher

import (
	"time"

	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/db"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/tealeg/xlsx/v3"
	"gitlab.com/ykgk/kgo/controller/errors"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"gt=0"`
	Password string `json:"password" form:"password" query:"password" validate:"gt=0"`
}

/* 管理员帐号登录 */
func Login(c echo.Context) error {
	ctx := c.(*context.Context)
	req := LoginRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	teacher, err := kindergartendb.GetKindergartenTeacherByUsername(req.Username)
	if err != nil {
		return ctx.InternalServerError()
	} else if teacher == nil || teacher.Deleted { /* 用户不存在 */
		return ctx.Fail(errors.ErrKindergartenTeacherNotFound, nil)
	}

	if !teacher.Auth(req.Password) { /* 密码错误 */
		return ctx.Fail(errors.ErrKindergartenTeacherPasswordIncorrect, nil)
	}

	expiresAt := time.Now().AddDate(0, 0, 1)
	token, err := kindergartendb.CreateKindergartenTeacherToken(teacher.ID, expiresAt, ctx.Request().UserAgent(), ctx.RealIP())
	if err != nil {
		return ctx.InternalServerError()
	}

	sess, _ := session.Get("teacher", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 3600 * 24,
	}
	sess.Values["token"] = token.ID
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.Success(teacher)
}

/* 注销登录 */
func Logout(c echo.Context) error {
	ctx := c.(*context.Context)

	sess, _ := session.Get("teacher", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 0,
	}
	sess.Values["token"] = ""
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.Success(nil)
}

func GetSelf(c echo.Context) error {
	ctx := c.(*context.Context)

	return ctx.Success(ctx.Teacher)
}

type UpdateSelfRequest struct {
	Name   string `json:"name" form:"name" query:"name" validate:"required"`
	Gender string `json:"gender" form:"gender" query:"gender" validate:"required"`
	Phone  string `json:"phone" form:"phone" query:"phone"`
}

func UpdateSelf(c echo.Context) error {
	ctx := c.(*context.Context)
// c.(*context.Context).Bind()
	req := UpdateSelfRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	teacher, err := kindergartendb.UpdateKindergartenTeacher(ctx.Teacher.ID, req.Name, req.Gender, req.Phone, ctx.Teacher.ClassID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(teacher)
}

type UpdateSelfPasswordRequest struct {
	Old string `json:"old" form:"old" query:"old" validate:"required"`
	New string `json:"new" form:"new" query:"new" validate:"gte=8"`
}

func UpdateSelfPassword(c echo.Context) error {
	ctx := c.(*context.Context)
	req := UpdateSelfPasswordRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	self := ctx.Teacher
	if !self.Auth(req.Old) {
		return ctx.Fail(errors.ErrKindergartenTeacherPasswordIncorrect, nil)
	}

	teacher, err := kindergartendb.UpdateKindergartenTeacherPassword(self.ID, req.New)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(teacher)
}

type FindKindergartenTeachersRequest struct {
	Query    string `json:"query" form:"query" query:"query"`
	Role     string `json:"role" form:"role" query:"role"`
	ClassID  int64  `json:"class_id" form:"class_id" query:"class_id"`
	Page     int    `json:"page" form:"page" query:"page"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindKindergartenTeachers(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenTeachersRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	teachers, total, err := kindergartendb.FindKindergartenTeachers(req.Query, ctx.Teacher.KindergartenID, req.ClassID, req.Role, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(teachers, req.Page, req.PageSize, total)
}

type CreateKindergartenTeacherRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"gte=6"`
	Password string `json:"password" form:"password" query:"password" validate:"gte=8"`
	Name     string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Phone    string `json:"phone" form:"phone" query:"phone"`
	Gender   string `json:"gender" form:"gender" query:"gender" validate:"required"`
	ClassID  int64  `json:"class_id" form:"class_id" query:"class_id"`
}

func CreateKindergartenTeacher(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenTeacherRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if !x.ValidateUsername(req.Username) {
		return ctx.BadRequest()
	}

	teacher, err := kindergartendb.GetKindergartenTeacherByUsername(req.Username)
	if err != nil {
		return ctx.InternalServerError()
	} else if teacher != nil {
		return ctx.Fail(errors.ErrKindergartenTeacherUsernameDuplicated, nil)
	}

	if req.ClassID > 0 {
		class, err := kindergartendb.GetKindergartenClass(req.ClassID)
		if err != nil {
			return ctx.InternalServerError()
		} else if class == nil {
			return ctx.BadRequest()
		}
		if class.KindergartenID != ctx.Teacher.KindergartenID {
			return ctx.BadRequest()
		}
	}

	/* 只有园长能创建老师 */
	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	teacher, err = kindergartendb.CreateKindergartenTeacher(req.Username, req.Password, req.Name, req.Gender, req.Phone, kindergartendb.KindergartenTeacherRoleTeacher, ctx.Teacher.KindergartenID, req.ClassID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(teacher)
}

type UpdateKindergartenTeacherRequest struct {
	Name    string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Phone   string `json:"phone" form:"phone" query:"phone"`
	Gender  string `json:"gender" form:"gender" query:"gender" validate:"required"`
	ClassID int64  `json:"class_id" form:"class_id" query:"class_id"`
}

func UpdateKindergartenTeacher(c echo.Context) error {
	ctx := c.(*context.Context)
	req := UpdateKindergartenTeacherRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	id := ctx.IntParam(`id`)

	teacher, err := kindergartendb.GetKindergartenTeacher(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if teacher == nil || teacher.KindergartenID != ctx.Teacher.KindergartenID {
		return ctx.NotFound()
	}

	if req.ClassID > 0 {
		class, err := kindergartendb.GetKindergartenClass(req.ClassID)
		if err != nil {
			return ctx.InternalServerError()
		} else if class == nil {
			return ctx.BadRequest()
		}
		if class.KindergartenID != teacher.KindergartenID {
			return ctx.BadRequest()
		}
	}

	/* 只有园长能编辑老师 */
	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	teacher, err = kindergartendb.UpdateKindergartenTeacher(id, req.Name, req.Gender, req.Phone, req.ClassID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(teacher)
}

/* 下载上传老师的模板 */
func DownloadKindergartenTeacherTemplate(c echo.Context) error {
	ctx := c.(*context.Context)

	headers := []string{"姓名", "电话", "性别", "班级", "用户名", "密码"}
	rows := make([][]interface{}, 0)

	return ctx.XLSX("老师模板", headers, rows)
}

type LoadKindergartenTeacherResult struct {
	Name      string `json:"name" xlsx:"姓名" validate:"required"`
	Phone     string `json:"phone" xlsx:"电话"`
	Gender    string `json:"gender" xlsx:"性别" validate:"required"`
	ClassName string `json:"class_name" xlsx:"班级"`
	ClassID   int64  `json:"class_id" validate:"gt=0"`

	Class    *kindergartendb.KindergartenClass `json:"class,omitempty"`
	Username string                            `json:"username" xlsx:"用户名" validate:"gte=6"`
	Password string                            `json:"password" xlsx:"密码" validate:"gte=8"`

	Status []string `json:"status"`
}

/* 解析XLSX文件，返回老师信息 */
func LoadKindergartenTeacher(c echo.Context) error {
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

	results := make([]*LoadKindergartenTeacherResult, 0)

	if err := x.ParseXLSXSheet(sheet, &results); err != nil {
		return ctx.BadRequest()
	}

	usernames := make(map[string]bool)
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
			}
		} else {
			r.Status = append(r.Status, "NoClass")
		}

		if r.Username == "" {
			r.Status = append(r.Status, "NoUsername")
		} else if !x.ValidateUsername(r.Username) {
			r.Status = append(r.Status, "InvalidUsername")
		} else if usernames[r.Username] {
			r.Status = append(r.Status, "DuplicatedUsername")
		} else {
			teacher, err := kindergartendb.GetKindergartenTeacherByUsername(r.Username)
			if err != nil {
				return ctx.InternalServerError()
			} else if teacher != nil {
				r.Status = append(r.Status, "DuplicatedUsername")
			}
		}

		if r.Password == "" {
			r.Status = append(r.Status, "NoPassword")
		} else if len(r.Password) < 8 {
			r.Status = append(r.Status, "ShortPassword")
		}

		if len(r.Status) == 0 {
			usernames[r.Username] = true
		}
	}
	return ctx.Success(results)
}

type CreateKindergartenTeacherLoadRequest struct {
	Teachers []*LoadKindergartenTeacherResult `json:"teachers" validate:"gt=0"`
}

func CreateKindergartenTeacherLoad(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenTeacherLoadRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	/* 只有园长能创建老师 */
	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	tx, err := db.Begin()
	if err != nil {
		return ctx.InternalServerError()
	}
	defer tx.Rollback()

	for _, t := range req.Teachers {
		if !x.ValidateUsername(t.Username) {
			return ctx.BadRequest()
		}

		teacher, err := kindergartendb.GetKindergartenTeacherByUsername(t.Username)
		if err != nil {
			return ctx.InternalServerError()
		} else if teacher != nil {
			return ctx.Fail(errors.ErrKindergartenTeacherUsernameDuplicated, nil)
		}

		if t.ClassID > 0 {
			class, err := kindergartendb.GetKindergartenClass(t.ClassID)
			if err != nil {
				return ctx.InternalServerError()
			} else if class == nil {
				return ctx.BadRequest()
			}
			if class.KindergartenID != ctx.Teacher.KindergartenID {
				return ctx.BadRequest()
			}
		}

		teacher, err = kindergartendb.CreateKindergartenTeacher(t.Username, t.Password, t.Name, t.Gender, t.Phone, kindergartendb.KindergartenTeacherRoleTeacher, ctx.Teacher.KindergartenID, t.ClassID)
		if err != nil {
			return ctx.InternalServerError()
		}
	}

	if err := tx.Commit(); err != nil {
		return ctx.InternalServerError()
	}

	return ctx.Success(nil)
}

func DeleteKindergartenTeacher(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	teacher, err := kindergartendb.GetKindergartenTeacher(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if teacher == nil {
		return ctx.NotFound()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager || ctx.Teacher.KindergartenID != teacher.KindergartenID {
		return ctx.Fail(errors.ErrKindergartenTeacherPermissionDenied, nil)
	}

	if teacher.Role == kindergartendb.KindergartenTeacherRoleManager {
		return ctx.Fail(errors.ErrKindergartenTeacherManager, nil)
	}

	teacher, err = kindergartendb.DeleteKindergartenTeacher(id)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(teacher)
}
