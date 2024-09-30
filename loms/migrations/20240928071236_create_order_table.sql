-- +goose Up
-- +goose StatementBegin
CREATE TABLE "order"
(
  id      serial PRIMARY KEY,
  status  text NOT NULL,
  "user"   int NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "order";
-- +goose StatementEnd
