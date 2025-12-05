-- Exemple 1 : post avec un title, un text et une image
BEGIN;

WITH post AS (
INSERT
INTO posts (slug, title)
VALUES ('first-post', 'Premier article')
ON CONFLICT (slug) DO
UPDATE SET title = EXCLUDED.title
    RETURNING id
    ),
    b_title AS (
INSERT
INTO blocks (type)
VALUES ('title') RETURNING id
    ), b_text AS (
INSERT
INTO blocks (type)
VALUES ('text') RETURNING id
    ), b_image AS (
INSERT
INTO blocks (type)
VALUES ('image') RETURNING id
    ), ins_title AS (
INSERT
INTO title_blocks (block_id, level, content)
SELECT id, 1, 'Bienvenue dans le premier article'
FROM b_title
    RETURNING block_id
    ), ins_text AS (
INSERT
INTO text_blocks (block_id, content)
SELECT id, 'Ceci est le contenu principal du premier article.'
FROM b_text
    RETURNING block_id
    ), ins_image AS (
INSERT
INTO image_blocks (block_id, url, title)
SELECT id, 'https://example.com/images/illustration1.jpg', 'Illustration 1'
FROM b_image
    RETURNING block_id
    )
INSERT
INTO contents (parent_table, parent_id, ord, type, block_id)
SELECT 'posts', p.id, s.ord, s.type, s.block_id
FROM post p
         JOIN (SELECT 0 AS ord, 'title'::text AS type, block_id
               FROM ins_title
               UNION ALL
               SELECT 1, 'text', block_id
               FROM ins_text
               UNION ALL
               SELECT 2, 'image', block_id
               FROM ins_image) s ON true;

COMMIT;


-- Exemple 2 : post avec un title, un code et un custom (JSON)
BEGIN;

WITH post AS (
INSERT
INTO posts (slug, title)
VALUES ('second-post', 'Second article')
ON CONFLICT (slug) DO
UPDATE SET title = EXCLUDED.title
    RETURNING id
    ),
    b_title AS (
INSERT
INTO blocks (type)
VALUES ('title') RETURNING id
    ), b_code AS (
INSERT
INTO blocks (type)
VALUES ('code') RETURNING id
    ), b_custom AS (
INSERT
INTO blocks (type)
VALUES ('custom') RETURNING id
    ), ins_title AS (
INSERT
INTO title_blocks (block_id, level, content)
SELECT id, 2, 'Sous-titre du second article'
FROM b_title
    RETURNING block_id
    ), ins_code AS (
INSERT
INTO code_blocks (block_id, language, content)
SELECT id, 'go', 'fmt.Println(\"Hello, world\")'
FROM b_code
    RETURNING block_id
    ), ins_custom AS (
INSERT
INTO custom_blocks (block_id, payload)
SELECT id, '{"widget":"chart", "data":[1,2,3]}'::jsonb
FROM b_custom
    RETURNING block_id
    )
INSERT
INTO contents (parent_table, parent_id, ord, type, block_id)
SELECT 'posts', p.id, s.ord, s.type, s.block_id
FROM post p
         JOIN (SELECT 0 AS ord, 'title'::text AS type, block_id
               FROM ins_title
               UNION ALL
               SELECT 1, 'code', block_id
               FROM ins_code
               UNION ALL
               SELECT 2, 'custom', block_id
               FROM ins_custom) s ON true;

COMMIT;
