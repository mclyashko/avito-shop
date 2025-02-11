# README

## Описание

Разработанный веб-сервис был создан на основе предложенного OpenAPI контракта версии 3.0.0, с учетом некоторых спорных моментов, которые возникли при его интерпретации. Несмотря на эти сложности, сервис был реализован с минимальными отклонениями от исходного контракта, при этом были сделаны несколько допущений, чтобы привести решение в соответствие с практическими требованиями.

### Основные допущения:
1. **Логин пользователя** является первичным ключом в таблице пользователей (`user`), что подразумевает использование его как уникального идентификатора для пользователя. Это решение принято, несмотря на возможную перспективу использования более гибких подходов (например, UUID).
   
2. **Названия предметов** в магазине являются первичными ключами в таблице товаров (`item`). В условиях задачи не было уточнений по поводу возможности изменения названия предметов, поэтому использовался этот подход.
   
3. **Баланс, размер перевода, стоимость предмета** — все эти поля были определены как целые числа, так как в описании API не были указаны более специфичные типы данных (например, `decimal` для денежных значений). Это решение было принято на основе ограничений контракта, что позволяет упростить работу с этими данными на уровне базы данных и API.

4. **Метод /api/buy/{item}** логичнее было бы реализовать как POST запрос, так как он связан с изменением состояния (покупка предмета). Однако, для полного соответствия описанию контракта, было принято решение использовать метод GET, как указано в спецификации OpenAPI.

5. **Транзакции** используются только там, где требуется консистентное изменение нескольких таблиц сразу. Это позволяет обеспечить атомарность операций, таких как перевод монет или покупка предмета, когда необходимо обновить несколько записей в базе данных одновременно, что минимизирует риск ошибок при выполнении операций.

## Миграция

Для обеспечения корректной работы сервиса была разработана миграция базы данных, которая учитывает структуру данных, описанную в предложенном OpenAPI контракте. Миграция включает создание таблиц для пользователей, товаров, покупок, а также истории переводов монет.

Миграционный файл можно найти в репозитории по пути `migrations/init.sql`.

### Особенности миграции:

1. **Таблица `user`** хранит баланс пользователя напрямую, а не вычисляет его на основе снапшотов и дельт. Это решение было принято с целью повышения производительности, особенно для эндпоинта `/api/info`, где требуется оперативный доступ к информации о пользователе, включая его баланс.

2. **Таблица `user_item`** хранит количество купленных предметов пользователем, что позволяет ускорить процесс получения информации о текущем количестве товаров в его наличии. Такой подход исключает необходимость вычислять эту информацию каждый раз при запросе, что также способствует улучшению производительности.

3. **Транзакции о покупке предметов** не отображаются в списке переводов в таблице `coin_transfer`, так как это не предусмотрено контрактом. Контракт описывает переводы монет между пользователями, и покупки товаров, соответственно, не должны фиксироваться в этой таблице. Это соответствует бизнес-логике, где покупка товара рассматривается как отдельная операция, не связанная с переводом монет.

Эти решения были приняты для улучшения производительности и соответствуют требованиям задачи. Они позволяют сервису эффективно обрабатывать запросы при высоких нагрузках.
