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
    "game_id"   uuid NOT NULL,
    "player_id" uuid NOT NULL,
    "team_id"   uuid NOT NULL,
    "bat_position" bigint           NOT NULL
);

CREATE TABLE "game_stat"
(
    "id"       uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "atbat_id" uuid,
    "type"     varchar
);

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

ALTER TABLE "game_participant"
    ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "game_stat"
    ADD FOREIGN KEY ("atbat_id") REFERENCES "atbat" ("id");
