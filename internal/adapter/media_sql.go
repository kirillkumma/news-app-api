package adapter

const (
	queryGetMediaByRegistrationNumber = `
SELECT ID_editor,
       Num_reg_media_r,
       Corp_name,
       Email_red,
       Editor_surname,
       Editor_name,
       Password
FROM media
WHERE Num_reg_media_r = $1
`

	queryGetMediaByName = `
SELECT ID_editor,
       Num_reg_media_r,
       Corp_name,
       Email_red,
       Editor_surname,
       Editor_name,
       Password
FROM media
WHERE Corp_name = $1
`

	queryGetMediaByEmail = `
SELECT ID_editor,
       Num_reg_media_r,
       Corp_name,
       Email_red,
       Editor_surname,
       Editor_name,
       Password
FROM media
WHERE Email_red = $1
`

	queryCreateMedia = `
INSERT INTO media (Num_reg_media_r, Corp_name, Email_red, Editor_surname, Editor_name, Password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING ID_editor, Num_reg_media_r, Corp_name, Email_red, Editor_surname, Editor_name, Password
`

	queryGetMediaByID = `
SELECT ID_editor,
       Num_reg_media_r,
       Corp_name,
       Email_red,
       Editor_surname,
       Editor_name,
       Password
FROM media
WHERE ID_editor = $1
`
)
