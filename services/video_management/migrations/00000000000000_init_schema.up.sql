-- videos: main metadata owned by Video Service
CREATE TABLE videos (
    id                  UUID            NOT NULL,
    uploader_id         UUID            NOT NULL,               -- references user ID from User Service
    title               VARCHAR(255)    NOT NULL,
    description         TEXT,
    status              VARCHAR(50)     DEFAULT 'processing',   -- processing, ready, failed
    original_file_url   TEXT,                                   -- s3://bucket/raw/video.mp4
    processed_at        TIMESTAMP,
    created_at          TIMESTAMP       DEFAULT NOW(),
    updated_at          TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id)
);

-- video_manifests
CREATE TABLE video_manifests (
    id                  UUID            NOT NULL,
    video_id            UUID            NOT NULL,
    manifest_url        TEXT            NOT NULL,            -- e.g., s3://bucket/processed/{video_id}/hls/master.m3u8
    size_bytes          BIGINT,
    created_at          TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_video_manifest FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
);

-- video_variants
CREATE TABLE video_variants (
    id                  UUID            NOT NULL,
    video_id            UUID            NOT NULL,
    quality             VARCHAR(50)     NOT NULL,           -- 480p, 720p, 1080p
    playlist_url        TEXT            NOT NULL,           -- e.g., s3://bucket/processed/{video_id}/hls/quality_0/index.m3u8
    total_segments      INT,
    total_duration      INT,                                -- total video duration in seconds
    created_at          TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_video_variant FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
);

-- video_thumbnails
CREATE TABLE video_thumbnails (
    id                  UUID            NOT NULL,
    video_id            UUID            NOT NULL,
    file_url            TEXT            NOT NULL,           -- e.g., s3://bucket/processed/{video_id}/thumbnail.jpg
    width               INT,
    height              INT,
    created_at          TIMESTAMP       DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_video_thumbnail FOREIGN KEY (video_id) REFERENCES videos (id) ON DELETE CASCADE
);

