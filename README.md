# Telegram bot shop
## Cli commands
tgbot - запускает бота и оставляет ввод Stdin открытым, 
таким образом написав "stop" вы сделаете "gracefully shutdown",
так же нужно определить 3 переменных окружения TG_BOT_DB, TG_BOT_LOGPATH, 
TG_BOT_TOKEN, url базы данных, путь к логам и токен бота соответсвенно.

tgbot init - инициализирует базу данных

tgbot droptb - удаляет все записи, таблици и свзяи в базе данных

## Telegram bot commands
### for all users
/start - после получения этой комманды пользователь региструется и 
заносится в базу данных

/make_admin key user_tg_id - делает из пользователя с telegram id 
= user_tg_id админа, у которого есть больший доступ к возможностям 
бота(key, определен в переменной, хотя возможно стоили перенести его
в переменные окружения)
### for admins
/add_product name, weight, cost, amount, kind, description - добавляет товар 
в каталог, также если к этом сообщение после прикрепить фотографии, то они будут 
отображаться при просмотре товара

/delete_prod name - удаляет товар
## work example
![Регистрация](https://i.ibb.co/ys43zsL/2023-07-06-17-13-47.png)
![](https://i.ibb.co/Rjj2RkH/2023-07-06-17-35-23.png)
![](https://i.ibb.co/tHzPqFb/2023-07-06-17-41-16.png)
![](https://i.ibb.co/GFpn0c5/2023-07-06-17-42-17.png)
![](https://i.ibb.co/17vGq99/2023-07-06-17-43-30.png)
![](https://i.ibb.co/HhCKZSF/2023-07-06-17-47-45.png)
