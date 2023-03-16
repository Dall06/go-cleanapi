package repository

// const of the stored procedure names from the db
const (
	spCreate = "CALL `go_cleanapi`.`sp_create_user`(?, ?, ?, ?);";
	spRead = "CALL `go_cleanapi`.`sp_read_user`(?);";
	spReadAll = "CALL `go_cleanapi`.`sp_read_users`();";
	spUpdate = "CALL `go_cleanapi`.`sp_update_user`(?, ?, ?, ?);";
	spDelete = "CALL `go_cleanapi`.`sp_delete_user`(?, ?);";
)

