#Сервис аккаунтов балансов
____
##Описание
Сервис позволяет управлять аккаунтами балансов пользователей. 
Отношение пользователь:баланс - 1:1. 
Баланс можно создать, удалить, пополнить, списать с него средства. 
Внутренней валютой является количество символов. 
Идея следующая: пользователь создает аккаунт баланса, пополняет его,
затем делает какую-то работу во внешнем сервисе и внешний сервис списывает с
аккаунта пользователя какое-то кол-во символов.
____
##Сборка сервиса
```
make
```