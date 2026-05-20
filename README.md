 # Estoty Go developer assignment

## Usage
Running the project:
```
docker compose up
```

To see how it works just visit the Nakama Console at `http://127.0.0.1:7351`, default credentials are `admin` and `password`.

To actually check things out we can use the [Api Explorer](http://127.0.0.1:7351/#/api-explorer), where our custom RPC methods should appear in the `Endpoint` field.

### 1. Updating account metadata
To update the metadata of an account, we will need it's ID, the easiest path is to just use `00000000-0000-0000-0000-000000000000`, it's a system user, but it will work.

Don't forget that this method requires a payload in the JSON format (it can be anything).

### 2. Requesting game config
Once again, we will need an account, but this time no payload is needed, you should see the full game config in the response as JSON.

### 3. Private (s2s) call
Private (a.k.a server to server) calls must never contain a user ID, you can including it and the request will fail, if you remove it everything will work just fine (the Console handles authorization automagically), the response will be empty since it's just a ping request.

## Notes
I've deliberately made the code structure as simple as possible, there are countless ways to make it more future-proof and "scalable", but for this specific task it would only make it needlessly complicated.

The first method doesn't have anything special, just make sure the caller is authorized, parse some JSON, make a quick update and that's it.

For the game config call I've decided to keep the data as a Storage object inside Nakama with the system user as the owner. The method itself is, once again, really simple.

Honestly, the most "complicated" part is the game config loading (located in `game_config.go`) simply because a hardcoded config was making me really uncomfortable.

The initial game config is stored in the `game_config.json`, it's being loaded using the `GAME_CONFIG_PATH` env variable (defined in `local.yml`), there's also a default path defined in code.

After locating the file we load the config if it's not yet present in Storage, or if the `version` field in our JSON file is greater than the one in the stored config.

If two versions match we don't do anything, same, if the version in our file is lower, compared to the stored config.

The private RPC method (s2s) is really simple, since the only requirement we have is the lack of a user ID.

## What I would change
- Of course, a flat structure only works in a project that's this small, if we were to add more RPC methods I would definitely add proper structuring.
- Handling the game config better. Right now the version is just a number, of course it's better to use semantic versioning (but it's just overkill for project this small). Also the idea of loading from a JSON file is a little scuffed, I mean, it works, but it's not as clean as I would've liked.
- Right now there's no caching (unless Nakama has something under the hood), just adding a simple Redis service would do the trick, but then again, no real point to overcomplicate things for something this small.
- Build a proper s2s auth layer. From my understanding, Nakama doesn't have first-class support for authorizing individual servers, there's just a single HTTP key that's shared between all servers, it would be nice to not have a single point of failure like that and especially valuable if you want to have different permissions.
- Testing. From what I've learned, testing in Nakama is a little awkward. I could've further separated the code with the sole purpose of covering it with unit tests, but again, at this size it feels really unnecessary. On the other hand, mocking Nakama's interfaces or writing integration tests for http handlers won't be pretty.
