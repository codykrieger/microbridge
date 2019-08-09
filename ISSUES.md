# issues

## Micro.blog support for Micropub

Micro.blog's support for Micropub has a number of limitations, which make this
bridge less functional than I'd like it to be. Broken down by action:

### authentication

- many requests made to `/micropub` succeed with 200 OK, even if an invalid
  bearer token (or no bearer token at all) is given in the request

### fetching items

- `GET /micropub?q=source` queries do not support filtering by URL
  (`?q=source&url=...`), so you always receive a list of all of your posts
- `GET /micropub?q=source` queries do not support specifying what properties the
  response should contain (`?q=source&properties=...` /
  `?q=source&properties[]=...&properties[]=...`)
- the `properties` object on items returned from `GET /micropub?q=source`
  queries does not contain a `categories` member, making it impossible to tell
  which categories an item is associated with
- it does not appear possible to retrieve a list of pages; just posts

### editing items

- the `url` property of items is fixed and does not respond to updates, making
  it impossible to change the permalink/slug of a post after it has been created
