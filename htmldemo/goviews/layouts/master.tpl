<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8"/>
        <title>{{.title}}</title>
    </head>
    <body>
        {{include "partials/ad"}}
        {{template "content" .}}
    </body>
</html>