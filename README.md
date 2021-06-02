# Projet
Projet universitaire de l'UE PC3R (Programmation Concurrente réactive et répartie) réalisé en binôme.\
L'objectif était de créer une application web utilisant une API. Nous avons décidé d'utiliser l'API de Deezer afin de créer une application web où il est possible de créer un compte et de s'y connecter, de commenter des musiques et d'aimer les musiques et commentaires, mais aussi aimer les commentaires.\ Le serveur est écrit en servlet en Go à l'aide du package "net/http", la base de donnée est une base de donnée MySQL. La gestion de la connexion s'effectue grâce à l'usage de cookies.

## Exécution
Pour que le site puisse tourner, il faut installer la base de donnée. Nous utilisons une base de donnée mysql.\
Pour connecter la votre, il suffit de modifier les informations de connexions à la bdd du fichier "database/database.go".

Par défaut, nous lançons le serveur sur le port 8080. Pour le modifier il suffit de modifier la variable "port" (ligne 118) de la fonction "main" du fichier "main.go".
