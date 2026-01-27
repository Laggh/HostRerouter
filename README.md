# Host Rerouter

O **Host Rerouter** é uma ferramenta simples e prática com interface gráfica (GUI) desenvolvida para facilitar a edição e gerenciamento do arquivo `hosts` do Windows. Ele permite criar redirecionamentos personalizados (overrides) para domínios específicos de forma visual, sem a necessidade de editar arquivos de sistema manualmente pelo bloco de notas.

## Tecnologias

- **Linguagem**: [Go (Golang)](https://go.dev/)
- **Interface Gráfica**: [Fyne](https://fyne.io/)

## Como funciona o arquivo Hosts?

Para entender o que este programa faz, imagine o arquivo `hosts` do seu computador (`C:\Windows\System32\drivers\etc\hosts`) como uma **lista telefônica** antiga.

Quando você digita um site no navegador (ex: `google.com`), seu computador precisa saber o "número de telefone" (Endereço IP) desse site para se conectar. Antes de perguntar para a internet (servidores DNS), o Windows olha primeiro nessa "lista telefônica" local.

Se você adicionar uma entrada nessa lista dizendo que `facebook.com` é igual a `127.0.0.1` (que é o endereço do seu próprio computador), você efetivamente bloqueia o acesso ao Facebook ou redireciona para onde você quiser.

O **Host Rerouter** automatiza a escrita nessa lista, adicionando tags especiais para gerenciar apenas as entradas criadas por ele, sem bagunçar o restante do seu arquivo.

## Pré-requisitos e Build

Para compilar e rodar este projeto localmente, você precisará configurar o ambiente de desenvolvimento Go e as dependências do Fyne.

### 1. Dependências do Fyne

O Fyne utiliza a GPU para renderização e requer um compilador C instalado no sistema (como o GCC).

**Siga o guia oficial de início rápido do Fyne para instalar os pré-requisitos no seu sistema:**
[https://docs.fyne.io/started/quick/](https://docs.fyne.io/started/quick/)

### 2. Compilando o Projeto

1. Clone este repositório.
2. Abra o terminal na pasta do projeto.
3. Instale as dependências do Go:
   ```bash
   go mod tidy
   ```
4. Para gerar o executável (`.exe`), você pode usar o script `build.bat` incluído ou rodar o comando manualmente:

   ```powershell
   go build -ldflags "-H=windowsgui" -o HostRerouter.exe .
   ```
   *(A flag `-H=windowsgui` serve para esconder a janela de terminal preta ao abrir o programa)*

## Como Usar

**Importante:** O arquivo `hosts` é protegido pelo sistema Windows. Portanto, para que o programa consiga salvar as alterações, você deve **Executá-lo como Administrador**.

1. Clique com o botão direito no `HostRerouter.exe`.
2. Selecione "Executar como administrador".
3. Adicione, edite ou remova seus redirecionamentos e clique em "Salvar".
