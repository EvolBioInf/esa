path = ~/WinHome/ownCloud/docs
NAME = esa
all: $(NAME)Doc.pdf
$(NAME)Doc.pdf: $(NAME)Doc.tex $(NAME).tex 
	latex $(NAME)Doc
	bibtex $(NAME)Doc
	latex $(NAME)Doc
	latex $(NAME)Doc
	latex $(NAME)Doc
	dvipdf -dALLOWPSTRANSPARENCY $(NAME)Doc
$(NAME).tex: ../$(NAME).org
	bash ../scripts/org2nw ../$(NAME).org | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/' > $(NAME).tex
publish: $(path)/$(NAME)Doc.pdf
$(path)/$(NAME)Doc.pdf: $(NAME)Doc.pdf
	if [ -d $(path) ]; then \
		cp $(NAME)Doc.pdf $(path); \
	fi

clean:
	rm -f $(NAME).tex $(NAME)_test.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
