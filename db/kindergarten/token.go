package kindergarten

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
)

func GetKindergartenTeacherToken(id string) (*KindergartenTeacherToken, error) {
	pg := db.PG()

	token := KindergartenTeacherToken{ID: id}
	if err := pg.Model(&token).Relation(`Teacher`).Relation(`Teacher.Kindergarten`).Relation(`Teacher.Class`).WherePK().Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &token, nil
}

func SetKindergartenTeacherTokenInvalid(teacherID int64) error {
	pg := db.PG()

	if _, err := pg.Model((*KindergartenTeacherToken)(nil)).Where(`"teacher_id"=?`, teacherID).Set(`"status"=?`, KindergartenTeacherTokenStatusInvalid).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

func CreateKindergartenTeacherToken(teacherID int64, expiresAt time.Time, device, ip string) (*KindergartenTeacherToken, error) {
	pg := db.PG()
	tx, err := pg.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	// if _, err := tx.Model((*KindergartenTeacherToken)(nil)).Where(`"teacher_id"=?`, teacherID).Where(`"status"=?`, KindergartenTeacherTokenStatusOK).Set(`"status"=?`, KindergartenTeacherTokenStatusInvalid).Update(); err != nil {
	// 	log.Errorf("DB Error: %v", err)
	// 	return nil, err
	// }

	token := KindergartenTeacherToken{
		TeacherID: teacherID,
		Device:    device,
		IP:        ip,
		Status:    KindergartenTeacherTokenStatusOK,
		ExpiresAt: expiresAt,
	}
	if _, err := tx.Model(&token).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &token, nil
}

func FindKindergartenTeacherTokens(teacherID int64, startTime, endTime time.Time, page, pageSize int) ([]*KindergartenTeacherToken, int, error) {
	pg := db.PG()

	tokens := make([]*KindergartenTeacherToken, 0)
	q := pg.Model(&tokens).Relation(`User`)
	if teacherID > 0 {
		q = q.Where(`"kindergarten_teacher_token"."teacher_id"=?`, teacherID)
	}
	if !startTime.IsZero() {
		q = q.Where(`"kindergarten_teacher_token"."created_at">=?`, startTime)
	}
	if !endTime.IsZero() {
		q = q.Where(`"kindergarten_teacher_token"."created_at"<?`, endTime)
	}

	if total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).Order(`kindergarten_teacher_token.created_at DESC`).SelectAndCountEstimate(100000); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	} else {
		return tokens, total, nil
	}
}
