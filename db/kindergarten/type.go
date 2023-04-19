package kindergarten

import (
	"time"

	systemModel "gitlab.com/ykgk/kgo/db/system"
)

type Kindergarten struct {
	tableName struct{} `pg:"kindergarten" json:"-"`
	ID        int64    `pg:"id,notnull,pk" json:"id"`

	Name   string `pg:"name,notnull" json:"name"`
	Remark string `pg:"remark,notnull" json:"remark"`

	NumberOfStudent int64 `pg:"number_of_student,notnull,use_zero" json:"number_of_student"`
	NumberOfTeacher int64 `pg:"number_of_teacher,notnull,use_zero" json:"number_of_teacher"`

	DistrictID string `pg:"district_id,notnull,use_zero" json:"district_id"`

	CreatedBy int64     `pg:"created_by,notnull,use_zero" json:"created_by"`
	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Deleted bool `pg:"deleted,notnull,use_zero" json:"-"`

	Manager  *KindergartenTeacher  `pg:"-" json:"manager,omitempty"`
	District *systemModel.District `pg:"rel:has-one" json:"district,omitempty"`
}

type KindergartenClass struct {
	tableName struct{} `pg:"kindergarten_class" json:"-"`
	ID        int64    `pg:"id,notnull,pk" json:"id"`

	Name           string `pg:"name,notnull" json:"name"`
	Remark         string `pg:"remark,notnull" json:"remark"`
	KindergartenID int64  `pg:"kindergarten_id,notnull" json:"kindergarten_id"`

	NumberOfStudent int64 `pg:"number_of_student,notnull,use_zero" json:"number_of_student"`
	NumberOfTeacher int64 `pg:"number_of_teacher,notnull,use_zero" json:"number_of_teacher"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Deleted bool `pg:"deleted,notnull,use_zero" json:"-"`

	Kindergarten *Kindergarten `pg:"rel:has-one" json:"kindergarten,omitempty"`
}

const (
	KindergartenTeacherRoleManager = "manager"
	KindergartenTeacherRoleTeacher = "teacher"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"
)

type KindergartenTeacher struct {
	tableName struct{} `pg:"kindergarten_teacher" json:"-"`
	ID        int64    `pg:"id,notnull,pk" json:"id"`

	Username string `pg:"username,notnull" json:"username"`
	Name     string `pg:"name,notnull" json:"name"`
	Role     string `pg:"role,notnull" json:"role"`
	Gender   string `pg:"gender,notnull" json:"gender"`
	Phone    string `pg:"phone,notnull,use_zero" json:"phone"`

	KindergartenID int64 `pg:"kindergarten_id,notnull" json:"kindergarten_id"`
	ClassID        int64 `pg:"class_id,notnull" json:"class_id"`

	Salt     string `pg:"salt,notnull,use_zero" json:"-"`
	PType    string `pg:"ptype,notnull,use_zero" json:"-"`
	Password string `pg:"password,notnull,use_zero" json:"-"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Deleted bool `pg:"deleted,notnull,use_zero" json:"-"`

	Kindergarten *Kindergarten      `pg:"rel:has-one" json:"kindergarten,omitempty"`
	Class        *KindergartenClass `pg:"rel:has-one" json:"class,omitempty"`
}

func (t *KindergartenTeacher) TeacherClassID(classID int64) int64 {
	if t.Role == KindergartenTeacherRoleManager {
		return classID
	} else if t.ClassID > 0 {
		return t.ClassID
	} else {
		return -1
	}
}

const (
	KindergartenTeacherTokenStatusOK      = "OK"
	KindergartenTeacherTokenStatusInvalid = "Invalid"
)

type KindergartenTeacherToken struct {
	tableName struct{} `pg:"kindergarten_teacher_token" json:"-"`

	ID        string `pg:"id,pk" json:"id"`
	TeacherID int64  `pg:"teacher_id,notnull" json:"teacher_id"`
	Device    string `pg:"device,notnull,use_zero" json:"device"`
	IP        string `pg:"ip,notnull,use_zero" json:"ip"`
	Status    string `pg:"status,notnull,use_zero" json:"status"`

	ExpiresAt time.Time `pg:"expires_at" json:"expires_at"`
	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Teacher *KindergartenTeacher `pg:"rel:has-one" json:"teacher,omitempty"`
}

type KindergartenStudent struct {
	tableName struct{} `pg:"kindergarten_student" json:"-"`
	ID        int64    `pg:"id,notnull,pk" json:"id"`

	NO string `pg:"no,notnull,use_zero" json:"no"`

	Name     string    `pg:"name,notnull" json:"name"`
	Gender   string    `pg:"gender,notnull" json:"gender"`
	Birthday time.Time `pg:"birthday,notnull" json:"birthday"`
	Device   string    `pg:"device,notnull" json:"device"`
	Remark   string    `pg:"remark,notnull,use_zero" json:"remark"`

	KindergartenID int64 `pg:"kindergarten_id,notnull" json:"kindergarten_id"`
	ClassID        int64 `pg:"class_id,notnull,use_zero" json:"class_id"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Deleted bool `pg:"deleted,notnull,use_zero" json:"-"`

	Kindergarten *Kindergarten      `pg:"rel:has-one" json:"kindergarten,omitempty"`
	Class        *KindergartenClass `pg:"rel:has-one" json:"class,omitempty"`
}

func (s *KindergartenStudent) GenderName() string {
	switch s.Gender {
	case "female":
		return "女"
	case "male":
		return "男"
	}
	return ""
}

func (s *KindergartenStudent) Age() float64 {
	now := time.Now()

	months := 12*(now.Year()-s.Birthday.Year()) + (int(now.Month()) - int(s.Birthday.Month()))

	year := months / 12
	month := months % 12
	var age float64 = float64(year)
	if month >= 4 && month <= 9 {
		age += 0.5
	} else if month > 9 {
		age += 1
	}
	return age
}
