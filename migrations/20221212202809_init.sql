-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    ID_user    BIGSERIAL PRIMARY KEY,
    Login      VARCHAR(32) NOT NULL UNIQUE,
    Password   VARCHAR(32),
    FIO_user   VARCHAR(255),
    Email_user VARCHAR(16)
);

CREATE TABLE media
(
    ID_editor       BIGSERIAL PRIMARY KEY,
    Num_reg_media_r INT UNIQUE NOT NULL,
    Corp_name       VARCHAR(32),
    Email_red       VARCHAR(16),
    Editor_surname  VARCHAR(16),
    Editor_name     VARCHAR(16),
    Password        VARCHAR(32)
);

CREATE TABLE news
(
    ID_news            BIGSERIAL PRIMARY KEY,
    Num_reg_media_news INT REFERENCES media (Num_reg_media_r),
    Title              VARCHAR(32),
    Text_content       VARCHAR(8000),
    Video_content      VARCHAR(512),
    Audio_content      VARCHAR(512),
    Release            TIMESTAMPTZ
);

CREATE TABLE feed
(
    ID_news INT REFERENCES news (ID_news),
    ID_user INT REFERENCES "user" (ID_user)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed;

DROP TABLE news;

DROP TABLE media;

DROP TABLE "user";
-- +goose StatementEnd
