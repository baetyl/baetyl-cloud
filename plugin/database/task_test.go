package database

import (
	"fmt"
)

var (
	taskTables = []string{
		`
CREATE TABLE baetyl_task
(
    trace_id          varchar(36)  NOT NULL DEFAULT '' PRIMARY KEY,
    namespace         varchar(64)  NOT NULL DEFAULT '',
    node              varchar(128)   NOT NULL DEFAULT '',
    type              varchar(32)  NOT NULL DEFAULT '',
    state             varchar(16)       NOT NULL DEFAULT '0',
    step              text   NOT NULL,
    old_version       varchar(36)   NOT NULL DEFAULT '',
    new_version       varchar(36)     NOT NULL DEFAULT '',
    create_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *DB) MockCreateTaskTable() {
	for _, sql := range taskTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}