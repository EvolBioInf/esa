NAME = esa
all: $(NAME)Doc.pdf
	latex $(NAME)Doc
	bibtex $(NAME)Doc
	latex $(NAME)Doc
	latex $(NAME)Doc
	latex $(NAME)Doc
	dvips $(NAME)Doc -o -q
	ps2pdf $(NAME)Doc.ps
$(NAME)Doc.pdf: $(NAME)Doc.tex $(NAME).tex 
$(NAME).tex: ../$(NAME).org
	bash ../scripts/org2nw ../$(NAME).org       | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/' > $(NAME).tex
clean:
	rm -f $(NAME).tex $(NAME)_test.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
