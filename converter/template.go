package converter

var HtmlTemplate = `<html>
<head><style>%s</style></head>
<body style="">%s</body>
</html>`

var Style = `
.article {
	margin: auto;
	max-width: 800px;
	font-family: Verdana, Candara, Arial, Helvetica, Microsoft YaHei, sans-serif;
    line-height: 1.6;
    padding: 16px 16px;
    overflow-wrap: break-word;
    word-break: break-word;
    font-size: larger;
}

.article img {
    max-width: 80%;
	display: block;
	margin: auto;
}


.article a {
    color: #368CCB;
    text-decoration: none;
}

.article a:hover {
    color: #368CCB;
    text-decoration: none;
}

.article h2,
.article h3,
.article h4,
.article h5,
.article h6 {
    font-weight: 700;
    line-height: 1.5;
    margin: 20px 0 15px;
    margin-block-start: 1em;
    margin-block-end: 0.2em;
}

.article h1 {
    font-size: 1.7em
}

.article h2 {
    font-size: 1.6em
}

.article h3 {
    font-size: 1.45em
}

.article h4 {
    font-size: 1.25em;
}

.article h5 {
    font-size: 1.1em;
}
.article h6 {
    font-size: 1em;
    font-weight: bold
}

@media screen and (max-width: 960px) {
    .article h1 {
        font-size: 1.5em
    }

    .article h2 {
        font-size: 1.35em
    }

    .article h3 {
        font-size: 1.3em
    }

    .article h4 {
        font-size: 1.2em;
    }
}

.article p {
    margin-top: 0;
    margin-bottom: 1.25rem;
}

.article table {
    margin: auto;
    border-collapse: collapse;
    border-spacing: 0;
    vertical-align: middle;
    text-align: left;
    min-width: 66%;
}

.article table td,
.article table th {
    padding: 5px 8px;
    border: 1px solid #bbb;
}

.article q {
    margin-left: 0;
    padding: 0 1em;
    border-left: 5px solid #ddd;
}

.article q::before, .article q::after {
    content: none;
}

.article code {
    color: #343d46;
    padding: .065em .4em;
	font-family: Consolas, 'Courier New', monospace;
}

.article pre {
    overflow-x: auto;
    padding: 8px;
    font-size: 16px;
    margin: 12px 0;
	background-color: #343d46;
}

.article pre code {
	color: white;
}

.article ol {
    text-decoration: none;
    padding-inline-start: 40px;
    margin-bottom: 1.25rem;
}

.article input[type='checkbox'] {
	margin-right: 8px;
	font-size: larger;
}

`
