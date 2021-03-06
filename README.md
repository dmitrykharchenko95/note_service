# Note service

Note service предоставляет функции хранения заметок в памяти. Хранение заметок каждого пользователя
осуществляется в виде файлов с именем пользователя в определенной директории. Данные хранятся в формате `json`.
Работа сервиса производится через консольный ввод/вывод.

Для компиляции бинаргоного файла, запуска программы и тестов воспользуйтесь Makefile.

## Пакеты

Программа сосотоит из трех пакетов:
* store - отвечает за работу с данными (сохранение, загрузка)
* service - реализует взаимодействие пользользователя с хранилищем (добавление, удаление, показ заметок)
* collector - обеспечивает удаление заметок с прошедшим временем жизни

## Запуск программы

При запуске программа принимает следующие флаги:
* `-period` - период запуска коллектора. С флагом передается временной промежуток, через который будет запускаться
  note_collector: подпрограмма, отвечающая за удаление старых заметок. Дефолтное значение - `10m`.
* `-lifetime` - устанавливает время жизни создаваемых заметок. Дефолтное значение - `24h`. Также время жизни заметки можно
  установить непосредственно при создании конкретной заметки.
* `-dir` - путь к директории, в которой будут храниться данные. При отсутствии директории программа создаст ee
  автоматически. Дефолтное значение - `./notes_data`.

Пример запуска:

```bash
$ note_service -period=5m -lifetime=36h -dir=./my_dir
```


## Использование

При запуске note_service необходимо ввести логин, под которым будут храниться заметки. Для нового пользователя
создастся новый файл в хранилище.

После авторизации note_service принимает следующие команды (не чувствительны к регистру):
* `add [lifetime] <note_data>` - добавление новой заметки с временем жизни `lifetime` (по умолчанию - `24h`, если не
  было изменено при запуске)
* `del <note_id>` - удаление заметки с номером `note_id`
* `get` - показать все заметки пользователя
* `last` - показать последнюю заметку пользователя
* `old` - показать самую старую заметку пользователя
* `out` - вернуться к этапу авторизации

На всех этапах работы note_service можно ввести команду `q`, что приведет к завершению программы.
