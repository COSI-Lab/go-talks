# Markdown Usage

Talk descriptions support Markdown formatting. This allows you to add links, images, and other formatting to your descriptions. For the nitty details read the source code, but here are some examples.

## Links

Links are automatically detected, but the advanced syntax is also supported.

```markdown
[link text](https://example.com)
```

[link text](https://example.com)

```markdown
<https://example.com>
```

<https://example.com>

```markdown
https://example.com
```

<https://example.com>

Add the `!` prefix to a link to make it an image.

```markdown
![alt text](https://talks.cosi.clarkson.edu/static/construction_tux.png)
```

![alt text](/static/construction_tux.png)

## Formatting

```markdown
**bold**
```

**bold**

```markdown
*italic*
```

*italic*

```markdown
~~strikethrough~~
```

~~strikethrough~~

```markdown
> blockquote
```

> blockquote

```markdown
1/2 3/4 5/6 7/8
```

1/2 3/4 5/6 7/8

## Code

Inline code is with a single backtick.

```markdown
`code`
```

`code`

Code blocks are also supported, without syntax highlighting.

````markdown
```python
print("Hello World!")
```
````

\

```python
print("Hello World!")
```

## Headers

Headers are also supported.

```markdown
# Header 1
## Header 2
### Header 3
#### Header 4
##### Header 5
###### Header 6
```

Here's how they look in context.

![Headers](/static/markdown_headers.png)
