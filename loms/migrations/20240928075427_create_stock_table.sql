-- +goose Up
-- +goose StatementBegin
CREATE TABLE stock
(
  stock_id      serial PRIMARY KEY,
  sku   integer NOT NULL ,
  total_count integer NOT NULL ,
  reserved integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
-- +goose StatementEnd
