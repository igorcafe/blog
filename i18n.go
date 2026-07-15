package main

import "log"

var i18n = map[string]string{
	"Page not found": "Página não encontrada",
	"Can't find the page you are looking for...": "Não encontrei a página que você procura...",
	"Post not found": "Post não encontrado",
	"Can't find the post you are looking for...": "Não encontrei o post que você procura...",
	"Welcome!": "Bem-vindo!",
	"Here I talk about programming, (good and bad) experiences, Go, Linux, Emacs, and anything else that interests me.": "Aqui eu falo sobre programação, experiências (boas e ruins), Go, Linux, Emacs, e o que der na telha.",
	"Results for filter by tags: ": "Resultados para filtro por tags: ",
	"No post found...":             "Nenhum post encontrado...",
}

func translator(lang string) func(string) string {
	if lang == "en" {
		return func(key string) string {
			return key
		}
	}

	return func(key string) string {
		val, ok := i18n[key]
		if !ok {
			log.Panicf("missing translation: '%s'", key)
		}
		return val
	}
}
