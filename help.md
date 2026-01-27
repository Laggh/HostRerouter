# Ajuda do Host Rerouter

Este aplicativo permite que você modifique facilmente o arquivo hosts do Windows para redirecionar ou bloquear sites.
# 1. O que são overrides

Overrides são regras que interceptam o tráfego de internet para sites específicos. Eles dizem ao seu computador para ir a um endereço IP diferente do real quando você tenta acessar um domínio. Isso é comumente usado para bloquear sites (redirecionando para seu próprio computador `127.0.0.1`) ou para desenvolvimento.

# 2. Criando overrides

Clique no botão **Novo** na barra lateral. Uma nova linha aparecerá na tabela com valores de exemplo (`example.com` apontando para `127.0.0.1`).

# 3. Editando overrides

Você pode editar diretamente as células na tabela:
- **Original**: O domínio que você quer afetar (ex: `youtube.com`).
- **Override**: O endereço IP numérico (IPv4 ou IPv6) para onde ele deve ir (ex: `127.0.0.1` ou `::1`).
- **Enabled**: Use a caixa de seleção para ligar ou desligar a regra temporariamente.

**Importante**: Clique em **Salvar** para aplicar as mudanças. O Windows pode pedir permissão de administrador.

# 4. Apagando overrides

Para remover uma regra permanentemente:
1. Clique no **número da linha** (botão à esquerda da tabela) para selecioná-la. O botão ficará destacado.
2. Clique em **Apagar** na barra lateral.
3. Clique em **Salvar** para atualizar o arquivo do sistema.

# 5. Edição Manual

O botão **Editar** abre uma janela com o conteúdo bruto do arquivo hosts. Use esta opção para:
- Colar grandes blocos de configurações.
- Ver comentários ou linhas ignoradas pelo Rerouter.
- Fazer ajustes finos manualmente.

Ao clicar em Salvar nesta janela, a tabela principal será recarregada automaticamente.

# 6. Botão Descartar/Atualizar

Este botão muda de função dependendo do estado:
- **Descartar**: Aparece quando você fez alterações na tabela mas ainda não salvou. Clicar nele reverte tudo para como estava no arquivo.
- **Atualizar**: Aparece quando não há alterações pendentes. Serve para reler o arquivo do disco (útil se você editou o arquivo por outro programa).

v1.0

