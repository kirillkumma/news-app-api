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

	queryGetMediaList = `
SELECT ID_editor,
       Num_reg_media_r,
       Corp_name,
       Email_red,
       Editor_surname,
       Editor_name,
       (SELECT COUNT(*) FROM subscription WHERE media_id = ID_editor)
FROM media
LIMIT $1 OFFSET $2
`

	queryCountMedia = `
SELECT COUNT(*)
FROM media
`

	queryIsSubscriptionExists = `
SELECT EXISTS(SELECT 1 FROM subscription WHERE media_id = $1 AND user_id = $2)
`

	queryCreateSubscription = `
INSERT INTO subscription (media_id, user_id)
VALUES ($1, $2)
`

	queryDeleteSubscription = `
DELETE FROM subscription WHERE media_id = $1 AND user_id = $2
`
)
