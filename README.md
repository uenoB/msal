# Simple MSAL public client

## How to Build

```sh
go build
```

To build without interactive authentication:

```sh
go build -tags=simple
```

## How to Use

1. Create `config.json` file with the following format:
``` js
{
  "clientId": "<public client id>",
  "scopes": ["<scope1>", "<scope2>", /* ... */],
  "redirectURI": "<redirect uri>"
}
```

2. Run `./msal config.json` to get refresh token.
   If succeeded, the token obtained is added to `config.json`.

3. Run `./msal config.json` again to get access token.
   If needed, the token information in `config.json` is updated.
