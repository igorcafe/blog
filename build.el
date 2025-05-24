(require 'ox-publish)
(require 'org-roam-export)
;; (package-install 'htmlize)
;; (require 'htmlize)

(defun org-export-index-page ()
  (let ((posts (org-roam-ql-nodes
                '(tags "post")))
        (concepts (org-roam-ql-nodes
                   '(and (not (title "Blog | Igor Melo" t))
                         (not (tags "post")))
                   "title")))

    (save-excursion
      (goto-char (point-max))
      (insert "* Posts\n")
      (if (null posts)
          (insert "\nNo posts yet.")
        (mapcar (lambda (node)
                  (insert
                   (format "- [[id:%s][%s]]\n"
                           (org-roam-node-id node)
                           (org-roam-node-title node))))
                posts))
      (insert "\n")

      (insert "* Concepts\n")
      (if (null concepts)
          (insert "\nNothing here yet.")
        (mapcar (lambda (node)
                  (insert
                   (format "- [[id:%s][%s]]\n"
                           (org-roam-node-id node)
                           (org-roam-node-title node))))
                concepts))
      (insert "\n")
      )))

(defun org-export-insert-node-links ()
  (let* ((node (with-current-buffer (replace-regexp-in-string "<[0-9]+>$" "" (buffer-name))
                 (org-roam-node-at-point)))
         (links (and node (org-roam-ql-nodes `(backlink-from ,node))))
         (backlinks (and node (org-roam-ql-nodes `(backlink-to ,node)))))
    (when node
      (save-excursion
        (goto-char (point-max))
        (insert "* Links\n")
        (if (null links)
            (insert "\nNo links.")
          (mapcar (lambda (node)
                    (insert
                     (format "- [[id:%s][%s]]\n"
                             (org-roam-node-id node)
                             (org-roam-node-title node))))
                  links))
        (insert "\n")

        (insert "* Backlinks\n")
        (if (null backlinks)
            (insert "\nNo backlinks.")
          (mapcar (lambda (node)
                    (insert
                     (format "- [[id:%s][%s]]\n"
                             (org-roam-node-id node)
                             (org-roam-node-title node))))
                  backlinks))))))

(defun org-before-processing (&rest args)
  (if (string-prefix-p "index.org" (buffer-name))
      (org-export-index-page)
    (org-export-insert-node-links)))

(add-hook 'org-export-before-processing-hook 'org-before-processing)

(setq org-roam-db-location "~/Git/blog/org-roam.db")
(setq org-roam-directory "~/Git/blog/")

(setq org-export-with-broken-links t
      org-html-validation-link nil
      org-html-head-include-scripts nil
      org-html-head-include-default-style nil
      org-html-head nil
      org-html-head "
<link rel=\"stylesheet\" href=\"https://cdn.simplecss.org/simple.min.css\" />
<link rel=\"stylesheet\" href=\"style.css\" />"
      )

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
             :recursive nil
             :base-directory "."
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
             ;;            (concat "* Artigos\n"
             ;;                (mapconcat (lambda (link)
             ;;                       (format "- [[%s][%s]]"
             ;;                           (car link) (cdr link)))
             ;;                     list "\n")))
             :html-preamble header
             ;; :sitemap-title "Igor Melo"
             ;; :sitemap-filename "index.org"
             :with-title t
             :sitemap-sort-files 'anti-chronologically
             :publishing-function 'org-html-publish-to-html)))

(org-publish-all t)
