
# Nome da imagem Docker
IMAGE = n64go

# Nome do container (já rodando)
ENV_NAME = n64suply

# Caminhos
WORKDIR = /go64
BUILD_DIR = build

# Arquivos
ELF = $(BUILD_DIR)/n64go.elf
ROM = $(BUILD_DIR)/n64go.z64

#  ------------------------------------------------
#  ____        _ _     _
# | __ ) _   _(_) | __| |
# |  _ \| | | | | |/ _` |
# | |_) | |_| | | | (_| |
# |____/ \__,_|_|_|\__,_|

# --- Alvo padrão
all: build rom

build:
	mkdir -p $(BUILD_DIR)
	docker run --rm \
		--platform linux/amd64 \
		-v $(PWD):$(WORKDIR) \
		-w $(WORKDIR) \
		$(IMAGE) \
		sh -c "go build -o $(ELF) ."

rom: build
	docker run --rm \
		--platform linux/amd64 \
		-v $(PWD):$(WORKDIR) \
		-w $(WORKDIR) \
		$(IMAGE) \
		sh -c "n64go rom $(ELF)"

#  ------------------------------------------------
#  ____
# |  _ \  _____   __
# | | | |/ _ \ \ / /
# | |_| |  __/\ V /
# |____/ \___| \_/

build-devenv:
	docker build -t $(IMAGE) .

# Sobe o container e deixa pronto
start-devenv:
	docker run -d --rm \
		--platform linux/amd64 \
		-v $(PWD):$(WORKDIR) \
		-w $(WORKDIR) \
		--name $(ENV_NAME) \
		$(IMAGE) sleep infinity 


stop-devenv:
	docker stop $(ENV_NAME)
# Gera build dentro do container ativo
build-dev:
	mkdir -p $(BUILD_DIR)
	docker exec -w $(WORKDIR) $(ENV_NAME) sh -c \
		"go build -o $(ELF) ."

# Gera ROM dentro do container ativo
rom-dev: 
	docker exec -w $(WORKDIR) $(ENV_NAME) sh -c \
		"n64go rom $(ELF)"

build-n-rom-dev: build-dev rom-dev
# --- CLEAN ---------------------------------------------------------

clean:
	rm -rf $(BUILD_DIR)

.PHONY: all build rom clean start-devenv build-local rom-local
