package adapter

const (
	queryGetUserByLogin = `
SELECT ID_user,
       LoginUser,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE LoginUser = $1
`

	queryGetUserByEmail = `
SELECT ID_user,
       LoginUser,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE Email_user = $1
`

	queryCreateUser = `
INSERT INTO "user" (LoginUser, Password, FIO_user, Email_user)
VALUES ($1, $2, $3, $4)
RETURNING ID_user, LoginUser, Password, FIO_user, Email_user
`

	queryGetUserByID = `
SELECT ID_user,
       LoginUser,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE ID_user = $1
`
)
