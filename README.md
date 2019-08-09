# microbridge

`microbridge` is a WordPress XML-RPC-to-Micropub bridge.

## status

As far as end-to-end XML-RPC-to-Micropub functionality is concerned, at present,
`microbridge` only fully supports:

- Getting the list of categories
- Getting the list of posts

WIP/partial/stubbed support is available for:

- Creating posts
- Editing posts
- Creating categories
- Uploading images/media

## purpose

I use [MarsEdit][marsedit] to create and manage [Micro.blog][microblog] posts.
MarsEdit can speak to Micro.blog via the MetaWeblog API, but MetaWeblog does not
support the concept of drafts. Micro.blog also supports the Micropub API, which
_does_ support drafts, but MarsEdit does not yet understand how to speak
Micropub.

So, until such a time as MarsEdit natively supports Micropub, it can instead
speak to `microbridge`, which implements enough of the WordPress XML-RPC API to
make MarsEdit happy.

[marsedit]: https://www.red-sweater.com/marsedit/
[microblog]: https://micro.blog
