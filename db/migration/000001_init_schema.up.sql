CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users"
(
    "id"                  uuid PRIMARY KEY    NOT NULL DEFAULT (uuid_generate_v4()),
    "username"            varchar(100)        NOT NULL,
    "first_name"          varchar(100)        NOT NULL,
    "last_name"           varchar(100)        NOT NULL,
    "email"               varchar(100) UNIQUE NOT NULL,
    "is_email_verified"   boolean             NOT NULL DEFAULT false,
    "hashed_password"     varchar             NOT NULL,
    "password_changed_at" timestamptz         NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"          timestamptz         NOT NULL DEFAULT (now()),
    "updated_at"          timestamptz         NOT NULL DEFAULT (now())
);

CREATE TABLE "verify_emails"
(
    "id"          uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "user_id" varchar(100) NOT NULL,
    "email"   varchar(100) NOT NULL,
    "secret_code" varchar          NOT NULL,
    "is_used"     boolean          NOT NULL DEFAULT false,
    "created_at"  timestamptz      NOT NULL DEFAULT (now()),
    "expired_at"  timestamptz      NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "teams"
(
    "id"         uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name" varchar(100) NOT NULL,
    "created_at" timestamptz      NOT NULL DEFAULT (now()),
    "updated_at" timestamptz      NOT NULL DEFAULT (now())
);

CREATE TABLE "team_members"
(
    "id"               uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "number"           bigint           NOT NULL,
    "primary_position" varchar(100) NOT NULL,
    "user_id"          uuid             NOT NULL,
    "team_id"          uuid             NOT NULL,
    "created_at"       timestamptz      NOT NULL DEFAULT (now()),
    "updated_at"       timestamptz      NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users" ("username");

CREATE UNIQUE INDEX ON "teams" ("name");

ALTER TABLE "team_members"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "team_members"
    ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");
