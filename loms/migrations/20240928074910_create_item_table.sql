-- +goose Up
-- +goose StatementBegin
CREATE TABLE item
(
  id      serial PRIMARY KEY,
  sku     integer NOT NULL,
  count   integer NOT NULL,
  order_id INT NOT NULL REFERENCES "order" (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE item;
-- +goose StatementEnd
