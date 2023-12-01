DELETE
FROM team_members
WHERE ctid NOT IN (SELECT MAX(ctid)
                   FROM team_members
                   GROUP BY user_id, team_id);

CREATE UNIQUE INDEX ON "team_members" ("user_id", "team_id");
