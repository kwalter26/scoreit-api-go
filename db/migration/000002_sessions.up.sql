CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY,
    "user_id"    uuid        NOT NULL,
    "refresh_token" varchar(700) NOT NULL,
    "user_agent"    varchar(100) NOT NULL,
    "client_ip"     varchar(100) NOT NULL,
    "is_blocked" boolean     NOT NULL DEFAULT false,
    "expires_at"    timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");