CREATE TABLE IF NOT EXISTS organizations (
		id TEXT PRIMARY KEY,
		name TEXT,
    slug TEXT
);

CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT,
    slug TEXT,
    status TEXT,

		organization_id TEXT UNSIGNED NOT NULL,


		CONSTRAINT fk__projects__organizations
			FOREIGN KEY (organization_id) REFERENCES organizations (id)
			ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS issues (
		id TEXT PRIMARY KEY,
		title TEXT,
    permalink TEXT,
    has_seen BOOL,
    first_seen TIMESTAMP,
    last_seen TIMESTAMP,
    user_count INT,
    level TEXT,
    status TEXT,
    type TEXT,

		project_id TEXT UNSIGNED NOT NULL,


		CONSTRAINT fk__issues__projects
			FOREIGN KEY (project_id) REFERENCES projects (id)
			ON DELETE CASCADE
);
