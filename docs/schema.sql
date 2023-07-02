-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-07-02T13:18:37.065Z

CREATE TABLE "users"
(
    "id"                  uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "username"            varchar          NOT NULL,
    "first_name"          varchar          NOT NULL,
    "last_name"           varchar          NOT NULL,
    "email"               varchar UNIQUE   NOT NULL,
    "is_email_verified"   boolean          NOT NULL DEFAULT false,
    "hashed_password"     varchar          NOT NULL,
    "password_changed_at" timestamptz      NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"          timestamptz      NOT NULL DEFAULT (now()),
    "updated_at"          timestamptz      NOT NULL DEFAULT (now())
);

CREATE TABLE "verify_emails"
(
    "id"          uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "user_id"     varchar          NOT NULL,
    "email"       varchar          NOT NULL,
    "secret_code" varchar          NOT NULL,
    "is_used"     boolean          NOT NULL DEFAULT false,
    "created_at"  timestamptz      NOT NULL DEFAULT (now()),
    "expired_at"  timestamptz      NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "teams"
(
    "id"         uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"       varchar          NOT NULL,
    "created_at" timestamptz      NOT NULL DEFAULT (now()),
    "updated_at" timestamptz      NOT NULL DEFAULT (now())
);

CREATE TABLE "user_teams"
(
    "id"               uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "number"           bigint           NOT NULL,
    "primary_position" varchar          NOT NULL,
    "user_id"          uuid             NOT NULL,
    "team_id"          uuid             NOT NULL,
    "created_at"       timestamptz      NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY,
    "user_id"       uuid        NOT NULL,
    "refresh_token" varchar     NOT NULL,
    "user_agent"    varchar     NOT NULL,
    "client_ip"     varchar     NOT NULL,
    "is_blocked"    boolean     NOT NULL DEFAULT false,
    "expires_at"    timestamptz,
    "created_at"    timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users" ("username");

CREATE UNIQUE INDEX ON "verify_emails" ("secret_code");

CREATE UNIQUE INDEX ON "teams" ("name");

ALTER TABLE "verify_emails"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_teams"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_teams"
    ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
