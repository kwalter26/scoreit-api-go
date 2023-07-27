CREATE TABLE "user_roles"
(
    "id"         uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "name"       varchar          NOT NULL,
    "user_id"    uuid             NOT NULL,
    "created_at" timestamptz      NOT NULL DEFAULT (now()),
    "updated_at" timestamptz      NOT NULL DEFAULT (now())
);

ALTER TABLE "user_roles"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

