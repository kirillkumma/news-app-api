package adapter

const (
	queryGetUserByLogin = `
SELECT ID_user,
       Login,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE Login = $1
`

	queryGetUserByEmail = `
SELECT ID_user,
       Login,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE Email_user = $1
`

	queryCreateUser = `
INSERT INTO "user" (Login, Password, FIO_user, Email_user)
VALUES ($1, $2, $3, $4)
RETURNING ID_user, Login, Password, FIO_user, Email_user
`

	queryGetUserByID = `
SELECT ID_user,
       Login,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE ID_user = $1
`

	queryGetSubscriptionList = `
SELECT ID_editor,
       Num_reg_media_r,
       Corp_name,
       Email_red,
       Editor_surname,
       Editor_name,
       (SELECT COUNT(*) FROM subscription WHERE media_id = ID_editor)
FROM media
WHERE id_editor IN (SELECT media_id FROM subscription WHERE user_id = $1)
LIMIT $2 OFFSET $3
`

	queryCountSubscriptions = `
SELECT COUNT(*)
FROM media
WHERE id_editor IN (SELECT media_id FROM subscription WHERE user_id = $1)
`
)
