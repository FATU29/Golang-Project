package cron

import (
	"Authentication_Service/internal/repository/interface"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/robfig/cron/v3"
)

// StartRefreshTokenCleanup runs a cron job that deletes expired refresh tokens.
// schedule: cron expression, e.g. "0 * * * *" (every hour), "0 0 * * *" (daily at midnight).
// Schedule: 5-field cron (min hour day month dow), e.g. "0 * * * *" = hourly, "0 0 * * *" = daily midnight.
func StartRefreshTokenCleanup(repo _interface.IRefreshTokenRepository, schedule string) (*cron.Cron, error) {
	c := cron.New()

	_, err := c.AddFunc(schedule, func() {
		affected, err := repo.DeleteExpired()
		if err != nil {
			logger.Errorf("cron refresh_token cleanup: %v", err)
			return
		}
		if affected > 0 {
			logger.Infof("cron refresh_token cleanup: deleted %d expired token(s) at %s", affected, time.Now().Format(time.RFC3339))
		}
	})

	if err != nil {
		return nil, err
	}

	c.Start()
	logger.Infof("cron refresh_token cleanup started with schedule: %s", schedule)
	return c, nil
}
