CREATE TABLE favorite (
    user_id BIGINT NOT NULL REFERENCES "user" (id_user),
    news_id BIGINT NOT NULL REFERENCES news (id_news)
);