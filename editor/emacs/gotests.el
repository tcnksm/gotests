;;; gotests.el --- Integration of the gotests into emacs

;; Copyright (C) 2016 Taichi Nakashima

;; Author: Taichi Nakashima
;; Keywords: go, languages

(defun gotests()
  (interactive)    
  (let ((buf (get-buffer-create "*Gotests patch*"))        
        (coding-system-for-read 'utf-8)
        (coding-system-for-write 'utf-8))
    
    (with-current-buffer buf
      (erase-buffer))

    ;; Only run when buffer is test file
    (if (string-match-p "_test\\.go\\'" (buffer-file-name))
        (progn
          (call-process "gotests" nil buf nil "-w" "-r" (buffer-file-name))
          (revert-buffer :ignore-auto :noconfirm)))))
