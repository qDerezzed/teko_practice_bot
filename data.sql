DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS companys;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
	user_id integer PRIMARY KEY NOT NULL,
	user_name varchar(32),
	chat_position integer NOT NULL,
	current_company_id integer,
	current_expense_id integer
);
CREATE TABLE companys
(
	company_id bigserial PRIMARY KEY,
	company_name varchar(256) NOT NULL,
	creation_date date NOT NULL,
	fk_user_id integer REFERENCES users(user_id)
);
CREATE TABLE expenses
(
	expense_id bigserial PRIMARY KEY,
	sum integer NOT NULL,
	description text,
    email varchar(64),
    creation_date timestamp NOT NULL,
    fk_company_id integer REFERENCES companys(company_id) NOT NULL,
    fk_user_id integer REFERENCES users(user_id) NOT NULL
);

INSERT INTO users (user_id, chat_position, user_name)
VALUES 
('321', '1', 'artem'),
('5234', '2', 'danil'),
('423', '3', 'kirill');
INSERT INTO companys (company_name, creation_date, fk_user_id)
VALUES 
('Рога и копыта', '2018-08-16', '321'),
('ЧИХ-ПЫХ', '2022-01-16', '5234'),
('ООО ОБНАЛ', '2021-03-22', '423'),
('Я ЗДЕСЬ', '2021-03-22', '423');
INSERT INTO expenses (sum, description, email, creation_date, fk_company_id, fk_user_id)
VALUES
('1000000', '5к быков', 'qwerty@gmail.com', '2021-01-08 04:05:06','1', '321'),
('46283762', '12,5к коров', 'qwerty@gmail.com', NOW(),'1', '321'),
('43224', 'электронный испаритель', 'kuritlubly@gmail.com', '2022-01-12 14:31:00', '2', '5234'),
('100000', 'кальян', 'kuritlubly@gmail.com', NOW(), '2', '5234'),
('1312', 'kto zdes?', 'yazdes@gmail.com', '2021-12-12 00:00:00', '4', '423'),
('1500', 'YA ZDES', 'yazdes@gmail.com', NOW(), '4', '423');

