package adapter

const (
	queryCreateNews = `
INSERT INTO news (Num_reg_media_news, title, text_content, release)
SELECT Num_reg_media_r, $2, $3, NOW() FROM media WHERE ID_editor = $1
RETURNING ID_news, Num_reg_media_news, title, text_content, EXTRACT(EPOCH FROM release)::BIGINT
`

	queryAddNewsToFeed = `
INSERT INTO feed (ID_news, ID_user)
SELECT $1::BIGINT, user_id FROM subscription
WHERE media_id = (
	SELECT ID_editor FROM media WHERE Num_reg_media_r = (SELECT Num_reg_media_news FROM news WHERE ID_news = $1::BIGINT)
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
       EXTRACT(EPOCH FROM release)::BIGINT
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
       EXISTS(SELECT 1 FROM favorite WHERE user_id = $1 AND news_id = news.id_news),
       EXTRACT(EPOCH FROM news.release)::BIGINT
FROM feed
INNER JOIN news ON
    feed.id_news = news.id_news
INNER JOIN media ON
    media.num_reg_media_r = news.num_reg_media_news
WHERE id_user = $1
  AND ($2::BIGINT IS NULL OR EXTRACT(EPOCH FROM news.release)::BIGINT >= $2::BIGINT)
ORDER BY news.release
LIMIT $3 OFFSET $4
`

	queryCountFeedNews = `
SELECT COUNT(*)
FROM feed
INNER JOIN news ON
    news.id_news = feed.id_news
WHERE id_user = $1
  AND ($2::BIGINT IS NULL OR EXTRACT(EPOCH FROM news.release)::BIGINT >= $2::BIGINT)
`

	queryIsFavorite = `
SELECT EXISTS(SELECT 1 FROM favorite WHERE user_id = $1 AND news_id = $2)
`

	queryAddToFavorite = `
INSERT INTO favorite (user_id, news_id)
VALUES ($1, $2)
`

	queryRemoveFromFavorite = `
DELETE FROM favorite WHERE user_id = $1 AND news_id = $2
`

	queryGetFavoriteList = `
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
       EXISTS(SELECT 1 FROM favorite WHERE user_id = $1 AND news_id = news.id_news),
       EXTRACT(EPOCH FROM news.release)::BIGINT
FROM favorite
INNER JOIN news ON
    favorite.news_id = news.id_news
INNER JOIN media ON
    media.num_reg_media_r = news.num_reg_media_news
WHERE user_id = $1
LIMIT $2 OFFSET $3
`

	queryCountFavorites = `
SELECT COUNT(*)
FROM favorite
WHERE user_id = $1
`

	queryGetNewsList = `
SELECT news.id_news,
       media.id_editor,
       media.num_reg_media_r,
       media.corp_name,
       media.email_red,
       media.editor_name,
       media.editor_surname,
       (SELECT COUNT(*) FROM subscription WHERE media_id = $1),
       news.title,
       news.text_content,
       EXISTS(SELECT 1 FROM favorite WHERE user_id = $2 AND news_id = news.id_news),
       EXTRACT(EPOCH FROM news.release)::BIGINT
FROM news
INNER JOIN media ON
    media.num_reg_media_r = news.num_reg_media_news
WHERE media.id_editor = $1
ORDER BY news.release
LIMIT $3 OFFSET $4
`

	queryCountNews = `
SELECT COUNT(*)
FROM news
INNER JOIN media ON 
    news.num_reg_media_news = media.num_reg_media_r
WHERE media.id_editor = $1
`
)
