# tz_golangdev

Привет! Мы бы хотели, чтобы ты зарегистрировался на GetBlock.io, и используя API GetBlock, попробовал создать сервис, выводящий адрес, баланс которого изменился (в любую сторону) больше остальных за последние сто блоков.

- Получить номер последнего блока можно с помощью следующего метода: https://getblock.io/docs/available-nodes-methods/ETH/JSON-RPC/eth_blockNumber/
- А данные блока вместе с транзакциями через: https://getblock.io/docs/available-nodes-methods/ETH/JSON-RPC/eth_getBlockByNumber/

## Запуск сервера
```
make run
```

## Генерация свагера
```
make swag
```

## Запуск линтера
```
make lint
```
