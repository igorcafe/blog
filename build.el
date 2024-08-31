(require 'ox-publish)

(setq org-publish-project-alist
      (list
       (list "igormelo.org"
             :recursive t
             :base-directory "./src"
             :publishing-directory "./public"
             :publishing-function 'org-html-publish-to-html)))

(org-publish-all t)
