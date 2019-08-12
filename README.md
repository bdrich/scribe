# Scribe

Templater written in Golang based on sprig functions

Inspiriration from https://github.com/tsg/gotpl

## Installation

```bash
make
```

## Getting Started

Ex:

Given there is a file `input.yaml` of the form,

```yaml
bread:
  title: "bread_preferences.md"
  sourdough:
    - name: rye
      prefer: eat
    - name: wheat
      prefer: bake
coffee:
  title: "coffee_tasting.md"
  beans:
    - name: ethiopian
      taste: berries
    - name: columbian
      taste: chocolate
```

and a template in path `mytemplates/{{ .bread.title }}` of the form,

```markdown
# My Sourdough Bread Preferences
{{ range $bread := .bread.sourdough }}
{{- if eq $bread.prefer "eat" }}
- I like to eat {{ $bread.name }} bread.
{{- end }}
{{- if eq $bread.prefer "bake" }}
- I prefer to bake {{ $bread.name }} bread.
{{- end }}
{{- end }}
```

and a template in path `mytemplates/{{ .coffee.title }}` of the form,

```markdown
# My Impressions of Various Coffee
{{ range $coffee := .coffee.beans }}
{{- if eq $coffee.taste "chocolate" }}
- The {{ $coffee.name | title }} coffee tastes of chocolate.
{{- end }}
{{- if eq $coffee.taste "berries" }}
- The {{ $coffee.name | title }} coffee tasts of berries.
{{- end }}
{{- end }}

```

execution of scribe,

```go
scribe -in=input.yaml -templates=mytemplates -out=myoutput
```

produces a file at path `myoutput/bread_preferences.md` with the content

```markdown
# My Sourdough Bread Preferences

- I like to eat rye bread.
- I prefer to bake wheat bread.
```

and produces a file at path `myoutput/coffee_tasting.md` with the content

```markdown
# My Impressions of Various Coffee

- The Ethiopian coffee tasts of berries.
- The Columbian coffee tastes of chocolate.
```
