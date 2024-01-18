# Chat Backend

Offical documentation for chat-backend.

## Endpoints

- `/signup` body:

```json
{
    "username":"dtsf16",
    "password":"test123",
    "email":"test8@gmail.com"
}
```

- `/signin` body:

```json
{
    "password":"test123",
    "email":"test8@gmail.com"
}
```

- `/message` body:

```json
{
    "author_id":"1195205092224008192",
    "channel_id":"1195230073930645504",
    "content":"1st message"
}
```

- `/message/history` body:

```json
{
    "id":"1195205092224008192",
    "channel_id":"1195230073930645504",
    "page":1
}
```

- `/channel` body:

```json
{
    "id":"1195205092224008192"
}
```

- `/channel/list` body:

```json
{
    "id":"1195205092224008192"
}
```

> Note: in the moment, error messages from error.go and controller.Response are not related, so the whole error handling is a mess and does not work properly most of the time.
