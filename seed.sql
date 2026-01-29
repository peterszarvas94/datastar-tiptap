CREATE TABLE IF NOT EXISTS content (
  id INTEGER PRIMARY KEY,
  html TEXT NOT NULL
);

DELETE FROM content;

INSERT INTO content (id, html)
VALUES (
  1,
  '<h2>Welcome to the editor</h2>
<p>This content is stored in SQLite.</p>
<blockquote><p>Edit this and click Save.</p></blockquote>'
);
