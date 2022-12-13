package adapter

const (
	queryGetByLogin = `
SELECT ID_user,
       Login,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE Login = $1
`

	queryGetByEmail = `
SELECT ID_user,
       Login,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE Email_user = $1
`

	queryCreate = `
INSERT INTO "user" (Login, Password, FIO_user, Email_user)
VALUES ($1, $2, $3, $4)
RETURNING ID_user, Login, Password, FIO_user, Email_user
`

	queryGetByID = `
SELECT ID_user,
       Login,
       Password,
       FIO_user,
       Email_user
FROM "user"
WHERE ID_user = $1
`
)
