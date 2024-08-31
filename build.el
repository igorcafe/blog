(require 'ox-publish)

(setq org-export-with-broken-links t
	  org-html-validation-link nil
	  org-html-head-include-scripts nil
	  org-html-head-include-default-style nil
	  org-html-head "<link rel=\"stylesheet\" href=\"https://cdn.simplecss.org/simple.min.css\" />
					 <link rel=\"stylesheet\" href=\"style.css\" />")

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
			 :html-preamble "<nav><a href=\"/\">Inicio</a></nav>"
			 ;; :sitemap-title "Igor Melo"
			 ;; :sitemap-filename "index.org"
			 :with-title t
			 :sitemap-sort-files 'anti-chronologically
             :publishing-function 'org-html-publish-to-html)))

(org-publish-all t)
