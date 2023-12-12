# Aplicativo de Monitoramento de Ativos

## Introdução

Este repositório contém o código-fonte para o aplicativo de gerenciamento de ativos. Este aplicativo é projetado para monitorar e registrar informações críticas sobre os ativos de TI de uma empresa, incluindo o status das portas abertas, endereços IP, aplicativos instalados e, eventualmente, vulnerabilidades nos servidores, containers docker.

## Características

- **Monitoramento de Portas Abertas:** O aplicativo verifica continuamente as portas de rede em todos os servidores para identificar e registrar quaisquer portas abertas.
- **Rastreamento de IPs e Serviços:** Ele mantém um registro de todos os endereços IP dentro da rede da empresa, junto com os serviços que estão sendo executados em cada um.
- **Inventário de Aplicativos:** O aplicativo cataloga todos os softwares instalados nos servidores e dispositivos da rede, facilitando o gerenciamento de ativos de software.
- **Verificação de containers docker (Futuro):** O aplicativo identifica containers em execução, listando suas portas e imagens utilizadas.
- **Detecção de Vulnerabilidades (Futuro):** Em uma versão futura, planejamos implementar uma funcionalidade para identificar potenciais vulnerabilidades de segurança nos servidores.


## Uso

Após a instalação e configuração, o aplicativo começará automaticamente a monitorar e registrar os dados dos ativos de rede. A interface do usuário permite que os administradores visualizem e analisem esses dados, identifiquem problemas potenciais e planejem manutenções ou upgrades conforme necessário.
Atualmente a api central encontra-se em https://github.com/llucasramao/activeNowManager
