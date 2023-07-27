-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-07-27T04:17:04.228Z

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

CREATE TABLE "user_roles"
(
    "id"         uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"       varchar          NOT NULL,
    "user_id"    uuid             NOT NULL,
    "created_at" timestamptz      NOT NULL DEFAULT (now()),
    "updated_at" timestamptz      NOT NULL DEFAULT (now())
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

CREATE TABLE "team_members"
(
    "id"               uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "number"           bigint           NOT NULL,
    "primary_position" varchar          NOT NULL,
    "user_id"          uuid             NOT NULL,
    "team_id"          uuid             NOT NULL,
    "created_at"       timestamptz      NOT NULL DEFAULT (now()),
    "updated_at"       timestamptz      NOT NULL DEFAULT (now())
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

CREATE TABLE "game"
(
    "id"           uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "home_team_id" uuid             NOT NULL,
    "away_team_id" uuid             NOT NULL,
    "home_score"   bigint           NOT NULL,
    "away_score"   bigint           NOT NULL,
    "created_at"   timestamptz      NOT NULL DEFAULT (now()),
    "updated_at"   timestamptz      NOT NULL DEFAULT (now())
);

CREATE TABLE "inning"
(
    "id"            uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "game_id"       uuid,
    "number"        bigint           NOT NULL,
    "home_runs"     bigint           NOT NULL,
    "home_hits"     bigint           NOT NULL,
    "home_errors"   bigint           NOT NULL,
    "home_last_bat" uuid             NOT NULL,
    "away_runs"     bigint           NOT NULL,
    "away_hits"     bigint           NOT NULL,
    "away_errors"   bigint           NOT NULL,
    "away_last_bat" uuid             NOT NULL
);

CREATE TABLE "atbat"
(
    "id"          uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "inning_id"   uuid,
    "batter_id"   uuid,
    "pitcher_id"  uuid,
    "balls"       bigint           NOT NULL,
    "strikes"     bigint           NOT NULL,
    "init_bases"  bigint           NOT NULL,
    "total_bases" bigint           NOT NULL,
    "out"         boolean          NOT NULL
);

CREATE TABLE "game_participant"
(
    "id"           uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "game_id"      uuid,
    "player_id"    uuid,
    "home_team"    boolean          NOT NULL,
    "bat_position" bigint           NOT NULL
);

CREATE TABLE "game_stat"
(
    "id"       uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "atbat_id" uuid,
    "type"     varchar
);

CREATE UNIQUE INDEX ON "users" ("username");

CREATE UNIQUE INDEX ON "user_roles" ("name", "user_id");

CREATE UNIQUE INDEX ON "verify_emails" ("secret_code");

CREATE UNIQUE INDEX ON "teams" ("name");

ALTER TABLE "user_roles"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "verify_emails"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "team_members"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "team_members"
    ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "game"
    ADD FOREIGN KEY ("home_team_id") REFERENCES "teams" ("id");

ALTER TABLE "game"
    ADD FOREIGN KEY ("away_team_id") REFERENCES "teams" ("id");

ALTER TABLE "inning"
    ADD FOREIGN KEY ("game_id") REFERENCES "game" ("id");

ALTER TABLE "inning"
    ADD FOREIGN KEY ("home_last_bat") REFERENCES "game_participant" ("id");

ALTER TABLE "inning"
    ADD FOREIGN KEY ("away_last_bat") REFERENCES "game_participant" ("id");

ALTER TABLE "atbat"
    ADD FOREIGN KEY ("inning_id") REFERENCES "inning" ("id");

ALTER TABLE "atbat"
    ADD FOREIGN KEY ("batter_id") REFERENCES "game_participant" ("id");

ALTER TABLE "atbat"
    ADD FOREIGN KEY ("pitcher_id") REFERENCES "game_participant" ("id");

ALTER TABLE "game_participant"
    ADD FOREIGN KEY ("game_id") REFERENCES "game" ("id");

ALTER TABLE "game_participant"
    ADD FOREIGN KEY ("player_id") REFERENCES "users" ("id");

ALTER TABLE "game_stat"
    ADD FOREIGN KEY ("atbat_id") REFERENCES "atbat" ("id");
