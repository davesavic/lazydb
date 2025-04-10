INSERT INTO users (name, email) VALUES
    ('Test User', 'test@example.com'),
    ('Another User', 'another@example.com')
ON CONFLICT (email) DO NOTHING;
