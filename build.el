(require 'ox-publish)

(setq org-publish-project-alist
      (list
       (list "igormelo.org"
             :recursive t
             :base-directory "./src"
             :publishing-directory "./public"
			 :with-author nil
			 :section-numbers nil
			 :with-creator t
			 :with-toc nil
			 :time-stamp-file nil
             :publishing-function 'org-html-publish-to-html)))

(org-publish-all t)
