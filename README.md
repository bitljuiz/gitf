# gitf
Утилита для сбора статистики о git репозитории. Основана на аналогичной утилите [git-fame](https://pypi.org/project/git-fame/).

```
gitf -p=. -r=HEAD -l='kotlin' -f=tabular
Name          Lines Commits Files
Ilya Astafjev 600   1       15
```

## Доступные флаги:

``-c``, ``--exclude`` ограничивает список обрабатываемых файлов, с помощью [glob-шаблонов](https://ru.wikipedia.org/wiki/%D0%A8%D0%B0%D0%B1%D0%BB%D0%BE%D0%BD_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA%D0%B0).\
``-x``, ``--extensions`` Уменьшает количество обрабатываемых расширений. В качестве входных данных принимается список расширений, разделенный запятыми\
``-f``, ``--format``         Устанавливает формат вывода. Может быть tabular (по умолчанию, показан в примере), [csv](https://ru.wikipedia.org/wiki/CSV), [json](https://ru.wikipedia.org/wiki/JSON), [json-lines](https://jsonlines.org/)\
``-l``, ``--languages``      Уменьшает количество обрабатываемых языков. В качестве входных данных принимается список языков, разделенный запятыми.\
``-o``, ``--order-by``       Устанавливает ключ сортировки результатов. Это могут быть lines (по умолчанию), commits, files \
``-p``, ``--repository``     Возвращает путь к репозиторию Git. Текущий каталог по умолчанию \
``-r``, ``--restrict-to``    Уменьшает количество обрабатываемых файлов, исключая те, которые не удовлетворяют ни одному из glob-шаблонов \
``-v``, ``--revision``       Возвращает указатель на коммит. HEAD по умолчанию \
``-u``, ``--use-committer``  При сборе статистики меняет в коммите author (по умолчанию) на committer

