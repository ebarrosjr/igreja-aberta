# igreja-aberta

![Licença](https://img.shields.io/badge/license-MIT-blue)
![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow)
![Stacks](https://img.shields.io/badge/stacks-multilanguage-blue)

Plataforma de gestão para igrejas e instituições religiosas, desenvolvida como **produto real** e também como **laboratório multi-stack** para comparação arquitetural, estudo técnico e demonstração de capacidade de engenharia de software.

---

## Visão geral

O **igreja-aberta** nasce com dois objetivos principais:

1. **Resolver problemas reais de gestão institucional**, como cadastro de membros, agenda, finanças, eventos, ministérios e comunicação.
2. **Implementar o mesmo domínio de negócio em diferentes stacks**, permitindo comparar linguagens, frameworks, bancos de dados e abordagens arquiteturais.

Mais do que um CRUD, este projeto busca demonstrar:
- modelagem de domínio
- organização de monorepo
- separação de responsabilidades
- arquitetura backend/frontend
- integração com múltiplos bancos
- documentação e padronização
- capacidade de transitar entre diferentes tecnologias

---

## Objetivo do projeto

Criar um sistema flexível e escalável para apoiar a administração de igrejas, congregações e outras instituições religiosas, contemplando rotinas como:

- cadastro de membros
- cadastro de visitantes
- gestão de ministérios e departamentos
- agenda pastoral e eventos
- controle financeiro
- comunicação interna
- relatórios administrativos
- acompanhamento de atividades e participação

Ao mesmo tempo, o projeto servirá como base para implementações paralelas nas tecnologias com as quais tenho experiência, funcionando como portfólio técnico comparativo.

---

## Proposta técnica

Este repositório foi estruturado para abrigar múltiplas implementações do mesmo domínio de negócio.

A ideia central é manter uma **base conceitual única**, com regras de negócio equivalentes, variando a stack utilizada em cada camada.

### Camadas previstas

- **Backend**
  - PHP
  - Java
  - Python
  - .NET

- **Frontend**
  - Angular
  - Flutter
  - HTML/CSS/JS
  - Next.js

- **Banco de dados**
  - MySQL
  - Oracle
  - PostgreSQL
  - SQL Server

---

## Estrutura do repositório

```text
igreja-aberta/
├── backend/
│   ├── dotnet/
│   ├── java/
│   ├── php/
│   └── python/
│
├── database/
│   ├── mysql/
│   ├── oracle/
│   ├── postgresql/
│   └── sqlserver/
│
├── frontend/
│   ├── angular/
│   ├── flutter/
│   ├── html/
│   └── next/
│
├── docs/
│   ├── arquitetura/
│   ├── dominio/
│   ├── diagramas/
│   ├── requisitos/
│   └── roadmap/
│
├── .gitignore
└── README.md
```

## Sobre o autor

Everton Barros Jr — 35 anos de experiência em desenvolvimento de software, atuando com 
arquitetura, legados, migrações e times técnicos. Este projeto representa a consolidação 
de múltiplas stacks dominadas ao longo da carreira, aplicadas a um domínio real.

> "Não sou especialista em uma linguagem. Sou especialista em resolver problemas 
> com a ferramenta certa."

## Contato para oportunidades PJ

- LinkedIn: http://linkedin.com/in/ebarrosjr
- E-mail: ebarrosjr@gmail.com
- Disponibilidade: projetos remotos em todo o mundo ou híbridos no Rio de Janeiro

## Para quem está contratando (PJ)

**O que você leva comigo:**
- 35 anos de experiência comprovada (não teoria, não bootcamp)
- Capacidade de atuar em qualquer stack do projeto (back, front, banco, infra)
- Entendimento de legados, migrações e dívida técnica
- Código que não precisa ser refeito na próxima troca de equipe

**O que não levo:**
- "Fiz na última moda e sumi"
- Banco de dados que cai em produção na sexta
- Preço de júnior (porque entrega de júnior também não levo - e vocês não vão precisar reescrever tudo em 6 meses)

