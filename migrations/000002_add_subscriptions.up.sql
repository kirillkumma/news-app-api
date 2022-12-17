CREATE TABLE subscription (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user" (ID_user),
    media_id BIGINT NOT NULL REFERENCES media (ID_editor),
    UNIQUE (user_id, media_id)
);