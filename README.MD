

### Настройка Buf
[buf](https://github.com/bufbuild/buf) это protobuf линтер, детектор критических изменений, который обеспечивает совместимость на уровне исходного кода или проводного интерфейса (и не только). 
Buf можно скачать со [страницы релизов](https://github.com/bufbuild/buf/releases) или воспользоваться предложенными способами установки https://buf.build/docs/cli/installation. 

пример конфиг файла (должен лежать в корне проекта)
```yaml
version: v2

deps:
  - buf.build/googleapis/googleapis

lint:
  use:
    - DEFAULT  # базовый набор линтинга: стиль, порядок и т.п.

breaking:
  use:
    - FILE      # проверяет, что нельзя удалять/ломать файлы
    - PACKAGE   # проверяет изменения в package
    - WIRE_JSON # проверяет wire-совместимость (номера, типы и т.д.)
```

если у вас в прото используются зависимости
```yaml
import "google/api/annotations.proto";
```
и buf ругается 
`"google/api/annotations.proto": file does not exist` нужно зависимости добавить 
в конфиг `deps` и выполнить 
```
buf dep update
```
должен создасться файл `buf.lock` с вашими зависимостями.
Команде `buf dep update` может выдать 403 ошибку из-за [блокировок по ip](https://buf.build/docs/bsr/authentication/#ip-address-blocking) 
в таком случае нужен прокси, в powershell установить прокси можно так
```
$env:HTTPS_PROXY="http://localhost:8081"
```
после чего запрос `buf dep update` пойдет через прокси.




#### примеры команд:
запуск линтера
```
buf lint 
```

запуск проверки на критические изменения (работает с git)
```
buf breaking --against '.git#branch=master'
```

### protolint
[protolint](https://github.com/yoheimuta/protolint) еще один интересный линтер для protobuf и он куда проще в настройке чем buf.
поддерживает флаг `-fix` для автоматического исправления замечаний по стилю. 

пример:
```
protolint lint -fix .\proto
```
можно запускать из докера. Детали см в README репазитория protolint