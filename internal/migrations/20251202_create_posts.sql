-- sql
CREATE TABLE IF NOT EXISTS posts (
slug text PRIMARY KEY,
title text NOT NULL,
content jsonb NOT NULL,
created_at timestamptz NOT NULL DEFAULT now(),
updated_at timestamptz NOT NULL DEFAULT now()
);

-- unique index on slug implicit by PK; si vous voulez un index séparé :
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_slug ON posts (slug);

INSERT INTO posts (slug, title, content)
VALUES
('example-post', 'Exemple de post avec images et paragraphes',
'[
{"type":"img","content":{"url":"https://cdn.example.com/img1.jpg","title":"Première image"}},
{"type":"p","content":"Ceci est le premier paragraphe du post. Il introduit le sujet."},
{"type":"h2","content":"Titre secondaire du post"},
{"type":"img","content":{"url":"https://cdn.example.com/img2.jpg","title":"Deuxième image"}},
{"type":"p","content":"Ceci est le deuxième paragraphe, qui conclut ou développe davantage."}
]'::jsonb
)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO posts (slug, title, content)
VALUES
('second-post', 'Second Post', '[{"type":"p","content":"This is the content of the second post."}]'::jsonb)
ON CONFLICT (slug) DO NOTHING;
