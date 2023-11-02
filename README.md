# YouTrack
## Создавать и настраивать доски
https://comptest.youtrack.cloud/agiles/154-2/current
## Создавать отчеты
https://youtrack.tages.ru/reports
## Настраивать интеграцию YT - gitlab
https://comptest.youtrack.cloud/admin/vcs

# Git
## Уметь работать с историей git'a // TODO показать локально
## Схлопывать в один все коммиты в ветке
    - git merge --squash 
        - при мердже своей ветки в другую и указании флага скуаша, 
          все коммиты из твоей ветки собираются в самый свежи коммит и мерджится, как один коммит
    - git rebase -i HEAD~5 
        - открывает редактор, где есть последние 5 коммитов из ветки. Перед каждым коммитом уже будет стоять выржение pick.
          Это мы даём команду оставить этот коммит. Следом можно выбрать несколько коммитов и pick заменить на squash.
    - git commit --amend
        - При коммите с таким флагом, мы даём команду гиту найти последний коммит, использовать его сообщение и изменить 
          этот коммит локальными изменениями
## git cherry-pick
    * Он позволяет скопировать изменения из одного коммита в свою ветку. В отдельных случая он очень даже полозен, 
      но им нельзя перебарщивать, потому что дублирование коммита, может привести к конфликтам и смотрится в истории 
      гита не очень красиво
    * Синатксис 
        git cherry-pick <commit-hash>
## git submodule 
    - В целом это репозитории внутри репозитория, поэтому все команды гита будут работать и в сабмодуле
    - Добавление сабмодуля
        - git submodule add {ссылка на репозитории} - добавляет всю репу со всеми ветками
        - git submodule add -b {название ветки} {ссылка на репозитории} - добавляет версию репы с указанной ветки
    - Обновление локального сабмодуля удалённым сабмодулем
        - Надо перейти в папу сабмодуля и использовать git pull
    - Пуш локальных изменении сабмодуля
        - Надо перейти в папку сабмодуля и использовать git commit и git push
[.gitmodules](.gitmodules)

## git bare для полной копии
    - Обычно его удобно использовать на сервере, чтобы не хранить все файлы, что может сократить память
    - Синтаксис
        - git clone --bare {ссылка на репозитории} - это создаст .git папку с необходимыми данными, например, конфиг, 
          голову основной ветки и т.д.
# Инфраструктура
## Знать k8s на уровне Namespace/Deployment/Service/Pod
### Namespace
    Это область работы среды окружения
### Deployment
    - помогает:
        - разворачивать окружение
        - откатывать окружение при возникновении ошибки
        - обновление системы без простоя(стратегия RollingUpdate)
    - можно указать:
        - кол-во реплик
        - стратегию
### Service
    - даёт доступ к подам извне
    - занимается балансировкой между подами
### Pod
    - похож на контейнеры в докере, но в поде можно держать несколько контейнеров(в жизни не встречала такого)

## Уметь работать на уровне SELECT/INSERT/UPDATE/DELETE в SQL базах данных
[product.go](internal%2Frepository%2Fproduct.go)

## Уметь работать в RabbitMQ/Kafka на уровне чтения и отправки вручную сообщений
    - показать из локального order
# FulltextSearch
    Используется индекс GIN
    - Механизм похож на B-Tree, только в элементах индекса хранится не ссылка, а набор ссылок
[index_gin_fulltext_Search.sql](..%2F..%2F..%2F..%2FLibrary%2FApplication%20Support%2FJetBrains%2FGoLand2023.2%2Fconsoles%2Fdb%2F8459f2b0-4f4e-45ee-af12-e84b39a58b64%2Findex_gin_fulltext_Search.sql)

# Индексы в postgreSQL
## B-tree
    - Это тип индекса по умолчанию в Postgres
    - Поиск быстрый, но остальные действия замедляются, поэтому много индексов тоже плохо
![img.png](img.png)

## хэш - хеш таблицы
    - Идеальный пример это map в Golang
    - Коллизия - это верояность наслоение записи на уже существующую запись
    - Логика работы:
        - Перед тем, как записать в хеш-таблицу, ключ проходит через хеш-функцию
        - Хеш-функция вычисляет хеш
        - Ищет в таблице по этому хеш есть запись или нет. Если нет, то записывает. Если есть, то возникает коллизия
          Чтобы решить эту коллизию есть два метода: метод открытой адресации и метод цепочек
            - метод открытой адресации
                - при коллизии этот метод ищет следующую пустую ячейку в таблице и записывает туда
            - метод цепочек
                - при коллизии этот метод сохраняет эти данные в отдельную табличку и у каждой записи будет ссылка на 
                  сл. коллизионное значение
                
## GIN
    - Используется для полнотекстового поиска
    - Механизм похож на B-Tree, только в элементах индекса хранится не ссылка, а набор ссылок(документ, строчка и т.д)

# Логирование в pgx
[rest.go](internal%2Fhandler%2Frest.go)

# goroutine
    Goroutine это поток, запускаемый в потоке рантайма
    Особенности:
    - Легковесный - запускаются в потоке самой программы, а не в ОС
    - Отсутствие лимитов на кол-во горутин
    - Конкурентность
[goroutine.go](internal%2Fusecase%2Fgoroutine.go)

## Общение между goroutine-ами
    - Каналы - буферизованные и не буферизованные
        * Не буферизованные
            - Синтаксис 
                chn := make(int chan)
            - Если в канале уже есть запись, то записывающие горутины переходят в ожидание, пока читающие не прочитают из канала
        * Буферизованные
            - Синтаксис 
                chn := make(int chan, 5)
            - Если в канале уже есть 5 записей, то записывающие горутины переходят в ожидание, пока читающие не прочитают из канала
            - Иначе, пока кол-во записей в канале < 5, записывающие горутины могут без блокировки записывать
        * В процессе работы с каналами, может случится, что канал закрылся, а одна горутина всё ещё не закончил записывать.
            В таком случае го запаникует. Чтобы избежать паники, перед записью в канал надо проверять она закрыта или нет:
                chanValue, ok := <- chn
            Если ok = true, то канал ещё открыт и можно записывать, иначе завершаем горутину

    - wait group из пакета sync - нужен для завершения всех горутин до завершения основного
        * Из названия можно понять, что он нужен для синхронизации горутин, а именно группы горутин. 
          При использовании wait group, мы выставляем счётчик равны кол-ву всех горутин, и при завершении любого из них, 
          надо уменьшить счётчик. Таким образом мы даём шанс всем горутинам завершиться
    - mutex (Mutex, RWMutex)
        * Этот пакет позваляет блокировать область памяти(переменную), чтобы остальные горутины не вмешивались  не испортили хронологию
            - Mutex - он блокирует полностью - и для чтения и для записи
            - RWMutex - он блокирует только для записи, а читать из него можно будет
[goroutine.go](internal%2Fusecase%2Fgoroutine.go)

# Работать с композицией типов
    type HttpClient interface {
        Post()
        Get()
    }

    type Create struct {
        client HttpClient
    }

    func (c *Create) Post() { // some code }
    

    type Get struct {
        client HttpClient
    }

    func (g *Get) Get() {}

# ElasticSearch
## заполнять индекс elasticsearch на базе данных из Postgres
    Можно провести условную аналогию: индекс — это база данных, тип — таблица в этой БД, документ - запись в таблице
[elastic.go](internal%2Felastic%2Felastic.go)
## Писать сложные запросы с настройкой релевантности
    - поле _score(по умолчанию сортируются по этому полю)
    - релевантность считается по алгоритму BM25 на основе следующего:
        - Как часто термин появляется в документе — временная частота ( tf )
        - Насколько распространен термин для всех документов — частота обратного документа ( idf )
        - Документы, содержащие все или большинство условий запроса, оцениваются выше, чем документы, который содержат меньше условий
        - Нормализация основана на длине документа, более короткие документы оцениваются лучше, чем более длинные
    - чтобы повысить релевантность можно строить запросы следующим образом:
    {
        "query": {
            "function_score": {
                "query": {
                    "match": {
                        "product_name" : "куртка"
                    }
                },
                "field_value_factor": { -- помогает поднять на вверх продукты с лучшим рейтингом
                    "field": "rate"
                }
            }
        }
    }

# HTTP
## Автодокументация Swagger
[swagger.json](docs%2Fswagger.json)

## github.com/gorilla/mux
[rest.go](internal%2Fhandler%2Frest.go)

# gRPC
## Реализовывать Stream на базе gRPC
[stream.proto](qualification_proto%2Fproto%2Fstream.proto)

# Тестирование
## Пистать Unit тесты
## Проводить функциональное тестирование

# RabbitMQ
## exchange
    - это прослойка между клиентом и очередью
    - занимается роутингом сообщении между очередями
## routing key
    - это ключь по которму определяется в какую очередь положить сообщение
## exchange > queue relations:
### fanout
    - exchange просто пуляет сообщение во все очереди
### direct
    - по routin key он смотрит в какую очередь гнать сообщение
### topic
    - по аргументу или хедеру определяет в какую очередь гнать
## Publisher Confirms
    После паблиша сообщения, паблишер вешает слушатель на канал и ждёт там подтверждения от сервера, что сообщение получили
## RPC pattern
    - Это когда клиент отправляет запрос в очередь и указывает параметр reply_to, где может указать куда ссылать ответ от этого запроса
      Таким образом для отправки ответа не требуется доп. очередь
## DeadLetter
    - Это exchange, куда попадают сообщения если:
        - консюмер признал сообщение негативным(или некорректным) - nack и requeue = false
        - жизнь сообщения закончилось 
        - кол-во элементов в очереди максимальное
    - Консюмер отправляет обратно сообщение с тем же роутинг ки, с которым получил или присваивает новый в headers указав x-dead-letter-routing-key 
# Redis
## Применять inMemory Key/Value
    - основные методы это set <key> <value> и get <key> или keys <pattern>
    - в редисе также можно хранить хешированные данные:
            - hset <key> <field> <value> 
            - hget <key> <field>
            - hgetall <key>

# Логирование
## Работать с логированием в БД/Kafka и т.д.
## Сквозное прокидывание request_id по каждому запросу

# Devops
## Envoy(настройка ratelimit и авторизацию через микросервис)
    * Это прокси сервер, который который занимается балансировкой, маршрутизацией и т.д.
    * Отличия от nginx, которые я нашла:
        - файл конфигурации структурированнее, чем в nginx
        - очень много настроек идут под копотом, когда как в nginx надо отдельные пакеты затягивать
    * показать из order-dev

# Разделять запуск приложение с помощью команд и Обрабатывать флаги запуска
[goroutine.go](cmd%2Fgoroutine.go)
