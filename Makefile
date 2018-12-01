
DIRS = poly \
       goom \
       plots \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./cmd/$$dir $@; \
	done

format:
	goimports -w .

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./cmd/$$dir $@; \
	done
