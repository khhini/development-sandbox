-- Add up migration script here
-- sqlfluff:dialect:postgres
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    search_vector TSVECTOR GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(content, '')), 'B')
    ) STORED
);

INSERT INTO articles (title, content) VALUES
('The Great Gatsby', 'A novel by F. Scott Fitzgerald'),
('To Kill a Mockingbird', 'A novel by Harper Lee'),
('1984', 'A dystopian novel by George Orwell'),
('The Catcher in the Rye', 'A novel by J. D. Salinger'),
('Pride and Prejudice', 'A novel by Jane Austen'),
('The Hobbit', 'A fantasy novel by J. R. R. Tolkien'),
('Brave New World', 'A dystopian novel by Aldous Huxley'),
('The Lord of the Rings', 'A fantasy novel by J. R. R. Tolkien'),
('Animal Farm', 'An allegorical novella by George Orwell'),
('Fahrenheit 451', 'A dystopian novel by Ray Bradbury'),
('Jane Eyre', 'A novel by Charlotte Brontë'),
('Wuthering Heights', 'A novel by Emily Brontë'),
('Moby-Dick', 'A novel by Herman Melville');

