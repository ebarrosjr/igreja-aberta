# API em C# .Net Core

Estamos usando o .Net 10, então ele precisa estar instalado em sua máquina para poder executar o backend em .Net

Comando para rodar:

```bash
 git clone git@github.com:ebarrosjr/igreja-aberta.git
 cd igreja-aberta
 // Aqui considero que tenha o MySQL (usado neste exemplo) instalado e rodando
 mysql -uroot -p -e "CREATE SCHEMA igreja_aberta CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
 mysql -uroot -p igreja_aberta < database/mysql/initial_schema.sql
 cd backend/php/cakephp
composer install
bin/cake server -p 8765
```

Sua api estará funcionando em `http://localhost:8765/api`