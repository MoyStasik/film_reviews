DROP TABLE IF EXISTS film_category;
DROP TABLE IF EXISTS review;
DROP TABLE IF EXISTS film;


CREATE TYPE film_review_value AS ENUM ('1', '2', '3', '4', '5');
CREATE TYPE category AS ENUM ('боевик', 'ужасы', 'комедия', 'приключения', 'фантастика');

CREATE TABLE IF NOT EXISTS film_category(
 id int primary key NOT NULL,
 category_name category NOT NULL
);

CREATE TABLE IF NOT EXISTS film (
 id int NOT NULL GENERATED ALWAYS AS,
 film_category_id int not NULL,
 film_name text NOT NULL,
 img_path text NOT NULL,
 description text NOT NULL,
 reviews_value float,
 FOREIGN KEY (film_category_id) REFERENCES film_category (id)
);

CREATE TABLE IF NOT EXISTS review (
 id int primary key NOT NULL,
 user_id int NOT NULL,
 film_id int NOT NULL,
 review_content text NOT NULL,
 review_value film_review_value NOT NULL,
 FOREIGN KEY (film_id) REFERENCES film (id)
);

ALTER TABLE film_category ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY(
	START WITH 6
	INCREMENT BY 1
);

ALTER TABLE film ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY(
	START WITH 6
	INCREMENT BY 1
);

ALTER TABLE review ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY(
	START WITH 5
	INCREMENT BY 1
);

INSERT INTO film_category VALUES (
 1, 'боевик'
);
INSERT INTO film_category VALUES (
 2, 'ужасы'
);
INSERT INTO film_category VALUES (
 3, 'комедия'
);
INSERT INTO film_category VALUES (
 4, 'фантастика'
);
INSERT INTO film_category VALUES (
 5, 'приключения'
);

INSERT INTO film VALUES (
 1, 1, 'Аватар', '/img/avatar.jpg', 'Бывший морпех Джейк Салли прикован к инвалидному креслу. Несмотря на немощное тело, Джейк в душе по-прежнему остается воином. Он получает задание совершить путешествие в несколько световых лет к базе землян на планете Пандора, где корпорации добывают редкий минерал, имеющий огромное значение для выхода Земли из энергетического кризиса.',
 0
);
INSERT INTO film VALUES (
 2, 2, 'Остров проклятых', '/img/OstrovProklyatih.jpg', 'Два американских судебных пристава отправляются на один из островов в штате Массачусетс, чтобы расследовать исчезновение пациентки клиники для умалишенных преступников. При проведении расследования им придется столкнуться с паутиной лжи, обрушившимся ураганом и смертельным бунтом обитателей клиники.',
 0
);
INSERT INTO film VALUES (
 3, 3, 'Великий Гэтсби', '/img/Gatsby.jpg', 'Бывший морпех Джейк Салли прикован к инвалидному креслу. Несмотря на немощное тело, Джейк в душе по-прежнему остается воином. Он получает задание совершить путешествие в несколько световых лет к базе землян на планете Пандора, где корпорации добывают редкий минерал, имеющий огромное значение для выхода Земли из энергетического кризиса.',
 0
);
INSERT INTO film VALUES (
 4, 1, 'Начало', '/img/Inception.jpg', 'Кобб – талантливый вор, лучший из лучших в опасном искусстве извлечения: он крадет ценные секреты из глубин подсознания во время сна, когда человеческий разум наиболее уязвим. Редкие способности Кобба сделали его ценным игроком в привычном к предательству мире промышленного шпионажа, но они же превратили его в извечного беглеца и лишили всего, что он когда-либо любил.

И вот у Кобба появляется шанс исправить ошибки. Его последнее дело может вернуть все назад, но для этого ему нужно совершить невозможное – инициацию. Вместо идеальной кражи Кобб и его команда спецов должны будут провернуть обратное. Теперь их задача – не украсть идею, а внедрить ее. Если у них получится, это и станет идеальным преступлением.

Но никакое планирование или мастерство не могут подготовить команду к встрече с опасным противником, который, кажется, предугадывает каждый их ход. Врагом, увидеть которого мог бы лишь Кобб.',
 0
);
INSERT INTO film VALUES (
 5, 4, 'Побег из Шоушенка', '/img/shoushenk.jpg', 'Бухгалтер Энди Дюфрейн обвинён в убийстве собственной жены и её любовника. Оказавшись в тюрьме под названием Шоушенк, он сталкивается с жестокостью и беззаконием, царящими по обе стороны решётки. Каждый, кто попадает в эти стены, становится их рабом до конца жизни. Но Энди, обладающий живым умом и доброй душой, находит подход как к заключённым, так и к охранникам, добиваясь их особого к себе расположения.',
 0
);

INSERT INTO review VALUES (
 1, 4, 1, 'Хороший фильм, рекомендую!', '5'
);
INSERT INTO review VALUES (
 2, 3, 1, 'Крутой фильм, посмотрите обязательно!!!', '5'
);
INSERT INTO review VALUES (
 3, 2, 1, 'Хороший фильм, но есть косяки!', '4'
);
INSERT INTO review VALUES (
 4, 4, 2, 'Так круто, я сразу даже не понял, что произошло.', '5'
);