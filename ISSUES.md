# issues

## Micro.blog support for Micropub

Micro.blog's support for Micropub has a number of limitations, which make this
bridge less functional than I'd like it to be. Broken down by action:

### authentication

- Many requests made to `/micropub` succeed with 200 OK, even if an invalid
  bearer token (or no bearer token at all) is given in the request. This makes
  it challenging to know whether the user's credentials are valid or not. I'm
  working around this by checking for the presence of the `destination` property
  on the config object (`/micropub?q=config`).

### fetching items

- `GET /micropub?q=source` queries do not support filtering by URL
  (`?q=source&url=...`), so you always receive a list of all of your posts.
- `GET /micropub?q=source` queries do not support specifying what properties the
  response should contain (`?q=source&properties=...` /
  `?q=source&properties[]=...&properties[]=...`).
- The `properties` object on items returned from `GET /micropub?q=source`
  queries does not contain a `categories` member, making it impossible to tell
  which categories an item is associated with.
- It does not appear possible to retrieve a list of pages; just posts.
- Item URLs always seem to be prefixed with `http://` rather than `https://` as
  I'd expect, given that `config.destination[0].uid` correctly specifies
  `https://`.

### editing items

- The `url` property of items is fixed and does not respond to updates, making
  it impossible to change the permalink/slug of a post after it has been
  created.
