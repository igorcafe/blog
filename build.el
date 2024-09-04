(require 'ox-publish)

(setq org-export-with-broken-links t
	  org-html-validation-link nil
	  org-html-head-include-scripts nil
	  org-html-head-include-default-style nil
	  org-html-head "<link rel=\"stylesheet\" href=\"https://cdn.simplecss.org/simple.min.css\" />
					 <link rel=\"stylesheet\" href=\"style.css\" />")

(setq navbar-items (mapconcat
					(lambda (entry)
					  (format "<a href=\"%s\">%s</a>" (cdr entry) (car entry)))
					'(("Home" . "/")
					  ("LinkedIn" . "https://www.linkedin.com/in/igoracmelo/")
					  ("GitHub" . "https://github.com/igorcafe"))
					"\n"))

;;(setq navbar (concat "<nav>" navbar-items "</nav>"))
(setq header (concat "<header><nav>" navbar-items "</nav></header>"))

(setq org-publish-project-alist
      (list
       (list "igormelo.org"
             :recursive t
             :base-directory "./src"
             :publishing-directory "./public"
			 :with-author nil
			 :with-title nil
			 :section-numbers nil
			 :with-creator t
			 :with-toc nil
			 :time-stamp-file nil
			 :with-creator nil
			 ;; :auto-sitemap t	   
			 ;; :sitemap-function (lambda (title list)
			 ;; 					 (concat "* Artigos\n"
			 ;; 							 (mapconcat (lambda (link)
			 ;; 										  (format "- [[%s][%s]]"
			 ;; 												  (car link) (cdr link)))
			 ;; 										list "\n")))
			 :html-preamble header
			 ;; :sitemap-title "Igor Melo"
			 ;; :sitemap-filename "index.org"
			 :with-title t
			 :sitemap-sort-files 'anti-chronologically
             :publishing-function 'org-html-publish-to-html)))

(org-publish-all t)
