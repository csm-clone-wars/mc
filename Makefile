FPC = fpc
SRC = src/mc.pas
OUT_DIR = build
OUT = mc

.PHONY: all clean

all:
	mkdir -p $(OUT_DIR)
	$(FPC) -O2 -FU$(OUT_DIR) -FE. -o$(OUT) $(SRC)

clean:
	rm -f $(OUT) $(OUT_DIR)/*.o $(OUT_DIR)/*.ppu
