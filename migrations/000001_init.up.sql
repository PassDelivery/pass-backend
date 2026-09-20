CREATE TABLE IF NOT EXISTS users(
  id            UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
  email         TEXT        NOT NULL UNIQUE,
  password_hash TEXT        NOT NULL,
  name          TEXT        NOT NULL,
  phone         TEXT,
  role          TEXT        NOT NULL CHECK(role IN ('customer', 'restaurant_owner', 'courier', 'admin')),
  created_at    TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS restaurants(
  id          UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id    UUID        NOT NULL REFERENCES users(id),
  name        TEXT        NOT NULL,
  description TEXT,
  address     TEXT        NOT NULL,
  phone       TEXT,
  is_open     BOOLEAN     NOT NULL DEFAULT false,
  created_at  TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses(
  id         UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID        NOT NULL REFERENCES users(id),
  city       TEXT        NOT NULL,
  street     TEXT        NOT NULL,
  house      TEXT        NOT NULL,
  apartment  TEXT,
  comment    TEXT,
  is_default BOOLEAN     NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS orders(
  id               UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id          UUID        NOT NULL REFERENCES users(id),
  restaurant_id    UUID        NOT NULL REFERENCES restaurants(id),
  status           TEXT        NOT NULL CHECK(status IN ('new', 'accepted', 'cooking', 'ready', 'delivering', 'delivered', 'cancelled', 'rejected')),
  total_kopecks    BIGINT      NOT NULL,
  delivery_address TEXT        NOT NULL,
  comment          TEXT,
  created_at       TIMESTAMPTZ NOT NULL,
  updated_at       TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(id) UNIQUE,
  user_id UUID NOT NULL REFERENCES users(id),
  restaurant_id UUID NOT NULL REFERENCES restaurants(id),
  rating INT NOT NULL CHECK(rating BETWEEN 1 AND 5),
  body TEXT,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  order_id UUID REFERENCES orders(id),
  type TEXT NOT NULL,
  body TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS order_status_history(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(id),
  status TEXT NOT NULL,
  changed_by UUID NOT NULL REFERENCES users(id),
  changed_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS carts(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) UNIQUE,
  restaurant_id UUID REFERENCES restaurants(id),
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS dish_categories(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  restaurant_id UUID NOT NULL REFERENCES restaurants(id),
  name TEXT NOT NULL,
  sort_order INTEGER
);

CREATE TABLE IF NOT EXISTS dishes(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  restaurant_id UUID NOT NULL REFERENCES restaurants(id),
  category_id UUID NOT NULL REFERENCES dish_categories(id),
  name TEXT NOT NULL,
  description TEXT,
  price_kopecks BIGINT NOT NULL,
  cooking_time_minutes INTEGER,
  is_available BOOLEAN NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS cart_items(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cart_id UUID NOT NULL REFERENCES carts(id),
  dish_id UUID NOT NULL REFERENCES dishes(id),
  quantity INTEGER NOT NULL CHECK(quantity > 0),
  UNIQUE(cart_id, dish_id)
);

CREATE TABLE IF NOT EXISTS couriers(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) UNIQUE,
  transport_type TEXT NOT NULL,
  status TEXT NOT NULL CHECK(status IN ('offline', 'free', 'busy')) DEFAULT 'free',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(id),
  courier_id UUID NOT NULL REFERENCES couriers(id),
  status TEXT NOT NULL CHECK(status IN ('offered', 'accepted', 'refused', 'picked_up', 'delivered')),
  offered_at TIMESTAMPTZ,
  responded_at TIMESTAMPTZ,
  delivered_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS order_items(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(id),
  dish_id UUID NOT NULL REFERENCES dishes(id),
  dish_name TEXT NOT NULL,
  quantity INTEGER NOT NULL CHECK(quantity > 0),
  price_at_order_time BIGINT NOT NULL
);