# Cryptonews headless-shell

> Оболочка для запроса вебстраниц с защитой cloudflare и подобных

## Сборка и запуск
```bash
docker-compose build
docker-compose down && docker-compose up -d
```

## Использование
```bash
curl 127.0.0.1:4444/html/?url=…&selector=…
```
- `url` - ссылка на веб-страницу
- `selector` - селектор, по которому определяется загрузка нужной страницы.

`selector` необязательное поле. По-умолчанию `body`.
Страницы с защитой сначала загружают код для определения браузера и проверки работы javascript.
Если все пройдено загружается основной контент.

## Подробные логи
```bash
DEBUG=1 docker-compose up -d
```

## Смена порта
```bash
PORT=3333 docker-compose up -d
```