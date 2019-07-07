# gotg_md2html

Having played with telegram bots for a while, I've always hated how inconsistent the markdown support was.
It has never seemed stable, consistent, or supported.
Sending an "unterminated" entity would result in the message failing to send, even if it was an underscore in a url!

I have however had... better experiences with the HTML support, even if HTML5 symbols (such as &apos;) aren't supported.

This is the reason for this project!

I wanted a way to reliably parse messages from basic markdown, to HTML, in a way that telegram would understand.

## How to:

Simply use the MD2HTML function.

``` go
htmlText := tg_md2html.MD2HTML("_hello_ `there` *stranger*! I even support [links [with square brackets!]](github.com)")
```

I am also a fan of url buttons, so I added support for those too - using the following `buttonurl:` syntax:

``` go
htmlText, btnNames, btnLinks := tg_md2html.MD2HTMLButtons("_hello_ [this is a button](buttonurl:link.com)")
```

Simply prepending `buttonurl:` to any link will make the parser detect it as a button, and convert it appropriately.
The function will return two new lists; button names and their respective links, mapped 1:1
