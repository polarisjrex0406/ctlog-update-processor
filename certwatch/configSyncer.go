package certwatch

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/xoduxcrt/ctlog-update-processor/config"
	"bitbucket.org/xoduxcrt/ctlog-update-processor/utils"
	"github.com/jackc/pgx/v5"
)

func InitConfig() error {
	id := 1
	for _, op := range utils.CTLogList.Operators {
		if op.Logs != nil {
			for _, l := range op.Logs {
				isActive := false
				if l.State != nil && l.State.Usable != nil {
					isActive = true
				}

				_, err := connLogConfigSyncer.Exec(context.Background(),
					`INSERT INTO ct_log (ID, OPERATOR, URL, NAME, PUBLIC_KEY, IS_ACTIVE, LATEST_UPDATE, MMD_IN_SECONDS, BATCH_SIZE, REQUESTS_CONCURRENT)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
					id, op.Name, l.URL, l.Description, l.Key, isActive, time.Now().UTC(), l.MMD, 256, 8)
				if err == nil {
					id++
				}
			}
		}

		if op.TiledLogs != nil {
			for _, tl := range op.TiledLogs {
				isActive := false
				if tl.State != nil && tl.State.Usable != nil {
					isActive = true
				}

				_, err := connLogConfigSyncer.Exec(context.Background(),
					`INSERT INTO ct_log (ID, OPERATOR, TYPE, URL, SUBMISSION_URL, NAME, PUBLIC_KEY, IS_ACTIVE, 
										LATEST_UPDATE, MMD_IN_SECONDS, BATCH_SIZE, REQUESTS_CONCURRENT)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
					id, op.Name, "static", tl.MonitoringURL, tl.SubmissionURL, tl.Description, tl.Key, isActive,
					time.Now().UTC(), tl.MMD, 256, 8)
				if err == nil {
					id++
				}
			}
		}
	}

	return nil
}

func GetConfig() (pgx.Rows, error) {
	return connLogConfigSyncer.Query(context.Background(), fmt.Sprintf(`
	SELECT ctl.ID, ctl.PUBLIC_KEY, ctl.URL, coalesce(ctl.SUBMISSION_URL, ctl.URL), ctl.TYPE, ctl.MMD_IN_SECONDS,
			CASE WHEN ctl.TYPE = 'rfc6962' THEN coalesce(ctl.BATCH_SIZE, %d) ELSE 256 END,
			ctl.REQUESTS_THROTTLE, coalesce(ctl.REQUESTS_CONCURRENT, 8), coalesce(latest.ENTRY_ID, -1)
		FROM ct_log ctl
				LEFT JOIN LATERAL (
					SELECT max(ctle.ENTRY_ID) ENTRY_ID
						FROM ct_log_entry ctle
						WHERE ctle.CT_LOG_ID = ctl.ID
				) latest ON TRUE
		WHERE ctl.IS_ACTIVE
	`, config.Config.CTLogs.GetEntriesDefaultBatchSize))
}

func BeginUpdateConfig() (pgx.Tx, error) {
	return connLogConfigSyncer.Begin(context.Background())
}
