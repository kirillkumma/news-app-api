package adapter

const (
	queryCreateNews = `
INSERT INTO news (Num_reg_media_news, title, text_content, release)
SELECT Num_reg_media_r, $2, $3, NOW() FROM media WHERE ID_editor = $1
RETURNING ID_news, Num_reg_media_news, title, text_content, EXTRACT(EPOCH FROM release)::INT
`

	queryAddNewsToFeed = `
INSERT INTO feed (ID_news, ID_user)
SELECT $1::INT, user_id FROM subscription
WHERE media_id = (
	SELECT ID_editor FROM media WHERE Num_reg_media_r = (SELECT Num_reg_media_news FROM news WHERE ID_news = $1::INT)
)
`

	queryGetNews = `
SELECT id_news,
       media.id_editor,
       media.num_reg_media_r,
       media.corp_name,
       media.email_red,
       media.editor_name,
       media.editor_surname,
       (SELECT COUNT(*) FROM subscription WHERE media_id = media.id_editor),
       title,
       text_content,
       EXTRACT(EPOCH FROM release)::INT
FROM news
INNER JOIN media ON
    news.num_reg_media_news = media.num_reg_media_r
WHERE id_news = $1
`

	queryGetFeedNewsList = `
SELECT news.id_news,
       media.id_editor,
       media.num_reg_media_r,
       media.corp_name,
       media.email_red,
       media.editor_name,
       media.editor_surname,
       (SELECT COUNT(*) FROM subscription WHERE media_id = media.id_editor),
       news.title,
       news.text_content,
       EXTRACT(EPOCH FROM news.release)::INT
FROM feed
INNER JOIN news ON
    feed.id_news = news.id_news
INNER JOIN media ON
    media.num_reg_media_r = news.num_reg_media_news
WHERE id_user = $1
LIMIT $2 OFFSET $3
`

	queryCountFeedNews = `
SELECT COUNT(*)
FROM feed
WHERE id_user = $1
`
)
