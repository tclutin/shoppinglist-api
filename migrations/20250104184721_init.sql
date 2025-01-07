-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.users (
    user_id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    gender TEXT NOT NULL CHECK (gender in ('FEMALE', 'MALE', 'NONE')),
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS public.sessions (
    session_id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    refresh_token UUID UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.groups (
    group_id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS public.members (
    member_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('owner', 'member')),
    joined_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES public.users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES public.groups (group_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.categories (
    category_id BIGSERIAL PRIMARY KEY,
    name TEXT NULL
);

CREATE TABLE IF NOT EXISTS public.product_names (
    product_name_id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (category_id) REFERENCES public.categories (category_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.products (
    product_id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    product_name_id BIGINT NOT NULL,
    price decimal,
    status TEXT NOT NULL CHECK (status IN ('open', 'closed')) DEFAULT 'open',
    quantity INT NOT NULL DEFAULT 0,
    added_by BIGINT NOT NULL,
    bought_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (group_id) REFERENCES public.groups(group_id) ON DELETE CASCADE,
    FOREIGN KEY (product_name_id) REFERENCES public.product_names(product_name_id) ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES public.users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (bought_by) REFERENCES public.users(user_id) ON DELETE CASCADE
);

INSERT INTO public.categories (name)
VALUES
    ('молочные продукты'),
    ('мясные продукты'),
    ('рыбные продукты'),
    ('яйцо'),
    ('масложировая продукция'),
    ('хлебобулочные изделия'),
    ('кондитерские изделия'),
    ('продукты пчеловодства'),
    ('бакалейные товары'),
    ('безалкогольные напитки'),
    ('алкогольные напитки'),
    ('табачные изделия'),
    ('плодоовощная продукция'),
    ('прочие продовольственные товары');

INSERT INTO public.product_names (name, category_id)
VALUES
    ('молоко питьевое', 1),
    ('сливки питьевые', 1),
    ('молочный коктейль', 1),
    ('молоко сухое', 1),
    ('сливки сухие', 1),
    ('масло сливочное', 1),
    ('мороженое', 1),
    ('торт', 1),
    ('пирожное', 1),
    ('сгущенное молоко', 1),
    ('детское питание жидкое', 1),
    ('детское питание сухое', 1),
    ('кефир', 1),
    ('простокваша', 1),
    ('ряженка', 1),
    ('йогурт', 1),
    ('сметана', 1),
    ('творог', 1),
    ('творожный сыр', 1),
    ('сырки творожные глазированные', 1),
    ('вареники с творогом', 1),
    ('сырники', 1),
    ('сыр плавленый', 1),
    ('сыр мягкий', 1),
    ('сыр твердый', 1),
    ('говядина', 2),
    ('телятина', 2),
    ('свинина', 2),
    ('баранина', 2),
    ('пельмени', 2),
    ('колбаса вареная', 2),
    ('колбаса копченая', 2),
    ('сосиски', 2),
    ('мясные консервы', 2),
    ('курица', 2),
    ('цыпленок', 2),
    ('утка', 2),
    ('гусь', 2),
    ('индейка', 2),
    ('камбала', 3),
    ('скумбрия', 3),
    ('треска', 3),
    ('окунь', 3),
    ('тунец', 3),
    ('осетр', 3),
    ('лосось', 3),
    ('сельдь', 3),
    ('икра', 3),
    ('краб', 3),
    ('креветки', 3),
    ('рак', 3),
    ('мидии', 3),
    ('устрицы', 3),
    ('яйцо куриное', 4),
    ('яйцо перепелиное', 4),
    ('масло подсолнечное', 5),
    ('масло оливковое', 5),
    ('майонез', 5),
    ('кетчуп', 5),
    ('соус', 5),
    ('хлеб', 6),
    ('батон', 6),
    ('багет', 6),
    ('булка', 6),
    ('ватрушка', 6),
    ('слойка', 6),
    ('крендель', 6),
    ('плюшка', 6),
    ('рулет', 6),
    ('сухари', 6),
    ('баранки', 6),
    ('бублики', 6),
    ('соломка', 6),
    ('пирог', 6),
    ('пирожок', 6),
    ('пончик', 6),
    ('лаваш', 6),
    ('конфеты шоколадные', 7),
    ('мармелад', 7),
    ('ирис', 7),
    ('драже', 7),
    ('карамель', 7),
    ('пастила', 7),
    ('шоколад', 7),
    ('халва', 7),
    ('жевательная резинка', 7),
    ('сахарная вата', 7),
    ('печенье', 7),
    ('вафли', 7),
    ('пряники', 7),
    ('кексы', 7),
    ('торт бисквитный', 7),
    ('мед пчелиный', 8),
    ('мед в сотах', 8),
    ('мед искуственный', 8),
    ('пчелиная перга', 8),
    ('обножка', 8),
    ('маточное молочко', 8),
    ('пчелиный яд', 8),
    ('прополис', 8),
    ('мука', 9),
    ('рис', 9),
    ('манка', 9),
    ('овсянка', 9),
    ('пшено', 9),
    ('гречка', 9),
    ('горох', 9),
    ('фасоль', 9),
    ('бобы', 9),
    ('каша', 9),
    ('макароны', 9),
    ('соломка', 9),
    ('рожки', 9),
    ('вермишель', 9),
    ('чай зеленый', 9),
    ('чай черный', 9),
    ('чай красный', 9),
    ('чай желтый', 9),
    ('чай травяной', 9),
    ('кофе жареный в зернах', 9),
    ('кофе молотый', 9),
    ('кофе без добавок', 9),
    ('кофе растворимый', 9),
    ('какао', 9),
    ('специи', 9),
    ('уксус', 9),
    ('лимонная кислота', 9),
    ('соль', 9),
    ('сахар', 9),
    ('сахарная пудра', 9),
    ('кисель', 9),
    ('желе', 9),
    ('чипсы', 9),
    ('крахмал', 9),
    ('дрожжи', 9),
    ('желатин пищевой', 9),
    ('сода пищевая', 9),
    ('вода', 10),
    ('вода минеральная', 10),
    ('квас', 10),
    ('сок', 10),
    ('лимонад', 10),
    ('газировка', 10),
    ('водка', 11),
    ('ликер', 11),
    ('ром', 11),
    ('виски', 11),
    ('вино', 11),
    ('шампанское', 11),
    ('коньяк', 11),
    ('сидр', 11),
    ('пиво', 11),
    ('сигареты', 12),
    ('папиросы', 12),
    ('сигары', 12),
    ('табак', 12),
    ('картофель', 13),
    ('капуста белокочанная', 13),
    ('капуста краснокочанная', 13),
    ('капуста цветная', 13),
    ('брокколи', 13),
    ('морковь', 13),
    ('свекла', 13),
    ('редис', 13),
    ('редька', 13),
    ('лук репчатый', 13),
    ('чеснок', 13),
    ('огурец', 13),
    ('тыква', 13),
    ('кабачок', 13),
    ('помидор', 13),
    ('баклажан', 13),
    ('перец', 13),
    ('салат', 13),
    ('укроп', 13),
    ('петрушка', 13),
    ('сельдерей', 13),
    ('яблоко', 13),
    ('груша', 13),
    ('вишня', 13),
    ('черешня', 13),
    ('слива', 13),
    ('персик', 13),
    ('нектарин', 13),
    ('абрикос', 13),
    ('земляника', 13),
    ('клубника', 13),
    ('смородина', 13),
    ('малина', 13),
    ('крыжовник', 13),
    ('виноград', 13),
    ('апельсин', 13),
    ('лимон', 13),
    ('лайм', 13),
    ('мандарин', 13),
    ('грейпфрут', 13),
    ('ананас', 13),
    ('авокадо', 13),
    ('банан', 13),
    ('манго', 13),
    ('оливки', 13),
    ('хурма', 13),
    ('гранат', 13),
    ('киви', 13),
    ('брусника', 13),
    ('клюква', 13),
    ('смородина', 13),
    ('сухофрукты', 13),
    ('фундук', 13),
    ('грецкий орех', 13),
    ('миндаль', 13),
    ('арахис', 13),
    ('фисташки', 13),
    ('кедровый орех', 13),
    ('кешью', 13),
    ('семечки', 13),
    ('грибы', 13),
    ('мясные бульонные кубики', 14),
    ('горчица', 14),
    ('хрен', 14),
    ('соевый соус', 14),
    ('рыбий жир', 14),
    ('кошачий корм жидкий', 14),
    ('кошачий корм сухой', 14),
    ('собачий корм жидкий', 14),
    ('собачий корм сухой', 14);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS public.products;
DROP TABLE IF EXISTS public.product_names;
DROP TABLE IF EXISTS public.categories;
DROP TABLE IF EXISTS public.members;
DROP TABLE IF EXISTS public.sessions;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.groups;
