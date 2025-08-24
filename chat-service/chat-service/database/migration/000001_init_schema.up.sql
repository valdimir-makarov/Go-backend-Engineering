-- 20250822045501_init_schema.up.sql

-- Create tables
CREATE TABLE conversation (
    id UUID PRIMARY KEY,
    is_group BOOLEAN DEFAULT false,
    group_name TEXT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    last_message_id UUID NULL
);

CREATE TABLE messages (
    id UUID PRIMARY KEY,
    conversation_id UUID,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT now(),
    delivered BOOLEAN DEFAULT false,
    CONSTRAINT fk_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversation(id) ON DELETE CASCADE
);

CREATE TABLE conversation_participants (
    conversation_id UUID,
    user_id INT,
    joined_at TIMESTAMP DEFAULT now(),
    role TEXT CHECK (role IN ('member', 'admin', 'owner')),
    PRIMARY KEY (conversation_id, user_id),
    CONSTRAINT fk_conversation_participants FOREIGN KEY (conversation_id) REFERENCES conversation(id) ON DELETE CASCADE
);

-- Insert sample conversation
INSERT INTO conversation (id, is_group, group_name, created_at, updated_at, last_message_id)
VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    false,
    NULL,
    now(),
    now(),
    '162a7258-9a55-47dc-bd17-088fe8204d0f'::UUID
);

-- Insert conversation participants
INSERT INTO conversation_participants (conversation_id, user_id, joined_at, role)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 123, now(), 'member'),
    ('550e8400-e29b-41d4-a716-446655440000', 456, now(), 'member');

-- Insert a sample message
INSERT INTO messages (id, conversation_id, sender_id, receiver_id, content, created_at, delivered)
VALUES (
    '162a7258-9a55-47dc-bd17-088fe8204d0f'::UUID,
    '550e8400-e29b-41d4-a716-446655440000'::UUID,
    123,
    456,
    'Hello, this is a test message!',
    now(),
    false
);
