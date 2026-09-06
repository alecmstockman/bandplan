-- +goose Up

ALTER TABLE message_reactions
DROP CONSTRAINT message_reactions_message_id_user_id_reaction_key;

ALTER TABLE message_reactions
ADD CONSTRAINT message_reactions_message_id_user_id_key
UNIQUE (message_id, user_id);


-- +goose Down

ALTER TABLE message_reactions
DROP CONSTRAINT message_reactions_message_id_user_id_key;

ALTER TABLE message_reactions
ADD CONSTRAINT message_reactions_message_id_user_id_reaction_key
UNIQUE (message_id, user_id, reaction);
