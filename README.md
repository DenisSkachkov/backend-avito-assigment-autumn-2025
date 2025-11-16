<h1>Тестовое задание для стажёра Backend</h1>

<h2>Задача</h2>
Необходимо реализовать сервис, который назначает ревьюеров на PR из команды автора, позволяет выполнять переназначение ревьюверов и получать список PR’ов, назначенных конкретному пользователю, а также управлять командами и активностью пользователей. После merge PR изменение состава ревьюверов запрещено.

<h2>Запуск приложения и базы данных</h2> 
<pre>make docker-up</pre>

<h2>Остановка приложения и базы данных</h2> 
<pre>make docker-down</pre>

<h2>Описание методов</h2>
<h3>Создание команды с участниками</h3>
<h3>Пример запроса</h3>
<img width="792" height="732" alt="image" src="https://github.com/user-attachments/assets/889068f2-f516-4796-bc39-80b6c1da3139" />
<h3>Пример ответа</h3>
<img width="647" height="767" alt="image" src="https://github.com/user-attachments/assets/0e350e73-d53b-43b3-b987-29a130ed4ff6" />

<h3>Получить команду с участниками</h3>
<h3>Пример запроса</h3>
<img width="853" height="249" alt="image" src="https://github.com/user-attachments/assets/297ac207-443b-4f0c-8c42-03ae677e1ba0" />
<h3>Пример ответа</h3>
<img width="868" height="719" alt="image" src="https://github.com/user-attachments/assets/e8a6aa65-a63b-4a32-848d-dea9bd556f47" />
<h3>Установить флаг активности пользователя</h3>
<h3>Пример запроса</h3>
<img width="855" height="290" alt="image" src="https://github.com/user-attachments/assets/cf878418-8c1e-421c-9bda-aa4e15d9f806" />
<h3>Пример ответа</h3>
<img width="851" height="311" alt="image" src="https://github.com/user-attachments/assets/e629f95f-55aa-49b4-bb4e-acf0848067f1" />

<h3>Получить PR'ы, где пользователь назначен ревьювером</h3>
<h3>Пример запроса</h3>
<img width="852" height="303" alt="image" src="https://github.com/user-attachments/assets/9bbae8be-46ab-4270-971c-7aa79fb36b33" />
<h3>Пример ответа</h3>
<img width="855" height="563" alt="image" src="https://github.com/user-attachments/assets/82103565-e628-4274-913e-2a1426b7aa50" />

<h3>Создать PR и автоматически назначить до 2 ревьюверов из команды автора</h3>
<h3>Пример запроса</h3>
<img width="874" height="287" alt="image" src="https://github.com/user-attachments/assets/4cf5e2af-bb7c-41e2-b987-293de3feacd1" />
<h3>Пример ответа</h3>
<img width="857" height="323" alt="image" src="https://github.com/user-attachments/assets/9b443ba5-7e1e-4068-b85a-13adebe4431a" />

<h3>Пометить PR как MERGED</h3>
<h3>Пример запроса</h3>
<img width="859" height="332" alt="image" src="https://github.com/user-attachments/assets/ddfea658-8f74-45b7-a27e-a16c5de3c600" />
<h3>Пример ответа</h3>
<img width="859" height="411" alt="image" src="https://github.com/user-attachments/assets/2ecf2b5a-9f69-4978-bed3-a9c645c3cf20" />

<h3>Переназначить конкретного ревьювера на другого из его команды</h3>
<h3>Пример запроса</h3>
<img width="869" height="309" alt="image" src="https://github.com/user-attachments/assets/e9244217-f382-4a4b-9d5d-4a60aa091d7f" />
<h3>Пример ответа</h3>
<img width="854" height="340" alt="image" src="https://github.com/user-attachments/assets/be8aba09-b49a-4f06-be59-57ea100ff322" />
