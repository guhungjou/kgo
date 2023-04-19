package admin

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
)

func GetAdminToken(id string) (*AdminToken, error) {
	pg := db.PG()

	token := AdminToken{ID: id}
	if err := pg.Model(&token).Relation(`User`).WherePK().Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &token, nil
}

func SetAdminTokenInvalid(userID int64) error {
	pg := db.PG()

	if _, err := pg.Model((*AdminToken)(nil)).Where(`"user_id"=?`, userID).Set(`"status"=?`, AdminTokenStatusInvalid).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return err
	}
	return nil
}

func CreateAdminToken(userID int64, expiresAt time.Time, device, ip string) (*AdminToken, error) {
	pg := db.PG()
	tx, err := pg.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.Model((*AdminToken)(nil)).Where(`"user_id"=?`, userID).Where(`"status"=?`, AdminTokenStatusOK).Set(`"status"=?`, AdminTokenStatusInvalid).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	token := AdminToken{
		UserID:    userID,
		Device:    device,
		IP:        ip,
		Status:    AdminTokenStatusOK,
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

func FindAdminTokens(userID int64, startTime, endTime time.Time, page, pageSize int) ([]*AdminToken, int, error) {
	pg := db.PG()

	tokens := make([]*AdminToken, 0)
	q := pg.Model(&tokens).Relation(`User`)
	if userID > 0 {
		q = q.Where(`"admin_token"."user_id"=?`, userID)
	}
	if !startTime.IsZero() {
		q = q.Where(`"admin_token"."created_at">=?`, startTime)
	}
	if !endTime.IsZero() {
		q = q.Where(`"admin_token"."created_at"<?`, endTime)
	}

	if total, err := q.Limit(pageSize).Offset((page - 1) * pageSize).Order(`admin_token.created_at DESC`).SelectAndCountEstimate(100000); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	} else {
		return tokens, total, nil
	}
}
