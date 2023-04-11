DROP DATABASE IF EXISTS db_go_cleanapi;
CREATE DATABASE db_go_cleanapi;

USE db_go_cleanapi;

CREATE TABLE users (
	id_user VARCHAR(64) NOT NULL UNIQUE,
    user_email VARCHAR(128) NOT NULL,
    user_phone VARCHAR(16),
    user_password VARCHAR(64) NOT NULL
);

DELIMITER $$
CREATE DEFINER=`root`@`localhost` FUNCTION `fn_validate_user`(
	in_u_id VARCHAR(36),
    in_u_pass VARCHAR(128)
) RETURNS BOOL
    READS SQL DATA
    DETERMINISTIC
BEGIN
	RETURN IF( EXISTS(
		SELECT * FROM `users` WHERE `user_id` = in_u_id AND `user_pass` = SHA2(in_u_pass, 512)), 1, 0);
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_read_user`(
	p_id_user VARCHAR(64)
)
BEGIN
	SELECT  (`id_user`,
	`user_email`,
	`user_phone`) FROM users WHERE id_user = p_id_user;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_read_users`()
BEGIN
	SELECT 
    (`id_user`,
	`user_email`,
	`user_phone`)
    FROM users;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_create_user`(
	p_id_user VARCHAR(64),
    p_user_email VARCHAR(128),
    p_user_phone VARCHAR(16),
    p_user_password VARCHAR(64)
)
BEGIN
	INSERT INTO `db_go_cleanapi`.`users`
	(`id_user`,
	`user_email`,
	`user_phone`,
	`user_password`)
	VALUES
	(p_id_user,
	p_user_email,
	p_user_phone,
	SHA2(p_user_password, 512));
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_update_user`(
	p_id_user VARCHAR(64),
    p_user_email VARCHAR(128),
    p_user_phone VARCHAR(16),
    p_user_password VARCHAR(64)
)
BEGIN
    UPDATE `db_go_cleanapi`.`users`
	SET
		`user_email` = p_user_email,
		`user_phone` = p_user_phone,
		`user_password` = p_user_password
	WHERE `id_user` = p_id_user;
END$$
DELIMITER ;

DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_delete_user`(
	p_id_user VARCHAR(64),
    p_user_password VARCHAR(64)
)
BEGIN
	DECLARE is_auth TINYINT;
    
	SELECT `db_go_cleanapi`.`fn_validate_user`(p_id_user, p_user_password) INTO is_auth;
    IF is_auth = FALSE THEN 
		SIGNAL SQLSTATE '40400' SET MESSAGE_TEXT = 'not authorized (worng credentials)';
    END IF;
    
    DELETE FROM `db_go_cleanapi`.users WHERE id_user = p_id_user AND user_password = SHA2(p_user_password, 512);
END$$
DELIMITER ;