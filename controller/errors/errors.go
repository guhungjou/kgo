package errors

const (
	ErrAdminUserNotFound          = 10001
	ErrAdminUserPasswordIncorrect = 10002
	ErrAdminUserStatusInvalid     = 10003

	ErrKindergartenTeacherNotFound           = 20001
	ErrKindergartenTeacherPasswordIncorrect  = 20002
	ErrKindergartenTeacherPermissionDenied   = 20003
	ErrKindergartenTeacherUsernameDuplicated = 20004
	ErrKindergartenTeacherManager            = 20005

	ErrKindergartenStudentDeviceDuplicated = 21001
	ErrKindergartenStudentDeviceInvalid    = 21002
	ErrKindergartenStudentNODuplicated     = 21003
	ErrKindergartenStudentNONotFound       = 21004
	ErrKindergartenStudentNameNotFound     = 21005
	ErrKindergartenStudentNameDuplicated   = 21006

	ErrKindergartenClassNameDuplicated = 22001
)
