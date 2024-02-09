.name "pierino"
.description "stay alive"

    fork %2048
    # after exec, two twin process should run
a:	live %-1
    ld %0, r2
	zjmp %:a

