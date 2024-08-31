(require 'ox-publish)

(setq org-html-validation-link nil            ;; Don't show validation link
      org-html-head-include-scripts nil       ;; Use our own scripts
      org-html-head-include-default-style nil ;; Use our own styles
      org-html-head "<link rel=\"stylesheet\" href=\"https://cdn.simplecss.org/simple.min.css\" />")

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
             :publishing-function 'org-html-publish-to-html)))

(org-publish-all t)
