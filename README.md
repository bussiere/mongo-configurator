# mongo-configurator

Provide configuration fo MongoDB

# How to USE

```bash
go get github.com/PxyUp/mongo-configurator
go install github.com/PxyUp/mongo-configurator
```

```bash
mongo-configurator {YML fileName}
```

```bash
mongo-configurator myConfig.yml
```

```yml
databases:
  - urlConnect: mongodb://127.0.0.1:27017/testdb
    name: testdb
    collections:
    - name: User
      indexes:
        - username
        - usernameCanonical
    - name: Users
      indexes:
        - username
        - usernameCanonical
  - urlConnect: mongodb://127.0.0.1:27017/testdbdb
    name: testdbdb
    collections:
    - name: User
      indexes:
        - username
        - usernameCanonical
    - name: Users
      indexes:
        - username
        - usernameCanonical

```
